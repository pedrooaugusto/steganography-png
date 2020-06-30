package png

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestInsertBytesIntoScanline(t *testing.T) {
	const bytesPerScanline = 12
	const scanlineIndex = 0
	const bitLoss = 4
	var availablePerScanline = uint32(math.Floor((bytesPerScanline - 1) * 1))

	var data []byte = []byte{'P', 'E', 'D'}
	var scanlines []byte = getScanlines(bytesPerScanline, 2, 0)

	insertBytesIntoScanline(data, &scanlines, scanlineIndex, availablePerScanline, bitLoss)

	expected := []byte{1, 5, 0, 0, 0, 4, 0, 5, 0, 4, 0, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(expected, scanlines) {
		log := fmt.Sprintf("\nData: %08b\nBitLoss: %d\n", data, bitLoss)

		log += printScanlines(byteArrayToIntArray(scanlines), bytesPerScanline, true, 0, true)

		t.Log(log)
		t.Log(scanlines)

		t.Fail()
	}
}

func TestInsertBytesIntoScanlineMultipleScanlines(t *testing.T) {
	const bytesPerScanline = 12
	const scanlineIndex = 0
	const bitLoss = 4
	var availablePerScanline = uint32(math.Floor((bytesPerScanline - 1) * 1))

	var data []byte = []byte{'P', 'E', 'D'}
	var scanlines []byte = getScanlines(bytesPerScanline, 2, 0)

	insertBytesIntoScanline(data, &scanlines, scanlineIndex*bytesPerScanline, availablePerScanline, bitLoss)
	insertBytesIntoScanline([]byte{'R', 'O'}, &scanlines, (scanlineIndex+1)*bytesPerScanline, availablePerScanline, bitLoss)

	expected := []byte{1, 5, 0, 0, 0, 4, 0, 5, 0, 4, 0, 4, 1, 5, 0, 0, 2, 0, 0, 4, 0, 0, 15, 0}
	if !reflect.DeepEqual(expected, scanlines) {
		log := fmt.Sprintf("\nData: %08b\nBitLoss: %d\n", data, bitLoss)

		log += printScanlines(byteArrayToIntArray(scanlines), bytesPerScanline, true, 0, true)

		t.Log(log)
		t.Log(scanlines)

		t.Fail()
	}
}

func TestInsertBytesIntoScanlineOneBit(t *testing.T) {
	const bytesPerScanline = 20
	const scanlineIndex = 0
	const bitLoss = 1
	var availablePerScanline = uint32(math.Floor((bytesPerScanline - 1) * 1))

	var data []byte = []byte{'4', '2'}
	var scanlines []byte = getScanlines(bytesPerScanline, 2, 0)

	insertBytesIntoScanline(data, &scanlines, scanlineIndex*bytesPerScanline, availablePerScanline, bitLoss)

	expected := []byte{1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(expected, scanlines) {
		log := fmt.Sprintf("\nData: %08b\nBitLoss: %d\n", data, bitLoss)

		log += printScanlines(byteArrayToIntArray(scanlines), bytesPerScanline, true, 0, true)

		t.Log(log)
		t.Log(scanlines)

		t.Fail()
	}
}

func TestGetBytesFromScanlineInsertBytesIntoScanline(t *testing.T) {
	hideAndRetrieveBytes([]byte("OLA MEU NOME"), 32, 3, 4, t)
	hideAndRetrieveBytes([]byte("GODDAMNIT"), 132, 1, 1, t)
	hideAndRetrieveBytes([]byte("FORTY-TWO"), 256, 3, 6, t)
	hideAndRetrieveBytes([]byte("HELLO, "), 100, 3, 8, t)
	hideAndRetrieveBytes([]byte("PEDRO"), 21, 3, 3, t)
}

func TestWriteDataReadData(t *testing.T) {
	// will fail: scanline width (5) is smaller than min (8)
	hideAndRetrieveData([]byte("0b00101010"), 5, 25, 1, true, t)
	// will fail: data to hide is bigger than scanlines combined
	hideAndRetrieveData([]byte("Get out you big piece of amazingness!"), 10, 2, 1, true, t)
	hideAndRetrieveData([]byte("May I have your attention, please!"), 14, 15, 4, false, t)
	hideAndRetrieveData([]byte("Will the real slim shady please stand up ?!"), 14, 9, 8, false, t)
}

func TestCompressDecompress(t *testing.T) {
	original := bytes.NewBuffer([]byte("JUST LOSE IT"))
	copyy := make([]byte, original.Len())
	copy(copyy, original.Bytes())

	compressed, err := compress(original)
	if err != nil {
		t.Fail()
	}

	k, err := decompress(&compressed)
	if err != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(copyy, k.Bytes()) {
		fmt.Print(original.Bytes())
		fmt.Print(k.Bytes())
		t.Errorf("Error values are not equivalent")
	}
}

func hideAndRetrieveData(data []byte, bytesPerScanline, nScanlines, bitloss int, shouldFail bool, t *testing.T) {
	buffer := bytes.NewBuffer(getScanlines(bytesPerScanline, nScanlines, -1))
	data2 := make([]byte, len(data))

	errW := WriteData(buffer, &data, bitloss, uint32(nScanlines))

	if !shouldFail && errW != nil {
		t.Errorf("Failed to write data\nReasson: %s", errW)
	}

	errR := ReadData(buffer, &data2, bitloss, uint32(nScanlines))

	if !shouldFail && errR != nil {
		t.Errorf("Failed to read data\nReasson: %s", errR)
	}

	if shouldFail && (errW == nil || errR == nil) {
		t.Errorf("Function was expected to fail but it did not!")
		t.FailNow()
	}

	if !shouldFail && !reflect.DeepEqual(data, data2) {
		log := printScanlines(byteArrayToIntArray(buffer.Bytes()), bytesPerScanline, false, 0, false)
		t.Errorf("\nValues are not equal:\n"+
			"\tExpected: %d\n\tGot     : %d\n"+
			"\n"+
			"Bitloss: %d\n"+
			"BytesPerScanline: %d"+
			"\n\n"+
			"LOG:\n%s", data, data2, bitloss, bytesPerScanline, log)
	}
}

func hideAndRetrieveBytes(data []byte, bytesPerScanline, nScanlines, bitloss int, t *testing.T) {
	var availablePerScanline = uint32(math.Floor(float64(bytesPerScanline-1) * 1))
	var scanlines []byte = getScanlines(bytesPerScanline, nScanlines, -1)

	insertBytesIntoScanline(data, &scanlines, 0, availablePerScanline, bitloss)

	var data2 []byte = make([]byte, len(data))

	getBytesFromScanline(data2, &scanlines, 0, availablePerScanline, bitloss)

	if !reflect.DeepEqual(data2, data) {
		t.Errorf("\nValues are not equal:\n\tExpected: %d\n\tGot     : %d\n\nBitloss: %d\nBytesPerScanline: %d", data, data2, bitloss, bytesPerScanline)
	}
}

func getScanlines(size, n int, fill int) []byte {
	d := make([]byte, size*n)

	for i := 0; i < len(d); i++ {
		if fill == -1 {
			d[i] = byte(rand.Intn(256))
		} else {
			if i%size == 0 {
				d[i] = 1
			} else {
				d[i] = byte(fill) // 46
			}
		}
	}

	return d
}

func printScanlines(scanlines []int, scanlineSize int, hiddenOnly bool, nullV int, byteRep bool) string {
	s := ""

	for i, j := 0, 0; i < len(scanlines); i++ {
		el := scanlines[i]

		if i%scanlineSize == 0 {
			s += "\n\t" + strconv.Itoa(j)

			if byteRep {
				s += " TTTTTTTT "
			} else {
				s += " TTT "
			}

			j++
		} else {
			if hiddenOnly && (el == nullV) {
				if byteRep {
					s += "........ "
				} else {
					s += "... "
				}
			} else {
				if byteRep {
					s += fmt.Sprintf("%08b ", (el))
				} else {
					s += fmt.Sprintf("%03d ", (el))
				}
			}
		}
	}

	return "\nScanlines:" + s
}

func byteArrayToIntArray(b []byte) []int {
	n := make([]int, len(b))
	for i := 0; i < len(n); i++ {
		n[i] = int(b[i])
	}

	return n
}
