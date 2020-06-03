package png

import _ "fmt"
import "bytes"
import "io"
import "compress/zlib"
import "math"
import "errors"

// decompress decompress data stored in the zlib format
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

// PRESERVE At each scanline we can only compromise 20% of the bytes
const PRESERVE float32 = 0.2

// ReplaceData Replaces some bytes of bytess with bytes of data
func replaceData(buffer *bytes.Buffer, data *[]byte, height uint32) error {
	return findData(buffer, data, height, true)
}

// ReadData Replaces some bytes of bytess with bytes of data
func readData(buffer *bytes.Buffer, data *[]byte, height uint32) error {
	return findData(buffer, data, height, false)
}

func findData(buffer *bytes.Buffer, data *[]byte, height uint32, replace bool) error {
	bytess := buffer.Bytes()
	buffer = bytes.NewBuffer([]byte{})

	bytesPerScanline := uint32(len(bytess)) / height
	compromisedBytes := uint32(float32(bytesPerScanline-1) * PRESERVE) // first bit of the scanliine is the filter type

	if compromisedBytes < 1 { // Less than a bit per scanline is zero!
		return ErrDataTooSmall
	}

	dataSize := uint32(len(*data))

	if dataSize > (compromisedBytes * height) { // Data to hide is larger than image itself
		return ErrDataTooSmall
	}

	var scanline, added uint32 = 0, 0
	var neededScanlines float64 = math.Ceil(float64(uint32(dataSize / compromisedBytes)))
	var step uint32 = (bytesPerScanline - 1) / dataSize

	for i := 0; uint32(i) < dataSize; i++ {
		if added == compromisedBytes {
			added = 0
			scanline = scanline + (height / uint32(neededScanlines))
		}
		index := scanline * bytesPerScanline
		index = index + (1 + added*step)

		if replace {
			bytess[index] = (*data)[i] // 1 byte is too much image might get buzzy
		} else {
			(*data)[i] = bytess[index]
		}

		added = added + 1
	}

	_, err := buffer.Read(bytess)

	if err != io.EOF {
		return err
	}

	return nil
}
