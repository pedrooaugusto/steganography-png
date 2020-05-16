package chunk

import (
	"fmt"
	"encoding/binary"
	"hash/crc32"
	"bytes"
	"compress/zlib"
	"errors"
	"math"
)

// Chunk represents one of the many chunks of a png file
type Chunk struct {
	id int // chunk index
	dataSize []byte // 4 bytes - data field size
	tipo []byte // 4 bytes - [a-Z] letters only
	Data []byte // $dataSize bytes
	crc []byte // 4 bytes - CRC algorithm
}

// GetDataSize returns dataSize in integer
func (r *Chunk) GetDataSize() uint32 {
	if len(r.dataSize) == 0 {
		return 0
	}

	return binary.BigEndian.Uint32(r.dataSize)
}

// GetType returns chunk type
func (r *Chunk) GetType() string {
	if len(r.dataSize) == 0 {
		return "undefined"
	}

	return string(r.tipo)
}

// GetCRC returns CRC in integer
func (r *Chunk) GetCRC() uint32 {
	if len(r.crc) == 0 {
		return 0
	}

	return binary.BigEndian.Uint32(r.crc)
}

// CalcCRC returns a recalculated version of the CRC
func (r *Chunk) CalcCRC() uint32 {
	if len(r.tipo) == 0 {
		return 0
	}

	return crc32.ChecksumIEEE(append(r.tipo, r.Data...))
}

// String String representation of chunk
func (r Chunk) String() string {
	s := fmt.Sprint("Data size: ", r.GetDataSize()) + "\n"
	s += "Type: " + r.GetType() + "\n"
	s += "Data: [...]\n"
	s += fmt.Sprintf("CRC:  %v", r.GetCRC()) + "\n"
	s += fmt.Sprintf("CRC': %v", r.CalcCRC()) + "\n"

	return s
}

// ToBytes Chunk byte representation
func (r *Chunk) ToBytes() []byte {
	raw := []byte{}

	raw = append(raw, r.dataSize...)
	raw = append(raw, r.tipo...)
	raw = append(raw, r.Data...)
	raw = append(raw, r.crc...)

	return raw
}

func decompressData(data []byte) []byte {
	// First of all decompress [inflate](zlib) data filed
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	// Write decompressed data to a buffer
	var readBuff bytes.Buffer
	if _, err := readBuff.ReadFrom(reader); err != nil {
		panic(err)
	}

	reader.Close()

	defer readBuff.Reset()

	return readBuff.Bytes()
}

func compressData(data []byte) []byte {
	var writeBuff bytes.Buffer
	writer, err := zlib.NewWriterLevel(&writeBuff, zlib.BestCompression)
	if err != nil {
		writer.Close()
		panic(err)
	}

	writer.Write(data)
	writer.Close()

	defer writeBuff.Reset()

	return writeBuff.Bytes()
}

// ErrTooSmallExeception Data is too small to hide anything in it
var ErrTooSmallExeception = errors.New("Scanline is too small to hide anything in it")

// PRESERVE At each scanline we can only compromise 20% of the bytes
const PRESERVE float32 = 0.2

func replaceData(bytess *[]byte, data []byte, height uint32) error {
	bytesPerScanline := uint32(len(*bytess)) / height
	compromisedBytes := uint32(float32(bytesPerScanline - 1) * PRESERVE) // first bit of the scanliine is the filter type

	if compromisedBytes < 1 { // Less than a bit per scanline is zero!
		return ErrTooSmallExeception
	}

	if uint32(len(data)) > (compromisedBytes * height) { // Data to hide is larger than image itself
		return ErrTooSmallExeception
	}

	var scanline, added uint32 = 0, 0
	var neededScanlines float64 = math.Ceil(float64(uint32(len(data)) / compromisedBytes))
	step := ((bytesPerScanline - 1) / compromisedBytes)

	for _, item := range data {
		if added == compromisedBytes {
			added = 0
			scanline = scanline + (height / uint32(neededScanlines))
		}
		index := scanline * bytesPerScanline // which scanline
		index = 1 + index + added * step

		(*bytess)[index] = item
		added = added + 1
	}

	return nil
}

// HideBytes HydeBytes Somewhere in the data array
func (r *Chunk) HideBytes(data []byte, height uint32) {
	// First of all decompress [inflate](zlib) data filed
	inflateBytes := decompressData(r.Data)

	// fmt.Println("Before")
	// fmt.Println(inflateBytes)

	// Function that given a raw byte array, something
	// to hide and bytes per scaline is capable of hiding
	// stuff in the better way possible inside the scanlines.
	// inflateBytes[3] = byte('A')
	err := replaceData(&inflateBytes, []byte{'P', 'E', 'D', 'R', 'O'}, height)

	if err != nil {
		panic(err)
	}

	fmt.Println("\nAfter")
	fmt.Printf("%s", inflateBytes)

	deflateBytes := compressData(inflateBytes)

	// Setting new data size
	newDataSize := make([]byte, 4)
	binary.BigEndian.PutUint32(newDataSize, uint32(len(deflateBytes)))
	r.dataSize = newDataSize

	// Setting new data
	r.Data = deflateBytes

	// Setting new CRC
	newCRC := make([]byte, 4)
	binary.BigEndian.PutUint32(newCRC, r.CalcCRC())
	r.crc = newCRC
}

// Parse converts a byte array into chunk
func Parse(index *uint32, data []byte) Chunk{
	chunk := Chunk{}

	chunk.dataSize = append(chunk.dataSize, data[*index : *index + 4]...)
	*index = *index + 4

	chunk.tipo = append(chunk.tipo, data[*index : *index + 4]...)
	*index = *index + 4

	chunk.Data = append(chunk.Data, data[*index : *index + chunk.GetDataSize()]...)
	*index = *index + chunk.GetDataSize()

	chunk.crc = append(chunk.crc, data[*index : *index + 4]...)
	*index = *index + 4

	return chunk
}

// CreateChunk create a new chunk with the desired options
func CreateChunk(data []byte, tipo []byte, ) Chunk {
	chunk := Chunk{}

	chunk.Data = data

	chunk.tipo = tipo

	newSize := make([]byte, 4)
	binary.BigEndian.PutUint32(newSize, uint32(len(data)))
	chunk.dataSize = newSize

	newCRC := make([]byte, 4)
	binary.BigEndian.PutUint32(newCRC, chunk.CalcCRC())
	chunk.crc = newCRC

	return chunk
}
