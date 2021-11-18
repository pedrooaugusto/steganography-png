class PortableNetWorkGraphicsWorker {
    public worker: Worker
    private listeners: [number, (val: any) => void, (err: Error) => void][] = []

    constructor() {
        this.worker = new Worker('go-worker.js')
        this.worker.onmessage = (event) => {
            if (event.data.type === 'ErrorLoadingWorker') {
                console.log('Killing Worker: ' + event.data.error)

                return this.worker.terminate()
            }

            if (event.data.type === 'OperationResponse') {
                const listener = this.listeners.find(a => a[0] === event.data.id)

                if (listener == null) return

                this.listeners = this.listeners.filter(a => a[0] !== event.data.id)

                if (event.data.error) {
                    listener[2](event.data.error)
                } else {
                    listener[1](event.data.payload)
                }

                return
            }
        }
    }

    public hideData(
        inputImage: Uint8Array,
        data: Uint8Array,
        dataType: string,
        bitLoss: 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | undefined
    ): Promise<{ data: Uint8Array, dataType?: string }> {
        return new Promise((res, rej) => {
            const id = +new Date()

            this.listeners.push([id, res, rej])
            this.worker.postMessage({ type: 'Operation', operationName: 'hideData', inputImage, data, dataType, bitLoss, id })
        })
    }

    public revealData(inputImage: Uint8Array): Promise<{ data: Uint8Array, dataType?: string }> {
        return new Promise((res, rej) => {
            const id = +new Date()

            this.listeners.push([id, res, rej])
            this.worker.postMessage({ type: 'Operation', operationName: 'revealData', inputImage, id })
        })
    }

    public toString(inputImage: Uint8Array): Promise<{ data: string }> {
        return new Promise((res, rej) => {
            const id = +new Date()

            this.listeners.push([id, res, rej])
            this.worker.postMessage({ type: 'Operation', operationName: 'toString', inputImage, id })
        })
    }
}

export default new PortableNetWorkGraphicsWorker()
