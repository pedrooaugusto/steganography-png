// go install steganographypng

// https://blog.logrocket.com/interesting-use-cases-for-javascript-bitwise-operators/
// https://www.interviewcake.com/concept/java/bit-shift
// http://www.sunshine2k.de/articles/coding/crc/understanding_crc.html
// http://www.libpng.org/pub/png/spec/1.2/PNG-Rationale.html#R.PNG-file-signature
// http://www.libpng.org/pub/png/spec/1.2/PNG-CRCAppendix.html
// http://www.libpng.org/pub/png/spec/1.2/PNG-Structure.html#CRC-algorithm

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"steganographypng/png"
)

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	// Chane to path.join
	return path + file
}

func main() {
	// Get image path
	byteArray, err := ioutil.ReadFile(getImage("/img.png"))

	if err != nil {
		panic(err)
	}

	parsedPng, err := png.Parse(byteArray)

	parsedPng.HideBytes([]byte{'A'})

	ioutil.WriteFile(getImage("/img2.png"), parsedPng.ToBytes(), 0644)

	if err != nil {
		panic(err)
	}

	fmt.Println(parsedPng.ToString())
}
