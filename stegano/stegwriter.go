package stegano

import (
	"fmt"
	"image"
)

type StegWriter struct {
	image *image.NRGBA
	index int64
}

func (w *StegWriter) writeByte(b byte) error {
	for bit := 0; bit < 8; bit++ {
		if w.index >= int64(len(w.image.Pix)) {
			return fmt.Errorf("Image does not have enough capacity to encode data")
		}
		w.image.Pix[w.index+int64(bit)] = (w.image.Pix[w.index+int64(bit)] & 0xFE) | ((b >> (7 - bit)) & 1)
	}
	w.index += 8
	return nil
}

func (w *StegWriter) Write(p []byte) (n int, err error) {
	for i := range p {
		err := w.writeByte(p[i])
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

func (w *StegWriter) WriteInt64(value int64) error {
	for i := 7; i >= 0; i-- {
		b := byte((value >> (i * 8)) & 0xFF)
		if err := w.writeByte(b); err != nil {
			return err
		}
	}
	return nil
}
