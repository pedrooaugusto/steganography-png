package chunk

import (
	"testing"
)

func TestCreateChunk(t *testing.T) {
	data := []byte("THIS IS THE CHUNK DATA")

	chunk := CreateChunk(data, []byte("IDAT"))

	if chunk.GetType() != "IDAT" {
		t.Errorf("\nComparison failed chunk type is IDAT\n")
	}

	if string(chunk.Data) != string(data) {
		t.Errorf("\nComparison failed wrong chunk data value\n")
	}
}

func TestParse(t *testing.T) {
	data := []byte{0, 0, 0, 5, 'I', 'D', 'A', 'T', 'A', 'B', 'C', 'D', 'E', 0, 0, 0, 0}

	var index uint32 = 0
	chunk := Parse(&index, data)

	if chunk.GetType() != "IDAT" {
		t.Errorf("\nComparison failed chunk type is IDAT\n")
	}

	if string(chunk.Data) != "ABCDE" {
		t.Errorf("\nComparison failed wrong chunk data value\n")
	}
}
