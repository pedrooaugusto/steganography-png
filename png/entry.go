package png

import "errors"
import "steganographypng/chunk"
import _ "fmt"
import "bytes"
import "encoding/binary"
import "math"

// PNG Represents a PNG file as described at www.png.org
type PNG struct {
	header []byte
	Chunks []chunk.Chunk
}

// String PNG converts into a string
func (r PNG) String() string {
	s := "PORTABLE NETWORK GRAPHICS\n\n"
	s += "Header: 137 PNG 13 10 26 10\n"

	for _, element := range r.Chunks {
		s += element.String()
		s += "\n"
	}

	return s
}

// ToBytes Reduces image to byte array
func (r *PNG) ToBytes() []byte {
	raw := []byte{}

	raw = append(raw, r.header...)

	for _, element := range r.Chunks {
		raw = append(raw, element.ToBytes()...)
	}

	return raw
}

// HideBytes HydeBytes Somewhere in the data array
func (r *PNG) HideBytes(data []byte) error {
	var compressedChunks bytes.Buffer
	defer compressedChunks.Reset()

	// Append all IDAT chunks
	var chunkSize uint32 = 0
	for _, element := range r.Chunks {
		if element.GetType() == "IDAT" {
			compressedChunks.Write(element.Data)
			chunkSize = uint32(math.Max(float64(chunkSize), float64(element.GetDataSize())))
		}
	}

	// Decompress IDAT chunks
	uncompressedChunks, err := decompress(&compressedChunks)
	if err != nil {
		return err
	}

	defer uncompressedChunks.Reset()

	// Write some data on IDAT chunks
	if err := replaceData(&uncompressedChunks, &data, r.GetHeight()); err != nil {
		return err
	}

	// Compress IDAT chunks
	compressedChunks, err = compress(&uncompressedChunks)
	if err != nil {
		return err
	}

	// Slice compressed big IDAT chunk into multiple smaller ones
	var chunks []chunk.Chunk = chunk.BuildIDATChunks(&compressedChunks, chunkSize)

	// Reorder chunks
	var chunks2 []chunk.Chunk
	for i := 0; i < len(r.Chunks); i++ {
		tipo := r.Chunks[i].GetType()
		if tipo != "IDAT" && tipo != "IEND" {
			chunks2 = append(chunks2, r.Chunks[i])
		}
	}

	chunks = append(chunks2, chunks...)
	chunks = append(chunks, r.Chunks[len(r.Chunks)-1])

	r.Chunks = chunks

	return nil
}

// UnhideBytes will look for hidden bytes into the image
func (r *PNG) UnhideBytes(data *[]byte) error {
	var compressedChunks bytes.Buffer
	defer compressedChunks.Reset()

	// Append all IDAT chunks
	for _, element := range r.Chunks {
		if element.GetType() == "IDAT" {
			compressedChunks.Write(element.Data)
		}
	}

	// Decompress IDAT chunks
	uncompressedChunks, err := decompress(&compressedChunks)
	if err != nil {
		return err
	}

	defer uncompressedChunks.Reset()

	// Write some data on IDAT chunks
	if err := readData(&uncompressedChunks, data, r.GetHeight()); err != nil {
		return err
	}

	return nil
}

// GetHeight returns the image height
func (r *PNG) GetHeight() uint32 {
	return binary.BigEndian.Uint32(r.Chunks[0].Data[4:8])
}

func (r *PNG) parseHeader(index *uint32, data []byte) error {
	arr := []byte{137, 80, 78, 71, 13, 10, 26, 10}

	res := bytes.Compare(arr, data[0:8])

	if res != 0 {
		return errors.New("this is simply not a PNG file, header does not contain the constant bytes")
	}

	r.header = append(r.header, data[0:8]...)

	*index = uint32(len(arr))

	return nil
}

// Parse Will parse a byte array to PNG structure
func Parse(file []byte) (PNG, error) {
	var png PNG = PNG{}
	var index uint32 = 0

	err := png.parseHeader(&index, file)

	if err != nil {
		return png, err
	}

	for {
		png.Chunks = append(png.Chunks, chunk.Parse(&index, file))
		if index == uint32(len(file)) {
			return png, nil
		} else if index > uint32(len(file)) {
			return png, errors.New("something went terrible wrong parsing the chunks")
		}
	}
}
