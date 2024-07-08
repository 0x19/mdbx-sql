package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected *SelectStatement
	}{
		{
			input: "SELECT name, age FROM users WHERE name = 'John' AND age = 30",
			expected: &SelectStatement{
				Columns:   []string{"name", "age"},
				TableName: "users",
				Conditions: []Condition{
					{Field: "name", Op: "=", Value: "John"},
					{Field: "age", Op: "=", Value: "30"},
				},
			},
		},
		{
			input: "SELECT users.name, users.age FROM users WHERE users.name = 'John'",
			expected: &SelectStatement{
				Columns:   []string{"users.name", "users.age"},
				TableName: "users",
				Conditions: []Condition{
					{Field: "users.name", Op: "=", Value: "John"},
				},
			},
		},
		{
			input: "SELECT user_id, user_name FROM users WHERE user_id = 1",
			expected: &SelectStatement{
				Columns:   []string{"user_id", "user_name"},
				TableName: "users",
				Conditions: []Condition{
					{Field: "user_id", Op: "=", Value: "1"},
				},
			},
		},
		{
			input: "SELECT orders.user_id, SUM(orders.amount) FROM orders JOIN users ON orders.user_id = users.id WHERE users.age > 30",
			expected: &SelectStatement{
				Columns:   []string{"orders.user_id"},
				TableName: "orders",
				Joins: []Join{
					{Table: "users", OnField1: "orders.user_id", OnField2: "users.id"},
				},
				Aggregates: []Aggregate{
					{Func: "SUM", Column: "orders.amount"},
				},
				Conditions: []Condition{
					{Field: "users.age", Op: ">", Value: "30"},
				},
			},
		},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		stmt := p.Parse()

		require.NotNil(t, stmt, "expected non-nil statement for input: %s", tt.input)

		selectStmt, ok := stmt.(*SelectStatement)
		require.True(t, ok, "expected *SelectStatement type for input: %s", tt.input)

		require.Equal(t, len(tt.expected.Columns), len(selectStmt.Columns), "expected columns length %d, got %d", len(tt.expected.Columns), len(selectStmt.Columns))

		for i, col := range selectStmt.Columns {
			require.Equal(t, tt.expected.Columns[i], col, "expected column %s, got %s", tt.expected.Columns[i], col)
		}

		require.Equal(t, tt.expected.TableName, selectStmt.TableName, "expected table name %s, got %s", tt.expected.TableName, selectStmt.TableName)
		require.Equal(t, len(tt.expected.Conditions), len(selectStmt.Conditions), "expected conditions length %d, got %d", len(tt.expected.Conditions), len(selectStmt.Conditions))

		for i, cond := range selectStmt.Conditions {
			require.Equal(t, tt.expected.Conditions[i].Field, cond.Field, "expected condition field %s, got %s", tt.expected.Conditions[i].Field, cond.Field)
			require.Equal(t, tt.expected.Conditions[i].Op, cond.Op, "expected condition op %s, got %s", tt.expected.Conditions[i].Op, cond.Op)
			require.Equal(t, tt.expected.Conditions[i].Value, cond.Value, "expected condition value %s, got %s", tt.expected.Conditions[i].Value, cond.Value)
		}

		require.Equal(t, len(tt.expected.Joins), len(selectStmt.Joins), "expected joins length %d, got %d", len(tt.expected.Joins), len(selectStmt.Joins))
		for i, join := range selectStmt.Joins {
			require.Equal(t, tt.expected.Joins[i].Table, join.Table, "expected join table %s, got %s", tt.expected.Joins[i].Table, join.Table)
			require.Equal(t, tt.expected.Joins[i].OnField1, join.OnField1, "expected join onField1 %s, got %s", tt.expected.Joins[i].OnField1, join.OnField1)
			require.Equal(t, tt.expected.Joins[i].OnField2, join.OnField2, "expected join onField2 %s, got %s", tt.expected.Joins[i].OnField2, join.OnField2)
		}

		require.Equal(t, len(tt.expected.Aggregates), len(selectStmt.Aggregates), "expected aggregates length %d, got %d", len(tt.expected.Aggregates), len(selectStmt.Aggregates))
		for i, agg := range selectStmt.Aggregates {
			require.Equal(t, tt.expected.Aggregates[i].Func, agg.Func, "expected aggregate func %s, got %s", tt.expected.Aggregates[i].Func, agg.Func)
			require.Equal(t, tt.expected.Aggregates[i].Column, agg.Column, "expected aggregate column %s, got %s", tt.expected.Aggregates[i].Column, agg.Column)
		}
	}
}
