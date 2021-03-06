package scanlines

import (
	"bytes"
	"math"
	"steganographypng/chunk"
)

// TODO: Scanlines bytes are filtered, that is why before changing this
// bytes first we need to unfilter the scanline
// https://www.w3.org/TR/PNG-Filters.html

// PRESERVE At each scaline we can compromise 70% of its bytes
const PRESERVE float32 = 0.9 //0.0075

// Scanliens Represensts the **parsed** union of all IDAT chunks.
//
// When talking about PNG files _the real_ image data lives inside the IDAT
// chunk, however this data is compressed and scattered in multiple IDAT chunks.
// This class is the actual representation of those IDAT chunks.
type Scanliens struct {
	length                    int
	bytesPerScanline          uint32
	bytesAvailablePerScanline uint32
	scalines                  [][]byte
}

// HideBytes Tries to hide some bytes inside the specified scanline
func (t *Scanliens) HideBytes(data []byte, bitloss int) error {
	return t.peekScanlines(data, bitloss, true)
}

// RevealBytes Tries to hide some bytes inside the specified scanline
func (t *Scanliens) RevealBytes(data []byte, bitloss int) error {
	return t.peekScanlines(data, bitloss, false)
}

// ToChunks Converts the scanlines into an array of IDAT chunks
func (t *Scanliens) ToChunks(chunkSize uint32) ([]chunk.Chunk, error) {
	// A basic Array.flat
	data := []byte{}

	for _, el := range t.scalines {
		data = append(data, el...)
	}

	data, err := compress(data)
	if err != nil {
		return nil, err
	}

	var chunks []chunk.Chunk = chunk.BuildIDATChunks(bytes.NewBuffer(data), chunkSize)

	return chunks, nil
}

// CanComport Determines wheather this scanlines can comport an x amount of data
//
// **@param** _dataSize dataSize_ Size of the data to be hide
func (t *Scanliens) canComport(dataSize uint32) bool {

	// At minumum each scanline should have at leas 8 bytes available
	if t.bytesAvailablePerScanline < 8 {
		return false
	}

	if dataSize > t.bytesAvailablePerScanline*uint32(t.length) {
		return false
	}

	return true
}

// scanlinesFor Returns an array with the best scanliens to hide data
func (t *Scanliens) scanlinesFor(data []byte, bitloss int) ([]int, error) {
	actualDataSize := uint32(len(data) * len(SectionsMap[bitloss-1]))

	// Less than a 8 bits per scanline is ultrageous!
	if !t.canComport(actualDataSize) {
		return nil, ErrDataTooSmall
	}

	scanlinesNeeded := int(math.Ceil(float64(actualDataSize) / float64(t.bytesAvailablePerScanline)))

	if scanlinesNeeded == 1 {
		scanlinesNeeded = 2
	}

	step := int((t.length - 1) / (scanlinesNeeded - 1))

	scanlines := []int{}
	for i, j := 0, 0; i < scanlinesNeeded; i++ {
		scanlines = append(scanlines, j)

		j += step
	}

	// for a := 0; a < len(t.scalines); a++ {
	// 	b := t.scalines[a]
	// 	fmt.Printf("\nType: %d\n", b[0])
	// }

	return scanlines, nil
}

// peekScanlines Select the best bytes on each scanline to hide info
func (t *Scanliens) peekScanlines(data []byte, bitloss int, replace bool) error {
	scanliens, err := t.scanlinesFor(data, bitloss)
	if err != nil {
		return err
	}

	// How many bytes we goona need to encode one byte
	var bytesPerByte = uint32(len(SectionsMap[bitloss-1]))

	// how many bytes we can **really** fit into a scanline.
	var maxFitted = int(math.Floor(float64(t.bytesAvailablePerScanline / bytesPerByte)))

	for i, j := 0, 0; i < len(scanliens); i++ {

		end := j + maxFitted

		if end > len(data) {
			end = len(data)
		}

		err = t.peekScanlineBytes(Slice{data: data, begin: j, end: end}, scanliens[i], bitloss, replace)
		if err != nil {
			return err
		}

		j = end
	}

	return nil
}

// peekScanlineBytes Will uniformly select some bytes from the selected scanline.
// If replace is true it will also edit those bytes.
//
// **@param** _data *[]byte_ Data to be distributed
//
// **@param** _index int_ Scanline index
//
// **@param** _bitloss int_ How many bits of a byte we should use to enconde information
//
// **@param** _replace bool_ Replace the selected bytes or just retrieve it
func (t *Scanliens) peekScanlineBytes(data Slice, index int, bitloss int, replace bool) error {
	scanline := t.scalines[index]

	bytesPerByte := len(SectionsMap[bitloss-1])

	actualSize := uint32(data.size() * bytesPerByte)
	step := (t.bytesAvailablePerScanline - 1) / (actualSize - 1)

	k := uint32(1) // first byte of scanline is the filter type

	for i := 0; i < data.size(); i++ {
		parts := Div(data.get(i), bitloss)

		for j := 0; j < len(parts); j++ {
			bits := parts[j]

			var mask byte = (((1 << (8 - bits[1])) - 1) << bits[1])

			if replace {
				scanline[k] = ((scanline)[k] & mask) | bits[0]
			} else {
				parts[j][0] = (scanline)[k] & ((1 << bits[1]) - 1)
			}

			k = k + step
		}

		if !replace {
			data.set(i, Unite(parts))
		}
	}

	return nil
}

// FromChunks Creates a new Scanlines instance from an array of chunks
func FromChunks(chunks []chunk.Chunk, height uint32) (Scanliens, uint32, error) {
	data, maxSize := assembleIDATChunks(chunks)

	data, err := decompress(data)
	if err != nil {
		return Scanliens{}, 0, err
	}

	bytesPerScanline := uint32(len(data)) / height

	var scanlines [][]byte
	for i, size := uint32(0), uint32(len(data)); i < size; {
		end := i + bytesPerScanline

		if end > size {
			end = size
		}

		scanlines = append(scanlines, data[i:end])

		i = end
	}

	return Scanliens{
		length:                    len(scanlines),
		bytesPerScanline:          bytesPerScanline,
		bytesAvailablePerScanline: uint32(float32(bytesPerScanline-1) * PRESERVE),
		scalines:                  scanlines,
	}, maxSize, nil
}

// Takes an Array of chunk.Chunk and appends all into a big IDAT chunk
func assembleIDATChunks(chunks []chunk.Chunk) ([]byte, uint32) {
	data := []byte{}

	var maxSize uint32 = 0

	for _, element := range chunks {
		if element.GetType() == "IDAT" {
			data = append(data, element.Data...)
			maxSize = uint32(math.Max(float64(maxSize), float64(element.GetDataSize())))
		}
	}

	return data, maxSize
}
