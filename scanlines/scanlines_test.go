package scanlines

import (
	"os"
	"reflect"
	"testing"

	"github.com/pedrooaugusto/steganography-png/chunk"
)

func TestFromChunks(t *testing.T) {
	chunks, err := getChunks(4)
	if err != nil {
		t.Error(err)
	}

	header := make(map[string]interface{})
	header["bpp"] = 4
	header["Height"] = uint32(4)

	scanlines, _, err := FromChunks(chunks, header)

	if err != nil {
		t.Errorf("\nError when converting to scanlines\n%s", err)
	}

	if scanlines.length != 4 {
		t.Errorf("\nWrong scanlines length\n 4 !== %d", scanlines.length)
	}
}

func TestCanComport(t *testing.T) {
	chunks, err := getChunks(2)
	if err != nil {
		t.Error(err)
	}

	header := make(map[string]interface{})
	header["Height"] = uint32(2)

	scanlines, _, err := FromChunks(chunks, header)

	if err != nil {
		t.Errorf("\nError when converting to scanlines\n%s", err)
	}

	if scanlines.canComport(500) {
		t.Errorf("\nThis scanaline should be unable to comport 500 bytes, since its size is only 110 bytes\n")
	}

	if !scanlines.canComport(2) {
		t.Errorf("\nThis scanline should be able to comport 2 bytes, since its size is 110 bytes\n")
	}
}

func TestToChunks(t *testing.T) {
	chunks, err := getChunks(2)
	if err != nil {
		t.Error(err)
	}

	header := make(map[string]interface{})
	header["Height"] = uint32(2)
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
	chunks, err := getChunks(4)
	if err != nil {
		t.Error(err)
	}

	header := make(map[string]interface{})
	header["Height"] = uint32(4)
	scanlines, size, err := FromChunks(chunks, header)

	if err != nil {
		t.Fatalf("\nError when converting to scanlines\n%s", err)
	}

	if err := scanlines.HideBytes([]byte("42"), []byte("text/plain"), 4); err != nil {
		t.Fatalf("\nError when hiding bytes in scanline\n%s", err)
	}

	chunks2, err := scanlines.ToChunks(size)
	if err != nil {
		t.Fatalf("\nError when converting back to chunks\n%s", err)
	}

	scanlines2, _, err := FromChunks(chunks2, header)

	if err != nil {
		t.Fatalf("\nError when converting to scanlines\n%s", err)
	}

	data2, dataType2, bitloss2, err := scanlines2.RevealBytes()
	if err != nil {
		t.Fatalf("\nError when hiding bytes in scanline\n%s", err)
	}

	if string(data2) != "42" || dataType2 != "text/plain" || bitloss2 != 4 {
		t.Fatalf("\nErrow hen retrieving data. Expected Result: 42\n")
	}
}

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}

// 268 bytes, 444
func getChunks(n int) ([]chunk.Chunk, error) {
	data := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum")

	c0, err := compress(data)
	if err != nil {
		return nil, err
	}

	step := len(c0) / n

	chunks := make([]chunk.Chunk, n)

	for i := 0; i < n; i++ {
		if (i+1)*step >= len(c0) {
			chunks[i] = chunk.CreateChunk(c0[i*step:], []byte("IDAT"))
		} else {
			chunks[i] = chunk.CreateChunk(c0[i*step:(i+1)*step], []byte("IDAT"))
		}
	}

	return chunks, nil
}
