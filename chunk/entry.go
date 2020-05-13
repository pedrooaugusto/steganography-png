package chunk

import "fmt"
import "encoding/binary"
import  "hash/crc32"

// Chunk represents one of the many chunks of a png file
type Chunk struct {
	id int // chunk index
	dataSize []byte // 4 bytes - data field size
	tipo []byte // 4 bytes - [a-Z] letters only
	data []byte // $dataSize bytes
	crc []byte // 4 bytes - CRC algorithm
	secret []byte // n bytes with make up a secret
}

func (r *Chunk) getDataSize() uint32 {
	return binary.BigEndian.Uint32(r.dataSize)
}

func (r *Chunk) getType() string {
	return string(r.tipo)
}

func (r *Chunk) getCRC() uint32 {
	return binary.BigEndian.Uint32(r.crc)
}

func (r *Chunk) calcCRC() uint32 {
	return crc32.ChecksumIEEE(append(r.tipo, r.data...))
}

// ToString String representation of chunk
func (r *Chunk) ToString() string {
	s := fmt.Sprint("Data size: ", r.getDataSize()) + "\n"
	s += "Type: " + r.getType() + "\n"
	s += "Data: [...]\n"
	s += fmt.Sprintf("CRC:  %v", r.getCRC()) + "\n"
	s += fmt.Sprintf("CRC': %v", r.calcCRC()) + "\n"

	return s
}

// ToBytes Chunk byte representation
func (r *Chunk) ToBytes() []byte {
	raw := []byte{}

	raw = append(raw, r.dataSize...)
	raw = append(raw, r.tipo...)
	raw = append(raw, r.data...)
	raw = append(raw, r.crc...)

	return raw
}

// HideBytes HydeBytes Somewhere in the data array
func (r *Chunk) HideBytes(data []byte) {
	step := int(r.getDataSize() / uint32(len(data)))
	index := 2

	for _, element := range data {
		r.data[index] = element
		index = index + step
	}

	newCRC := make([]byte, 4)
	binary.BigEndian.PutUint32(newCRC, r.calcCRC())
	r.crc = newCRC
}

// Parse converts a byte array into chunk
func Parse(index *uint32, data []byte) Chunk{
	chunk := Chunk{}

	chunk.dataSize = append(chunk.dataSize, data[*index : *index + 4]...)
	*index = *index + 4

	chunk.tipo = append(chunk.tipo, data[*index : *index + 4]...)
	*index = *index + 4

	chunk.data = append(chunk.data, data[*index : *index + chunk.getDataSize()]...)
	*index = *index + chunk.getDataSize()

	chunk.crc = append(chunk.crc, data[*index : *index + 4]...)
	*index = *index + 4

	return chunk
}
