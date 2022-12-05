# DataStax Workshop Builder

Bringing back to life the work created by [James Williams](https://www.linkedin.com/in/james-williams-b509341/) and [Peter Blum](https://www.linkedin.com/in/pblum/) and countless amazing Platform Architects at Pivotal. 

The purpose of this tool is to help create microsites used primarily to run customer workshops. 

## Quick Start

1. Download the correct `dscda` CLI binary from the releases tab.
    - *MAC OS Users Optional:* If you have the `brew tap dscda/tap` installed you can install the `dscda` CLI with `brew install dscda-cli`

1. Run `dscda init`.

1. Edit the `config.json`. The format should follow the `sampleConfig.json`.

1. Run `dscda build`. Notice the new `workshopGen` folder. This contains your new workshop.

1. **Optional** Run `dscda serve` to view your workshop. View local running site at http://localhost:1313

1. Deploy the static microsite built with [HUGO](https://gohugo.io/hosting-and-deployment/) at your environment of choice.

1. **Optional* Use our Netlify(https://app.netlify.com/teams/mborges-pivotal/overview) team to deploy. If you use this option, your workshop will be auto-deleted after 30 days.


## Notes

1. Content is pulled from the [workshop-content](https://github.com/datastax-cda/workshop-content) github repo. Feel free to add any content there that you can then use to build a workshop with `dscda build`

1. While `dscda` will build a generic homepage for your workshop you can setup a custom one by supplying a markdown file via the `workshopHomepage` field in the `config.json` file. This is not required.

## Build/Install workshop-builder manually
1. Download and install [go](https://golang.org/dl/)

1. Make sure you have something like this in your terminal profile

        export GOPATH=~/go
        export GOBIN=$GOPATH/bin

1. Fetch the latest workshop-builder source into the go workspace

        go get -u github.com/datastax-cda/workshop-builder

1. Open a terminal window to `workshop-builder` directory

        cd ${GOPATH}/src/datastax-cda/workshop-builder

1. Build binary by running:

        go install

1. You should have an executable binary in `$GOBIN/workshop-builder`.

1. [OPTIONAL] Rename `workshop-builder` to `dscda`:

        mv $GOBIN/workshop-builder $GOBIN/dscda

1. Test your new dscda install with:

        dscda -h
