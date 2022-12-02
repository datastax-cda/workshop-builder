package serve

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gohugoio/hugo/commands"
)

func ServeCmd() {

	fmt.Println("Checking workshopGen")
	if err := os.Chdir("workshopGen/"); err != nil {
		fmt.Println("Please `build` before `serve` to create the content. Error:" + err.Error())
		return
	}

	fmt.Println("Serving up a local version of your workshop!  Check your content at http://localhost:1313 ...")
	if err := serveHugo(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}
}

func serveHugo() error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	resp := commands.Execute([]string{"serve", "--quiet"})

	if resp.Err != nil {
		if resp.IsUserError() {
			resp.Cmd.Println("")
			resp.Cmd.Println(resp.Cmd.UsageString())
		}
		os.Exit(-1)
	}
	return nil
}
