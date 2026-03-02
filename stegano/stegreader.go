package stegano

import (
	"image"
	"io"
)

type StegReader struct {
	image *image.NRGBA
	index int64
}

// Reads a single byte from the image's pixel data by extracting the least significant bits.
func (r *StegReader) readByte() (b byte, err error) {
	if r.index >= int64(len(r.image.Pix)) {
		return 0, io.EOF
	}
	for bit := 0; bit < 8; bit++ {
		val := r.image.Pix[r.index+int64(bit)]
		b <<= 1
		b |= val & 1
	}
	r.index += 8
	return b, nil
}

func (r *StegReader) Read(p []byte) (n int, err error) {
	for i := range p {
		b, err := r.readByte()
		if err != nil {
			return n, err
		}
		p[i] = b
		n++
	}
	return n, nil
}

func (r *StegReader) Skip(n int) error {
	r.index += int64(n * 8)
	if r.index > int64(len(r.image.Pix)) {
		return io.EOF
	}
	return nil
}

func (r *StegReader) ReadInt64() (value int64, err error) {
	for i := 0; i < 8; i++ {
		b, err := r.readByte()
		if err != nil {
			return 0, err
		}
		value = (value << 8) | int64(b)
	}
	return value, nil
}
