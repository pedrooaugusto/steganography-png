package png

import (
	"bytes"
	"compress/zlib"
	"errors"
	_ "fmt"
	"math"
)

// decompress Decompress data stored in the zlib format
//
// **@param** _data *bytes.Buffer_ The compressed data
//
// **@return** _(*bytes.Buffer, error)_ The decompressed data and an error
func decompress(data *bytes.Buffer) (bytes.Buffer, error) {
	defer data.Reset()

	var readBuff bytes.Buffer

	reader, err := zlib.NewReader(data)
	if err != nil {
		return readBuff, err
	}

	defer reader.Close()

	// Write decompressed data to a buffer
	if _, err := readBuff.ReadFrom(reader); err != nil {
		return readBuff, err
	}

	return readBuff, nil
}

// compress compress data to the zlib format
func compress(data *bytes.Buffer) (bytes.Buffer, error) {
	defer data.Reset()

	var writeBuff bytes.Buffer

	writer, err := zlib.NewWriterLevel(&writeBuff, zlib.BestSpeed)
	if err != nil {
		return writeBuff, err
	}

	writer.Write(data.Bytes())
	writer.Close()

	return writeBuff, nil
}

// ErrDataTooSmall Data is too small to hide anything in it
var ErrDataTooSmall = errors.New("Scanline is too small to hide anything in it")

// PRESERVEROUGH At each scaline we can compromise 70% of its bytes
const PRESERVEROUGH float32 = 0.7

// ReadData Retrieves hidden data inside the buffer
func ReadData(buffer *bytes.Buffer, data *[]byte, bitloss int, height uint32) error {
	return peekData(buffer, data, bitloss, height, false)
}

// WriteData Hide data somewhere in the buffer
func WriteData(buffer *bytes.Buffer, data *[]byte, bitloss int, height uint32) error {
	return peekData(buffer, data, bitloss, height, true)
}

func peekData(buffer *bytes.Buffer, data *[]byte, bitloss int, height uint32, replace bool) error {
	bufferBytes := buffer.Bytes()

	bytesPerScanline := uint32(len(bufferBytes)) / height

	// first bit of the scanliine is the filter type
	maxCompromisedBytesPerScanLine := uint32(float32(bytesPerScanline-1) * PRESERVEROUGH)

	dataSize := uint32(len(*data))

	// Less than a 8 bits per scanline is ultrageous!
	if maxCompromisedBytesPerScanLine < 8 || dataSize > (maxCompromisedBytesPerScanLine*height) {
		return ErrDataTooSmall
	}

	// How many bytes we goona need to encode one byte
	var bytesPerByte = uint32(len(SectionsMap[bitloss-1]))

	// one byte of $data is equivalent to (1 * n) bytes being wasted on $bufferBytes
	// how many bytes we can **really** fit into a scanline
	var maxFitted = uint32(math.Floor(float64(maxCompromisedBytesPerScanLine / bytesPerByte)))

	var i, currentScanline uint32 = 0, 0
	for {
		if i >= dataSize {
			break
		}

		end := i + maxFitted

		if end > dataSize {
			end = dataSize
		}

		scanlineIndex := uint32(currentScanline * (bytesPerScanline))
		if replace {
			insertBytesIntoScanline((*data)[i:end], &bufferBytes, scanlineIndex, maxCompromisedBytesPerScanLine, bitloss)
		} else {
			getBytesFromScanline((*data)[i:end], &bufferBytes, scanlineIndex, maxCompromisedBytesPerScanLine, bitloss)
		}

		currentScanline++
		i = end
	}

	return nil
}

func insertBytesIntoScanline(bytess []byte, scanlines *[]byte, scanlineIndex uint32, scanlineSize uint32, bitLoss int) error {
	return peekScanlineBytes(bytess, scanlines, scanlineIndex, scanlineSize, bitLoss, true)
}

func getBytesFromScanline(bytess []byte, scanlines *[]byte, scanlineIndex uint32, scanlineSize uint32, bitLoss int) error {
	return peekScanlineBytes(bytess, scanlines, scanlineIndex, scanlineSize, bitLoss, false)
}

func peekScanlineBytes(bytess []byte, scanlines *[]byte, scanlineIndex uint32, scanlineSize uint32, bitLoss int, replace bool) error {
	bytesPerByte := len(SectionsMap[bitLoss-1])

	actualSize := uint32(len(bytess) * bytesPerByte)
	step := (scanlineSize - 1) / (actualSize - 1)

	index := scanlineIndex + 1
	for i := 0; i < len(bytess); i++ {
		parts := Div(bytess[i], bitLoss)

		for j := 0; j < len(parts); j++ {
			bits := parts[j]

			var mask byte = (((1 << (8 - bits[1])) - 1) << bits[1])

			if replace {
				(*scanlines)[index] = ((*scanlines)[index] & mask) | bits[0]
			} else {
				parts[j][0] = (*scanlines)[index] & ((1 << bits[1]) - 1)
			}

			index = index + step
		}

		if !replace {
			bytess[i] = Unite(parts)
		}
	}

	return nil
}

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

*/
