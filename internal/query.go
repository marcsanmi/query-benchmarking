package internal

// Query compose all the necessary params to run a PromQL expression query.
type Query struct {
	// Params contains all the query params like start, end and step.
	Params map[string]string
}
