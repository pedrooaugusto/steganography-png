import React from 'react'
import { Actions, State } from '../../duck'
import { Hex, Text, isInvalidState, PNG, PPNG } from './output-types'
import './style.scss'

type Props = {
    state: State,
    actions: Actions
}

export default function Output(props: Props) {
    const { state: { mode, output } } = props

    const download = () => {
        const a = document.createElement('a')
        const blob = new Blob([output.result as Uint8Array], { type: 'application/octet-stream' })
        a.href = URL.createObjectURL(blob)
        a.download = 'download-' + output.dataType
        a.click()
    }

    return (
        <div className="output-section">
            <div className="main-title">
                <div className="title">Output</div>
                <div className="subtitle">
                    {mode === 'HIDE' && (
                        <span>
                            This is the resultant image with the secret hidden deep down
                            in the pixels of each <i>scanline.</i> Higher values for <i>bit loss </i>
                            produces images with a high volume of noise.<br /><br />
                            <i className="fa fa-info-circle"></i>&nbsp;You can now download and use this image as
                            the input in the <i>Find secret</i> mode to retrieve your secret.
                        </span>
                    )}
                    {mode === 'FIND' && (
                        <span>
                            This is what we found after looking for a hidden secret inside this image.
                        </span>
                    )}
                </div>
            </div>
            <div className="result-file">
                <div className="output">
                    <Empty {...props.state} />
                    <Loading {...props.state} />
                    <Hex.OutputView {...props.state} />
                    <Text.OutputView {...props.state} />
                    <PNG.OutputView {...props.state} />
                    <PPNG.OutputView {...props.state} />
                </div>
                <div className="info">
                    <b>
                    {output.result && mode === 'HIDE' && (<span>New Image Length: {output?.result?.length} bytes;</span>)}
                    {output.result && mode === 'FIND' && (
                        <span>
                            Hidden Secret Length: {output?.result?.length} bytes;
                            Hidden Secret Type: {output?.dataType};
                        </span>
                    )}
                    </b>
                </div>
                <div className="view-options">
                    <PNG.Button {...props.state} setOutputView={props.actions.setOutputView} />
                    <PPNG.Button {...props.state} setOutputView={props.actions.setOutputView} />
                    <Text.Button {...props.state} setOutputView={props.actions.setOutputView} />
                    <Hex.Button {...props.state} setOutputView={props.actions.setOutputView} />
                    { !isInvalidState(props.state) && (
                        <button className={`btn`} onClick={download}>
                            Download Output
                        </button>
                    )}
                </div>
            </div>
        </div>
    )
}

const Empty = (props: State) => {
    if (isInvalidState(props) !== 'EMPTY') return null

    return (
        <div className="output-type empty">
            <p>Please, fill in the Configuration form.</p>
        </div>
    )

}

const Loading = (props: State) => {
    if (isInvalidState(props) !== 'LOADING') return null

    return (
        <div className="output-type loading">
            <h4>Loading please wait... <i className="fa-3x fa fa-spinner fa-spin"></i></h4><br />
            <h5>Why does it take so long ?</h5>
            <h5>I would say it's <a href="https://github.com/pedrooaugusto/steganography-png/blob/master/scanlines/utils.go#L115" target="_blank">Golang zlib compression level</a> setted to <i>BestCompression</i> instead of <i>BestSpeed.</i></h5>
        </div>
    )
}