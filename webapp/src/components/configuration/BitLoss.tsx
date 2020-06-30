import React from 'react';

export default function BitLoss(props: { setBitLoss: Function }) {
    return (
        <div className="config bit-loss">
            <div className="title">Bit Loss</div>
            <div className="subtitle">
                How many bits of the original file we should use to enconde 1 byte of the secret message ?
            </div>
            <div className="bit-loss-input">
                <select
                    name="bitloss"
                    defaultValue="8"
                    onChange={e => {
                        props.setBitLoss(parseInt(e.target.value))
                    }
                }>

                    <option value="1">1 bit</option>
                    <option value="2">2 bits</option>
                    <option value="4">4 bits</option>
                    <option value="6">6 bits</option>
                    <option value="8">8 bits</option>
                </select>
            </div>
        </div>
    )
}