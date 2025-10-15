package pipelines

// AppPipelines is the root container for all application business logic pipelines.
// It acts as a single dependency for the UI layers instead of listing every pipeline
// individually in the struct, although that's a valid way to do it too.

type Pipelines struct {
	LynchFairValuePipeline *LynchFairValuePipeline
	// Add new pipelines here in the future
}

func NewPipelines(client APIClient, writer ParquetWriter) *Pipelines {
	return &Pipelines{
		LynchFairValuePipeline: NewLynchFairValuePipeline(client, writer),
		// Add new pipelines here in the future
	}
}
