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

func sub(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))
	newScanlineData = append(newScanlineData, 1)

	for i := 0; i < len(current); i++ {
		prior := byte(0)

		if i-bpp >= 0 {
			prior = current[i-bpp]
		}

		newScanlineData = append(newScanlineData, (current[i]-prior)&0xff)
	}

	return newScanlineData
}

func undo_sub(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))

	for i := 0; i < len(current); i++ {
		prior := byte(0)

		if i-bpp >= 0 {
			prior = newScanlineData[i-bpp]
		}

		newScanlineData = append(newScanlineData, (current[i]+prior)&0xff)
	}

	newScanlineData = append([]byte{1}, newScanlineData...)

	return newScanlineData
}

// Sub Filter and Unfilter a byte array using the filter method 1 (Sub)
func Sub(current, previous []byte, undo bool, header map[string]interface{}) []byte {

	bpp := header["bpp"].(int)

	if undo {
		return undo_sub(current, previous, bpp)
	} else {
		return sub(current, previous, bpp)
	}
}
