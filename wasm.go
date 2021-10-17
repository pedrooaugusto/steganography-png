//+build ignore
// How to build this for wasm:
// GOOS=js GOARCH=wasm go build -o main.wasm wasm.go && mv main.wasm webapp/public/wasm

package main

import (
	"fmt"
	"steganographypng/png"
	"syscall/js"
)

func valueToByteArray(v js.Value) []byte {
	binImage := make([]byte, v.Length())
	js.CopyBytesToGo(binImage, v)

	return binImage
}

func hideData(this js.Value, args []js.Value) interface{} {
	inputImage := valueToByteArray(args[0])
	bytesToHide := valueToByteArray(args[1])
	bitloss := args[2].Int()
	callback := args[3]

	pngParsed, err := png.Parse(inputImage)
	if err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	if err := pngParsed.HideData(bytesToHide, bitloss); err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	rawBytes := pngParsed.ToBytes()

	dst := js.Global().Get("Uint8Array").New(len(rawBytes))

	js.CopyBytesToJS(dst, rawBytes)

	callback.Invoke(js.Null(), dst)

	return nil
}

func revealData(this js.Value, inputs []js.Value) interface{} {
	binImage := valueToByteArray(inputs[0])
	callback := inputs[1]

	pngParsed, err := png.Parse(binImage)
	if err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	dataSize, bitloss, err := pngParsed.GetParams()
	if err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	messsage := make([]byte, dataSize)
	if err := pngParsed.RevealData(messsage, bitloss); err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	dst := js.Global().Get("Uint8Array").New(len(messsage))

	js.CopyBytesToJS(dst, messsage)

	callback.Invoke(js.Null(), dst)

	return nil
}

func toString(this js.Value, inputs []js.Value) interface{} {
	binImage := valueToByteArray(inputs[0])
	callback := inputs[1]

	pngParsed, err := png.Parse(binImage)
	if err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	str := pngParsed.String()

	callback.Invoke(js.Null(), str)

	return nil
}

func main() {
	c := make(chan bool)

	fmt.Println("[Steganography PNG Module Loaded]")

	js.Global().Set("PNG", make(map[string]interface{}))

	js.Global().Get("PNG").Set("hideData", js.FuncOf(hideData))
	js.Global().Get("PNG").Set("revealData", js.FuncOf(revealData))
	js.Global().Get("PNG").Set("toString", js.FuncOf(toString))

	// Keep it open
	<-c
}
