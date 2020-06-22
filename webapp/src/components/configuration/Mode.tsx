import React from 'react';

export default function Mode(props: { setMode: (mode: 'HIDE' | 'FIND') => void }) {
    return (
        <div className="config mode">
            <div className="title">Mode</div>
            <div className="subtitle">
                You can either look for data hidden inside the input image or hide data inside
                the input image.
            </div>
            <div className="opts">
                <label htmlFor="hide">
                    <input
                        type="radio"
                        name="mode"
                        value="hide"
                        id="hide"
                        defaultChecked
                        onClick={() => props.setMode('HIDE')}
                    />
                    Hide data
                </label>
                <label htmlFor="find">
                    <input
                        type="radio"
                        name="mode"
                        value="find"
                        id="find"
                        onClick={() => props.setMode('FIND')}
                    />
                    Find data
                </label>
            </div>
        </div>
    )
}