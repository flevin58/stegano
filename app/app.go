package app

import (
	"fmt"
	"os"
)

func Run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Missing arguments")
	}
	for _, cmd := range commandList {
		if cmd.name == os.Args[1] {
			return cmd.handler()
		}
	}
	return fmt.Errorf("Bad command: %v\n", os.Args[1])
}
