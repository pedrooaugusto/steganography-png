# Steganography & PNG & Webassembly
> Hide <i>anything</i> inside a png image

<a href="https://github.com/pedrooaugusto/steganography-png/actions">
    <img alt="Toolkit audit status" src="https://github.com/pedrooaugusto/steganography-png/workflows/Run Tests and Build/badge.svg" />
</a>

### Demo
Online webassembly implementation: https://pedrooaugusto.github.io/steganography-png/

### What is Steganography ?
> Steganography is the practice of concealing a file, message, image, or video within another file, message, image, or video. The word steganography comes from Greek steganographia, which combines the words steganós, meaning "covered or concealed", and -graphia meaning "writing".
><br/><br/><i>From Wikipedia</i>

So basically is the practice of hiding and retrieving a file inside another file.

### How does this module work ?

#### PNG
This is a *steganography-png* module, meaning that the input file is restricted to png images only.
You can use it to hide **any type of file** inside a png image.

#### How does one hide information inside a PNG image ?
Every png image is structured as follows :
```
[Header] + [Chunk 1] + [Chunk 2] + ... + [Chunk n]
```
http://www.libpng.org/pub/png/spec/1.2/PNG-Structure.html

Each chunk has a *Type* and a *Data* property. Between the many types of chunks there is the **IDAT** chunk, this is where information about the pixels of the image
is stored (...) and is also the perfect spot to hide some piece of information.

That is how an IDAT chunk looks like:
```
Type: IDAT
Length: 16
Data: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16] (zlib compression)
Crc: 898323
```
https://www.w3.org/TR/PNG-Compression.html

The *Data* field is zipped, so in order to do any work with it, you must unzip it first. Once decompressed: the Data field is structured in N scanlines
each one containing M bytes. Hiding information inside a PNG image is a matter of editing some bytes in those scanlines.

This is how the decompressed *Data* field (scanlines) looks like:
```
1 234 212 121 098 069 032 035 067 012 000 000 000 012 011 001 000 (scanline 1)
2 034 099 253 018 012 232 110 073 120 002 111 000 120 111 063 001 (scanline 2)
2 007 002 042 015 013 234 100 038 000 010 200 120 110 188 037 010 (scanline 3)
```
*The first byte of each scanline is the filter type, you can not use it...*

If you want to hide the letter 'A' inside this image all you have to do is choose one scanline unfilter the bytes according with
the filter method and put the number 65 ('A') somewhere inside it.

https://www.w3.org/TR/PNG-Filters.html

### Usage

#### Web Client

https://pedrooaugusto.github.io/steganography-png/

Figma: https://www.figma.com/file/Mt6a7buAvM42YfVKpVIh8GYS/Steganography-Page?node-id=0%3A1

![Test Image 4](https://github.com/pedrooaugusto/steganography-png/blob/master/webapp/preview.png)

Video: https://drive.google.com/file/d/1_3SFULtktoeHTUg1sM8mLtlFUgMaQns0/view

#### CLI

Start the CLI with `go run .`

```
Usage ./steganography-png -o=[hide | reveal] -i=/path/to/png [OPTIONS]

 Where:
   -o  string: Operation to carry out (hide stuff or reveal stuff).
   -i  string: Path of the PNG image in which the operation will be done.
   -ss string: Secret plain text message to hide in the input image (Required if `-o=hide`).
   -sf Path  : Path of the secret file to hide inside the input image (Overrides `-ss`).
   -st string: Type of the content described by `-ss` or `-sf`. Eg: text/plain, text/html, audio/mp3 ... (Optional).
   -bl int   : How many BITS of the input image should be used to encode ONE BYTE of the secret. (Optional; Defaults to 8).

 Example:
   ./steganography-png -o=hide -i=./images/bisk.png -ss="Hello World!" // To hide 'Hello Wolrd' inside bisk.png
   ./steganography-png -o=hide -i=./images/bisk.png -sf=./pitou.jpg -st=image/jpeg -bl=1 //To hide the image 'pitou.jpg' inside bisk.png
   ./steganography-png -o=reveal -i=./images/killua.png // To search and reveal a hidden secret inside 'killua.png'

```
