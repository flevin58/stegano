package app

import (
	"fmt"
	"os"
)

func Run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Missing command: encode|decode")
	}
	for _, cmd := range commandList {
		if cmd.name == os.Args[1] {
			return cmd.handler()
		}
	}
	return fmt.Errorf("Bad command: %v\n", os.Args[1])
}

func Usage() {
	fmt.Println(`
Usage:
	stegano encode -text <string> <image.png>
	
	stegano decode <image.png>`)
}
