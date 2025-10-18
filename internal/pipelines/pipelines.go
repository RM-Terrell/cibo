package pipelines

// Pipelines is the root container for all application business logic pipelines.
// It acts as a single dependency for the UI layers instead of listing every pipeline
// individually in the struct

type Pipelines struct {
	LynchFairValue FairValuePipeline
	// Add new pipelines here in the future
}

func NewPipelines(client APIClient, writer ParquetWriter) *Pipelines {
	return &Pipelines{
		LynchFairValue: NewLynchFairValuePipeline(client, writer),
		// Add new pipelines here in the future
	}
}
