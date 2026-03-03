package stegano

import (
	"image"
	"io"
)

type StegWriter struct {
	image *image.NRGBA
	index int64
}

func (w *StegWriter) bitAtOffset(b byte, pos, offset uint) {
	index := w.index + int64(offset)
	w.image.Pix[index] = (w.image.Pix[index] & 0xFE) | ((b >> pos) & 1)
}

func (w *StegWriter) Pos() int64 {
	return w.index
}

// Each byte written occupies 3 pixels = 12 image bytes
func (w *StegWriter) writeByte(b byte) error {
	if w.index+BYTE_LEN >= int64(len(w.image.Pix)) {
		return io.EOF
	}
	// 1st pixel: R, G, B bit 1
	w.bitAtOffset(b, 7, 0)
	w.bitAtOffset(b, 6, 1)
	w.bitAtOffset(b, 5, 2)
	// 2nd pixel: R, G, B bit 1
	w.bitAtOffset(b, 4, 4)
	w.bitAtOffset(b, 3, 5)
	w.bitAtOffset(b, 2, 6)
	// 3rd pixel: R, G bit 1
	w.bitAtOffset(b, 1, 8)
	w.bitAtOffset(b, 0, 9)
	w.index += BYTE_LEN
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

// Writes the 8 bytes of an int64 starting with the MSB and down to the LSB
func (w *StegWriter) WriteInt64(value int64) error {
	for i := 7; i >= 0; i-- {
		b := byte((value >> (i * 8)) & 0xFF)
		if err := w.writeByte(b); err != nil {
			return err
		}
	}
	return nil
}

// Writes a string by first writing its length and then its content
func (w *StegWriter) WriteString(str string) error {
	length := int64(len(str))
	if err := w.WriteInt64(length); err != nil {
		return err
	}
	if _, err := w.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}
