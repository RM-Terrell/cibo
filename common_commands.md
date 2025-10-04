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

**CURRENT ISSUE**: There is currently a quirk in the unit testing extension where you must highlight the text name value (between the double quotes) of the unit test in order for the inline "run" and "debug" buttons to run the desired test. I think this is an issue with how I structured and named the tests causing the extension to be confused. The left side panel Testing bit works perfectly though.

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
