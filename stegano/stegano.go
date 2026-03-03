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

const MAGIC string = "STEG"
const BYTE_LEN int64 = 12
const INT64_LEN int64 = BYTE_LEN * 8

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
		// index: 0,
	}
}

func (s *Stegano) NewWriter() *StegWriter {
	return &StegWriter{
		image: s.image,
		index: 0,
	}
}

func (s *Stegano) IsEncoded() bool {
	reader := s.NewReader()
	return reader.ReadHeader() == MAGIC
}

func (s *Stegano) Encode(fileName string, data []byte) error {
	writer := s.NewWriter()

	// Write a magic number to identify the presence of hidden data
	_, err := writer.Write([]byte(MAGIC))
	if err != nil {
		return err
	}

	// Write the name of the file or "_TEXT_" to mean embedded message
	writer.WriteString(fileName)

	// Write the length of the hidden data (8 bytes, big-endian)
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

func (s *Stegano) Decode() (fileName string, data []byte, err error) {

	// Check the magic number to verify the presence of hidden data
	if !s.IsEncoded() {
		return fileName, data, fmt.Errorf("Image not encoded")
	}

	// Skip the header (we just tested it)
	reader := s.NewReader()
	reader.Skip(int64(len(MAGIC)))

	// Read the filename
	fileName, err = reader.ReadString()

	// Read the length of the hidden data (8 bytes, big-endian)
	length, err := reader.ReadInt64()
	if err != nil {
		return fileName, data, err
	}
	if length < 0 {
		return fileName, data, fmt.Errorf("Invalid data length: %d", length)
	}

	// Read the hidden data
	data = make([]byte, length)
	n, err := reader.Read(data)
	if err != nil {
		return fileName, data, err
	}
	if int64(n) != length {
		return fileName, data, fmt.Errorf("Data not fully read")
	}
	return fileName, data, nil
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
