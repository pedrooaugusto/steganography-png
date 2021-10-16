/*

Filter type 2: Sub

*/

package filters

func up(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))
	newScanlineData = append(newScanlineData, 2)

	for i := 0; i < len(current); i++ {
		prior := previous[i]
		newScanlineData = append(newScanlineData, (current[i]-prior)&0xff)
	}

	return newScanlineData
}

func undo_up(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))

	for i := 0; i < len(current); i++ {
		prior := previous[i]

		newScanlineData = append(newScanlineData, (current[i]+prior)&0xff)
	}

	newScanlineData = append([]byte{2}, newScanlineData...)

	return newScanlineData
}

// UpFilter Apply the Up algorithm to filter this array of bytes to better compression.
func UpFilter(scanlines [][]byte, currentScanline int, header map[string]interface{}) {
	if header["Color type"] == 3 {
		return
	}

	scanline := scanlines[currentScanline]
	newScanlineData := make([]byte, len(scanline))
	newScanlineData[0] = scanline[0]

	for i := 1; i < len(scanline); i++ {
		prior := byte(0)

		if currentScanline-1 >= 0 {
			prior = scanlines[currentScanline-1][i]
		}

		newScanlineData[i] = scanline[i] - prior
	}

	scanlines[currentScanline] = newScanlineData
}

// UpUnfilter Apply the Up algorithm to unfilter this array of filtered bytes.
func UpUnfilter(scanlines [][]byte, currentScanline int, header map[string]interface{}) {
	if header["Color type"] == 3 {
		return
	}

	scanline := scanlines[currentScanline]
	newScanlineData := make([]byte, len(scanline))
	newScanlineData[0] = scanline[0]

	for i := 1; i < len(scanline); i++ {
		prior := byte(0)

		if currentScanline-1 >= 0 {
			prior = scanlines[currentScanline-1][i]
		}

		newScanlineData[i] = scanline[i] + prior
	}

	scanlines[currentScanline] = newScanlineData
}

func Up(current, previous []byte, undo bool, header map[string]interface{}) []byte {

	bpp := header["bpp"].(int)

	if undo {
		return undo_up(current, previous, bpp)
	} else {
		return up(current, previous, bpp)
	}
}
