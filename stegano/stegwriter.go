package stegano

import (
	"image"
	"io"
)

type StegWriter struct {
	image *image.NRGBA
	index int64
}

func (w *StegWriter) bitAtOffset(b byte, pos int, offset int64) {
	index := w.index + int64(offset)
	w.image.Pix[index] = (w.image.Pix[index] & 0xFE) | ((b >> pos) & 1)
}

func (w *StegWriter) Pos() int64 {
	return w.index
}

// Each byte written occupies 3 pixels = 12 image bytes
// The algorithm used is defined in the OffsetForBit function, which determines how the bits of the byte are distributed across the RGB channels of the pixels.
func (w *StegWriter) writeByte(b byte) error {
	if w.index+BYTE_LEN >= int64(len(w.image.Pix)) {
		return io.EOF
	}

	for bit := 7; bit >= 0; bit-- {
		w.bitAtOffset(b, bit, OffsetForBit(bit))
	}

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
