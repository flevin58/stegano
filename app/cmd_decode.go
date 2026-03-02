package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/flevin58/stegano/stegano"
)

func cmdDecode() error {
	cnvt := flag.NewFlagSet("decode", flag.ExitOnError)
	cnvt.Parse(os.Args[2:])
	what := cnvt.Arg(0)
	if what == "" {
		return fmt.Errorf("Missing argument: what to decode")
	}
	img, err := stegano.Load(what)
	if err != nil {
		return err
	}
	encoded, err := img.IsEncoded()
	if err != nil {
		return err
	}
	if !encoded {
		return fmt.Errorf("Image is not encoded")
	}
	data, err := img.Decode()
	if err != nil {
		return err
	}
	fmt.Println("Hidden data:", string(data))
	return nil
}
