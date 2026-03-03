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
	// Check if we have enough pixels left to read a byte
	// We need 3 + 3 + 2 bits = 3 colors x 4 bytes each = 12 bytes
	if r.index+BYTE_LEN >= int64(len(r.image.Pix)) {
		return 0, io.EOF
	}
	// Extract the least significant bit from the R, G, and B channels of the current pixel
	rbit := r.image.Pix[r.index+0] & 1
	gbit := r.image.Pix[r.index+1] & 1
	bbit := r.image.Pix[r.index+2] & 1
	b = rbit<<2 | gbit<<1 | bbit
	rbit = r.image.Pix[r.index+4] & 1
	gbit = r.image.Pix[r.index+5] & 1
	bbit = r.image.Pix[r.index+6] & 1
	b = b<<3 | rbit<<2 | gbit<<1 | bbit
	// By choice we consider R and G as the last 2 bits
	rbit = r.image.Pix[r.index+8] & 1
	gbit = r.image.Pix[r.index+9] & 1
	b = b<<2 | rbit<<1 | gbit
	// Move to the next set of pixels for the next byte
	r.index += BYTE_LEN
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

func (r *StegReader) Pos() int64 {
	return r.index
}

func (r *StegReader) Reset() {
	r.index = 0
}

func (r *StegReader) Skip(n int64) error {
	r.index += n * BYTE_LEN
	if r.index > int64(len(r.image.Pix)) {
		return io.EOF
	}
	return nil
}

func (r *StegReader) ReadHeader() (header string) {
	var magic [4]byte
	r.Reset()
	n, err := r.Read(magic[:])
	if err == nil && n == 4 {
		header = string(magic[:])
	}
	return header
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

func (r *StegReader) ReadString() (str string, err error) {
	length, err := r.ReadInt64()
	if err != nil {
		return str, err
	}
	data := make([]byte, length)
	_, err = r.Read(data)
	if err != nil {
		return str, err
	}
	return string(data), nil
}
