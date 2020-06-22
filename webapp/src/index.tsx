import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';


declare global {
  interface Window {
    hideBytes: (
      input: Uint8Array,
      dataToHide: Uint8Array,
      bitLoss: 2 | 4 | 6 | 8 | undefined,
      callback: (err: Error | null, data: Uint8Array) => void
    ) => void
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
