import React from 'react'
import { State } from '../../duck'
import PNGWorker from '../../duck/go-worker'

export function isInvalidState (state: State): 'ERR' | 'LOADING' | 'EMPTY' | null {
    const {mode, output} = state

    if (!!output.loading) return 'LOADING'

    const thereIsOutput = !output.err && output.result

    if (!thereIsOutput) return 'EMPTY'

    const err = output.err

    if (err) return 'ERR'

    return null
}

interface OutputMode {
    Button: (props: State & { setOutputView: (type: 'PNG' | 'HEX' | 'PPNG' | 'PLAIN') => void }) => JSX.Element | null
    OutputView: (props: State) => JSX.Element | null
}

export const Hex: OutputMode = {
    Button: (props) => {
        if (isInvalidState(props)) return null
        const selected = props.output.viewType === 'HEX'
        const available = props.mode === 'FIND'

        return (
            <button
                className={`btn ${selected ? 'selected' : ''} ${available ? '' : 'disabled'}`}
                disabled={!available}
                onClick={() => props.setOutputView('HEX')}
            >
                Show as Hex
            </button>
        )
    },
    OutputView: (props) => {
        if (isInvalidState(props) || props.output.viewType !== 'HEX') return null

        const text = Array.from((props.output.result as Uint8Array)).map(item => item.toString(16)).join(" ")

        return (
            <div className="output-type hex" style={{ position: 'absolute' }}>
                <pre>{text}</pre>
            </div>
        )
    }
}

export const Text: OutputMode = {
    Button: (props) => {
        if (isInvalidState(props)) return null
        const selected = props.output.viewType === 'PLAIN'
        const available = props.mode === 'FIND'

        return (
            <button
                className={`btn ${selected ? 'selected' : ''} ${available ? '' : 'disabled'}`}
                disabled={!available}
                onClick={() => props.setOutputView('PLAIN')}
            >
                Show as Plain Text
            </button>
        )
    },
    OutputView: (props) => {
        if (isInvalidState(props) || props.output.viewType !== 'PLAIN') return null

        const text = new TextDecoder().decode((props.output.result as Uint8Array))

        return (
            <div className="output-type plain-text">
                {text.startsWith("#!HTML") ? <div dangerouslySetInnerHTML={{ __html: text }} /> : <pre>{text}</pre>}
            </div>
        )
    }
}

export const PNG: OutputMode & { isPng: Function } = {
    Button: (props) => {
        if (isInvalidState(props)) return null
        const selected = props.output.viewType === 'PNG'
        const available = (props.output.dataType?.search?.(/png|gif|jpg|jpeg/gi) || -1) >= 0

        return (
            <button
                className={`btn ${selected ? 'selected' : ''} ${available ? '' : 'disabled'}`}
                disabled={!available}
                onClick={() => props.setOutputView('PNG')}
            >
                Show as Image
            </button>
        )
    },
    OutputView: (props) => {
        const hide = isInvalidState(props) || props.output.viewType !== 'PNG'
        const { mode, output: { result, dataType } } = props

        const imageUrl = React.useMemo(() => {
            if (hide) return null

            let type = mode === 'HIDE' ? 'image/png' : (dataType || '')

            type = type.split(".")[0]

            const blob = new Blob([result as Uint8Array], { type })

            return URL.createObjectURL(blob)
    
        }, [hide, mode, result, dataType])

        if (hide || imageUrl === null) return null

        return (
            <div className="output-type png">
                <figure>
                    <img src={imageUrl} alt="Output file" />
                </figure>
            </div>
        )
    },
    isPng: (result: Uint8Array) => {
        const a = result.subarray(0, 8)
        const b = [137, 80, 78, 71, 13, 10, 26, 10]

        return a.every((v, i) => v === b[i])
    }
}

// Parsed PNG view
export const PPNG: OutputMode = {
    Button: (props) => {
        if (isInvalidState(props)) return null
        const selected = props.output.viewType === 'PPNG'
        const available = PNG.isPng(props.output.result) && (props.mode === 'FIND' || props.mode === 'HIDE')

        return (
            <button
                className={`btn ${selected ? 'selected' : ''} ${available ? '' : 'disabled'}`}
                disabled={!available}
                onClick={() => props.setOutputView('PPNG')}
            >
                Show as Parsed PNG
            </button>
        )
    },
    OutputView: (props) => {
        const invalid = isInvalidState(props) || props.output.viewType !== 'PPNG'
        const [text, setText] = React.useState('')

        // Please, dont do it here. Move it to somewhere else later.
        React.useEffect(() => {
            if (!invalid) {
                PNGWorker.toString(props.output.result as Uint8Array)
                    .then((payload) => setText(payload.data))
                    .catch((err) => console.error(err))
            }

        }, [invalid])

        if (invalid) return null

        return (
            <div className="output-type png-parsed">
                <pre>{text}</pre>
            </div>
        )
    }
}