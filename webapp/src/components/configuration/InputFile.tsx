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
    const [isImageLoading, setImageLoading] = React.useState(false);
    const inputFile = React.useRef<HTMLInputElement | null>(null);

    const onLoadFromFileSystem = async function () {
        setLoading(true);
        setImageLoading(true);

        if (inputFile.current?.files?.length) {
            const file = inputFile.current.files[0];

            if (!file.type.match(/png/gi)) {
                setErr(`Input file must be a png image!\n\tType "${file.type}" is not "image/png"`);
                setLoading(false);
                setImageLoading(false);
                setUrl(file.name);
                setStagedUrl(file.name);

                return;
            }

            props.setInputImage(new Uint8Array(await file.arrayBuffer()));
            setErr(null);

            console.log(file);
            const u = URL.createObjectURL(file);

            setUrl(u);
            setStagedUrl(u);
        }

        setLoading(false);
    }

    const onLoadFromUrl = (evt: React.MouseEvent) => {
        evt.preventDefault()

        setLoading(true);
        setImageLoading(true);

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
                setImageLoading(false);
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
                <form>
                    <input
                        list="images"
                        name="url" id="url"
                        autoComplete="off"
                        placeholder="Insert png url here"
                        value={stagedUrl}
                        onChange={evt => {
                            if (evt.target.value === '') {
                                setUrl('')
                                props.setInputImage(null)
                            }
                            setStagedUrl(evt.target.value)
                        }}
                    />
                    <datalist id="images">
                        <option value="https://raw.githubusercontent.com/pedrooaugusto/steganography-png/master/imagepack/suspicious-pitou.png"/>
                        <option value="https://raw.githubusercontent.com/pedrooaugusto/steganography-png/master/imagepack/funny_horse.png"/>
                        <option value="https://raw.githubusercontent.com/pedrooaugusto/steganography-png/master/imagepack/suspicious-bisky.png"/>
                    </datalist>
                    <button onClick={onLoadFromUrl} disabled={isEmpty}>Load</button>
                </form>
            </div>
            <div className={`preview-img ${(props.empty) ? 'empty' : ''}`}>
                <figure>
                    {/* Ugly! */}
                    {(() => {
                        if (isLoading) return <div className="loading">Loading... <i className="fa-3x fa fa-spinner fa-spin"></i></div>
                        if (isEmpty || url == null || url === '') return <div className="empty"><b>EMPTY PREVIEW -- NO IMAGE!</b></div>
                        if (err) return <div className="err"><span>{err}</span></div>

                        return <img src={url} alt="Input preview" onLoad={() => setImageLoading(false)}/>
                    })()}
                </figure>
            </div>
            <div className="load-file">
                <label htmlFor="file-upload-input-file" className="btn">
                    Or Load From File System
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