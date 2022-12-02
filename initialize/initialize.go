package initialize

import (
	"fmt"
	"os"

	"github.com/datastax-cda/workshop-builder/util"
)

func InitCmd() {

	fmt.Println("Generating default pace config.json")
	if err := createDefaultConfig(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Generating default cf push manifest.yml")
	if err := createDefaultManifest(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Generating default Staticfile.auth")
	if err := createDefaultAuthFile(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Pulling PACE workshop content...")
	if err := getWorkshopContent(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Sample Config, Manifest and Staticfile.auth have been generated. Edit the config, manifest and Staticfile.auth to your desire. Run `pace build` to build your first pace workshop!")
	fmt.Println("Adjusting workshop content locally can be done within the paceWorkshopContent folder. Once your content is ready to be shared with your fellow team members, commit it back to pace workshop content! ")
}

func createDefaultConfig() error {
	f, err := os.OpenFile("config.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error creating config.json")
	}
	defer f.Close()
	_, err = f.WriteString(util.DefaultConfig)
	if err != nil {
		return fmt.Errorf("error writing default config to config.json")
	}

	return nil
}

func createDefaultManifest() error {
	f, err := os.OpenFile("manifest.yml", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error creating manifest.yml")
	}
	defer f.Close()
	_, err = f.WriteString(util.DefaultManifest)
	if err != nil {
		return fmt.Errorf("error writing default manifest to manifest.yml")
	}

	return nil
}

func createDefaultAuthFile() error {
	f, err := os.OpenFile("Staticfile.auth", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error creating Staticfile.auth")
	}
	defer f.Close()
	_, err = f.WriteString(util.DefaultStaticFile)
	if err != nil {
		return fmt.Errorf("error writing default manifest to Staticfile.auth")
	}

	return nil
}

func getWorkshopContent() error {

	if _, err := os.Stat("paceWorkshopContent"); os.IsNotExist(err) {
		if err := util.CloneRepo("https://github.com/Pivotal-Field-Engineering/pace-workshop-content", "paceWorkshopContent"); err != nil {
			return err
		}
	}
	return nil
}
