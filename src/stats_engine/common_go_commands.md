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

> Go: Test Coverage in Current Package

And you should see coverage gutters in the UI of vscode along the side of your code by the line numbers.
