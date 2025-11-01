# Common commands used during development

## Compiling and running

To run the whole application `cd` into `/cmd/` and run

```bash
go run .
```

To run just the web server, also in `/cmd/` and (assuming the data file is present in the same directory)

```bash
go run . -webMode ./TICKER_GOES_HERE.parquet
```

## Debugging

To manually invoke the `dlv` debugger you can just run

```bash
dlv
```

Debugging the running application can be done via VS Code however, as defined in `.vscode/launch.json`. To debug the whole application just select the config desired in the Debug panel of VS Code and launch it via the green play button. That will launch both the dlv based vs code debugger, and the application itself in a terminal window. Interacting with the application will now trip break points.

## Unit testing

You can run unit tests via VS Code Tasks:

VS Code Command Palette -> Run Task -> "Run Go Tests <TEST_SUBSET>". This is configured in `.vscode/tasks.json`.

As part of dev container setup there should also be a unit testing extension installed that will both add a "Testing" flask icon on the left side of the editor from which you can run tests, and also add an overlay on unit test code for "run" and "debug" to run them right from the code with one click.

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
