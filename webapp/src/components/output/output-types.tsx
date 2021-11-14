import React from 'react'
import { State } from '../../duck'

export function isInvalidState (state: State): 'ERR' | 'LOADING' | 'EMPTY' | null {
    const {mode, output} = state

    const thereIsOutput = (!output.err && !output.loading && output.result)

    if (!thereIsOutput) return 'EMPTY'

    const isLoading = !!output.loading

    if (isLoading) return 'LOADING'

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

        const text = Array.from((props.output.result as Uint8Array)).map(item => item.toString(12)).join(" ")

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
                <pre>{text}</pre>
            </div>
        )
    }
}

export const PNG: OutputMode & { isPng: Function } = {
    Button: (props) => {
        if (isInvalidState(props)) return null
        const selected = props.output.viewType === 'PNG'
        const available = PNG.isPng(props.output.result) && (props.mode === 'FIND' || props.mode === 'HIDE')

        return (
            <button
                className={`btn ${selected ? 'selected' : ''} ${available ? '' : 'disabled'}`}
                disabled={!available}
                onClick={() => props.setOutputView('PNG')}
            >
                Show as PNG
            </button>
        )
    },
    OutputView: (props) => {
        const hide = isInvalidState(props) || props.output.viewType !== 'PNG'

        const imageUrl = React.useMemo(() => {
            if (hide) return null

            const blob = new Blob([props.output.result as Uint8Array], { type: "image/png" })

            return URL.createObjectURL(blob)
    
        }, [hide, props.mode, props.output.result])

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
                window.PNG.toString(props.output.result as Uint8Array, (err, str) => {
                    if (err) return console.error(err)

                    setText(str)
                })
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