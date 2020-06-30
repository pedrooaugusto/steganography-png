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
    console.log(props.state)

    return (
        <div className="config-section">
            <div className="main-title">
                <div className="title">Configuration</div>
                <div className="subtitle">
                    Here you can specify which parameters to apply during the proccess such as
                    input image, mode, bit loss, and data to be hidden.
                </div>
            </div>
            <InputFile setInputImage={props.actions.setInputImage} />
            <Mode setMode={props.actions.setMode} />
            {props.state.mode !== 'FIND' && (
                <React.Fragment>
                    <Secret setSecret={props.actions.setDataToHide} secret={props.state.dataToHide} />
                    <BitLoss setBitLoss={props.actions.setBitLoss} />
                </React.Fragment>
            )}
            <div className="submit-section">
                <button className="btn" onClick={props.actions.startProcess}>GO!</button>
            </div>
        </div>
    )
}
