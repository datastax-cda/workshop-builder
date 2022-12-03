package build

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"workshop-builder/util"

	"github.com/gohugoio/hugo/commands"
	cp "github.com/otiai10/copy"
)

var languages = [...]string{"en", "es", "fr", "pt"}

func BuildCmd() {

	config, err := util.DetermineConfig("config.json")
	if err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Cleaning up existing workshopGen...")
	if err := os.RemoveAll("workshopGen/"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}
	fmt.Println("Setting up base theme...")
	if err := util.CloneRepo("https://github.com/datastax-cda/workshop-base", "workshopGen"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}
	if err := util.RemoveGitMetadata("workshopGen"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	if err := setWorkshopTitle(config); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	if err := setWorkshopContent(config); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Building Static Website Content in /publicGen ...")
	_ = os.RemoveAll("publicGen/")
	runtime.GOMAXPROCS(runtime.NumCPU())
	resp := commands.Execute([]string{"-s", "workshopGen/", "-d", "../publicGen"})

	if resp.Err != nil {
		if resp.IsUserError() {
			resp.Cmd.Println("")
			resp.Cmd.Println(resp.Cmd.UsageString())
		}
		os.Exit(-1)
	}

	fmt.Println("Copying Staticfile.auth to /publicGen ...")
	if err := copyStaticfileAuth(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

}

func setWorkshopContent(config *util.WorkshopConfig) error {
	if _, err := os.Stat("paceWorkshopContent"); os.IsNotExist(err) {
		if err := util.CloneRepo("https://github.com/datastax-cda/workshop-content", "paceWorkshopContent"); err != nil {
			return err
		}
		fmt.Println("Adjusting workshop content locally can be done within the paceWorkshopContent folder. Once your content is ready to be shared with your fellow team members, commit it back to pace workshop content! ")
	}

	for _, module := range config.Modules {
		if (strings.Compare(module.Type, "concepts")) == 0 {
			if err := setWorkshopConcepts(module.Content); err != nil {
				return err
			}
		} else if module.Type == "demos" {
			if err := setWorkshopDemos(module.Content); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("config contains a module (%s) that is not of type demos or concepts. This is not allowed", module.Type)
		}
	}
	return nil
}

func setWorkshopDemos(contents []util.ContentConfig) error {
	for order, content := range contents {
		err := setWorkshopExtras(content, "demos")
		if err != nil {
			return err
		}
		for _, language := range languages {
			fileName := strings.Split(content.Filename, "/")
			pageFile := "workshopGen/content/demos/" + fileName[len(fileName)-1] + "." + language + ".md"
			err := createPage(pageFile, content.Name, order)

			if err != nil {
				return err
			}

			contentPath := "paceWorkshopContent/" + content.Filename
			err = addMarkdown(pageFile, contentPath+"."+language+".md", language)
			if err != nil {
				fmt.Printf("cannot add specified demo markdown to file, %s, %+v", fileName[len(fileName)-1]+"."+language+".md", err)
				return err
			}
		}
	}
	return nil
}

func setWorkshopConcepts(contents []util.ContentConfig) error {
	for order, content := range contents {
		err := setWorkshopExtras(content, "concepts")
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, language := range languages {
			fileName := strings.Split(content.Filename, "/")
			pageFile := "workshopGen/content/concepts/" + fileName[len(fileName)-1] + "." + language + ".md"
			err := createPage(pageFile, content.Name, order)

			if err != nil {
				return err
			}

			contentPath := "paceWorkshopContent/" + content.Filename
			err = addMarkdown(pageFile, contentPath+"."+language+".md", language)
			if err != nil {
				fmt.Printf("cannot add specified content markdown to file, %s, %+v", fileName[len(fileName)-1]+"."+language+".md", err)
			}
		}
	}
	return nil
}

func setWorkshopExtras(curContent util.ContentConfig, contType string) error {

	var (
		destination string
		source      string
	)

	contentPath := strings.Split(curContent.Filename, "/")
	folders := contentPath[:len(contentPath)-1]
	folderPath := strings.Join(folders, "/")

	source = "paceWorkshopContent/" + folderPath + "/"

	if contType == "demos" {
		destination = "workshopGen/content/demos/" + contentPath[len(contentPath)-1] + "/"
		_ = os.MkdirAll(destination, os.FileMode(0777))
	} else if contType == "concepts" {
		destination = "workshopGen/content/concepts/" + contentPath[len(contentPath)-1] + "/"
		_ = os.MkdirAll(destination, os.FileMode(0777))
	} else {
		return fmt.Errorf("%s content is not of demos or concepts types", contType)
	}

	fds, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(source, fd.Name())
		dstfp := path.Join(destination, fd.Name())

		if !fd.IsDir() {
			if filepath.Ext(strings.TrimSpace(fd.Name())) != ".md" {

				srcfd, err := os.Open(srcfp)
				if err != nil {
					return err
				}
				defer srcfd.Close()

				dstfd, err := os.Create(dstfp)
				if err != nil {
					return err
				}
				defer dstfd.Close()

				if _, err = io.Copy(dstfd, srcfd); err != nil {
					return err
				}
			}
		} else {
			err := cp.Copy(srcfp, dstfp)
			fmt.Println(err)
		}
	}

	return nil
}

func addMarkdown(existingFile string, additionalMarkDown string, lang string) error {
	additionalMarkDownWriter, err := os.Open(additionalMarkDown)
	if err != nil {
		if lang == "en" {
			fmt.Printf("%s not found!\n", additionalMarkDown)
		}
		os.Remove(existingFile)
		return nil
	}
	defer additionalMarkDownWriter.Close()
	existingFileWriter, err := os.OpenFile(existingFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing %s", err)
	}
	defer existingFileWriter.Close()
	_, err = io.Copy(existingFileWriter, additionalMarkDownWriter)
	if err != nil {
		log.Fatalln("failed to append files:", err)
	}

	return nil
}

func createPage(file string, title string, order int) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("cannot create file, %s, %+v", file, err)
	}
	order = order + 3
	//header := fmt.Sprintf("+++\ntitle = \"\"\nmenuTitle = \"%s\"\nchapter = false\nweight = %d\ndescription = \"\"\ndraft = false\n+++\n", title, order)
	header := fmt.Sprintf("+++\ntitle = \"%s\"\nweight = %d\ndescription = \"\"\ndraft = false\n+++\n", title, order)
	_, err = f.WriteString(header)
	if err != nil {
		return fmt.Errorf("cannot write string %s, %+v", header, err)
	}
	return nil
}

