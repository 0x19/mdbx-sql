package mdbxsql

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func randomUserGo(rand *rand.Rand, _ int) *UserGo {
	return &UserGo{
		ID:   int32(rand.Intn(1000000)),
		Name: randomString(rand),
		Age:  int32(rand.Intn(100)),
	}
}

func randomString(rand *rand.Rand) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := rand.Intn(20) + 1 // Random length between 1 and 20
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func TestBatchInsert(t *testing.T) {
	ctx := context.Background()
	db, err := NewDb(ctx, "testdb", 10)
	require.NoError(t, err)
	defer db.Close()

	schema := NewSchema(db)
	userTable, err := schema.CreateTable("users", "ID")
	require.NoError(t, err)

	// Generate 10,000 random UserGo records
	users := make([]Model, 10000)
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10000; i++ {
		users[i] = randomUserGo(randSource, i)
	}

	// Test Batch Insert
	start := time.Now()
	err = BatchInsert(userTable, users)
	log.Printf("Batch insert operation completed in %v", time.Since(start))
	require.NoError(t, err)
}
