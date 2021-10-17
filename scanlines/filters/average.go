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

// Average Filter and Unfilter a byte array using the filter method 3 (average)
func Average(current, previous []byte, undo bool, header map[string]interface{}) []byte {

	bpp := header["bpp"].(int)

	if undo {
		return undo_average(current, previous, bpp)
	} else {
		return average(current, previous, bpp)
	}
}
