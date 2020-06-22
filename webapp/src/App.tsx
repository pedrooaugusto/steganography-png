import React from 'react';
import Header from './components/header';
import Configuration from './components/configuration';
import Output from './components/output';
import reducer, { makeActions, initialState } from './duck';

function App() {
	const [state, actions] = makeActions(React.useReducer(reducer, initialState));

	return (
		<div className="App">
			<Header />
			<main>
				<section>
					<Configuration state={state} actions={actions} />
				</section>
				<section>
					<Output state={state} actions={actions} />
				</section>
			</main>
		</div>
	);
}

export default App;
