/*

Filter type 1: Sub

The Sub filter transmits the difference between each byte and the
value of the corresponding byte of the prior pixel.
To compute the Sub filter, apply the following formula to each byte
of the scanline:

Sub(x) = Raw(x) - Raw(x-bpp)

More at: https://github.com/pedrooaugusto/steganography-png/issues/15

*/

package filters

import "fmt"

func sub(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))
	newScanlineData = append(newScanlineData, 1)

	// fmt.Println(current)

	for i := 0; i < len(current); i++ {
		prior := byte(0)

		if i-bpp >= 0 {
			prior = current[i-bpp]
		}

		newScanlineData = append(newScanlineData, (current[i]-prior)&0xff)
	}

	// fmt.Println(newScanlineData)

	return newScanlineData
}

func undo_sub(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))

	fmt.Println(current)

	for i := 0; i < len(current); i++ {
		prior := byte(0)

		if i-bpp >= 0 {
			prior = newScanlineData[i-bpp]
		}

		newScanlineData = append(newScanlineData, (current[i]+prior)&0xff)
	}

	newScanlineData = append([]byte{1}, newScanlineData...)

	fmt.Println(newScanlineData)

	return newScanlineData
}

// SubFilter Apply the sub algorithm to filter this array of bytes to better compression.
func SubFilter(scanlines [][]byte, currentScanline int, header map[string]interface{}) {
	if header["Color type"] == 3 {
		return
	}

	bpp := header["bpp"].(int)
	scanline := scanlines[currentScanline]
	newScanlineData := make([]byte, len(scanline))
	newScanlineData[0] = scanline[0]

	for i := 1; i < len(scanline); i++ {
		prior := byte(0)

		if i-bpp >= 0 {
			prior = scanline[i-bpp]
		}

		newScanlineData[i] = scanline[i] - prior
	}

	scanlines[currentScanline] = newScanlineData
}

// SubUnfilter Apply the sub algorithm to unfilter this array of filtered bytes.
func SubUnfilter(scanlines [][]byte, currentScanline int, header map[string]interface{}) {
	if header["Color type"] == 3 {
		return
	}

	bpp := header["bpp"].(int)
	scanline := scanlines[currentScanline]
	newScanlineData := make([]byte, len(scanline))
	newScanlineData[0] = scanline[0]

	for i := 1; i < len(scanline); i++ {
		prior := byte(0)

		if i-bpp >= 0 {
			prior = newScanlineData[i-bpp]
		}

		newScanlineData[i] = scanline[i] + prior
	}

	scanlines[currentScanline] = newScanlineData
}

func Sub(current, previous []byte, undo bool, header map[string]interface{}) []byte {

	bpp := header["bpp"].(int)

	if undo {
		return undo_sub(current, previous, bpp)
	} else {
		return sub(current, previous, bpp)
	}
}
