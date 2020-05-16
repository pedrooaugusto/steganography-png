// go install steganographypng

// PNG Stuff
// General: https://blog.logrocket.com/interesting-use-cases-for-javascript-bitwise-operators/
// General: https://www.interviewcake.com/concept/java/bit-shift
// General: http://www.sunshine2k.de/articles/coding/crc/understanding_crc.html
// General: http://www.libpng.org/pub/png/spec/1.2/PNG-Rationale.html#R.PNG-file-signature
// General: http://www.libpng.org/pub/png/spec/1.2/PNG-CRCAppendix.html
// General: http://www.libpng.org/pub/png/spec/1.2/PNG-Structure.html#CRC-algorithm
// Chunks: http://www.libpng.org/pub/png/spec/1.2/PNG-Chunks.html

// CHUNK IDAT and Compression algorithm
// http://www.libpng.org/pub/png/spec/1.2/PNG-Compression.html
// https://tools.ietf.org/html/rfc1950

// We can just insert random bytes on the IDAT chunk beacause of two reassons
// 1. IDAT data is compressed using the ZILIB algorithm
// 2. The data is filtered, thorough  a filter algorithm to ensure
//    max compression.
//
// So in order to create a PNG file you must have a raw
// stream of bytes called "scanlines" then you filter that stream
// and finally compress using zlib.
//
// After decompression you can divide the array into $HEIGHT scanlines
// each containing [N / $HEIGHT] bytes.
//
// To alter anything we must do the reverse proccess


package main

import (
	_ "fmt"
	"io/ioutil"
	_ "io"
	"os"
	"steganographypng/png"
	_ "bytes"
	_ "compress/zlib"
)

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	// Change to path.join
	return path + file
}

func main() {
	// Get image path
	byteArray, err := ioutil.ReadFile(getImage("/test3.png"))

	if err != nil {
		panic(err)
	}

	parsedPng, err := png.Parse(byteArray)

	if err:= parsedPng.HideBytes([]byte("PEDRO")); err!=nil {
		panic(err)
	}

	ioutil.WriteFile(getImage("/suspicous.png"), parsedPng.ToBytes(), 0644)

}
