package png

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"steganographypng/scanlines"
	"testing"
)

func TestHideData(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/test3.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	data := []byte("Hello, Doctor!")
	data2 := make([]byte, len(data))

	err = pngParsed.HideData(data, 1)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	dataSize, bitloss, err := pngParsed.GetParams()

	err = pngParsed.RevealData(data2, 1)
	if err != nil {
		t.Errorf("\nError when retrieving hidden data\n%s", err)
	}

	if !reflect.DeepEqual(data, data2) {
		t.Errorf("\n%d\nis not equal to\n%d", data2, data)
	}

	if dataSize != 14 || bitloss != 1 {
		t.Errorf("Could not retrieve params stored in the image")
	}
}

func TestHideDataRevealData(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/test3.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	data := []byte("Hello, Doctor!")

	err = pngParsed.HideData(data, 1)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	newImage := pngParsed.ToBytes()

	pngParsed, err = Parse(newImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	data2 := make([]byte, len(data))

	err = pngParsed.RevealData(data2, 1)

	if err != nil {
		t.Errorf("\nError when retrieving data\n%s", err)
	}

	if !reflect.DeepEqual(data, data2) {
		t.Errorf("\nHidden data is not equal to the retrieved data.\n")
	}
}

func TestScanlinesType(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/pitou.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	err = pngParsed.HideData([]byte("Hello my name is Pedro"), 1)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	ioutil.WriteFile(getImage("/../imagepack/suspicous.png"), pngParsed.ToBytes(), 0644)
}

func TestFilterAndUnfilter(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/pitou.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	s, _, err := scanlines.FromChunks(pngParsed.Chunks, pngParsed.GetHeader())

	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	original := make([][]byte, len(s.GetScanlines()))
	copy(original, s.GetScanlines())

	s.ToggleFilter(true, nil)

	if reflect.DeepEqual(original, s.GetScanlines()) {
		t.Errorf("\nUnfilter was not sucessful: filtered and unfiltered are equal!")
	}

	s.ToggleFilter(false, nil)

	if !reflect.DeepEqual(original, s.GetScanlines()) {
		t.Errorf("\nFilter was not sucessful: original filtered and new filtered are not equal!")

		for i := 0; i < len(original); i++ {
			if !reflect.DeepEqual(original[i], s.GetScanlines()[i]) {

				fmt.Printf("Index: %v\n", i)
				fmt.Printf("Size: %v\n", len(original[i])-len(s.GetScanlines()[i]))

				fmt.Println("Original")
				fmt.Println(original[i])
				fmt.Println("New")
				fmt.Println(s.GetScanlines()[i])
			}
		}
	}
}

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}
