import * as vscode from 'vscode';
import * as cp from 'child_process';
import getPort = require('get-port');
import net = require('net');
import path = require('path');
import {
	CloseAction,
	ErrorAction,
	LanguageClient,
	LanguageClientOptions,
	Message,
	StreamInfo,
} from 'vscode-languageclient/node';
import {hideWhalelintStatusBar, initWhalelintStatusBar, showWhalelintStatusInStatusBar, whalelintEnvStatusbarItem} from './status';


export var whalelintClient: LanguageClient;
export let diagnosticCollection: vscode.DiagnosticCollection;
export let whalelintLsp: cp.ChildProcess;
export let whalelintLspPort: number;
export let whalelintLspSocket: net.Socket;

export async function activate(context: vscode.ExtensionContext) {

	let publisherName = "TamasGBarna.whalelint";
	let extension = vscode.extensions.getExtension(publisherName);
	if(extension === undefined) {
		vscode.window.showErrorMessage(`WhaleLint extension path not found!`);
		return;
	}
	let binaryPath = path.join(extension.extensionPath, "whalelint");

	// get an open port, preferably 18888, if not, from the range of 18000 - 65535.
	whalelintLspPort = await getPort({port: [18888, 18000, 65535]});

	// launch server
	whalelintLsp = cp.exec(binaryPath + " lsp --port " + whalelintLspPort, function (error, stdout, stderr) {
		if (error) {
			vscode.window.showErrorMessage(error.message, "For additional details, please see the console.");
			vscode.window.showErrorMessage(stderr);
			console.log(error.stack);
			console.log('Error code: ' + error.code);
			console.log('Signal received: ' + error.signal);
			console.log('Child Process STDOUT: ' + stdout);
			console.log('Child Process STDERR: ' + stderr);
		}
	});

	await sleep(2000); // Wait for server activation	
	
	whalelintLspSocket = net.connect(whalelintLspPort); //.on("error", onError);

	let serverOptions = () => {
		// Connect to language server via socket
		// let socket = net.connect(whalelintLspPort);
		let result: StreamInfo = {
			writer: whalelintLspSocket,
			reader: whalelintLspSocket
		};
		
		return Promise.resolve(result);
	};

	const traceOutputChannel = vscode.window.createOutputChannel(
		"Whalelint Language Server Trace"
	);

	// Client side configurations
	let clientOptions: LanguageClientOptions = {
		// js is used to trigger things	
		documentSelector: [{ scheme: 'file', language: 'dockerfile' }],
		errorHandler: {
			error: (error: Error, message: Message, count: number): ErrorAction => {
				vscode.window.showErrorMessage(
					`Error communicating with the language server: ${error}: ${message}.`
				);
				return ErrorAction.Shutdown;
			},
			closed: (): CloseAction => {
				// Allow 5 crashes before shutdown.
				vscode.window.showErrorMessage(
					`Error the language server may have crashed.`
				);
				return CloseAction.DoNotRestart;
			},
		},
		traceOutputChannel,
		diagnosticCollectionName: "WhaleLint",
		// middleware: {
		// 	handleDiagnostics: (uri, diagnostics, next) => {
		// 		// assert.equal(uri, "uri:/test.ts");
		// 		// console.log("Diagnostics length:", diagnostics.length);

		// 		let d = -1;
		// 		diagnostics.forEach((diag, index) => {
		// 			if (!diag.range.isSingleLine) {
		// 				d = diagnostics.indexOf(diag);
		// 				let start = diag.range.start.line;
		// 				let end   = diag.range.end.line;
		// 				for (var _i = start; _i <= end; _i++) {
		// 					let line = vscode.window.activeTextEditor?.document.lineAt(_i);
		// 					let s = new vscode.Position(_i, 0);
		// 					let e = new vscode.Position(_i, 50);
		// 					if (line) {
		// 						s = new vscode.Position(_i, line.firstNonWhitespaceCharacterIndex);
		// 						// line.text.replace("\S+\\", "");
		// 						e = new vscode.Position(_i, line.text.replace("\\", "").trimRight().length);
		// 					}
		// 					let range = new vscode.Range(s, e);
							
		// 					let derivedDiag = new vscode.Diagnostic(range, diag.message, diag.severity);
		// 					derivedDiag.source = "WhaleLint";

		// 					diagnostics.push(derivedDiag);
		// 				}
		// 			}
		// 		});
	
		// 		diagnostics.splice(d, 1);
				
		// 		next(uri, diagnostics);
		// 	}
		// }
	};

	whalelintClient = new LanguageClient(
		'WhalelintLanguageServer',
		'Whalelint Language Server',
		serverOptions,
		clientOptions
	);

	whalelintClient.onReady().then(() => {
		// vscode.window.showInformationMessage(
		// 	`WhaleLint language client is finally ready!`
		// );
		initWhalelintStatusBar(whalelintLspPort);
	});

	// Start the client side, and at the same time also start the language server
	whalelintClient.start();

	// vscode.window.showInformationMessage(
	// 	`WhaleLint language server started: ${whalelintClient.initializeResult}`
	// );

	diagnosticCollection = vscode.languages.createDiagnosticCollection('WhaleLint');
	context.subscriptions.push(diagnosticCollection);

	// context.subscriptions.push(vscode.commands.registerCommand('whalelint.status', showWhalelintStatusInStatusBar(port)));
}

export function deactivate(): Thenable<void> | undefined {
	hideWhalelintStatusBar();

	if (whalelintClient) {
		return whalelintClient.stop().then(() => {
			if (whalelintLsp !== undefined) {
				whalelintLsp.kill("SIGKILL");
			}
		});
	}
}

async function sleep(ms: number) {
	return new Promise(resolve => setTimeout(resolve, ms));
}
