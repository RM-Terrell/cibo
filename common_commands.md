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

## Unit testing

As part of dev container setup there should be a unit testing extension installed that will both add a "Testing" flask icon on the left side of the editor from which you can run tests, and also add an overlay on unit test code for "run" and "debug" to run them right from the code with one click.

To run the unit test at your cursor use `Command` + `;` followed by `C`, to run all tests in the file use `Command` + `;` followed by `F` (assuming default vs code shortcuts on Mac).

## Checking coverage

In vscode make sure you have

```json
    "go.coverageDecorator": {
        "type": "gutter",
    },
```

Then run `Command` + `shift` + `P`, then select:

> Go: Toggle Test Coverage in Current Package

And you should see coverage gutters in the UI of vscode along the side of your code by the line numbers.

## Data Wrangler

The first time you rebuild a container you may need to kick Data Wrangler into action by using the command palette and searching for "Data Wrangler: Open File" and selecting the file you want to open. That should initialize the Python virtual env interpreter and launch the extension.

If Data Wangler is ever ditched, Python and Juptyer extensions can be removed from the container if not being used by other extensions.
