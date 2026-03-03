package expenseswiring

import (
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	coredb "github.com/Hughzu/trackstack/apps/server/internal/core/db"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
	expensesdb "github.com/Hughzu/trackstack/apps/server/internal/modules/expenses/adapters/db"
	"github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
)

type ExpensesDependencies struct {
	Service *expenses.Service
	Close   func() error
}

func BuildExpenses(cfg config.Config) (ExpensesDependencies, error) {
	dsn, err := common.BuildTursoDSN(common.TursoConfig{
		Mode:    cfg.DBConnectionMode,
		URL:     cfg.TursoExpensesURL,
		URLHTTP: cfg.TursoExpensesURLHTTP,
		URLWS:   cfg.TursoExpensesURLWS,
		Token:   cfg.TursoExpensesToken,
	})
	if err != nil {
		return ExpensesDependencies{}, fmt.Errorf("expenses dsn: %w", err)
	}

	db, err := coredb.OpenLibSQL(dsn)
	if err != nil {
		return ExpensesDependencies{}, fmt.Errorf("expenses db: %w", err)
	}

	store := expensesdb.NewExpensesStore(db)
	service := expenses.NewService(store)

	return ExpensesDependencies{
		Service: service,
		Close:   db.Close,
	}, nil
}
