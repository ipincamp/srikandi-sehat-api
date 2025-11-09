package graphql

// It's the entry point for injecting core services into the adapter.
func NewResolver() *Resolver {
	// When services are added, they would be injected here.
	// Example:
	// return &Resolver{userService: us}
	return &Resolver{}
}
