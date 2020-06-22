import React from 'react';
import bisk from '..//configuration/test3.png';
import { Actions, State } from '../../duck';
import './style.scss';

type Props = {
    state: State,
    actions: Actions
};

export default function Output(props: Props) {
    const { state: { mode, output } } = props
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
                    Here you can specify which parameters to apply during the proccess such as
                    input image, mode, bit loss, and data to be hidden.
                </div>
            </div>
            <div className="result-file">
                <div className="output">
                    <figure>
                        <img src={imageUrl} alt="Output file" />
                    </figure>
                </div>
                <div className="info">
                    Hidden File Length: 12; Total Time: 134ms; Hidden file was detected as being another PNG
                </div>
                <div className="view-options">
                    <button className="btn selected">Show as PNG</button>
                    <button className="btn">Show as Parsed PNG</button>
                    <button className="btn disabled">Show as Plain Text</button>
                    <button className="btn">Show as Hex</button>
                </div>
            </div>
        </div>
    )
}