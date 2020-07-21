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
	"flag"
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

	return path + file
}

// Usage:
//
// steganographypng -o=[hide | reveal] -i=/path/to/png [OPTIONS]
//
// OPTIONS (*available for hide*):
//
// **-ss**: Secret string to hide in the input image, must be expressed as quoted string.
//  eg: steganographypng -o=hide -i=/path/hi.png -ss="My Secret"
//
// **-sf**: Secret file to hide in the input image, must be expressed as path to a file.
//  eg: steganographypng -o=hide -i=/path/hi.png -sf=/path/to/secret/song.mp3
//
// **-bl**: Bit loss, how many bits use to encode a byte, must be a number between [1-8]. Default: 8
//  eg: steganographypng -o=hide -i=/path/hi.png -sf=/path/to/secret/song.mp3 -bl=8
//
func main() {
	var operation *string = flag.String("o", "", "Operation to carry out. Eg: 'hide' or 'reveal'")
	var input *string = flag.String("i", "", "an input image path")
	var secretS *string = flag.String("ss", "", "a secret message")
	var secretF *string = flag.String("sf", "", "a secret file")
	var bitloss *int = flag.Int("bl", 8, "bit loss value [1-8]")

	flag.Parse()

	if *operation != "hide" && *operation != "reveal" {
		fmt.Printf("\nInvalid operation type '%s'. Options are 'hide' or 'reveal'.\n", *operation)
		return
	}

	binImage, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("\nError opening input image '%s'.\nErr: %s\n", *input, err)
		return
	}

	pngParsed, err := png.Parse(binImage)
	if err != nil {
		fmt.Printf("\nError when parsing input image.\nErr: %s\n", err)
		return
	}

	if *operation == "hide" {
		var bytesToHide []byte
		if *secretF != "" {
			bytesToHide, err = ioutil.ReadFile(*secretF)
			if err != nil {
				fmt.Printf("\nError opening secret file '%s'.\nErr: %s\n", *secretF, err)
				return
			}
		} else if *secretS != "" {
			bytesToHide = []byte(*secretS)
		} else {
			fmt.Println("When using the 'hide' operation options '-ss' or '-sf' are required.")
			return
		}

		if err := pngParsed.HideData(bytesToHide, *bitloss); err != nil {
			panic(err)
		}

		ioutil.WriteFile(getImage("/new-image.png"), pngParsed.ToBytes(), 0644)

		fmt.Printf("\nCompleted! New image with hidden secret at %s\n", getImage("/new-image.png"))

	} else {
		dataSize, bitloss, err := pngParsed.GetParams()

		if err != nil {
			fmt.Printf("\nThere seems to be no hidden secret in this image.\nErr: %s\n", err)
			return
		}

		messsage := make([]byte, dataSize)
		if err := pngParsed.RevealData(messsage, bitloss); err != nil {
			fmt.Printf("\nError when revealing hidden data.\nErr: %s\n", err)
			return
		}

		ioutil.WriteFile(getImage("/secret-data.txt"), messsage, 0644)

		fmt.Printf("\nCompleted!\nSecret retrieved from image was saved at %s. File extension may differ...\n\n", getImage("/secret-data.txt"))

		n := 50

		if n > len(messsage) {
			n = len(messsage) - 1
		}

		fmt.Printf("\nPreview: %s\n", messsage[0:n])
	}
}
