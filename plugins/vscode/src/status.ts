
import * as vscode from 'vscode';

// statusbar item for switching the Go environment
export let whalelintEnvStatusbarItem: vscode.StatusBarItem;
export const languageServerIcon = '$(zap)';

export async function initWhalelintStatusBar(port: number) {
	if (!whalelintEnvStatusbarItem) {
		whalelintEnvStatusbarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 0);
	}

	// get extension version
	let extension = vscode.extensions.getExtension("TamasGBarna.whalelint");
	let extensionVersion = "unknown";
	if (extension) {
		extensionVersion = extension.packageJSON.version;
	}

	// set Go version and command
	const name = "whlsp";

	whalelintEnvStatusbarItem.text = languageServerIcon + name;
	whalelintEnvStatusbarItem.tooltip = "WhaleLint v" + extensionVersion + " | Serving on port " + port;

	whalelintEnvStatusbarItem.show();
}

export async function hideWhalelintStatusBar() {
	if (whalelintEnvStatusbarItem) {
		whalelintEnvStatusbarItem.hide();
		whalelintEnvStatusbarItem.dispose();
	}	
}

export async function showWhalelintStatusInStatusBar(port: number) {
	
}

