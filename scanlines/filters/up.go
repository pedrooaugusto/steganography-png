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

// Up Filter and Unfilter a byte array using the filter method 2 (up)
func Up(current, previous []byte, undo bool, header map[string]interface{}) []byte {

	bpp := header["bpp"].(int)

	if undo {
		return undo_up(current, previous, bpp)
	} else {
		return up(current, previous, bpp)
	}
}
