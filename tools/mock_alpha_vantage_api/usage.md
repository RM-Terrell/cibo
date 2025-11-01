# Using the mock server

In one window start the mock API server via

```bash
go run .
```

Or via its VS Code Task.

Then in another terminal window ,launch the CLI application with the flag `-mockAPI`

```bash
go run . -mockAPI
```

and the full application should now be hitting the fake API server, which you should be able to confirm by seeing server logs go by in its terminal window.
