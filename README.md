# DataStax Workshop Builder

[![asciicast](https://asciinema.org/a/213808.svg)](https://asciinema.org/a/213808?autoplay=1)

## Quick Start

1. Download the correct `pace` CLI binary from the releases tab.
    - *MAC OS Users Optional:* If you have the `brew tap pivotal/tap` installed you can install the `pace` CLI with `brew install pace-cli`

1. Run `pace init`.

1. Edit the `config.json`. The format should follow the `sampleConfig.json`.

1. Run `pace build`. Notice the new `workshopGen` folder. This contains your new workshop.

1. **Optional** Run `pace serve` to view your workshop. View local running site at http://localhost:1313

1. Update your `manifest.yml` file settings. You may want to change the application `name` and turn off `random-route`

1. **Optional* Run `cf target -o Southeast -s pace`. This space is reserved for PACE workshop pushes. If you use this
space, your workshop will be auto-deleted after 30 days.

1. Run `cf push`. This will push your workshop to the PCF environment that you have specified in `cf target`. 


## Notes

1. Content is pulled from the [pace-workshop-content](https://github.com/Pivotal-Field-Engineering/pace-workshop-content) github repo. Feel free to add any content there that you can then use to build a workshop with `pace build`

1. While `pace` will build a generic homepage for your workshop you can setup a custom one by supplying a markdown file via the `workshopHomepage` field in the `config.json` file. This is not required.

1. `pace push` will automatically generate a random hostname for your workshop but you can specifiy a custom one with the use of the `workshopHostname` attribute inside your `config.json`.

## Build/Install pace-builder manually
1. Download and install [go](https://golang.org/dl/)

1. Make sure you have something like this in your terminal profile

        export GOPATH=~/go
        export GOBIN=$GOPATH/bin

1. Fetch the latest pace-builder source into the go workspace

        go get -u github.com/Pivotal-Field-Engineering/pace-builder

1. Open a terminal window to `pace-builder` directory

        cd ${GOPATH}/src/github.com/Pivotal-Field-Engineering/pace-builder

1. Build binary by running:

        go install

1. You should have an executable binary in `$GOBIN/pace-builder`.

1. [OPTIONAL] Rename `pace-builder` to `pace`:

        mv $GOBIN/pace-builder $GOBIN/pace

1. Test your new pace install with:

        pace -h
