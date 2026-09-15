import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import terser from '@rollup/plugin-terser';
import image from '@rollup/plugin-image';
import babel from '@rollup/plugin-babel';
import polyfill from 'rollup-plugin-polyfill';
import rollupPluginCssOnly from 'rollup-plugin-css-only';
import postcss from 'postcss';
import cssvariables from 'postcss-css-variables';

const production = !process.env.ROLLUP_WATCH;

function serve() {
	let server;

	function toExit() {
		if (server) server.kill(0);
	}

	return {
		writeBundle() {
			if (server) return;
			server = require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
				stdio: ['ignore', 'inherit', 'inherit'],
				shell: true
			});

			process.on('SIGTERM', toExit);
			process.on('exit', toExit);
		}
	};
}

export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: 'public/build/bundle.js',
		inlineDynamicImports: true,
	},
	plugins: [
		image(),
		svelte({
			compilerOptions: {
				dev: !production,
			},
			preprocess: {
				style: ({ content }) => {
					const plugins = [
						cssvariables({
							preserve: true,
							variables: {
								"--bg-color": "#2e3440",
								"--bg-color2": "#434c5e",
								"--bg-color3": "#4c566a",
								"--bg-inverse-color": "#eceffa",
								"--bg-inverse-color2": "#e5e9f0",
								"--bg-inverse-color3": "#d8dee9",
								"--bg-input-color": "#242933",

								"--border-color": "#3b4252",

								"--text-color": "#eceff4",
								"--text-color2": "#d8dee9",
								"--text-color3": "#bcc7d9",
								"--text-inverse-color": "#2e3440",

								"--primary-color": "#88c0d0",
								"--accent-color": "#8fbcbb",
								"--accent-color2": "#81a1c1",
								"--accent-color3": "#5e81ac",

								"--red-color": "#bf616a",
								"--orange-color": "#d08770",
								"--yellow-color": "#ebcb8b",
								"--green-color": "#a3be8c",
								"--purple-color": "#b48ead",

								"--padding": "12px",
								"--border": "1px solid #3b4252",
								"--font-size": "10pt",

								"--height": "40px",
								"--background": "var(--bg-input-color)",
								"--borderRadius": "0",
								"--borderFocusColor": "var(--border-color)",
								"--borderHoverColor": "var(--border-color)",
								"--inputColor": "var(--text-color)",
								"--listBackground": "var(--bg-inverse-color)",
								"--listBorderRadius": "0",
								"--itemColor": "var(--text-inverse-color)",
								"--itemHoverBG": "var(--bg-inverse-color2)",
								"--itemIsActiveColor": "var(--text-inverse-color)",
								"--itemIsActiveBG": "var(--bg-inverse-color3)",
								"--placeholderColor": "var(--border-color)",
								"--indicatorColor": "var(--border-color)",
								"--clearSelectRight": "36px",
								"--clearSelectBottom": "0",
								"--clearSelectTop": "3px",
								"--clearSelectColor": "var(--border-color)",
								"--clearSelectFocusColor": "var(--border-color)",
							}
						})
					];
					return postcss(plugins)
						.process(content, {
							from: 'src',
							map: { inline: false }
						})
						.then((result) => ({
							code: result.css.toString(),
							map: result.map.toString()
						}));
				}
			}
		}),

		resolve({
			browser: true,
			dedupe: ['svelte', 'svelte/transition', 'svelte/internal']
		}),
		commonjs(),

		rollupPluginCssOnly({
			output: 'bundle.css'
		}),

		!production && serve(),
		!production && livereload('public'),

		babel({
			babelHelpers: 'runtime',
			extensions: ['.js', '.jsx', '.es6', '.es', '.mjs', '.svelte', '.html'],
			exclude: ['node_modules/@babel/**', 'node_modules/core-js/**'],
			presets: [
				[
					'@babel/preset-env',
					{
						targets: 'defaults',
						modules: false,
						useBuiltIns: 'usage',
						corejs: 3,
					},
				]
			],
			plugins: [
				'@babel/plugin-syntax-dynamic-import',
				[
					'@babel/plugin-transform-runtime',
					{
						useESModules: true
					}
				]
			]
		}),
		polyfill(['@webcomponents/webcomponentsjs']),

		production && terser()
	],
	watch: {
		clearScreen: false
	}
};
