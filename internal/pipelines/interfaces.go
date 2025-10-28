package pipelines

import (
	"cibo/internal/types"
	"io"
)

/*
	The main point of these interfaces is to make unit testing easy. By wrapping a
	pile of code like the API Client or IO in an interface, we can swap it out in code that
	uses it when testing so we don't have to fight concrete implementations of it.

	When in doubt, if it sucks to unit test on its own, wrap it in an interface here so
	the suck doesn't double. I generally dislike endless layers of abstraction without
	doing real work, but in this case it's well justified. I even ran into it in the
	wild and wrote about it here: https://www.dashdashforce.dev/posts/golang-interface-testing
*/

/*
	If you wind up here via a "go to definition" trying to find the source code of a function, you'll want
	to use the vscode right click "go to implementation" to find the source of the concrete implementation,
	instead of command click / "go to definition".

	If using Goland as an editor you should see a "see definitions" or "see implementations" option to jump to it.
	This initial confusion of traversing imports is an unfortunate side effect of interfaces in Go.
*/

type APIClient interface {
	FetchDailyPrice(ticker string) ([]byte, error)
	FetchEarnings(ticker string) ([]byte, error)
}

type ParquetWriter interface {
	WriteCombinedPriceDataToParquet(records []types.CombinedPriceRecord, writer io.WriteCloser) (string, error)
}

type FairValuePipeline interface {
	RunPipeline(input LynchFairValueInputs) (*LynchFairValueOutputs, error)
}
