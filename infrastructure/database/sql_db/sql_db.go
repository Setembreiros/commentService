package sql_db

import (
	model "commentservice/internal/model/domain"
	"database/sql"

	"github.com/rs/zerolog/log"
)

type SqlDatabase struct {
	Client *sql.DB
}

func NewDatabase(connStr string) (*SqlDatabase, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't open a connection with the database")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Database is not reachable")
		return nil, err
	}

	return &SqlDatabase{
		Client: db,
	}, nil
}

func (sd *SqlDatabase) CreateComment(comment *model.Comment) (uint64, error) {
	tx, err := sd.Client.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `
		INSERT INTO comments (
			postId, 
			username, 
			content, 
			createdAt, 
		) VALUES (?, ?, ?, ?)
		RETURNING id
	`
	var commentId uint64

	err = tx.QueryRow(
		query,
		comment.PostId,
		comment.Username,
		comment.Content,
		comment.CreatedAt).Scan(&commentId)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Insert comment failed")
		return 0, err
	}

	return commentId, nil
}
