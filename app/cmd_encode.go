package app

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/flevin58/stegano/stegano"
)

func cmdEncode() error {
	cnvt := flag.NewFlagSet("encode", flag.ExitOnError)
	text := cnvt.String("text", "", "Text to encode into the image")
	file := cnvt.String("file", "", "File to encode into the image")
	cnvt.Parse(os.Args[2:])
	what := cnvt.Arg(0)
	if what == "" {
		return fmt.Errorf("Missing argument: what to encode")
	}
	img, err := stegano.Load(what)
	if err != nil {
		return err
	}

	if img.IsEncoded() {
		return fmt.Errorf("Image is already encoded")
	}

	if len(*text) > 0 && len(*file) > 0 {
		return fmt.Errorf("Flags text and file cannot be used together")
	}

	if len(*text) > 0 {
		err = img.Encode("_TEXT_", []byte(*text))
		if err != nil {
			return err
		}
	}

	if len(*file) > 0 {
		fileInfo, err := os.Stat(*file)
		fileName := path.Base(*file)
		if err != nil {
			return err
		}
		fileSize := fileInfo.Size()
		inp, err := os.Open(*file)
		if err != nil {
			return err
		}
		data := make([]byte, fileSize)
		_, err = inp.Read(data)
		err = img.Encode(fileName, data)
		if err != nil {
			return err
		}
	}

	err = img.Save()
	if err != nil {
		return err
	}
	fmt.Println("Data encoded and image saved")
	return nil
}
