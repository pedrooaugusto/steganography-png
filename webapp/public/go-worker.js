self.importScripts('wasm/wasm_exec.js')

onmessage = function(event) {
    if (event.data.operationName === 'hideData') {
        const { inputImage, data, dataType, bitLoss, id } = event.data

        self.PNG.hideData(inputImage, data, dataType, bitLoss, function(err, data, dataType = "") {
            self.postMessage({ type: 'OperationResponse', id, error: err, payload: { data, dataType }})
        })
    } else if (event.data.operationName === 'revealData') {
        self.PNG.revealData(event.data.inputImage, function(err, data, dataType = "") {
            self.postMessage({ type: 'OperationResponse', id: event.data.id, error: err, payload: { data, dataType }})
        })
    } else if (event.data.operationName === 'toString') {
        self.PNG.toString(event.data.inputImage, function(err, data) {
            self.postMessage({ type: 'OperationResponse', id: event.data.id, error: err, payload: { data }})
        })
    }
}

onerror = function(event) {
    console.error(event.message)
}

try {
    const go = new Go()

    WebAssembly.instantiateStreaming(fetch('wasm/main.wasm'), go.importObject).then((result) => {
        go.run(result.instance)
    }).catch(err => {
        postMessage({ type: 'ErrorLoadingWorker', error: err })
    })
} catch(err) {
    postMessage({ type: 'ErrorLoadingWorker', error: err })
}
