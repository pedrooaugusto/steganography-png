package png

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestHideData(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/test3.png"))

	if err != nil {
		t.Fail()
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Fail()
	}

	data := []byte("Hello, Doctor!")
	data2 := make([]byte, len(data))

	err = pngParsed.HideData2(data, 1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	dataSize, bitloss, err := pngParsed.GetParams()

	err = pngParsed.RevealData2(data2, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(data, data2) {
		t.Errorf("\n%d\nis not equal to\n%d", data2, data)
	}

	if dataSize != 14 || bitloss != 1 {
		t.Errorf("Could not retrieve params stored in the image")
	}
}

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}
