# CIBO

A stock price analysis assistant written in Go with a Bubbletea based terminal UI and a React based web UI. Currently in VERY early phase development.

## For local building and development

- Git
- VSCode
- Docker

After cloning the application, open the file directory in VS Code, then click the prompt to reopen in a devcontainer. If you don't get the UI prompt you can find the command via the Command Palette. Once the container builds all needed software should be installed that is required to run code and work on the app including debugging, data exploration, compilation, etc.

## If you want to explore raw data from the parquet files

You'll need the following extensions installed in VS Code, installed in the devcontainer (not needed to for compilation or execution of the application).

- Data Wrangler
- Python extension (Data Wrangler dependency)
- Jupyter extension (Data Wrangler dependency)

These SHOULD all auto install and configure based on the setup in the `.devcontainer/Dockerfile` and `.devcontainer/devcontainer.json` files. See `common_commands.md` for more info on them.
