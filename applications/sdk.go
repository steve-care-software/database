package appliations

import (
	"github.com/steve-care-software/database/applications/queries"
	"github.com/steve-care-software/database/applications/transactions"
	"github.com/steve-care-software/cryptography/domain/hash"
)

// NewApplication creates a new application instance
func NewApplication(
	query queries.Application,
	trx transactions.Application,
) Application {
	return createApplication(query, trx)
}

// Builder represents the application builder
type Builder interface {
	Create() Builder
	WithApplication(application hash.Hash) Builder
	Now() (Application, error)
}

// Application represents a database application
type Application interface {
	Query() queries.Application
	Transaction() transactions.Application
}
