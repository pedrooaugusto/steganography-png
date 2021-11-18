import PNGWorker from './go-worker'

export type State = {
    imageBuf: Uint8Array | null,
    mode: 'HIDE' | 'FIND',
    dataToHide?: Uint8Array | string | null,
    bitLoss?: 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8,
    output: {
        viewType: 'PNG' | 'PPNG' | 'PLAIN' | 'HEX',
        result: Uint8Array | null,
        dataType?: string,
        err: Error | null,
        loading: boolean
    }
}

type Action = {
    type: 'SET_IMAGE_BUFF' | 'SET_MODE' | 'SET_DATA_TO_HIDE' | 'SET_BITLOSS' |
        'PROCCESS' | 'CHANGE_OUTPUT_VIEW_MODE' | 'SET_OUTPUT_VIEW_TYPE',
    data: any
}

export const initialState: State = {
    imageBuf: null,
    mode: 'HIDE',
    dataToHide: '',
    bitLoss: 8,
    output: {
        viewType: 'PNG',
        result: null,
        err: null,
        loading: false
    }
}

export default function reducer(state = initialState, action: Action): State {
    switch(action.type) {
        case 'SET_IMAGE_BUFF':
            return { ...state, imageBuf: action.data, output: { ...initialState.output } }
        case 'SET_MODE':
            return { ...state, mode: action.data, output: { ...initialState.output } }
        case 'SET_DATA_TO_HIDE':
            return { ...state, dataToHide: action.data }
        case 'SET_BITLOSS':
            return { ...state, bitLoss: action.data }
        case 'SET_OUTPUT_VIEW_TYPE':
            return { ...state, output: { ...state.output, viewType: action.data } }
        case 'PROCCESS':
            return { ...state, output: { ...state.output, ...action.data } }

        default:
            return state
    }
}

export interface Actions {
    setInputImage: (buf: Uint8Array | null) => void,
    setMode: (mode: 'HIDE' | 'FIND') => void,
    setDataToHide: (buf: Uint8Array | null | string) => void,
    setBitLoss: (bitLoss: 2 | 4 | 6 | 8) => void,
    setOutputView: (type: 'PNG' | 'PPNG' | 'PLAIN' | 'HEX') => void,
    startProcess: () => void
}

export function makeActions([state, dispatch]: [State, (action: Action) => void]): [State, Actions] {
    return [state, {
        setInputImage(buf: Uint8Array | null) {
            dispatch({ type: 'SET_IMAGE_BUFF', data: buf })
        },
        setMode(mode: 'HIDE' | 'FIND') {
            dispatch({ type: 'SET_MODE', data: mode })
        },
        setDataToHide(buf: Uint8Array | null | string) {
            dispatch({ type: 'SET_DATA_TO_HIDE', data: buf })
        },
        setBitLoss(bitLoss: 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8) {
            dispatch({ type: 'SET_BITLOSS', data: bitLoss })
        },
        setOutputView(veiwType) {
            dispatch({ type: 'SET_OUTPUT_VIEW_TYPE', data: veiwType })
        },
        startProcess(){
            dispatch({ type: 'PROCCESS', data: {
                result: null,
                err: null,
                loading: true
            }})

            const myHandle = (err: null | Error, data: Uint8Array, dataType: string = "") => handle(err, data, dataType, state.mode, dispatch)

            if (state.mode === 'HIDE') {
                // @ts-ignore
                let dataType = state.dataToHide?.type

                if (typeof state.dataToHide === 'string') {
                    dataType = state.dataToHide.startsWith("#!HTML") ? 'text/html.html' : 'text/plain.txt'
                }

                PNGWorker.hideData(state.imageBuf!, toUint8Array(state.dataToHide), dataType, state.bitLoss).then(res => {
                    myHandle(null, res.data, res.dataType)
                }).catch(err => {
                    myHandle(err, new Uint8Array(), "")
                })
            } else {
                PNGWorker.revealData(state.imageBuf!).then(res => {
                    myHandle(null, res.data, res.dataType)
                }).catch(err => {
                    myHandle(err, new Uint8Array(), "")
                })
                // window.PNG.revealData(state.imageBuf!, myHandle)
            }

            if (matchMedia('screen and (max-width: 860px)').matches) setTimeout(() => window.scrollTo(0, document.body.scrollHeight), 100)
            else window.scrollTo(0, 0)
        }
    }]
}

const handle = (err: null | Error, data: Uint8Array, dataType: string = "", mode: string, dispatch: (action: Action) => void): void => {
    if (err) {
        alert(err)
        return dispatch({
            type: 'PROCCESS',
            data: {
                result: null,
                err: err,
                loading: false
            }
        })
    }

    // Force display as PNG if the mode is hide.
    const isImage = dataType.search(/png|gif|jpg|jpeg/gi) >= 0 || mode === 'HIDE'
    const isText = dataType.search(/text/gi) >= 0

    dispatch({
        type: 'PROCCESS',
        data: {
            result: data,
            err: null,
            loading: false,
            viewType: isImage ? 'PNG' : isText ? 'PLAIN' : 'HEX',
            dataType: mode === 'HIDE' ? 'image-png.png' : dataType
        }
    })
}

function toUint8Array(raw: string | Uint8Array | undefined | null): Uint8Array {
    if ((raw as Uint8Array).buffer) return (raw as Uint8Array)

    return new TextEncoder().encode(raw as string)
}