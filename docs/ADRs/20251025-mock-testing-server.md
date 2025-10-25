# The mock testing server

In `/tools/` there's a directory now that contains another web server written in Go. This is meant for dev / testing purposes, as with Alpha Vantage there is a limit of 25 calls per day on a free account. When testing the whole application heavily its very easy to hit such a low limit.

The intent of that server is to be an API server that returns the same (albeit unchanging and static) data that the application use to generate a data file and render it. Using it was originally intended to be done in a separate terminal window. So in one window start the mock API server, then in another terminal window launch the CLI application itself. Minimal and simple with low performance overhead.

Any changes to the alpha vantage api and the way the application consumes it will obviously require changes to the mock server.