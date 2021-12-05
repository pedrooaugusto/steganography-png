//+build ignore
// How to build this for wasm (just run):
// GOOS=js GOARCH=wasm go build -o main.wasm wasm.go && mv main.wasm webapp/public/wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/pedrooaugusto/steganography-png/png"
)

func valueToByteArray(v js.Value) []byte {
	binImage := make([]byte, v.Length())
	js.CopyBytesToGo(binImage, v)

	return binImage
}

func hideData(this js.Value, args []js.Value) interface{} {
	inputImage := valueToByteArray(args[0])
	dataToHide := valueToByteArray(args[1])
	dataType := args[2].String()
	bitloss := args[3].Int()
	callback := args[4]

	pngParsed, err := png.Parse(inputImage)
	if err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	if err := pngParsed.HideData(dataToHide, dataType, bitloss); err != nil {
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

	messsage, dataType, _, err := pngParsed.RevealData()

	if err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	dst := js.Global().Get("Uint8Array").New(len(messsage))

	js.CopyBytesToJS(dst, messsage)

	callback.Invoke(js.Null(), dst, dataType)

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

	fmt.Println("[Wasm Steganography PNG Module Loaded]")
	fmt.Println("['go-worker' Scope: window.PNG]")

	js.Global().Set("PNG", make(map[string]interface{}))

	js.Global().Get("PNG").Set("hideData", js.FuncOf(hideData))
	js.Global().Get("PNG").Set("revealData", js.FuncOf(revealData))
	js.Global().Get("PNG").Set("toString", js.FuncOf(toString))

	// Keep it open (channels??)
	<-c
}
