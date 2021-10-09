package chunk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// Chunk represents one of the many chunks of a png file
type Chunk struct {
	id       int    // chunk index
	dataSize []byte // 4 bytes - data field size
	tipo     []byte // 4 bytes - [a-Z] letters only (chunk type)
	Data     []byte // $dataSize bytes
	crc      []byte // 4 bytes - CRC algorithm
}

// GetDataSize returns dataSize in integer
func (r *Chunk) GetDataSize() uint32 {
	if len(r.dataSize) == 0 {
		return 0
	}

	return binary.BigEndian.Uint32(r.dataSize)
}

// SetDataSize set chunk dataSize field
func (r *Chunk) SetDataSize(size []byte) {
	r.dataSize = size
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

// SetCRC sets the CRC
func (r *Chunk) SetCRC(ncrc []byte) {
	r.crc = ncrc
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

// Parse converts a byte array into chunk
func Parse(index *uint32, data []byte) Chunk {
	chunk := Chunk{}

	chunk.dataSize = append(chunk.dataSize, data[*index:*index+4]...)
	*index = *index + 4

	chunk.tipo = append(chunk.tipo, data[*index:*index+4]...)
	*index = *index + 4

	chunk.Data = append(chunk.Data, data[*index:*index+chunk.GetDataSize()]...)
	*index = *index + chunk.GetDataSize()

	chunk.crc = append(chunk.crc, data[*index:*index+4]...)
	*index = *index + 4

	return chunk
}

// CreateChunk create a new chunk with the desired options
func CreateChunk(data []byte, tipo []byte) Chunk {
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

// BuildIDATChunks will create a bunch of IDAT chunks from a compressed rawBytes array
func BuildIDATChunks(bytesBuff *bytes.Buffer, chunkSize uint32) []Chunk {
	defer bytesBuff.Reset()

	var rawBytes []byte = bytesBuff.Bytes()
	var pointer uint32 = 0
	var chunks []Chunk
	var length uint32 = uint32(len(rawBytes) - 1)
	var tipo []byte = []byte("IDAT")

	for {
		if pointer >= length {
			break
		}

		b := pointer
		e := pointer + chunkSize

		if e > length {
			e = length + 1
		}

		chunks = append(chunks, CreateChunk(rawBytes[b:e], tipo))

		pointer = e
	}

	return chunks
}
