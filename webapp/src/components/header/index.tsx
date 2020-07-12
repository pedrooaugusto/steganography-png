import React from 'react';
import './style.scss';

export default function Header(props: any) {

    return (
        <header className="app-header">
            <div className="icons">
                <div className="icon github">
                    <a href="https://github.com/pedrooaugusto/steganography-png" target="_blank">
                        <i className="fa fa-github" /> Github
                    </a>
                </div>
                <div className="icon medium">
                    <a href="#dois">
                        <i className="fa fa-medium" /> Article
                    </a>
                </div>
            </div>
            <div className="main-title">
                <h1>Portable Network Graphics & Steganography</h1>
                <div className="subtitle">Hiding and retrieving secret files inside PNG files</div>
            </div>
        </header>
    )
}