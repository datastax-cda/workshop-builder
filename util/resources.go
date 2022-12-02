package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var DefaultConfig = `{
    "workshopSubject":"PACE",
    "workshopHomepage":"",
    "modules": [
    {
        "type": "concepts",
        "content": [
            {
            "name":"example-slide",
            "filename":"example/example-slide"
            }
        ]
    },
    {
        "type": "demos",
        "content": [
            {
            "name":"example-demo",
            "filename":"example/example-demo"
            }
        ]
    }
  ]
}`

var DefaultManifest = `---
applications:
- name: my-pace-workshop
  memory: 64M
  instances: 1
  buildpacks: 
  - staticfile_buildpack
  random-route: true
  path: publicGen/`

var DefaultStaticFile = `guest:$apr1$oM3ne/Oz$86q6.UWNEb0Nfv3xbSiiB0`

type WorkshopConfig struct {
	WorkshopHomepage string `json:"workshopHomepage"`
	WorkshopSubject  string `json:"workshopSubject"`
	WorkshopHostname string `json:"workshopHostname"`
	Modules          []struct {
		Type    string          `json:"type"`
		Content []ContentConfig `json:"content"`
	} `json:"modules"`
}

type ContentConfig struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

func DetermineConfig(path string) (*WorkshopConfig, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config not found")
	}
	var config WorkshopConfig
	err = json.Unmarshal(configFile, &config)
	return &config, nil
}
