package postgres

import (
	"errors"
	"github.com/go-jet/jet/internal/jet"
	"strconv"
	"strings"
)

var Dialect = NewDialect()

func NewDialect() jet.Dialect {

	serializeOverrides := map[string]jet.SerializeOverride{}
	serializeOverrides["REGEXP_LIKE"] = postgres_REGEXP_LIKE_function

	dialectParams := jet.DialectParams{
		Name:                "PostgreSQL",
		PackageName:         "postgres",
		CastOverride:        castFunc,
		SerializeOverrides:  serializeOverrides,
		AliasQuoteChar:      '"',
		IdentifierQuoteChar: '"',
		ArgumentPlaceholder: func(ord int) string {
			return "$" + strconv.Itoa(ord)
		},
		SetClause:         postgresSetClause,
		SupportsReturning: true,
	}

	return jet.NewDialect(dialectParams)
}

func castFunc(expression jet.Expression, castType string) jet.SerializeFunc {
	return func(statement jet.StatementType, out *jet.SqlBuilder, options ...jet.SerializeOption) error {
		if err := jet.Serialize(expression, statement, out, options...); err != nil {
			return err
		}
		out.WriteString("::" + castType)
		return nil
	}
}

func postgresSetClause(columns []jet.IColumn, values []jet.Clause, out *jet.SqlBuilder) (err error) {
	if len(columns) > 1 {
		out.WriteString("(")
	}

	err = jet.SerializeColumnNames(columns, out)

	if err != nil {
		return
	}

	if len(columns) > 1 {
		out.WriteString(")")
	}

	out.WriteString("=")

	if len(values) > 1 {
		out.WriteString("(")
	}

	err = jet.SerializeClauseList(jet.UpdateStatementType, values, out)

	if err != nil {
		return
	}

	if len(values) > 1 {
		out.WriteString(")")
	}

	return
}

func postgres_REGEXP_LIKE_function(expressions ...jet.Expression) jet.SerializeFunc {
	return func(statement jet.StatementType, out *jet.SqlBuilder, options ...jet.SerializeOption) error {
		if len(expressions) < 2 {
			return errors.New("jet: invalid number of expressions for operator")
		}

		if err := jet.Serialize(expressions[0], statement, out, options...); err != nil {
			return err
		}

		caseSensitive := false

		if len(expressions) >= 3 {
			if stringLiteral, ok := expressions[2].(jet.LiteralExpression); ok {
				matchType := stringLiteral.Value().(string)

				caseSensitive = !strings.Contains(matchType, "i")
			}
		}

		if caseSensitive {
			out.WriteString("~")
		} else {
			out.WriteString("~*")
		}

		if err := jet.Serialize(expressions[1], statement, out, options...); err != nil {
			return err
		}

		return nil
	}
}
