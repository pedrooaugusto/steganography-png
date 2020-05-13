package png

import "errors"
import "steganographypng/chunk"
import _ "fmt"
import "bytes"

// PNG Represents a PNG file as described at www.png.org
type PNG struct {
	header []byte
	chunks []chunk.Chunk
}

// ToString PNG converts into a string
func (r PNG) ToString() string {
	s := "PORTABLE NETWORK GRAPHICS\n\n"
	s += "Header: 137 PNG 13 10 26 10\n"

	for _, element := range r.chunks {
		s += element.ToString()
		s += "\n"
	}

	return s
}

// ToBytes Reduces image to byte array
func (r PNG) ToBytes() []byte {
	raw := []byte{}

	raw = append(raw, r.header...)

	for _, element := range r.chunks {
		raw = append(raw, element.ToBytes()...)
	}

	return raw
}

// HideBytes HydeBytes Somewhere in the data array
func (r PNG) HideBytes(data []byte) {
	r.chunks[4].HideBytes(data)
}

func (r *PNG) parseHeader(index *uint32, data []byte) error {
	arr := []byte{137, 80, 78, 71, 13, 10, 26, 10}

	res := bytes.Compare(arr, data[0 : 8])

	if res != 0 {
		return errors.New("This is simply not a PNG file, header does not contain the constant bytes")
	}

	r.header = append(r.header, data[0 : 8]...)

	*index = uint32(len(arr))

	return nil
}

// Parse Will parse a byte array to PNG structure
func Parse(file []byte) (PNG, error) {
	var png PNG = PNG{}
	var index uint32 = 0

	err := png.parseHeader(&index, file)

	if err != nil {
		panic(err)
	}

	for ;; {
		png.chunks = append(png.chunks, chunk.Parse(&index, file))
		if index == uint32(len(file)) {
			return png, nil
		} else if index > uint32(len(file)) {
			return png, errors.New("Something went terrible worng")
		}
	}
}
