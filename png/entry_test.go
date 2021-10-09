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
	t.Skip()

	binImage, err := ioutil.ReadFile(getImage("/../imagepack/pitou.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	fmt.Print("\n\n===NEXT ONE===\n\n\n")

	err = pngParsed.HideData([]byte("CAESAR DIED IN THE IDLES OF MARCH AND AFTER THAT HIS BODY WAS BURNED IN A BIG FIRE IN THE CENTER OF ROME"), 1)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	ioutil.WriteFile(getImage("/../imagepack/suspicous.png"), pngParsed.ToBytes(), 0644)

	myscanlines, _, err := scanlines.FromChunks(pngParsed.Chunks, pngParsed.GetHeight())

	if err != nil {
		t.Errorf("\nError parsing scalines from chunks")
	}

	t.Log(scanlines.ToJson(pngParsed.GetHeader()))

	unfiltered := myscanlines.Unfilter()[0]

	// for i := 0; i < len(unfiltered); i++ {
	// 	t.Logf("%v === %v ", unfiltered[i], myscanlines.Get(0)[i])
	// }

	filtered := myscanlines.Filter([][]byte{unfiltered})[0]

	for i := 0; i < len(filtered); i++ {
		t.Logf("%v === %v ", filtered[i], myscanlines.Get(0)[i])
	}

	// t.Fail()
}

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}
