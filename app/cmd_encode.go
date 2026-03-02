package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/flevin58/stegano/stegano"
)

func cmdEncode() error {
	cnvt := flag.NewFlagSet("encode", flag.ExitOnError)
	text := cnvt.String("text", "This is hidden data!", "Text to encode into the image")
	cnvt.Parse(os.Args[2:])
	what := cnvt.Arg(0)
	if what == "" {
		return fmt.Errorf("Missing argument: what to encode")
	}
	img, err := stegano.Load(what)
	if err != nil {
		return err
	}
	encoded, err := img.IsEncoded()
	if err != nil {
		return err
	}
	if encoded {
		return fmt.Errorf("Image is already encoded")
	}
	err = img.Encode([]byte(*text))
	if err != nil {
		return err
	}
	err = img.Save()
	if err != nil {
		return err
	}
	fmt.Println("Data encoded and image saved")
	return nil
}
