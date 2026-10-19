package problems

import (
	"context"
	"fmt"

	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	sqlQuery := `
		DELETE FROM problems
		WHERE id = $1;
	`

	conn := dbpostgres.ExecutorFromContext(ctx, r.pool)
	tag, err := conn.Exec(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.NotFound(fmt.Sprintf("problem with id=%s not found", id))
	}

	return nil
}
