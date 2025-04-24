package integration_test_arrange

import (
	"commentservice/cmd/provider"
	database "commentservice/internal/db"
	model "commentservice/internal/model/domain"
	integration_test_assert "commentservice/test/integration_test_common/assert"
	"context"
	"testing"
)

func CreateTestDatabase(ctx context.Context) *database.Database {
	provider := provider.NewProvider("test", "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable")
	sqlDb, err := provider.ProvideDb()
	if err != nil {
		panic(err)
	}
	return database.NewDatabase(sqlDb)
}

func AddComment(t *testing.T, comment *model.Comment) uint64 {
	provider := provider.NewProvider("test", "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable")
	sqlDb, err := provider.ProvideDb()
	if err != nil {
		panic(err)
	}

	commentId, err := sqlDb.CreateComment(comment)
	if err != nil {
		panic(err)
	}

	integration_test_assert.AssertCommentExists(t, database.NewDatabase(sqlDb), commentId, comment)

	return commentId
}

func GetNextCommentId() uint64 {
	provider := provider.NewProvider("test", "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable")
	sqlDb, err := provider.ProvideDb()
	if err != nil {
		panic(err)
	}
	return sqlDb.GetNextCommentId()
}