func copyStaticfileAuth() error {
	destination, err := os.Create("./publicGen/Staticfile.auth")

	defer destination.Close()
	source, err := os.Open("./Staticfile.auth")
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("cannot create Staticfile.auth", err)
	}
	return nil
}

func setWorkshopTitle(config *util.WorkshopConfig) error {
	workshopTitle := fmt.Sprintf("%s Workshop", config.WorkshopSubject)
	workshopToml := fmt.Sprintf("+++\ntitle = \"%s\"\nchapter = true\nweight = 1\n+++\n\n", workshopTitle)
	workshopHomepageContent := workshopToml
	if config.WorkshopHomepage != "" {
		homepageContent, err := ioutil.ReadFile(config.WorkshopHomepage)
		if err != nil {
			fmt.Printf("%s not found!\n", config.WorkshopHomepage)
			return err
		}
		workshopHomepageContent = workshopHomepageContent + string(homepageContent)

	} else {
		workshopHomepageContent = workshopHomepageContent + `<p style="font-family: Novacento Sans Wide, Helvetica, Tahoma, Geneva, Arial, sans-serif;
    text-align: center;
    text-transform: uppercase;
    color: #222;
    font-weight: 200;
	font-size: 3rem;">` + workshopTitle + `
</p>

<div class="text" style="background-color: #4fb2a3; border-radius: 15px; padding: 30px;">
<div style="width: 145px; height: 45px; color: #ffffff">
<svg width="396.8" height="56.962" viewBox="0 0 396.8 56.962" xmlns="http://www.w3.org/2000/svg"><g id="svgGroup" stroke-linecap="round" fill-rule="evenodd" font-size="9pt" stroke="#000" stroke-width="0.5mm" fill="#000" style="stroke:#000;stroke-width:0.5mm;fill:#000"><path d="M 250.56 16.961 L 250.56 17.921 L 234 17.921 L 234 17.601 A 5.345 5.345 0 0 0 233.474 15.215 A 5.381 5.381 0 0 0 232.24 13.601 Q 230.48 12.001 226.88 12.001 A 19.16 19.16 0 0 0 224.887 12.097 Q 222.858 12.31 221.569 12.992 A 5.502 5.502 0 0 0 221.48 13.041 A 5.364 5.364 0 0 0 220.625 13.618 Q 219.6 14.478 219.6 15.601 A 2.996 2.996 0 0 0 220.894 18.106 Q 221.371 18.463 222.023 18.743 A 7.416 7.416 0 0 0 222.16 18.801 Q 224.72 19.841 230.4 20.961 A 128.764 128.764 0 0 1 234.868 21.954 Q 236.955 22.459 238.735 22.98 A 56.456 56.456 0 0 1 241.32 23.801 Q 245.6 25.281 248.8 28.641 A 11.465 11.465 0 0 1 251.615 33.911 A 17.041 17.041 0 0 1 252.08 37.761 A 22.807 22.807 0 0 1 251.403 43.496 A 15.039 15.039 0 0 1 245.48 52.241 A 23.466 23.466 0 0 1 238.045 55.706 Q 233.533 56.961 227.84 56.961 Q 219.422 56.961 213.447 55.115 A 25.396 25.396 0 0 1 207.8 52.641 A 13.989 13.989 0 0 1 201.486 44.111 Q 200.772 41.594 200.661 38.519 A 31.958 31.958 0 0 1 200.64 37.361 L 217.36 37.361 Q 217.36 41.521 219.52 42.921 Q 221.68 44.321 226.24 44.321 A 25.929 25.929 0 0 0 228.505 44.228 Q 230.391 44.062 231.8 43.601 A 3.808 3.808 0 0 0 232.876 43.066 Q 233.851 42.351 233.98 41.052 A 4.151 4.151 0 0 0 234 40.641 A 2.784 2.784 0 0 0 232.82 38.343 Q 232.303 37.945 231.56 37.641 Q 229.706 36.881 226.074 36.052 A 104.747 104.747 0 0 0 223.6 35.521 Q 216.88 34.081 212.48 32.521 Q 208.08 30.961 204.8 27.361 A 12.104 12.104 0 0 1 202.02 22.023 Q 201.559 20.164 201.523 17.974 A 22.702 22.702 0 0 1 201.52 17.601 A 18.31 18.31 0 0 1 202.31 12.064 A 13.781 13.781 0 0 1 208.52 4.281 Q 214.754 0.469 223.94 0.052 A 50.74 50.74 0 0 1 226.24 0.001 Q 236.8 0.001 243.6 4.281 A 14.593 14.593 0 0 1 248.63 9.394 Q 250.478 12.629 250.56 16.961 Z M 380.96 34.001 L 396.8 56.001 L 378.24 56.001 L 370.24 42.721 L 369.92 42.721 L 361.84 56.001 L 344.16 56.001 L 360.4 33.601 L 345.76 13.761 L 364.48 13.761 L 371.2 24.961 L 371.52 24.961 L 378.24 13.761 L 395.84 13.761 L 380.96 34.001 Z M 103.84 26.561 L 103.84 43.041 A 4.743 4.743 0 0 0 103.937 44.029 Q 104.078 44.688 104.42 45.196 A 2.964 2.964 0 0 0 104.48 45.281 Q 105.12 46.161 106.4 46.161 L 109.28 46.161 L 109.28 55.361 A 1.361 1.361 0 0 1 109.14 55.439 Q 108.823 55.599 108.04 55.881 Q 107.137 56.206 105.548 56.531 A 40.864 40.864 0 0 1 105.2 56.601 A 18.25 18.25 0 0 1 103.352 56.859 Q 102.392 56.946 101.315 56.959 A 31.172 31.172 0 0 1 100.96 56.961 A 25.386 25.386 0 0 1 97.837 56.78 Q 96.225 56.58 94.873 56.158 A 12.406 12.406 0 0 1 93.32 55.561 A 10.079 10.079 0 0 1 91.299 54.33 A 7.094 7.094 0 0 1 89.2 51.681 Q 86.16 54.081 82.4 55.521 A 20.86 20.86 0 0 1 78.182 56.623 Q 76.179 56.939 73.913 56.959 A 34.861 34.861 0 0 1 73.6 56.961 A 27.447 27.447 0 0 1 68.169 56.472 Q 59.747 54.764 58.832 47.082 A 16.587 16.587 0 0 1 58.72 45.121 A 18.248 18.248 0 0 1 59.069 41.433 Q 59.524 39.231 60.571 37.55 A 9.743 9.743 0 0 1 62.04 35.721 A 14.306 14.306 0 0 1 65.868 33.096 Q 68.334 31.909 71.6 31.281 Q 77.84 30.081 87.92 30.081 L 87.92 28.001 A 5.672 5.672 0 0 0 87.717 26.437 A 4.058 4.058 0 0 0 86.2 24.241 A 6.211 6.211 0 0 0 83.913 23.197 Q 82.929 22.961 81.76 22.961 A 11.743 11.743 0 0 0 79.72 23.129 Q 78.492 23.346 77.48 23.841 A 3.601 3.601 0 0 0 76.497 24.525 Q 75.75 25.276 75.686 26.422 A 3.911 3.911 0 0 0 75.68 26.641 L 75.68 26.961 L 60 26.961 A 3.264 3.264 0 0 1 59.958 26.679 Q 59.925 26.384 59.921 25.977 A 12.144 12.144 0 0 1 59.92 25.841 A 10.244 10.244 0 0 1 63.687 17.736 A 15.277 15.277 0 0 1 65.64 16.321 Q 71.332 12.818 81.896 12.801 A 64.048 64.048 0 0 1 82 12.801 A 47.714 47.714 0 0 1 88.23 13.183 Q 93.77 13.914 97.76 16.041 A 11.547 11.547 0 0 1 101.697 19.34 Q 103.354 21.582 103.73 24.675 A 15.642 15.642 0 0 1 103.84 26.561 Z M 191.36 26.561 L 191.36 43.041 A 4.743 4.743 0 0 0 191.457 44.029 Q 191.598 44.688 191.94 45.196 A 2.964 2.964 0 0 0 192 45.281 Q 192.64 46.161 193.92 46.161 L 196.8 46.161 L 196.8 55.361 A 1.361 1.361 0 0 1 196.66 55.439 Q 196.343 55.599 195.56 55.881 Q 194.657 56.206 193.068 56.531 A 40.864 40.864 0 0 1 192.72 56.601 A 18.25 18.25 0 0 1 190.872 56.859 Q 189.912 56.946 188.835 56.959 A 31.172 31.172 0 0 1 188.48 56.961 A 25.386 25.386 0 0 1 185.357 56.78 Q 183.745 56.58 182.393 56.158 A 12.406 12.406 0 0 1 180.84 55.561 A 10.079 10.079 0 0 1 178.819 54.33 A 7.094 7.094 0 0 1 176.72 51.681 Q 173.68 54.081 169.92 55.521 A 20.86 20.86 0 0 1 165.702 56.623 Q 163.699 56.939 161.433 56.959 A 34.861 34.861 0 0 1 161.12 56.961 A 27.447 27.447 0 0 1 155.689 56.472 Q 147.267 54.764 146.352 47.082 A 16.587 16.587 0 0 1 146.24 45.121 A 18.248 18.248 0 0 1 146.589 41.433 Q 147.044 39.231 148.091 37.55 A 9.743 9.743 0 0 1 149.56 35.721 A 14.306 14.306 0 0 1 153.388 33.096 Q 155.854 31.909 159.12 31.281 Q 165.36 30.081 175.44 30.081 L 175.44 28.001 A 5.672 5.672 0 0 0 175.237 26.437 A 4.058 4.058 0 0 0 173.72 24.241 A 6.211 6.211 0 0 0 171.433 23.197 Q 170.449 22.961 169.28 22.961 A 11.743 11.743 0 0 0 167.24 23.129 Q 166.013 23.346 165 23.841 A 3.601 3.601 0 0 0 164.017 24.525 Q 163.27 25.276 163.206 26.422 A 3.911 3.911 0 0 0 163.2 26.641 L 163.2 26.961 L 147.52 26.961 A 3.264 3.264 0 0 1 147.478 26.679 Q 147.445 26.384 147.441 25.977 A 12.144 12.144 0 0 1 147.44 25.841 A 10.244 10.244 0 0 1 151.207 17.736 A 15.277 15.277 0 0 1 153.16 16.321 Q 158.852 12.818 169.416 12.801 A 64.048 64.048 0 0 1 169.52 12.801 A 47.714 47.714 0 0 1 175.75 13.183 Q 181.29 13.914 185.28 16.041 A 11.547 11.547 0 0 1 189.217 19.34 Q 190.874 21.582 191.25 24.675 A 15.642 15.642 0 0 1 191.36 26.561 Z M 338 26.561 L 338 43.041 A 4.743 4.743 0 0 0 338.097 44.029 Q 338.238 44.688 338.58 45.196 A 2.964 2.964 0 0 0 338.64 45.281 Q 339.28 46.161 340.56 46.161 L 343.44 46.161 L 343.44 55.361 A 1.361 1.361 0 0 1 343.3 55.439 Q 342.983 55.599 342.2 55.881 Q 341.297 56.206 339.708 56.531 A 40.864 40.864 0 0 1 339.36 56.601 A 18.25 18.25 0 0 1 337.512 56.859 Q 336.552 56.946 335.475 56.959 A 31.172 31.172 0 0 1 335.12 56.961 A 25.386 25.386 0 0 1 331.997 56.78 Q 330.385 56.58 329.033 56.158 A 12.406 12.406 0 0 1 327.48 55.561 A 10.079 10.079 0 0 1 325.459 54.33 A 7.094 7.094 0 0 1 323.36 51.681 Q 320.32 54.081 316.56 55.521 A 20.86 20.86 0 0 1 312.342 56.623 Q 310.339 56.939 308.073 56.959 A 34.861 34.861 0 0 1 307.76 56.961 A 27.447 27.447 0 0 1 302.329 56.472 Q 293.907 54.764 292.992 47.082 A 16.587 16.587 0 0 1 292.88 45.121 A 18.248 18.248 0 0 1 293.229 41.433 Q 293.684 39.231 294.731 37.55 A 9.743 9.743 0 0 1 296.2 35.721 A 14.306 14.306 0 0 1 300.028 33.096 Q 302.494 31.909 305.76 31.281 Q 312 30.081 322.08 30.081 L 322.08 28.001 A 5.672 5.672 0 0 0 321.877 26.437 A 4.058 4.058 0 0 0 320.36 24.241 A 6.211 6.211 0 0 0 318.073 23.197 Q 317.089 22.961 315.92 22.961 A 11.743 11.743 0 0 0 313.88 23.129 Q 312.653 23.346 311.64 23.841 A 3.601 3.601 0 0 0 310.657 24.525 Q 309.91 25.276 309.846 26.422 A 3.911 3.911 0 0 0 309.84 26.641 L 309.84 26.961 L 294.16 26.961 A 3.264 3.264 0 0 1 294.118 26.679 Q 294.085 26.384 294.081 25.977 A 12.144 12.144 0 0 1 294.08 25.841 A 10.244 10.244 0 0 1 297.847 17.736 A 15.277 15.277 0 0 1 299.8 16.321 Q 305.492 12.818 316.056 12.801 A 64.048 64.048 0 0 1 316.16 12.801 A 47.714 47.714 0 0 1 322.39 13.183 Q 327.93 13.914 331.92 16.041 A 11.547 11.547 0 0 1 335.857 19.34 Q 337.514 21.582 337.89 24.675 A 15.642 15.642 0 0 1 338 26.561 Z M 23.84 56.001 L 0 56.001 L 0 0.961 L 23.84 0.961 A 43.762 43.762 0 0 1 35.093 2.278 Q 52.72 6.981 52.72 28.481 A 37.656 37.656 0 0 1 51.238 39.538 Q 46.177 56.001 23.84 56.001 Z M 132.32 13.761 L 141.28 13.761 L 141.28 24.561 L 132.32 24.561 L 132.32 40.641 A 13.711 13.711 0 0 0 132.394 42.118 Q 132.558 43.623 133.084 44.538 A 3.384 3.384 0 0 0 133.28 44.841 A 2.828 2.828 0 0 0 134.525 45.785 Q 135.023 45.991 135.656 46.084 A 7.898 7.898 0 0 0 136.8 46.161 L 141.28 46.161 L 141.28 55.521 A 18.151 18.151 0 0 1 139.755 55.95 Q 138.268 56.304 136.32 56.561 Q 133.28 56.961 131.04 56.961 Q 124 56.961 120.2 54.401 A 8.083 8.083 0 0 1 116.98 49.944 Q 116.545 48.551 116.436 46.849 A 18.337 18.337 0 0 1 116.4 45.681 L 116.4 24.561 L 110.48 24.561 L 110.48 13.761 L 117.04 13.761 L 120.48 0.961 L 132.32 0.961 L 132.32 13.761 Z M 278.96 13.761 L 287.92 13.761 L 287.92 24.561 L 278.96 24.561 L 278.96 40.641 A 13.711 13.711 0 0 0 279.034 42.118 Q 279.198 43.623 279.724 44.538 A 3.384 3.384 0 0 0 279.92 44.841 A 2.828 2.828 0 0 0 281.165 45.785 Q 281.663 45.991 282.296 46.084 A 7.898 7.898 0 0 0 283.44 46.161 L 287.92 46.161 L 287.92 55.521 A 18.151 18.151 0 0 1 286.395 55.95 Q 284.907 56.304 282.96 56.561 Q 279.92 56.961 277.68 56.961 Q 270.64 56.961 266.84 54.401 A 8.083 8.083 0 0 1 263.62 49.944 Q 263.185 48.551 263.076 46.849 A 18.337 18.337 0 0 1 263.04 45.681 L 263.04 24.561 L 257.12 24.561 L 257.12 13.761 L 263.68 13.761 L 267.12 0.961 L 278.96 0.961 L 278.96 13.761 Z M 17.68 14.161 L 17.68 42.801 L 23.52 42.801 Q 32.35 42.801 34.169 35.284 A 18.782 18.782 0 0 0 34.64 30.881 L 34.64 26.081 A 18.484 18.484 0 0 0 34.117 21.472 Q 32.589 15.556 26.595 14.427 A 16.634 16.634 0 0 0 23.52 14.161 L 17.68 14.161 Z M 87.92 41.201 L 87.92 37.601 Q 81.12 37.601 77.88 39.081 Q 76.123 39.883 75.319 40.98 A 3.403 3.403 0 0 0 74.64 43.041 A 4.106 4.106 0 0 0 75.041 44.918 Q 75.912 46.636 78.673 46.963 A 11.971 11.971 0 0 0 80.08 47.041 A 9.789 9.789 0 0 0 83.392 46.494 A 8.85 8.85 0 0 0 85.56 45.361 Q 87.92 43.681 87.92 41.201 Z M 175.44 41.201 L 175.44 37.601 Q 168.64 37.601 165.4 39.081 Q 163.643 39.883 162.839 40.98 A 3.403 3.403 0 0 0 162.16 43.041 A 4.106 4.106 0 0 0 162.561 44.918 Q 163.432 46.636 166.193 46.963 A 11.971 11.971 0 0 0 167.6 47.041 A 9.789 9.789 0 0 0 170.912 46.494 A 8.85 8.85 0 0 0 173.08 45.361 Q 175.44 43.681 175.44 41.201 Z M 322.08 41.201 L 322.08 37.601 Q 315.28 37.601 312.04 39.081 Q 310.283 39.883 309.479 40.98 A 3.403 3.403 0 0 0 308.8 43.041 A 4.106 4.106 0 0 0 309.201 44.918 Q 310.072 46.636 312.833 46.963 A 11.971 11.971 0 0 0 314.24 47.041 A 9.789 9.789 0 0 0 317.552 46.494 A 8.85 8.85 0 0 0 319.72 45.361 Q 322.08 43.681 322.08 41.201 Z" vector-effect="non-scaling-stroke"/></g></svg>
</div>
	<h1 style="color: #09243c; font-family:"Proxima Nova", sans-serif; font-size: 70px;">The real-time data cloud.</h1>
	<h2 style="color: #ffffff; font-family:"Proxima Nova", sans-serif; font-size: 45px;">Leading companies build and run their most important applications with a limitless open data cloud. Learn how our 
open data stack, tools, and methodology help you deliver exceptional real-time user experiences.</h2>
	<p style="color: #ffffff; font-family:"Proxima Nova", sans-serif; font-size: 25px;">DataStax Cloud Data Architects helps companies learn how to design 
and deploy real-time solutions. We encourage you to explore our workshops. 
Build the future with Pivotal!</p>
</div>

`
	}

	workshop, err := os.OpenFile("workshopGen/content/_index.en.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}

	workshop, err = os.OpenFile("workshopGen/content/_index.es.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}

	workshop, err = os.OpenFile("workshopGen/content/_index.fr.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}

	workshop, err = os.OpenFile("workshopGen/content/_index.pt.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}
	return nil
}
