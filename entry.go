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
	_ "bytes"
	"fmt"
	"io/ioutil"
	"os"
	"steganographypng/png"
	"strconv"
)

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}

/*
 * Hide file *f inside png file *i
 * steganographypng hide -i /path/to/png -f /path/to/file
 *
 * Hide string *s inside png file *i
 * steganographypng hide -i /path/to/png -s "hello, how are you?"
 *
 * Look for some hideen content inside *i with lenght of *l
 * steganographypng unhide -i /path/to/png -l 300 -s output.png
 *x
 */

func main() {
	var operation string = os.Args[1]
	var imagePath string = os.Args[3]

	binImage, err := ioutil.ReadFile(imagePath)
	if err != nil {
		panic(err)
	}

	pngParsed, err := png.Parse(binImage)
	if err != nil {
		panic(err)
	}

	if operation == "hide" {
		var bytesToHide []byte

		if os.Args[4] == "-f" {
			bytesToHide, err = ioutil.ReadFile(os.Args[5])
			if err != nil {
				panic(err)
			}

		} else {
			bytesToHide = []byte(os.Args[5])
			fmt.Println("Hiding: ")
			fmt.Println(string(bytesToHide))
		}

		if err := pngParsed.HideData(bytesToHide, 8); err != nil {
			panic(err)
		}

		ioutil.WriteFile(getImage("/suspicous.png"), pngParsed.ToBytes(), 0644)

		fmt.Println("\nTotal Size: " + strconv.Itoa(len(bytesToHide)))
		fmt.Println("File name: suspicous.png")

	} else {
		length, err := strconv.Atoi(os.Args[5])
		if err != nil {
			panic(err)
		}

		messsage := make([]byte, length)

		if err := pngParsed.RevealData(messsage, 8); err != nil {
			panic(err)
		}

		if len(os.Args)-1 <= 5 {
			fmt.Println("Bytes representaion")
			fmt.Println(messsage)
			fmt.Println("\nChar representation")
			fmt.Println(string(messsage))
		} else {
			fmt.Println("Writing to file: " + os.Args[7])
			ioutil.WriteFile(os.Args[7], messsage, 0644)
		}
	}

}
