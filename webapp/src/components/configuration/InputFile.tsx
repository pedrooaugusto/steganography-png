import React from 'react';

type InputFileProps = {
    setInputImage: (input: Uint8Array | null) => void,
    empty: boolean
}

export default function InputFile(props: InputFileProps) {
    const [url, setUrl] = React.useState(initialUrl);
    const [stagedUrl, setStagedUrl] = React.useState(initialUrl);

    const [err, setErr] = React.useState<string | null>(null);
    const [isLoading, setLoading] = React.useState(false);
    const inputFile = React.useRef<HTMLInputElement | null>(null);

    const onLoadFromFileSystem = async function () {
        setLoading(true);

        if (inputFile.current?.files?.length) {
            console.log(inputFile.current.files[0])
            props.setInputImage(new Uint8Array(await inputFile.current.files[0].arrayBuffer()));
            setErr(null);

            const u = URL.createObjectURL(inputFile.current.files[0]);

            setUrl(u);
            setStagedUrl(u);
        }

        setLoading(false);
    }

    const onLoadFromUrl = () => {
        setLoading(true);

        fetch(stagedUrl, { method: 'GET' })
            .then(res => {
                if (res.status !== 200) throw new FailedToLoadImage(`Request response was not ok:\n\t'${res.statusText}`);

                const type = res.headers.get('content-type') || '';

                if (!type.toLocaleLowerCase().includes('png'))
                    throw new FileTypeNotSupported(`Input file must be a png image!\n\tType "${type}" is not "image/png"`);

                return res;
            })
            .then((res: any) => res.arrayBuffer())
            .then((buff: ArrayBuffer) => {
                setErr(null);
                setUrl(stagedUrl);
                props.setInputImage(new Uint8Array(buff));
            })
            .catch((err: Error) => {
                setErr(err.toString());
                setUrl(stagedUrl);
                props.setInputImage(null);
            }).finally(() => {
                setLoading(false);
            });
    }

    const isEmpty = (stagedUrl == null || stagedUrl === '')

    return (
        <div className="config input-file">
            <div className="title">Input image</div>
            <div className="subtitle">
                The input file must be a PNG image, you can either load from the file system
                or from an external URL. This is the image in which the secret is hidden or
                the secret will be hidden (depending on the mode).
            </div>
            <div className="load-url-input">
                <input
                    type="text"
                    value={stagedUrl}
                    onChange={evt => {
                        if (evt.target.value === '') {
                            setUrl('')
                            props.setInputImage(null)
                        }
                        setStagedUrl(evt.target.value)
                    }}
                    name="url"
                    placeholder="Insert png url here"
                />
                <button onClick={onLoadFromUrl} disabled={isEmpty}>Load</button>
            </div>
            <div className={`preview-img ${(props.empty) ? 'empty' : ''}`}>
                <figure>
                    {/* Ugly! */}
                    {(() => {
                        if (isEmpty || url == null || url === '') return <div className="empty"><b>EMPTY PREVIEW -- NO IMAGE!</b></div>
                        if (err) return <div className="err"><span>{err}</span></div>
                        if (isLoading) return <div className="loading">Loading...</div>

                        return <img src={url} alt="Input preview" />
                    })()}
                </figure>
            </div>
            <div className="load-file">
                <label htmlFor="file-upload-input-file" className="btn">
                    Or Load From File
                </label>
                <input
                    id="file-upload-input-file"
                    type="file"
                    accept=".png"
                    onChange={onLoadFromFileSystem}
                    ref={inputFile}
                />
            </div>
        </div>
    )
}

const initialUrl = '' ?? 'https://vignette.wikia.nocookie.net/anicrossbr/images/2/20/109_-_Neferpitou_portrait.png/revision/latest/scale-to-width-down/340?cb=20160308215759&path-prefix=pt-br';

type Error = {
    message: string
};

class FailedToLoadImage extends Error {
    constructor(message: string){
        super('Failed to load image:\n\t' + message);
    }
}

class FileTypeNotSupported extends Error {
    constructor(message: string){
        super('File type not supported:\n\t' + message);
    }
}