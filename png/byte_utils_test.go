package png

import (
	"reflect"
	"testing"
)

func TestDiv(t *testing.T) {
	expected := [][2]byte{
		[2]byte{0, 2},
		[2]byte{2, 2},
		[2]byte{2, 2},
		[2]byte{2, 2},
	}

	got := Div(0b00101010, 2)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nDiv(0b00101010, 2) = %d; expected %d", got, expected)
	}
}

func TestDiv2(t *testing.T) {
	expected := [][2]byte{
		[2]byte{0b00000, 5},
		[2]byte{0b100, 3},
	}

	number := byte(0b00000100)
	parts := 5
	got := Div(number, parts)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nDiv(%b, %d) = %b; expected %b", number, parts, got, expected)
	}
}

func TestDiv3(t *testing.T) {
	expected := [][2]byte{
		[2]byte{0b001, 3},
		[2]byte{0b010, 3},
		[2]byte{0b10, 2},
	}

	number := byte(0b00101010)
	parts := 3
	got := Div(number, parts)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nDiv(%b, %d) = %b; expected %b", number, parts, got, expected)
	}
}

func TestDiv4(t *testing.T) {
	expected := [][2]byte{
		[2]byte{0b1000000, 7},
		[2]byte{0b1, 1},
	}

	number := byte(0b10000001)
	parts := 7
	got := Div(number, parts)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nDiv(%b, %d) = %b; expected %b", number, parts, got, expected)
	}
}

func TestUnite(t *testing.T) {
	var number, got byte
	var b int

	number = 65
	b = 5
	if got = Unite(Div(number, b)); got != number {
		t.Errorf("\nUnite(?) = %d; expected %d", got, number)
	}

	number = 12
	b = 6
	if got = Unite(Div(number, b)); got != number {
		t.Errorf("\nUnite(?) = %d; expected %d", got, number)
	}

	number = 128
	b = 2
	if got = Unite(Div(number, b)); got != number {
		t.Errorf("\nUnite(?) = %d; expected %d", got, number)
	}

	number = 255
	b = 1
	if got = Unite(Div(number, b)); got != number {
		t.Errorf("\nUnite(?) = %d; expected %d", got, number)
	}

	number = 0
	b = 8
	if got = Unite(Div(number, b)); got != number {
		t.Errorf("\nUnite(?) = %d; expected %d", got, number)
	}

	number = 42
	b = 1
	if got = Unite(Div(number, b)); got != number {
		t.Errorf("\nUnite(?) = %d; expected %d", got, number)
	}

	if got = Unite([][2]byte{[2]byte{0b0010, 4}, [2]byte{0b1010, 4}}); got != 42 {
		t.Errorf("\nUnite(?) = %d; expected %d", got, 42)
	}
}
