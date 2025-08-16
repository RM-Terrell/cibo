# Common commands used during development

## Compiling and running

```bash
go run .
```

VS Code Command Palette -> Run Task -> Run Go Tests (data_engine). This is configured in `.vscode/tasks.json`.

## Delv Debugger

```bash
dlv
```

## Checking coverage

In vscode make sure you have

```json
    "go.coverageDecorator": {
        "type": "gutter",
    },
```

Then run command shift P, then select:

> Go: Toggle Test Coverage in Current Package

And you should see coverage gutters in the UI of vscode along the side of your code by the line numbers.

## Data Wrangler

The first time you rebuild a container you may need to kick Data Wrangler into action by using the command palette and searching for "Data Wrangler: Open File" and selecting the file you want to open. That should initialize the Python virtual env interpreter and launch the extension.

If Data Wangler is ever ditched, Python and Juptyer extensions can be removed from the container if not being used by other extensions.
