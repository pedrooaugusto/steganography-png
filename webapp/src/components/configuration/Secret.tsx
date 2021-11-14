import React from 'react';

type Props = {
    setSecret: (data: Uint8Array | null | string) => void,
    secret: string | Uint8Array | null | undefined,
    empty: boolean
}

export default function Secret({ secret, setSecret, empty }: Props) {
    const secretFile = React.useRef<HTMLInputElement | null>(null);
    const isBuffer = typeof secret !== 'string' && secret?.byteLength != null
    const text = isBuffer
        ? (secret as Uint8Array).slice(0, 100).join(' ') + '...'
        : !!secret ? secret : ''

    return (
        <div className="config secret">
            <div className="title">Secret to be hidden</div>
            <div className="subtitle">
                The secret can be a plain text message or a file loaded from the file system.
            </div>
            <div className="plain-text">
                <textarea
                    className={(empty && (secret == null || secret === '')) ? 'empty' : ''}
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
                            const v = new Uint8Array(await secretFile.current.files[0].arrayBuffer())

                            const ext = secretFile.current.files[0].name.includes('.') && secretFile.current.files[0].name.split('.').pop()

                            // @ts-ignore
                            v.type = secretFile.current.files[0].type + (ext ? '.' + ext : '')

                            setSecret(v);
                        } else {
                            setSecret(null);
                        }
                    }}
                />
            </div>
        </div>
    )
}