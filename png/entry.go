package png

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/pedrooaugusto/steganography-png/chunk"
	scls "github.com/pedrooaugusto/steganography-png/scanlines"
)

// PNG Represents a PNG file as described at www.png.org
type PNG struct {
	header []byte
	Chunks []chunk.Chunk
}

// String Text representation of the PNG
func (r PNG) String() string {
	s := "PORTABLE NETWORK GRAPHICS\n\n"
	s += "Header: 137 PNG 13 10 26 10\n\n"

	s += "Parsed IDHR:\n"
	parsedidhr, _ := json.MarshalIndent(r.GetHeader(), "", "  ")
	s += string(parsedidhr) + "\n\n"

	for _, element := range r.Chunks {
		s += element.String()
		s += "\n"
	}

	return s
}

// ToBytes Reduces this PNG structure to a byte array
func (r *PNG) ToBytes() []byte {
	raw := []byte{}

	raw = append(raw, r.header...)

	for _, element := range r.Chunks {
		raw = append(raw, element.ToBytes()...)
	}

	return raw
}

// HideData Hide some data in this png file
func (r *PNG) HideData(data []byte, dataType string, bitloss int) error {
	scanlines, maxSize, err := scls.FromChunks(r.Chunks, r.GetHeader())
	if err != nil {
		return err
	}

	if err := scanlines.HideBytes(data, []byte(dataType), bitloss); err != nil {
		return err
	}

	chunks, err := scanlines.ToChunks(maxSize)
	if err != nil {
		return err
	}

	r.setIdatChunks(chunks)

	return nil
}

// RevealData Reveal hidden data in this png image
func (r *PNG) RevealData() (data []byte, dataType string, bitloss int, err error) {
	scanlines, _, err := scls.FromChunks(r.Chunks, r.GetHeader())

	if err != nil {
		return nil, "", 0, err
	}

	return scanlines.RevealBytes()
}

// setIdatChunks Replaces the current IDAT chunks with diffrent ones
func (r *PNG) setIdatChunks(chunks []chunk.Chunk) {
	// First we need to reorder the chunks
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
}

// GetHeight returns the image height
func (r *PNG) GetHeight() uint32 {
	return binary.BigEndian.Uint32(r.Chunks[0].Data[4:8])
}

// GetHeader returns the image header (IDHR chunk)
func (r *PNG) GetHeader() map[string]interface{} {
	header := make(map[string]interface{})

	header["Width"] = binary.BigEndian.Uint32(r.Chunks[0].Data[0:4])
	header["Height"] = binary.BigEndian.Uint32(r.Chunks[0].Data[4:8])
	header["Bit depth"] = uint32(r.Chunks[0].Data[8])
	header["Color type"] = uint32(r.Chunks[0].Data[9])
	header["Compression method"] = uint32(r.Chunks[0].Data[10])
	header["Filter method"] = uint32(r.Chunks[0].Data[11])
	header["Interlace method"] = uint32(r.Chunks[0].Data[11])
	header["bpp"] = calculateBPP(ColorType[header["Color type"].(uint32)], header["Bit depth"].(uint32))

	return header
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

// ColorType Mapping of color type to number of samples
var ColorType = []uint32{
	0: 1,
	2: 3,
	3: 3,
	4: 2,
	6: 4,
}

func calculateBPP(colorType, bitDepth uint32) int {
	bpp := int(colorType * bitDepth / 8)

	if bpp <= 0 {
		bpp = 1
	}

	return bpp
}
