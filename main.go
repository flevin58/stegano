package main

import "github.com/flevin58/stegano/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
