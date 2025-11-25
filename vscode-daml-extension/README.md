# DAML Syntax Highlighting for VS Code

This extension provides basic syntax highlighting for DAML smart contract files (`.daml`).

## Installation

### Install from VSIX (Recommended)

This method packages the extension into a standard installer file.

1.  **Install the packaging tool** (requires Node.js):
    ```bash
    npm install -g vsce
    ```

2.  **Package the extension**:
    Run this command from inside the `vscode-daml-extension` directory:
    ```bash
    vsce package
    ```
    This will generate a file like `daml-syntax-0.0.1.vsix`.

3.  **Install in Editor**:
    - Open VS Code / Cursor.
    - Go to the **Extensions** view (Cmd+Shift+X).
    - Click the **...** (Views and More Actions) menu at the top of the pane.
    - Select **Install from VSIX...**.
    - Choose the generated `.vsix` file.
