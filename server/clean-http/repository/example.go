package repository

type exampleRepository struct {
	config ExampleRepositoryConfig
}

// ExampleRepositoryConfig represents the configuration for example repository
type ExampleRepositoryConfig struct {
}

// NewExampleRepository creates a new example repository
func NewExampleRepository(c ExampleRepositoryConfig) ExampleRepository {
	return &exampleRepository{
		config: c,
	}
}
