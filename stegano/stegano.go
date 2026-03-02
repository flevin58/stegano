package stegano

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type Stegano struct {
	filePath string
	image    *image.NRGBA
}

func Load(filePath string) (*Stegano, error) {
	img, err := LoadNRGBA(filePath)
	if err != nil {
		return nil, err
	}
	return &Stegano{
		filePath: filePath,
		image:    img,
	}, nil
}

func (s *Stegano) Save() error {
	outPath := strings.TrimSuffix(s.filePath, filepath.Ext(s.filePath)) + "_stegano.png"
	return SaveNRGBA(s.image, outPath)
}

func (s *Stegano) SaveAs(filePath string) error {
	if filepath.Ext(filePath) != ".png" {
		return fmt.Errorf("Output file must have a .png extension")
	}
	return SaveNRGBA(s.image, filePath)
}

func (s *Stegano) NewReader() *StegReader {
	return &StegReader{
		image: s.image,
		index: 0,
	}
}

func (s *Stegano) NewWriter() *StegWriter {
	return &StegWriter{
		image: s.image,
		index: 0,
	}
}

func (s *Stegano) IsEncoded() (bool, error) {
	var magic [4]byte
	reader := s.NewReader()
	n, err := reader.Read(magic[:])
	if err != nil {
		return false, err
	}
	if n != 4 {
		return false, fmt.Errorf("Magic number not fully read")
	}
	return string(magic[:]) == "STEG", nil
}

func (s *Stegano) Encode(data []byte) error {
	writer := s.NewWriter()

	// Write a magic number to identify the presence of hidden data
	_, err := writer.Write([]byte("STEG"))
	if err != nil {
		return err
	}

	// Write the length of the hidden data (4 bytes, big-endian)
	length := int64(len(data))
	if err := writer.WriteInt64(length); err != nil {
		return err
	}

	// Write the hidden data
	if _, err := writer.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *Stegano) Decode() (data []byte, err error) {
	var magic [4]byte
	reader := s.NewReader()

	// Read the magic number to verify the presence of hidden data
	n, err := reader.Read(magic[:])
	if err != nil {
		return data, err
	}
	if n != 4 {
		return data, fmt.Errorf("Magic number not fully read")
	}
	if string(magic[:]) != "STEG" {
		return data, fmt.Errorf("Invalid magic number: %s", string(magic[:]))
	}

	// Read the length of the hidden data (4 bytes, big-endian)
	length, err := reader.ReadInt64()
	if err != nil {
		return data, err
	}
	if length < 0 {
		return data, fmt.Errorf("Invalid data length: %d", length)
	}

	// Read the hidden data
	data = make([]byte, length)
	n, err = reader.Read(data)
	if err != nil {
		return data, err
	}
	if int64(n) != length {
		return data, fmt.Errorf("Data not fully read")
	}
	return data, nil
}

// Loads an image from the specified file path and returns it as an NRGBA image.
func LoadNRGBA(filePath string) (*image.NRGBA, error) {
	// Open the image file
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	// Decode the image
	var img image.Image
	switch filepath.Ext(filePath) {
	case ".png":
		img, err = png.Decode(imgFile)
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(imgFile)
	default:
		return nil, fmt.Errorf("Unsupported format: %v", filepath.Ext(filePath))
	}
	if err != nil {
		return nil, err
	}
	// Convert the image to NRGBA format if it's not already
	if img.ColorModel() != color.NRGBAModel {
		newimg := image.NewNRGBA(img.Bounds())
		draw.Draw(newimg, newimg.Bounds(), img, image.Point{}, draw.Src)
		img = newimg
	}

	return img.(*image.NRGBA), nil
}

// Saves an NRGBA image to the specified file path in PNG format.
func SaveNRGBA(img *image.NRGBA, filePath string) error {
	// Create the output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode the image to PNG format and save it
	err = png.Encode(outFile, img.SubImage(img.Rect))
	if err != nil {
		return err
	}

	return nil
}
