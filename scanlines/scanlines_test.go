package scanlines

import (
	"io/ioutil"
	"os"
	"reflect"
	"steganographypng/chunk"
	"steganographypng/png"
	"testing"
)

func TestFromChunks(t *testing.T) {
	chunks, err := getChunks()
	if err != nil {
		t.Error(err)
	}

	scanlines, _, err := FromChunks(chunks, 2)

	if err != nil {
		t.Errorf("\nError when converting to scanlines\n%s", err)
	}

	if scanlines.length != 2 {
		t.Errorf("\nWrong scanlines length\n 2 !== %d", scanlines.length)
	}
}

func TestCanComport(t *testing.T) {
	chunks, err := getChunks()
	if err != nil {
		t.Error(err)
	}

	header := make(map[string]interface{})
	header["Header"] = 2

	scanlines, _, err := FromChunks(chunks, header)

	if err != nil {
		t.Errorf("\nError when converting to scanlines\n%s", err)
	}

	if scanlines.canComport(500) {
		t.Errorf("\nThis scanlien should be unable to comport 500 bytes, since its size is only 110 bytes\n")
	}

	if !scanlines.canComport(2) {
		t.Errorf("\nThis scanlien should be able to comport 2 bytes, since its size is 110 bytes\n")
	}
}

func TestToChunks(t *testing.T) {
	chunks, err := getChunks()
	if err != nil {
		t.Error(err)
	}

	header := make(map[string]interface{})
	header["Header"] = 2
	scanlines, size, err := FromChunks(chunks, header)

	if err != nil {
		t.Errorf("\nError when converting to scanlines\n%s", err)
	}

	chunks2, err := scanlines.ToChunks(size)

	if err != nil {
		t.Errorf("\nUnable to convert scanlines to chunk\n%s", err)
	}

	if len(chunks) != len(chunks2) {
		t.Errorf("\nScanlines chunk data is not the same as original chunk data\n")
	}

	if !reflect.DeepEqual(chunks[0].Data, chunks2[0].Data) {
		t.Errorf("\nScanlines chunk data is not the same as original chunk data\n")
	}
}

func TestHideBytesRevealBytes(t *testing.T) {
	chunks, err := getChunks()
	if err != nil {
		t.Error(err)
	}

	header := make(map[string]interface{})
	header["Header"] = 2
	scanlines, size, err := FromChunks(chunks, header)

	if err != nil {
		t.Errorf("\nError when converting to scanlines\n%s", err)
	}

	if err := scanlines.HideBytes([]byte("42"), 4); err != nil {
		t.Errorf("\nError when hiding bytes in scanline\n%s", err)
	}

	chunks2, err := scanlines.ToChunks(size)
	if err != nil {
		t.Errorf("\nError when converting back to chunks\n%s", err)
	}

	scanlines2, size, err := FromChunks(chunks2, header)

	if err != nil {
		t.Errorf("\nError when converting to scanlines\n%s", err)
	}

	data2 := []byte{0, 0}
	if err := scanlines2.RevealBytes(data2, 4); err != nil {
		t.Errorf("\nError when hiding bytes in scanline\n%s", err)
	}

	if string(data2) != "42" {
		t.Errorf("\nErrow hen retrieving data. Expected Result: 42\n")
	}
}

func TestFilterAndUnfilter(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/pitou.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := png.Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	s, _, err := FromChunks(pngParsed.Chunks, pngParsed.GetHeader())

	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	original := make([][]byte, len(s.scalines))
	copy(original, s.scalines)

	s.Unfilter()

	if reflect.DeepEqual(original, s.scalines) {
		t.Errorf("\nUnfilter was not sucessful: filtered and unfiltered are equal!")
	}

	s.Filter()

	if !reflect.DeepEqual(original, s.scalines) {
		t.Errorf("\nFilter was not sucessful: original filtered and new filtered are not equal!")
	}
}

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}

func getChunks() ([]chunk.Chunk, error) {
	data := []byte("The Quick Brown Fox Jumped Over the Zeu's Fence. Padin How does one count sheep before numbers were invented ?")
	c0, err := compress(data)
	if err != nil {
		return nil, err
	}

	d1 := chunk.CreateChunk(c0[0:55], []byte("IDAT"))
	d2 := chunk.CreateChunk(c0[55:], []byte("IDAT"))

	return []chunk.Chunk{d1, d2}, nil
}
