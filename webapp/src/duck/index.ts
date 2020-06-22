export type State = {
    imageBuf: Uint8Array | null,
    mode: 'HIDE' | 'FIND',
    dataToHide?: Uint8Array | string | null,
    bitLoss?: 2 | 4 | 6 | 8,
    output: {
        viewType: 'PNG' | 'PPNG' | 'PLAIN' | 'HEX',
        result: Uint8Array | null,
        err: Error | null,
        loading: boolean
    }
}

type Action = {
    type: 'SET_IMAGE_BUFF' | 'SET_MODE' | 'SET_DATA_TO_HIDE' | 'SET_BITLOSS' |
        'PROCCESS' | 'CHANGE_OUTPUT_VIEW_MODE',
    data: any
}

export const initialState: State = {
    imageBuf: null,
    mode: 'HIDE',
    dataToHide: 'Hello, Doctor!',
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
            return { ...state, imageBuf: action.data }
        case 'SET_MODE':
            return { ...state, mode: action.data }
        case 'SET_DATA_TO_HIDE':
            return { ...state, dataToHide: action.data }
        case 'SET_BITLOSS':
            return { ...state, bitLoss: action.data }
        case 'PROCCESS':
            return { ...state, output: { ...state.output, ...action.data } }

        default:
            return state
    }
}

export type Actions = {
    setInputImage: (buf: Uint8Array | null) => void,
    setMode: (mode: 'HIDE' | 'FIND') => void,
    setDataToHide: (buf: Uint8Array | null | string) => void,
    setBitLoss: (bitLoss: 2 | 4 | 6 | 8) => void,
    startProcess: () => void
}

export function makeActions([state, dispatch] : [State, (action: Action) => void]): [State, Actions] {
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
        setBitLoss(bitLoss: 2 | 4 | 6 | 8) {
            dispatch({ type: 'SET_BITLOSS', data: bitLoss })
        },
        startProcess(){
            dispatch({ type: 'PROCCESS', data: {
                result: null,
                err: null,
                loading: true
            }})

            if (state.mode === 'HIDE') {
                const buf: Uint8Array = toUint8Array(state.dataToHide)
                const handle = (err: null | Error, data: Uint8Array): void => {
                    if (err) {
                        console.error(err)
                        return dispatch({
                            type: 'PROCCESS',
                            data: {
                                result: null,
                                err: err,
                                loading: false
                            }
                        })
                    }

                    dispatch({
                        type: 'PROCCESS',
                        data: {
                            result: data,
                            err: null,
                            loading: false
                        }
                    })
                }

                window.hideBytes(state.imageBuf as Uint8Array, buf, state.bitLoss, handle)
            }
        }
    }]
}

function toUint8Array(raw: string | Uint8Array | undefined | null): Uint8Array {
    if ((raw as Uint8Array).buffer) return (raw as Uint8Array)

    return new TextEncoder().encode(raw as string)
}