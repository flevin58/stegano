package main

//go:generate go run generate/main.go

import (
	"fmt"

	"github.com/flevin58/stegano/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
