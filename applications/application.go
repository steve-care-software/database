package appliations

import (
	"github.com/steve-care-software/database/applications/queries"
	"github.com/steve-care-software/database/applications/transactions"
)

type application struct {
	query queries.Application
	trx   transactions.Application
}

func createApplication(
	query queries.Application,
	trx transactions.Application,
) Application {
	out := application{
		query: query,
		trx:   trx,
	}

	return &out
}

// Query returns the query application
func (app *application) Query() queries.Application {
	return app.query
}

// Transaction returns the transaction application
func (app *application) Transaction() transactions.Application {
	return app.trx
}
