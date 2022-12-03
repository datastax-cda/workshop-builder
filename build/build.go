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
<svg viewBox="0 0 291.93 68.4" xmlns="http://www.w3.org/2000/svg">
<g style="fill: white">
<path style="fill: white" d="M131.51,61.59H121.14v-10h10.37v10Zm0,57.21h-10.4V69h10.4V118.8Z" transform="translate(-66.72 -51.56)"></path>
<path style="fill: white" d="M191.48,69L177,112c-2.49,7-6.92,7.94-10.52,7.94-5.33,0-8.54-2.45-10.44-7.93L144,76.27h-4.77V69h12.45l12.44,40c0.54,1.71.88,2.73,2.39,2.73s1.88-1,2.38-2.73l12.69-40h9.92Z" transform="translate(-66.72 -51.56)"></path>
<path style="fill: white" d="M217.08,69c13.45,0,22.83,8.87,22.83,21.59v7.83c0,12.7-9.38,21.59-22.83,21.59s-22.84-8.89-22.84-21.59V90.54c0-12.71,9.4-21.59,22.84-21.59m0,42.8c8,0,13-6.08,13-13.38V90.54c0-7.3-4.94-13.39-13-13.39-8.52,0-13,6.08-13,13.39v7.83c0,7.3,4.7,13.38,13,13.38" transform="translate(-66.72 -51.56)"></path>
<path style="fill: white" d="M322.9,70.55a88.57,88.57,0,0,0-19.25-2.3c-13.65,0-22.12,8.61-22.12,22.48v5.46c0,13.86,8.47,22.6,22.12,22.6,0.32,0,2.74,0,3.85-.1v-8.37c-0.42,0-3.54.11-3.85,0.11-7.42,0-12.42-5.72-12.42-14.24V90.73c0-8.52,5-14.24,12.42-14.24a72.93,72.93,0,0,1,10.18.64l0.56,0.12V118.8h10.39V72.27c0-.89,0-1.23-1.89-1.73" transform="translate(-66.72 -51.56)"></path>
<rect style="fill: white" height="67.24" width="10.39" x="268.07"></rect>
<path style="fill: white" d="M85,51.56H66.72v67.23H77.53V60.91h6.34c1.35,0,2.49.07,3.64,0.09,9.37,0.18,14,3.9,14,11.18,0,0.29,0,.49,0,0.79,0,6.74-3.7,11.21-13.92,11.21-1,0-2,0-3,0,0,2.58,0,7.36,0,9,1.05,0.05,2,.09,3.05.09,14.66,0,25-5.76,25-20.23,0-.28,0-0.58,0-0.87,0-15-11.28-20.64-27.61-20.64" transform="translate(-66.72 -51.56)"></path>
<path style="fill: white" d="M258.37,58.21V69h16.93V77H258.37v29c0,4.56,2.91,4.68,7.13,4.68h9.8v8.06H262c-9.81,0-14.18-3.93-14.18-12.74V59.65Z" transform="translate(-66.72 -51.56)"></path>
</g>
<path style="stroke: white;" d="M350,114.44a4.34,4.34,0,1,1,4.34,4.36A4.36,4.36,0,0,1,350,114.44Zm4.35-3.75a3.73,3.73,0,1,0,3.71,3.73A3.73,3.73,0,0,0,354.31,110.7Zm-1,6.05h-0.5V112h1.33c1.2,0,1.78.55,1.78,1.43a1.39,1.39,0,0,1-1,1.35l1.22,1.95h-0.6l-1.13-1.83-1.13.07v1.76Zm0.85-2.24a1.09,1.09,0,0,0,1.26-1c0-.67-0.42-1-1.33-1h-0.8v2.07Z" transform="translate(-66.72 -51.56)"></path>
</svg>
</div>
	<h1 style="color: #09243c; font-family:"Proxima Nova", sans-serif; font-size: 70px;">The way the future gets built.</h1>
	<h2 style="color: #ffffff; font-family:"Proxima Nova", sans-serif; font-size: 45px;">Leading companies build and run their most important applications on Pivotal. Learn how our 
platform, tools, and methodology help you deliver exceptional user experiences.</h2>
	<p style="color: #ffffff; font-family:"Proxima Nova", sans-serif; font-size: 25px;">Pivotal Platform Architecture helps companies learn how to solve IT
and engineering challenges. We encourage you to explore our workshops. 
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
