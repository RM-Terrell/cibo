
# VS Code configs and settings notes

## launch.json and debugging

The line
```json
"console": "integratedTerminal"
```

is needed to launch a "real" terminal that can be interacted with, which is needed being a TUI application. Youll get
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

without it.