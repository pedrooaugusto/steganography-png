/*

Filter type 3: Average

*/

package filters

func average(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))
	newScanlineData = append(newScanlineData, 3)

	for i := 0; i < len(current); i++ {
		prior := previous[i]
		rawBpp := byte(0)

		if i-bpp >= 0 {
			rawBpp = current[i-bpp]
		}

		newScanlineData = append(newScanlineData, (current[i]-((rawBpp+prior)/2))&0xff)
	}

	return newScanlineData
}

func undo_average(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))

	for i := 0; i < len(current); i++ {
		prior := previous[i]
		rawBpp := byte(0)

		if i-bpp >= 0 {
			rawBpp = newScanlineData[i-bpp]
		}

		newScanlineData = append(newScanlineData, (current[i]+((rawBpp+prior)/2))&0xff)
	}

	newScanlineData = append([]byte{3}, newScanlineData...)

	return newScanlineData
}

// AverageFilter Apply the Average algorithm to filter this array of bytes to better compression.
func AverageFilter(scanlines [][]byte, currentScanline int, header map[string]interface{}) {
	if header["Color type"] == 3 {
		return
	}

	bpp := header["bpp"].(int)
	scanline := scanlines[currentScanline]
	newScanlineData := make([]byte, len(scanline))
	newScanlineData[0] = scanline[0]

	for i := 1; i < len(scanline); i++ {
		prior, rawBpp := byte(0), byte(0)

		if currentScanline-1 >= 0 {
			prior = scanlines[currentScanline-1][i]
		}

		if i-bpp >= 0 {
			rawBpp = scanline[i-bpp]
		}

		sum := (rawBpp + prior) / 2
		newScanlineData[i] = scanline[i] - sum
	}

	scanlines[currentScanline] = newScanlineData
}

// AverageUnfilter Apply the Average algorithm to unfilter this array of filtered bytes.
func AverageUnfilter(scanlines [][]byte, currentScanline int, header map[string]interface{}) {
	if header["Color type"] == 3 {
		return
	}

	bpp := header["bpp"].(int)
	scanline := scanlines[currentScanline]
	newScanlineData := make([]byte, len(scanline))
	newScanlineData[0] = scanline[0]

	for i := 1; i < len(scanline); i++ {
		prior, rawBpp := byte(0), byte(0)

		if currentScanline-1 >= 0 {
			prior = scanlines[currentScanline-1][i]
		}

		if i-bpp >= 0 {
			rawBpp = newScanlineData[i-bpp]
		}

		sum := (rawBpp + prior) / 2
		newScanlineData[i] = scanline[i] + sum
	}

	scanlines[currentScanline] = newScanlineData
}

func Average(current, previous []byte, undo bool, header map[string]interface{}) []byte {

	bpp := header["bpp"].(int)

	if undo {
		return undo_average(current, previous, bpp)
	} else {
		return average(current, previous, bpp)
	}
}
