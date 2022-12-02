package clean

import (
	"fmt"
	"os"
)

func CleanCmd() {
	if err := os.RemoveAll("paceWorkshopContent/"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}
	if err := os.RemoveAll("workshopGen/"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}
}
