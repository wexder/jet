package sqlbuilder

import "bytes"

type SelectStatementTable struct {
	statement SelectStatement
	columns   []Column
	alias     string
}

func (s *SelectStatementTable) Columns() []Column {
	return s.columns
}

func (s *SelectStatementTable) RefIntColumnName(name string) Column {
	intColumn := NewIntegerColumn(name, NotNullable)
	intColumn.setTableName(s.alias)

	return intColumn
}

func (s *SelectStatementTable) RefIntColumn(column Column) *IntegerColumn {
	intColumn := NewIntegerColumn(column.TableName()+"."+column.Name(), NotNullable)
	intColumn.setTableName(s.alias)

	return intColumn
}

func (s *SelectStatementTable) RefStringColumn(column Column) *StringColumn {
	strColumn := NewStringColumn(column.Name(), NotNullable)
	strColumn.setTableName(column.TableName())
	return strColumn
}

func (s *SelectStatementTable) SerializeSql(out *bytes.Buffer) error {
	out.WriteString("( ")
	statementStr, err := s.statement.String()

	if err != nil {
		return err
	}

	out.WriteString(statementStr)

	out.WriteString(" ) AS ")
	out.WriteString(s.alias)

	return nil
}

// Generates a select query on the current tableName.
func (s *SelectStatementTable) SELECT(projections ...Projection) SelectStatement {
	return newSelectStatement(s, projections)
}

// Creates a inner join tableName expression using onCondition.
func (s *SelectStatementTable) INNER_JOIN(table ReadableTable, onCondition BoolExpression) ReadableTable {
	return InnerJoinOn(s, table, onCondition)
}

//func (s *SelectStatementTable) InnerJoinUsing(table ReadableTable, col1 Column, col2 Column) ReadableTable {
//	return INNER_JOIN(s, table, col1.Eq(col2))
//}

// Creates a left join tableName expression using onCondition.
func (s *SelectStatementTable) LeftJoinOn(table ReadableTable, onCondition BoolExpression) ReadableTable {
	return LeftJoinOn(s, table, onCondition)
}

// Creates a right join tableName expression using onCondition.
func (s *SelectStatementTable) RightJoinOn(table ReadableTable, onCondition BoolExpression) ReadableTable {
	return RightJoinOn(s, table, onCondition)
}

func (s *SelectStatementTable) FULL_JOIN(table ReadableTable, onCondition BoolExpression) ReadableTable {
	return FullJoin(s, table, onCondition)
}

func (s *SelectStatementTable) CrossJoin(table ReadableTable) ReadableTable {
	return CrossJoin(s, table)
}
