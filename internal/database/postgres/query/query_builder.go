package query_builder

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	count    int
	args     []any
	sqlQuery string
}

func New(queryExists string) QueryBuilder {
	queryExists = strings.TrimSpace(queryExists)
	queryExists = strings.TrimSuffix(queryExists, ";")

	return QueryBuilder{
		count:    0,
		args:     []any{},
		sqlQuery: fmt.Sprintf("%s WHERE 1=1", queryExists),
	}
}

func (qb *QueryBuilder) Get() (string, []any) {
	return qb.sqlQuery, qb.args
}

func (qb *QueryBuilder) AddString(column, value string) {
	if value == "" {
		return
	}

	qb.args = append(qb.args, value)
	qb.count++
	qb.sqlQuery += fmt.Sprintf(" AND %s = $%d", column, qb.count)
}

func (qb *QueryBuilder) AddTypeArray(column string, values []string, typeName string) {
	if len(values) == 0 {
		return
	}

	placeholders := make([]string, len(values))
	for i := range values {
		qb.count++
		placeholders[i] = fmt.Sprintf("$%d", qb.count)
		qb.args = append(qb.args, values[i])
	}

	qb.sqlQuery += fmt.Sprintf(
		" AND %s && ARRAY[%s]::%s[]",
		column,
		strings.Join(placeholders, ","),
		typeName,
	)
}
