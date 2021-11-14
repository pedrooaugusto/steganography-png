import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';


declare global {
  interface PNG {
    hideData: (
      input: Uint8Array,
      dataToHide: Uint8Array,
      type: string,
      bitLoss: 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | undefined,
      callback: (err: Error | null, data: Uint8Array, dataType?: string) => void
    ) => void,
    revealData: (
      input: Uint8Array,
      callback: (err: Error | null, data: Uint8Array, dataType?: string) => void
    ) => void,
    toString: (
      input: Uint8Array,
      callback: (err: Error, str: string) => void
    ) => void
  }

  interface Window {
    PNG: PNG
  }
}

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
