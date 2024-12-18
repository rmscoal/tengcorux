package tracing

func (t *tracing) Version() string {
	return Version()
}

func Version() string {
	return "v0.1.1"
}
