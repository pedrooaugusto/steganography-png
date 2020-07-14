import React from 'react';
import Header from './components/header';
import Configuration from './components/configuration';
import Output from './components/output';
import reducer, { makeActions, initialState } from './duck';
import { match } from 'assert';

function App() {
	const [state, actions] = makeActions(React.useReducer(reducer, initialState));
	const matches = useMatchMedia('screen and (max-width: 860px)');

	const showOutput = !matches || (state.output.result);

	return (
		<div className="App">
			<Header />
			<main>
				<section>
					<Configuration state={state} actions={actions} />
				</section>
				{showOutput && (
					<section>
						<Output state={state} actions={actions} />
					</section>
				)}
			</main>
		</div>
	);
}

// We can move it to duck
const useMatchMedia = (str: string) => {
	const [matches, setMatches]	= React.useState(matchMedia(str).matches);
	const fn = React.useCallback(() => {
		console.log('nn')
		setMatches(matchMedia(str).matches);
	}, [str]);

	React.useEffect(() => {
		window.addEventListener("resize", fn);
		console.log('dd')

		return () => {
			window.removeEventListener("resize", fn);
		}
	}, [fn])

	return matches;
}

export default App;
