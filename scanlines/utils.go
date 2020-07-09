package scanlines

import (
	"bytes"
	"compress/zlib"
	"errors"
)

// SectionsMap ways of dividing a 8 bit number
var SectionsMap = [][]byte{
	[]byte{1, 1, 1, 1, 1, 1, 1, 1},
	[]byte{2, 2, 2, 2},
	[]byte{3, 3, 2},
	[]byte{4, 4},
	[]byte{5, 3},
	[]byte{6, 2},
	[]byte{7, 1},
	[]byte{8},
}

// ErrDataTooSmall Data is too small to hide anything in it
var ErrDataTooSmall = errors.New("Scanline is too small to hide anything in it")

// Div Divides a 8 bit number $N into sections, each section contaning at most $M bits
//
// Usage:
//	b = Div(0b00101010, 3) == [[0b001, 3], [0b010, 3], [0b10, 2]]
//	b[0][0] is the first 3 bits of the number
//	b[0][1] is the length of b[0][0]
//
// **@param** _number byte_ The byte to be divided
//
// **@param** _parts int_ How many parts should it be divided into
//
// **@return** _[][2]byte_ An array containg all parts of the divided number into $N parts
//
func Div(number byte, parts int) [][2]byte {
	sections := SectionsMap[parts-1]
	na := [][2]byte{}
	r := byte(8)
	for i := 0; i < len(sections); i++ {
		n := sections[i]
		r = r - n

		d := number >> r

		na = append(na, [2]byte{d, n})

		number = number & ((1 << r) - 1)
	}

	return na
}

// Unite is the exact opposite of Div.
// It the takes the broken byte pieces and reassemble it, into a full byte.
func Unite(parts [][2]byte) byte {
	var n byte = 0
	var r int = 8
	for i := 0; i < len(parts); i++ {
		v, s := parts[i][0], parts[i][1]

		n = n + (v << (r - int(s)))
		r -= int(s)
	}
	return n
}

// decompress Decompress data stored in the zlib format
//
// **@param** _data []byte_ The compressed data
//
// **@return** _([]byte, error)_ The decompressed data and an error
func decompress(data []byte) ([]byte, error) {
	dataBuf := bytes.NewBuffer(data)
	defer dataBuf.Reset()

	var readBuff bytes.Buffer

	reader, err := zlib.NewReader(dataBuf)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	// Write decompressed data to a buffer
	if _, err := readBuff.ReadFrom(reader); err != nil {
		return nil, err
	}

	return readBuff.Bytes(), nil
}

// compress Compress data to the zlib format
//
// **@param** _data []byte_ Data to be compressed
//
// **@return** _([]byte, error)_ The compressed data and an error
func compress(data []byte) ([]byte, error) {
	var writeBuff bytes.Buffer

	writer, err := zlib.NewWriterLevel(&writeBuff, zlib.BestSpeed)
	if err != nil {
		return nil, err
	}

	writer.Write(data)
	writer.Close()

	return writeBuff.Bytes(), nil
}

// Generare mask (((1 << n) - 1) << n).toString(2)
