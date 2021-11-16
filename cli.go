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

// We cant just insert random bytes on the IDAT chunk beacause of two reasons:
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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pedrooaugusto/steganography-png/png"
)

func getImage(file string) string {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path + file
}

/*
Usage ./steganography-png -o=[hide | reveal] -i=/path/to/png [OPTIONS]

 Where:
   -o  string: Operation to carry out (hide stuff or reveal stuff).
   -i  string: Path of the PNG image in which the operation will be done.
   -ss string: Secret plain text message to hide in the input image (Required if `-o=hide`).
   -sf Path  : Path of the secret file to hide inside the input image (Overrides `-ss`).
   -st string: Type of the content described by `-ss` or `-sf`.
   			Eg: text/plain, text/html, audio/mp3 ... (Optional).
   -bl int   : How many BITS of a byte of the input image should be used to encode ONE BYTE of the secret.
   			(Optional; Defaults to 8).

 Example:
   ./steganography-png -o=hide -i=./images/bisk.png -ss="Hello World!" // To hide 'Hello Wolrd' inside bisk.png
   ./steganography-png -o=hide -i=./images/bisk.png -sf=./pitou.jpg -st=image/jpeg -bl=1 //To hide the image 'pitou.jpg' inside bisk.png
   ./steganography-png -o=reveal -i=./images/killua.png // To search and reveal a hidden secret inside 'killua.png'

*/
func main() {
	var operation *string = flag.String("o", "", "Operation to carry out. Eg: 'hide' or 'reveal'")
	var input *string = flag.String("i", "", "an input image path")
	var secretS *string = flag.String("ss", "", "a plain text message that will be used as the secret")
	var secretF *string = flag.String("sf", "", "a path to the file that will be used as the secret")
	var secretFT *string = flag.String("st", "", "the secret file type (type/subtype eg: text/html, text/plain, image/png)")
	var bitloss *int = flag.Int("bl", 8, "bit loss value [1-8]")

	flag.Parse()
	checkArgs(operation, input, secretS, secretF, secretFT, bitloss)

	binImage, _ := ioutil.ReadFile(*input)

	pngParsed, err := png.Parse(binImage)
	if err != nil {
		fmt.Printf("\nError when parsing input image.\nErr: %s\n", err)
		os.Exit(1)
	}

	if *operation == "hide" {
		var bytesToHide []byte

		if *secretF != "" {
			bytesToHide, _ = ioutil.ReadFile(*secretF)
			*secretFT = *secretFT + path.Ext(*secretF)
		} else if *secretS != "" {
			bytesToHide = []byte(*secretS)
			if *secretFT == "" {
				*secretFT = "text/plain.txt"
			}
		}

		*secretFT = strings.ReplaceAll(*secretFT, "/", "-")

		if err := pngParsed.HideData(bytesToHide, *secretFT, *bitloss); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ioutil.WriteFile(getImage("/new-image.png"), pngParsed.ToBytes(), 0644)

		fmt.Printf("\nCompleted! New image with hidden secret at %s\n", getImage("/new-image.png"))

	} else {
		dataSize, dataType, bitloss, err := pngParsed.GetParams()

		if err != nil {
			fmt.Printf("\nThere seems to be no hidden secret in this image.\nErr: %s\n", err)
			return
		}

		messsage := make([]byte, dataSize)
		if err := pngParsed.RevealData(messsage, bitloss); err != nil {
			fmt.Printf("\nError when revealing hidden data.\nErr: %s\n", err)
			return
		}

		ioutil.WriteFile(getImage("/secret-data-"+dataType), messsage, 0644)

		fmt.Printf("\nCompleted!\nSecret retrieved from image was saved at %s. File extension may differ...\n\n", getImage("/secret-data-"+dataType))

		n := 50

		if n > len(messsage) {
			n = len(messsage) - 1
		}

		fmt.Printf("\nPreview: %s\n", messsage[0:n])
	}
}

func checkArgs(operation *string, input *string, secretS *string, secretF *string, secretFT *string, bitloss *int) {

	flag.Usage = func() {
		fmt.Println("Usage ./steganography-png -o=[hide | reveal] -i=/path/to/png [OPTIONS]")
		fmt.Println("\n Where:")
		fmt.Println("   -o  string: Operation to carry out (hide stuff or reveal stuff).")
		fmt.Println("   -i  string: Path of the PNG image in which the operation will be done.")
		fmt.Println("   -ss string: Secret plain text message to hide in the input image (Required if `-o=hide`).")
		fmt.Println("   -sf Path  : Path of the secret file to hide inside the input image (Overrides `-ss`).")
		fmt.Println("   -st string: Type of the content described by `-ss` or `-sf`. Eg: text/plain, text/html, audio/mp3 ... (Optional).")
		fmt.Println("   -bl int   : How many BITS of the input image should be used to encode ONE BYTE of the secret. (Optional; Defaults to 8).")
		fmt.Println("\n Example:")
		fmt.Println("   ./steganography-png -o=hide -i=./images/bisk.png -ss=\"Hello World!\" /*To hide 'Hello Wolrd' inside bisk.png*/")
		fmt.Println("   ./steganography-png -o=hide -i=./images/bisk.png -sf=./pitou.jpg -st=image/jpeg -bl=1 /*To hide the image 'pitou.jpg' inside bisk.png*/")
		fmt.Println("   ./steganography-png -o=reveal -i=./images/killua.png /*To search and reveal a hidden secret inside 'killua.png' */")
	}

	if *operation != "hide" && *operation != "reveal" {
		flag.Usage()
		fmt.Printf("\nInvalid operation type '%s'. Options are 'hide' or 'reveal'.\n", *operation)
		os.Exit(1)
	}

	if _, err := os.Stat(*input); os.IsNotExist(err) {
		flag.Usage()
		fmt.Printf("\nError opening input image '%s'.\nErr: %s\n", *input, err)
		os.Exit(1)
	}

	if *operation == "hide" && (*secretS == "" && *secretF == "") {
		flag.Usage()
		fmt.Println("\nWhen using the 'hide' operation options '-ss' or '-sf' are required.")
		os.Exit(1)
	}

	if *operation == "hide" && *secretF != "" {
		if _, err := os.Stat(*secretF); os.IsNotExist(err) {
			flag.Usage()
			fmt.Printf("\nError opening secret file '%s'.\nErr: %s\n", *secretF, err)
			os.Exit(1)
		}
	}
}
