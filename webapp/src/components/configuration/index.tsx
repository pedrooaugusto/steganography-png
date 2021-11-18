import React from 'react';
import InputFile from './InputFile';
import Mode from './Mode';
import Secret from './Secret';
import BitLoss from './BitLoss';
import { Actions, State } from '../../duck';
import './style.scss';

type Props = {
    state: State,
    actions: Actions
};

export default function Configuration(props: Props) {
    const [formError, setFormError] = React.useState(false)

    const startProcess = () => {
        if(props.state.imageBuf == null) return setFormError(true)

        if(props.state.mode === 'HIDE' && (props.state.dataToHide == null || props.state.dataToHide === '' )) {
            return setFormError(true)
        }

        setFormError(false)

        props.actions.startProcess()
    }

    return (
        <div className="config-section">
            <div className="main-title">
                <div className="title">Configuration</div>
                <div className="subtitle">
                    Here you can specify which parameters to apply during the proccess, such as
                    the input image, mode (hide secret or reveal secret), bit loss and the secret to be hidden.
                </div>
            </div>
            <InputFile setInputImage={props.actions.setInputImage} empty={formError && props.state.imageBuf === null} />
            <Mode setMode={props.actions.setMode} />
            {props.state.mode !== 'FIND' && (
                <React.Fragment>
                    <Secret
                        setSecret={props.actions.setDataToHide}
                        secret={props.state.dataToHide}
                        empty={formError}
                    />
                    <BitLoss setBitLoss={props.actions.setBitLoss} />
                </React.Fragment>
            )}
            <div className="submit-section">
                <button
                    className={`btn ${!!props.state?.output?.loading ? 'disabled loading' : ''}`}
                    disabled={!!props.state?.output?.loading}
                    onClick={startProcess}
                >
                    GO!
                </button>
                {formError && (
                    <span className="subtitle">Pleae fill in all the required fields! </span>
                )}
            </div>
        </div>
    )
}
