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
	CreateComment(data *model.Comment) (uint64, error)
}

func NewDatabase(client DatabaseClient) *Database {
	return &Database{
		Client: client,
	}
}
