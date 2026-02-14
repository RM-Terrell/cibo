# Common Commands Used During Development

## Prerequisites

- Go (see `go.mod` for version)
- Node.js (see `react_ui/package.json` for engine compatibility)
- GNU Make

After cloning the repo, install the React UI dependencies:

```bash
make install-ui
```

## Quick Reference

Run `make` or `make help` from the project root to see all available targets:

```bash
make help
```

| Task | Command |
| :--- | :--- |
| Build everything (UI + Go binary) | `make build` |
| Launch the TUI (builds UI first) | `make run-tui` |
| Launch the TUI with mock API | `make run-tui-mock` |
| Launch standalone web mode (sample data) | `make run-web` |
| Start Vite dev server (hot reload, no Go rebuild) | `make run-ui-dev` |
| Run all tests (Go + React) | `make test` |
| Run Go tests only | `make test-go` |
| Run React tests only | `make test-ui` |
| Run all linters | `make lint` |
| Clean build artifacts | `make clean` |

## Compiling and Running

All build and run commands are executed from the **project root** using `make`. The Makefile handles `cd`-ing into the correct directories and building dependencies in the right order.

To launch the full TUI application:

```bash
make run-tui
```

To launch in standalone web mode using a sample data file:

```bash
make run-web
```

This uses `sample_data/AAPL.parquet` by default. To use a different file, run the Go command directly:

```bash
cd react_ui &amp;&amp; npm run build
cd cmd &amp;&amp; go run . -webMode ../sample_data/TICKER_GOES_HERE.parquet
```

### React UI Development

When iterating on the frontend, use Vite's dev server for hot module replacement. This does not require rebuilding the Go binary:

```bash
make run-ui-dev
```

Note: The Vite dev server proxies API requests to the Go backend, so you'll need the Go server running separately for `/api/data` to work.

## Unit Testing

### From the Command Line

Run all tests (Go and React):

```bash
make test
```

Or run them independently:

```bash
make test-go
make test-ui
```

### From VS Code

- **Tasks:** VS Code Command Palette → Run Task → "Run Go Tests \<TEST\_SUBSET\>". This is configured in `.vscode/tasks.json`.
- **Test Explorer:** The dev container includes a testing extension that adds a flask icon in the sidebar. You can run and debug individual tests from there.
- **Keyboard shortcuts (Mac defaults):**
    - Run test at cursor: `Cmd` + `;` then `C`
    - Run all tests in file: `Cmd` + `;` then `F`

## Linting

Run all linters (Go + React):

```bash
make lint
```

Or independently:

```bash
make lint-go    # staticcheck
make lint-ui    # ESLint
```

## Debugging

To manually invoke the `dlv` debugger:

```bash
dlv debug ./cmd
```

Debugging via VS Code is configured in `.vscode/launch.json`. Select the desired configuration in the Debug panel and launch with the green play button. This starts both the `dlv`-based debugger and the application in a terminal window. Breakpoints will be hit as you interact with the application.

## Unit testing

You can run unit tests via VS Code Tasks:

VS Code Command Palette -> Run Task -> "Run Go Tests <TEST_SUBSET>". This is configured in `.vscode/tasks.json`.

As part of dev container setup there should also be a unit testing extension installed that will both add a "Testing" flask icon on the left side of the editor from which you can run tests, and also add an overlay on unit test code for "run" and "debug" to run them right from the code with one click.

To run the unit test at your cursor use `Command` + `;` followed by `C`, to run all tests in the file use `Command` + `;` followed by `F` (assuming default vs code shortcuts on Mac).

## Checking Coverage

In your VS Code `settings.json`, ensure you have:

```json
    "go.coverageDecorator": {
        "type": "gutter",
    },
```

Then run `Cmd` + `Shift` + `P` and select:

> Go: Toggle Test Coverage in Current Package

Coverage gutters will appear alongside line numbers in the editor.

## Data Wrangler

The first time you rebuild a dev container you may need to initialize Data Wrangler by using the command palette and searching for **"Data Wrangler: Open File"**, then selecting the Parquet file you want to inspect. This sets up the Python virtual environment and launches the extension.

If Data Wrangler is ever ditched, the Python and Jupyter extensions can also be removed from the container if they aren't used by other extensions or features.
