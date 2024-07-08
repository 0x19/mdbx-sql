package mdbxsql

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchemaCreateAndGetTable(t *testing.T) {
	ctx := context.Background()
	db, err := NewDb(ctx, "testdb", 10)
	require.NoError(t, err)
	defer db.Close()

	schema := NewSchema(db)

	// Test CreateTable
	userTable, err := schema.CreateTable("users", "ID")
	require.NoError(t, err)
	require.NotNil(t, userTable)
	require.Equal(t, "users", userTable.Name)
	require.Equal(t, "ID", userTable.Primary)

	// Test GetTable
	retrievedTable, err := schema.GetTable("users")
	require.NoError(t, err)
	require.NotNil(t, retrievedTable)
	require.Equal(t, userTable, retrievedTable)

	// Test GetTable with a non-existing table
	_, err = schema.GetTable("nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "table not found")
}
