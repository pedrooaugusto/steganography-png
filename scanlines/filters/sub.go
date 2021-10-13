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

// SubFilter Apply the sub algorithm to filter this array of bytes to better compression.
func SubFilter(scanlines [][]byte, currentScanline int, header map[string]interface{}) {
	if header["Color type"] == 3 {
		return
	}

	bpp := header["bpp"].(int)
	scanline := scanlines[currentScanline]
	newScanlineData := make([]byte, len(scanline))

	for i := 0; i < len(scanline); i++ {
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

	for i := 0; i < len(scanline); i++ {
		prior := byte(0)

		if i-bpp >= 0 {
			prior = newScanlineData[i-bpp]
		}

		newScanlineData[i] = scanline[i] + prior
	}

	scanlines[currentScanline] = newScanlineData
}
