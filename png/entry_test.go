package png

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/pedrooaugusto/steganography-png/scanlines"
)

func TestHideData(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/bisky.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	data := []byte("Hello, Doctor!")

	err = pngParsed.HideData(data, "plain/text", 1)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	data2, dataType, bitloss, err := pngParsed.RevealData()

	if err != nil {
		t.Errorf("\nError when retrieving hidden data\n%s", err)
	}

	if !reflect.DeepEqual(data, data2) {
		t.Errorf("\n%d\nis not equal to\n%d", data2, data)
	}

	if len(data2) != 14 || bitloss != 1 || dataType != "plain/text" {
		t.Errorf("Could not retrieve params stored in the image")
	}
}

func TestHideDataRevealData(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/bisky.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	data := []byte("Hello, Doctor!")

	err = pngParsed.HideData(data, "plain/text", 1)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	newImage := pngParsed.ToBytes()

	pngParsed, err = Parse(newImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	data2, _, _, err := pngParsed.RevealData()

	if err != nil {
		t.Errorf("\nError when retrieving data\n%s", err)
	}

	if !reflect.DeepEqual(data, data2) {
		t.Errorf("\nHidden data is not equal to the retrieved data.\n")
	}
}

func TestHideStringAndSaveFile(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/pitou.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	text := `#!HTML
	
	<br/><br/>
	<span>
		Using this tool you can hide images, text and random binary data inside PNG images (provided there is enough space).<br/>
		If your plain text message starts with "#!HTML" it will be interpreted as such.<br/>
		<b>And now the hidden message inside the input image: </b><br/><br/>
		<iframe width="530" height="280" src="https://www.youtube.com/embed/uRQ7ecvU56k?start=12" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
	</span>`
	err = pngParsed.HideData([]byte(text), "text/html.html", 8)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	ioutil.WriteFile(getImage("/../imagepack/suspicious-pitou.png"), pngParsed.ToBytes(), 0644)
}

func TestRevealStringAndSaveFile(t *testing.T) {
	binImage, err := ioutil.ReadFile(getImage("/../imagepack/suspicious-pitou.png"))

	if err != nil {
		t.Errorf("\nError when opening file\n%s", err)
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	// dataSize, _, bitloss, err := pngParsed.GetParams()

	// text := make([]byte, dataSize)

	text, _, _, err := pngParsed.RevealData()

	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	expectedText := `#!HTML
	
	<br/><br/>
	<span>
		Using this tool you can hide images, text and random binary data inside PNG images (provided there is enough space).<br/>
		If your plain text message starts with "#!HTML" it will be interpreted as such.<br/>
		<b>And now the hidden message inside the input image: </b><br/><br/>
		<iframe width="530" height="280" src="https://www.youtube.com/embed/uRQ7ecvU56k?start=12" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
	</span>`

	if string(text) != expectedText {
		t.Errorf("\n%s\nis not equal to\n%s", string(text), expectedText)
	}
}

func TestHideImageAndSaveFile(t *testing.T) {
	inputImage, err1 := ioutil.ReadFile(getImage("/../imagepack/bisky.png"))
	secret, err2 := ioutil.ReadFile(getImage("/../imagepack/pitou.png"))

	if err1 != nil || err2 != nil {
		t.Errorf("\nError when opening file\n%s\n%s", err1, err2)
	}

	pngParsed, err := Parse(inputImage)
	if err != nil {
		t.Errorf("\nError when parsing png file\n%s", err)
	}

	err = pngParsed.HideData(secret, "image/png.png", 1)
	if err != nil {
		t.Errorf("\nError when hiding data\n%s", err)
	}

	ioutil.WriteFile(getImage("/../imagepack/suspicious-bisky.png"), pngParsed.ToBytes(), 0644)
}

func TestFilterAndUnfilterAreOpposites(t *testing.T) {
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

func TestFilterAndUnfilter(t *testing.T) {
	unfiltered := getScanlinesForImage("/../imagepack/unfiltered.png")

	// Paeth filter
	scanlines_paeth := getScanlinesForImage("/../imagepack/paeth_filtered.png")
	scanlines_paeth.ToggleFilter(true, nil)
	assertEqualScanlines(scanlines_paeth, unfiltered)

	// Average filter
	scanlines_average := getScanlinesForImage("/../imagepack/average_filtered.png")
	scanlines_average.ToggleFilter(true, nil)
	assertEqualScanlines(scanlines_average, unfiltered)

	// Sub filter
	scanlines_sub := getScanlinesForImage("/../imagepack/sub_filtered.png")
	scanlines_sub.ToggleFilter(true, nil)
	assertEqualScanlines(scanlines_sub, unfiltered)

	// Sup filter
	scanlines_sup := getScanlinesForImage("/../imagepack/sup_filtered.png")
	scanlines_sup.ToggleFilter(true, nil)
	assertEqualScanlines(scanlines_sup, unfiltered)
}

func assertEqualScanlines(a, b scanlines.Scanliens) {
	if len(a.GetScanlines()) != len(b.GetScanlines()) {
		panic("Different number of scanlines")
	}

	for i := 0; i < len(a.GetScanlines()); i++ {
		if !reflect.DeepEqual(a.Get(i)[1:], b.Get(i)[1:]) {
			panic("Scanlines are not equal")
		}
	}
}

func getScanlinesForImage(img string) scanlines.Scanliens {
	binImage, err := ioutil.ReadFile(getImage(img))

	if err != nil {
		panic(fmt.Sprintf("\nError when opening file\n%s", err))
	}

	pngParsed, err := Parse(binImage)
	if err != nil {
		panic(fmt.Sprintf("\nError when parsing png file\n%s", err))
	}

	s, _, err := scanlines.FromChunks(pngParsed.Chunks, pngParsed.GetHeader())

	if err != nil {
		panic(fmt.Sprintf("\nError when parsing png file\n%s", err))
	}

	return s
}

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}
