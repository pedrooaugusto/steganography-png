//+build ignore
// GOOS=js GOARCH=wasm go build -o main.wasm wasm.go | mv main.wasm webapp/public/wasm

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

func hideBytes(this js.Value, args []js.Value) interface{} {
	inputImage := valueToByteArray(args[0])
	bytesToHide := valueToByteArray(args[1])
	bitloss := args[2].Int()
	callback := args[3]

	pngParsed, err := png.Parse(inputImage)
	if err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	if err := pngParsed.HideBytes(bytesToHide, bitloss); err != nil {
		callback.Invoke(err.Error(), js.Null())
		return nil
	}

	rawBytes := pngParsed.ToBytes()

	dst := js.Global().Get("Uint8Array").New(len(rawBytes))

	js.CopyBytesToJS(dst, rawBytes)

	callback.Invoke(js.Null(), dst)

	return nil
}

func unhideBytes(this js.Value, inputs []js.Value) interface{} {
	binImage := valueToByteArray(inputs[0])
	callback := inputs[2]

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
	if err := pngParsed.UnhideBytes(&messsage, bitloss); err != nil {
		panic(err)
	}

	dst := js.Global().Get("Uint8Array").New(len(messsage))

	js.CopyBytesToJS(dst, messsage)

	callback.Invoke(js.Null(), dst)

	return nil
}

func main() {
	c := make(chan bool)

	fmt.Println("Hi")
	js.Global().Set("hideBytes", js.FuncOf(hideBytes))
	js.Global().Set("unhideBytes", js.FuncOf(unhideBytes))

	<-c
}
