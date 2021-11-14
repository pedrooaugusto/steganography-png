/*

Filter type 4: Paeth

The Paeth filter computes a simple linear function of the three
neighboring pixels (left, above, upper left), then chooses as
predictor the neighboring pixel closest to the computed value.
This technique is due to Alan W. Paeth [PAETH].

To compute the Paeth filter, apply the following formula to each
byte of the scanline:

	Paeth(x) = Raw(x) - PaethPredictor(Raw(x-bpp), Prior(x), Prior(x-bpp))

*/

package filters

func paeth(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))
	newScanlineData = append(newScanlineData, 4)

	for i := 0; i < len(current); i++ {
		rawBpp, prior, priorBpp := byte(0), previous[i], byte(0)

		if i-bpp >= 0 {
			rawBpp = current[i-bpp]
			priorBpp = previous[i-bpp]
		}

		sub := int(current[i]) - paethPredictor(int(rawBpp), int(prior), int(priorBpp))

		newScanlineData = append(newScanlineData, byte(sub&0xff))
	}

	return newScanlineData
}

func undo_paeth(current, previous []byte, bpp int) []byte {
	newScanlineData := make([]byte, 0, len(previous))

	for i := 0; i < len(current); i++ {
		rawBpp := byte(0)
		prior := previous[i]
		priorBpp := byte(0)

		if i-bpp >= 0 {
			rawBpp = newScanlineData[i-bpp]
			priorBpp = previous[i-bpp]
		}

		add := int(current[i]) + paethPredictor(int(rawBpp), int(prior), int(priorBpp))

		newScanlineData = append(newScanlineData, byte(add&0xff))
	}

	newScanlineData = append([]byte{4}, newScanlineData...)

	return newScanlineData
}

// Paeth Filter and Unfilter a byte array using the filter method 4 (Paeth)
func Paeth(current, previous []byte, undo bool, header map[string]interface{}) []byte {

	bpp := header["bpp"].(int)

	if undo {
		return undo_paeth(current, previous, bpp)
	} else {
		return paeth(current, previous, bpp)
	}
}

/**
  As described in: https://www.w3.org/TR/PNG-Filters.html

  function PaethPredictor (a, b, c)
  begin
       ; a = left, b = above, c = upper left
       p := a + b - c        ; initial estimate
       pa := abs(p - a)      ; distances to a, b, c
       pb := abs(p - b)
       pc := abs(p - c)
       ; return nearest of a,b,c,
       ; breaking ties in order a,b,c.
       if pa <= pb AND pa <= pc then return a
       else if pb <= pc then return b
       else return c
  end
*/

func paethPredictor(a, b, c int) int {
	p := (a + b - c)

	pa := abs(p - a)
	pb := abs(p - b)
	pc := abs(p - c)

	if pa <= pb && pa <= pc {
		return a
	} else if pb <= pc {
		return b
	} else {
		return c
	}
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}

	return n
}
