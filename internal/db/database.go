package database

import (
	model "commentservice/internal/model/domain"

	_ "github.com/lib/pq"
)

//go:generate mockgen -source=database.go -destination=test/mock/database.go

type Database struct {
	Client DatabaseClient
}

type DatabaseClient interface {
	Clean()
	CreateComment(data *model.Comment) (uint64, error)
	GetCommentById(id uint64) (*model.Comment, error)
}

func NewDatabase(client DatabaseClient) *Database {
	return &Database{
		Client: client,
	}
}
