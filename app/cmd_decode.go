package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/flevin58/stegano/stegano"
)

func cmdDecode() error {
	cnvt := flag.NewFlagSet("decode", flag.ExitOnError)
	as := cnvt.String("as", "", "Saves the embedded file as the new given name")
	cnvt.Parse(os.Args[2:])
	what := cnvt.Arg(0)
	if what == "" {
		return fmt.Errorf("Missing argument: what to decode")
	}
	img, err := stegano.Load(what)
	if err != nil {
		return err
	}

	if !img.IsEncoded() {
		return fmt.Errorf("Image is not encoded")
	}

	fileName, data, err := img.Decode()
	if *as != "" {
		fileName = *as
	}
	if err != nil {
		return err
	}
	if fileName == "_TEXT_" {
		fmt.Println("Hidden data:", string(data))
	} else {
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		f.Write(data)
		fmt.Printf("Hidden file '%v' saved to disk.\n", fileName)
	}
	return nil
}
