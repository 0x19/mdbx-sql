package mdbxsql

import (
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestGetWithSQLParser(t *testing.T) {
	ctx := context.Background()
	db, err := NewDb(ctx, "testdb", 10)
	require.NoError(t, err)
	defer db.Close()

	schema := NewSchema(db)
	userTable, err := schema.CreateTable("users", "ID")
	require.NoError(t, err)

	// Insert a record for testing
	user := &UserGo{ID: 1, Name: "John Doe", Age: 30}
	start := time.Now()
	err = Insert(userTable, user)
	require.NoError(t, err)
	log.Printf("Insert operation completed in %v", time.Since(start))

	// Test Get with SQL parser
	sqlQuery := "SELECT name, age FROM users WHERE id=1"
	retrievedUser := &UserGo{}
	start = time.Now()
	err = Get(db, sqlQuery, retrievedUser)
	log.Printf("Get operation completed in %v", time.Since(start))
	require.NoError(t, err)
	require.Equal(t, user.ID, retrievedUser.ID)
	require.Equal(t, user.Name, retrievedUser.Name)
	require.Equal(t, user.Age, retrievedUser.Age)
	log.Printf("Retrieved User: %+v", retrievedUser)
}
