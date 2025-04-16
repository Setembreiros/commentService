package integration_test_arrange

import (
	"commentservice/cmd/provider"
	database "commentservice/internal/db"
	"context"
)

func CreateTestDatabase(ctx context.Context) *database.Database {
	provider := provider.NewProvider("test", "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable")
	sqlDb, err := provider.ProvideDb()
	if err != nil {
		panic(err)
	}
	return database.NewDatabase(sqlDb)
}

func GetNextCommentId() uint64 {
	provider := provider.NewProvider("test", "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable")
	sqlDb, err := provider.ProvideDb()
	if err != nil {
		panic(err)
	}
	return sqlDb.GetNextCommentId()
}
