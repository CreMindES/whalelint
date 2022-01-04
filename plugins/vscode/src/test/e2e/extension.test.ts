import * as assert from "assert";
import * as path from "path";
import * as vscode from "vscode";

const testFolderLocation = "/../../../src/test/e2e/examples/";

suite("Dockerfile diagnostics", () => {
  vscode.window.showInformationMessage('Start end to end tests.');

  test("should have at least 1 diagnostic result", async () => {
    const diagnosticArray = await getDecorationsFromExample("Dockerfile");
    assert.notEqual(diagnosticArray.filter((d) => d.source === 'WhaleLint').length, 0);

    vscode.commands.executeCommand("workbench.action.closeActiveEditor");
  });
});

async function getDecorationsFromExample(exampleName: string): Promise<vscode.Diagnostic[]> {
  const uri = vscode.Uri.file(path.join(__dirname + testFolderLocation + exampleName));
  const document = await vscode.workspace.openTextDocument(uri);
  const editor = await vscode.window.showTextDocument(document);
  
  // TODO: remove hard coded wait
  await sleep(15000);
  let diagnosticArray = vscode.languages.getDiagnostics(uri);

  return diagnosticArray;
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}