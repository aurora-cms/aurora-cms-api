package common

// Mapper defines a generic interface for mapping between domain entities and persistence models
type Mapper[E any, M any] interface {
	// ToModel converts a domain entity to a persistence model
	ToModel(entity E) (M, error)

	// ToDomain converts a persistence model to a domain entity
	ToDomain(model M) (E, error)

	// ToDomains converts a slice of persistence models to a slice of domain entities
	ToDomains(models []M) ([]E, error)

	// ToModels converts a slice of domain entities to a slice of persistence models
	ToModels(entities []E) ([]M, error)
}
