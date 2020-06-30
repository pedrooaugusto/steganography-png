package png

// SectionsMap ways of dividing a 8 bit number
var SectionsMap = [][]byte{
	[]byte{1, 1, 1, 1, 1, 1, 1, 1},
	[]byte{2, 2, 2, 2},
	[]byte{3, 3, 2},
	[]byte{4, 4},
	[]byte{5, 3},
	[]byte{6, 2},
	[]byte{7, 1},
	[]byte{8},
}

// Div Divides a 8 bit number $N into sections, each section contaning at most $M bits
//
// Usage:
//	b = Div(0b00101010, 3) == [[0b001, 3], [0b010, 3], [0b10, 2]]
//	b[0][0] is the first 3 bits of the number
//	b[0][1] is the length of b[0][0]
//
// **@param** _number byte_ The byte to be divided
//
// **@param** _parts int_ How many parts should it be divided into
//
// **@return** _[][2]byte_ An array containg all parts of the divided number into $N parts
//
func Div(number byte, parts int) [][2]byte {
	sections := SectionsMap[parts-1]
	na := [][2]byte{}
	r := byte(8)
	for i := 0; i < len(sections); i++ {
		n := sections[i]
		r = r - n

		d := number >> r

		na = append(na, [2]byte{d, n})

		number = number & ((1 << r) - 1)
	}

	return na
}

// Unite is the exact opposite of Div.
// It the takes the broken byte pieces and reassemble it, into a full byte.
func Unite(parts [][2]byte) byte {
	var n byte = 0
	var r int = 8
	for i := 0; i < len(parts); i++ {
		v, s := parts[i][0], parts[i][1]

		n = n + (v << (r - int(s)))
		r -= int(s)
	}
	return n
}

// Generare mask (((1 << n) - 1) << n).toString(2)
