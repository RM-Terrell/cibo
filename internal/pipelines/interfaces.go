package pipelines

import (
	"cibo/internal/types"
	"io"
)

// The main point of these interfaces is to make unit testing easy. By wrapping a
// pile of code like the API Client or IO in an interface, we can swap it out in code that
// uses it when testing so we dont have to fight concrete implementations of it.

// When in doubt, if it sucks to unit test on its own, wrap it in an interface here so
// the suck doesnt double. I generally dislike endless layers of abstraction without
// doing real work, but in this case it's well justified. I even ran into it in the
// wild and wrote about it here: https://www.dashdashforce.dev/posts/golang-interface-testing

/*
	If you wind up here trying to find a definition of a function in the code, you'll want
	to look in the concrete implementation passed in higher up the code stack (probably main.go),
	and find the function with the matching name from the interface, or locate it with a whole file system search.
	If using Goland as an editor you might be able to click "see definitions" or "see implementations" to jump to it.
	this initial confusion following imports is an unfortunate side effect of interfaces in Go so aim to
	keep them at a minimum.
*/

type APIClient interface {
	FetchDailyPrice(ticker string) ([]byte, error)
	FetchEarnings(ticker string) ([]byte, error)
}

type ParquetWriter interface {
	WriteCombinedPriceDataToParquet(records []types.CombinedPriceRecord, writer io.WriteCloser) error
}

type FairValuePipeline interface {
	RunPipeline(input LynchFairValueInputs) (*LynchFairValueOutputs, error)
}
