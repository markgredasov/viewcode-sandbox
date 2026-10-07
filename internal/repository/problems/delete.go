package problems

import (
	"context"
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		DELETE FROM problems
		WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.NotFound(fmt.Sprintf("problem with id=%s not found", id))
	}

	return nil
}
