
# VS Code configs and settings notes

## launch.json and debugging

[Official docs here](https://code.visualstudio.com/docs/debugtest/debugging)

Of note for this application, the line

```json
"console": "integratedTerminal"
```

is needed to launch a "real" terminal that can be interacted with, which is needed being a TUI application. Without it you'll get:

```console
Starting: /go/bin/dlv dap --listen=127.0.0.1:33735 --log-dest=3 from /workspaces/cibo/cmd
DAP server listening at: 127.0.0.1:33735
Type 'dlv help' for list of commands.
Successfully loaded configuration from: /workspaces/cibo/internal/statistics/keys/api_keys.toml
2025/10/24 13:34:29 There's been an error: could not open a new TTY: open /dev/tty: no such device or address
Process 10525 has exited with status 1
Detaching
dlv dap (8237) exited with code: 0
```

## Tasks

[Official docs here](https://code.visualstudio.com/docs/debugtest/tasks)

The `tasks.json` file is where I've scripted out common commands. Of note, the mock API server can be launched via a task, and that task is also launched as part of a specific debug config for debugging the application without hitting the live server. If you find yourself runnin things a lot that you think "i could shell script this" do it via a Task.
