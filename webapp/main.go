package main

import "fmt"
import "syscall/js"
import "steganographypng/png"

func valueToByteArray(v js.Value) []byte{
    binImage := make([]byte, v.Length())
    js.CopyBytesToGo(binImage, v)

    return binImage
}

func hideBytes(this js.Value, inputs []js.Value) interface{} {
    binImage := valueToByteArray(inputs[0])
    bytesToHide := valueToByteArray(inputs[1])
    callback := inputs[2]

    pngParsed, err := png.Parse(binImage)
	if err != nil {
        callback.Invoke(err.Error(), js.Null())
        return nil
    }

    if err := pngParsed.HideBytes(bytesToHide); err != nil {
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
    length := inputs[1].Int()
    callback := inputs[2]

    pngParsed, err := png.Parse(binImage)
	if err != nil {
        callback.Invoke(err.Error(), js.Null())
        return nil
    }

    messsage := make([]byte, length)
    if err := pngParsed.UnhideBytes(&messsage); err != nil {
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