import React from 'react';
import { Actions, State } from '../../duck';
import './style.scss';

type Props = {
    state: State,
    actions: Actions
};

export default function Output(props: Props) {
    const { state: { mode, output } } = props
    const thereIsOutput = (!output.err && !output.loading && output.result)
    const isLoading = !!output.loading
    const err = output.err
    const [showAs, setShowAs] = React.useState({ png: null, pPng: null, text: null, hex: null })

    const imageUrl = React.useMemo(() => {
        if (mode === 'HIDE') {
            if (!output.err && !output.loading && output.result) {
                const blob = new Blob([output.result], { type: "image/png" })

                return URL.createObjectURL(blob)
            }
        }

        return 'fake_url'

    }, [mode, output.result])

    return (
        <div className="output-section">
            <div className="main-title">
                <div className="title">Output</div>
                <div className="subtitle">
                    {mode === 'HIDE' && (
                        <span>
                            This is the resultant image with the data you have chosen hidden deep down
                            in the pixels of each <i>scanline.</i> If you picked a high value for 
                            <i> bit loss</i>, is possible to spot some differencies betwenn the
                            resultant image and the original one.
                        </span>
                    )}
                    {mode === 'FIND' && (
                        <span>
                            This is what we found after searching deep down in the bits of the
                            input image: <b>Please fill the configuration form first</b>
                        </span>
                    )}
                </div>
            </div>
            <div className="result-file">
                <div className="output">
                    {!isLoading && !err && !thereIsOutput && (
                        <h4>Please fill the configuration form first!</h4>
                    )}
                    {!isLoading && !err && thereIsOutput && showAs.png && (
                        <figure>
                            <img src={imageUrl} alt="Output file" />
                        </figure>
                    )}
                </div>
                <div className="info">
                    Hidden File Length: 12; Total Time: 134ms; Hidden file was detected as being another PNG
                </div>
                <div className="view-options">
                    <button
                        className={`btn ${showAs.png == true ? 'selected' : ''} ${showAs.png === null ? 'disabled' : ''}`}
                        disabled={showAs.pPng === null}
                    >
                        Show as PNG
                    </button>
                    <button
                        className={`btn ${showAs.pPng == true ? 'selected' : ''} ${showAs.pPng === null ? 'disabled' : ''}`}
                        disabled={showAs.pPng === null}
                    >
                        Show as Parsed PNG
                    </button>
                    <button
                        className={`btn ${showAs.text == true ? 'selected' : ''} ${showAs.text === null ? 'disabled' : ''}`}
                        disabled={showAs.text === null}
                    >
                        Show as Plain Text
                    </button>
                    <button
                        className={`btn ${showAs.hex == true ? 'selected' : ''} ${showAs.hex === null ? 'disabled' : ''}`}
                        disabled={showAs.hex === null}
                    >
                        Show as Hex
                    </button>
                </div>
            </div>
        </div>
    )
}