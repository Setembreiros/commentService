package sql_db

import (
	model "commentservice/internal/model/domain"
	"database/sql"
	"fmt"

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

func (sd *SqlDatabase) Clean() {
	tx, err := sd.Client.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Clean each table
	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		_, err = tx.Exec(query)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Failed to clean table %s", table)
		}
	}

	log.Info().Msg("Database cleaned successfully")
	return
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
        	createdAt
    	) VALUES ($1, $2, $3, $4)
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

func (sd *SqlDatabase) GetCommentById(id uint64) (*model.Comment, error) {
	query := `
		SELECT 
			id,
			postId, 
			username, 
			content, 
			createdAt
		FROM comments 
		WHERE id = $1
	`

	var comment model.Comment
	err := sd.Client.QueryRow(query, id).Scan(
		&comment.Id,
		&comment.PostId,
		&comment.Username,
		&comment.Content,
		&comment.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No comment found with the given ID
		}
		log.Error().Stack().Err(err).Msgf("Get comment by id %d failed", id)
		return nil, err
	}

	return &comment, nil
}

func (sd *SqlDatabase) GetNextCommentId() uint64 {
	query := `
		SELECT nextval('comments_id_seq')
	`

	var lastId uint64
	err := sd.Client.QueryRow(query).Scan(&lastId)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to get next comment id")
		return 0
	}

	return lastId + uint64(1)
}
