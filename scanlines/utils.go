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

/**
[10001000, 10001000, 10001000, 10001000, 10001000, 10001000, 10001000, 10001000, 10001000]
[11 10 00 01, 11 10 00 01]

1 = 4

step = (len(buff) - 1) / (len(data) - 1)  [ (8 - 1) / (4 - 1) = 2]

-------------------------------------------------
step: 2
00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
xx    xx    xx    xx    xx    xx    xx    xx

(7)*2 = 14
(0)*2 = 0
(4)*2 = 8
---- ok !
-------------------------------------------------


step: 19/4 = 4
00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20
   xx          xx          xx          xx          xx

(7)*2 = 14
(0)*2 = 0
(4)*2 = 8
---- ok !



[0a, 1, 2a, 3, 4a, 5, 6a, 7]


maxPerScanLine = 4
ola
neededScanLines = (3 * (8 / 1)) / 4 = 6
[
    [0, 1o, 2, 3o, 4, 5o, 6, 7o, 8],
    [0, 1o, 2, 3o, 4, 5o, 6, 7o, 8],
    [0, 1l, 2, 3l, 4, 5l, 6, 7l, 8],
    [0, 1l, 2, 3l, 4, 5l, 6, 7l, 8],
]


// remember
var sectionsMap = new Array(9)
sectionsMap[1] = [1, 1, 1, 1, 1, 1, 1, 1]
sectionsMap[2] = [2, 2, 2, 2]
sectionsMap[3] = [3, 3, 2]
sectionsMap[4] = [4, 4]
sectionsMap[5] = [5, 3]
sectionsMap[6] = [6, 2]
sectionsMap[7] = [7, 1]
sectionsMap[8] = [8]

function div(number, parts) {
	console.time('ola')
    let sections = sectionsMap[parts]
    const na = new Array()
    for (let i = 0; i < sections.length; i++) {
        const n = sections[i]
        const d = (number >> 0) & ((1 << n) - 1)

        na.splice(0, 0, d.toString(2).padStart(n, "0"))

        number = number >> n
	}

	console.timeEnd('ola')

    return na;
}

bytesPerScanline: 30
available: 15

data: 32

(data * (8 | 4 | 3 | 2 | *1) )/available = 3 scanlines
(data * (*8 | 4 | 3 | 2 | 1) )/available = 18 scanlines

[8], [7]
[1]



[[11, 3], [0, 3], [1, 2]] --> 01100001

0b01100000 + 0b00000 + 0b1

-----------


0 1 2 3 4 5 6 7 8 9
e 1 2 e 4 5 e 7 8 e


10 / 3
3

*/
