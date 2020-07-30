# Steganography & PNG & Webassembly
> Hide <i>anything</i> inside a png image

<a href="https://github.com/pedrooaugusto/steganography-png/actions">
    <img alt="Toolkit audit status" src="https://github.com/pedrooaugusto/steganography-png/workflows/Run Tests and Build/badge.svg" />
</a>

### What is Steganography ?
> Steganography is the practice of concealing a file, message, image, or video within another file, message, image, or video. The word steganography comes from Greek steganographia, which combines the words steganós, meaning "covered or concealed", and -graphia meaning "writing".
><br/><br/><i>From Wikipedia</i>

So basically is the practice of hiding and retrieving a file inside another file.

### How does this module works ?

#### PNG
This is a *steganography-png* module, meaning that the input file is restricted to png images only.
You can use it to hide **any type of file** inside a png image.

#### How to hide information inside a PNG file ?
Every png image is strucured as follows :
```
[Header] + [Chunk 1] + [Chunk 2] + ... + [Chunk n]
```
http://www.libpng.org/pub/png/spec/1.2/PNG-Structure.html

Each chunk has a *Type* and a *Data* property. Between the many types of chunks is the **IDAT** chunk, this is where information about the pixels of the image
are stored (...) and is also perfect place to hide information.

#### The IDAT Chunk

**Demo at: https://pedrooaugusto.github.io/live/steganography/**

![Test Image 4](https://github.com/pedrooaugusto/steganography-png/blob/master/webapp/preview.png)

