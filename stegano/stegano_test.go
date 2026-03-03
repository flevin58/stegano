package stegano

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initStegano() Stegano {
	rect := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 100, Y: 100},
	}
	return Stegano{
		filePath: "",
		image:    image.NewNRGBA(rect),
	}
}

func initStreams() (r *StegReader, w *StegWriter) {
	ste := initStegano()
	return ste.NewReader(), ste.NewWriter()
}

func TestHeader(t *testing.T) {
	r, w := initStreams()
	n, err := w.Write([]byte(MAGIC))
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	head := r.ReadHeader()
	assert.Equal(t, MAGIC, head)
	assert.Equal(t, w.Pos(), r.Pos())
	expected_pos := BYTE_LEN * int64(len(MAGIC))
	assert.Equal(t, expected_pos, r.Pos())
}

func TestByte(t *testing.T) {
	r, w := initStreams()
	w.writeByte(69)
	w.writeByte(42)
	assert.Equal(t, int64(24), w.Pos())
	b1, _ := r.readByte()
	assert.Equal(t, byte(69), b1)
	b2, _ := r.readByte()
	assert.Equal(t, byte(42), b2)
}

func TestInt64(t *testing.T) {
	r, w := initStreams()
	var expected int64 = 32045820348502
	w.WriteInt64(expected)
	actual, err := r.ReadInt64()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, INT64_LEN, r.Pos())
	assert.Equal(t, INT64_LEN, w.Pos())
}

func TestReadString(t *testing.T) {
	const TESTSTR = "Pippo"
	r, w := initStreams()
	w.WriteString(TESTSTR)
	str, _ := r.ReadString()
	assert.Equal(t, TESTSTR, str)
	assert.Equal(t, w.Pos(), r.Pos())
	var expected_pos int64 = INT64_LEN + BYTE_LEN*int64(len(TESTSTR))
	assert.Equal(t, expected_pos, r.Pos())
}

func TestMultipleInt64(t *testing.T) {
	expected := []int64{6456234, 34563456, 245234, 23452345}
	r, w := initStreams()
	for _, i64 := range expected {
		w.WriteInt64(i64)
	}
	actual := []int64{}
	for i := 0; i < len(expected); i++ {
		i64, _ := r.ReadInt64()
		actual = append(actual, i64)
	}
	assert.Equal(t, expected, actual)
	expected_pos := int64(len(expected)) * INT64_LEN
	assert.Equal(t, expected_pos, w.Pos())
}

func TestRealScenarioOnMemory(t *testing.T) {
	expected_str := "This is a hidden message inside a PNG file!"
	r, w := initStreams()
	w.Write([]byte(MAGIC))
	w.WriteString("_TEXT_")
	w.WriteString(expected_str)
	head := r.ReadHeader()
	assert.Equal(t, MAGIC, head)
	fileName, err := r.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "_TEXT_", fileName)
	str, err := r.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, w.Pos(), r.Pos())
	assert.Equal(t, expected_str, str)
}

func TestRealScenarioOnFile(t *testing.T) {
	expected_str := "This is a hidden message inside a PNG file!"
	img, err := Load("../assets/test.png")
	assert.NoError(t, err)
	assert.False(t, img.IsEncoded())
	w := img.NewWriter()
	w.Write([]byte(MAGIC))
	assert.True(t, img.IsEncoded())
	w.WriteString("_TEXT_")
	w.WriteString(expected_str)
	img.Save()

	img, err = Load("../assets/test_stegano.png")
	assert.NoError(t, err)
	assert.True(t, img.IsEncoded())
	r := img.NewReader()
	r.ReadHeader()
	str, err := r.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "_TEXT_", str)
	actual_str, err := r.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, expected_str, actual_str)
}

func TestEncode(t *testing.T) {
	expected_str := "This is a hidden message inside a PNG file!"
	img := initStegano()
	img.Encode("_TEXT_", []byte(expected_str))
	assert.True(t, img.IsEncoded())
	r := img.NewReader()
	r.ReadHeader()
	str, err := r.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "_TEXT_", str)
	actual_str, err := r.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, expected_str, actual_str)
}

func TestDecode(t *testing.T) {
	expected_str := "This is a hidden message inside a PNG file!"
	img := initStegano()
	img.Encode("_TEXT_", []byte(expected_str))
	assert.True(t, img.IsEncoded())
	fname, actual_data, err := img.Decode()
	assert.NoError(t, err)
	assert.Equal(t, fname, "_TEXT_")
	assert.Equal(t, expected_str, string(actual_data))
}
