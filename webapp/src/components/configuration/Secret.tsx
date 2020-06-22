import React from 'react';

type Props = {
    setSecret: (data: Uint8Array | null | string) => void,
    secret: string | Uint8Array | null | undefined
}

export default function Secret({ secret, setSecret }: Props) {
    const secretFile = React.useRef<HTMLInputElement | null>(null);
    const isBuffer = typeof secret !== 'string' && secret?.byteLength != null
    const text = isBuffer
        ? (secret as Uint8Array).slice(0, 100).join(' ') + '...'
        : !!secret ? secret : ''

    return (
        <div className="config secret">
            <div className="title">Data to be hidden</div>
            <div className="subtitle">
                You can hide a secret message in plain text or hide another file itself.
            </div>
            <div className="plain-text">
                <textarea
                    placeholder="Type here the secret message to hide inside the input image"
                    value={text.toString()}
                    readOnly={isBuffer}
                    disabled={isBuffer}
                    title={isBuffer ? 'You cannot edit content loaded directly from a file!' : ''}
                    onChange={(evt) => setSecret(evt.target.value)}
                ></textarea>
                <div className="footer">
                    <div className="clear-all">
                        <button onClick={() => setSecret(null)}>Clear</button>
                    </div>
                    <div className="info">
                        Data Length: {isBuffer ? (secret as Uint8Array).length : text.length}
                    </div>
                </div>
            </div>
            <div className="load-file">
                <label htmlFor="file-upload-secret" className="btn">
                    Or Load From File
                </label>
                <input
                    id="file-upload-secret"
                    type="file"
                    accept="*"
                    ref={secretFile}
                    onChange={async function(this: HTMLInputElement) {
                        if (secretFile.current?.files?.length) {
                            setSecret(new Uint8Array(await secretFile.current.files[0].arrayBuffer()));
                        } else {
                            setSecret(null);
                        }
                    }}
                />
            </div>
        </div>
    )
}