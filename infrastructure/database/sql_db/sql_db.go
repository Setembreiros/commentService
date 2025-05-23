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
		query := fmt.Sprintf("DELETE FROM commentservice.%s", table)
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
		INSERT INTO commentservice.comments (
        	postId, 
        	username, 
        	content, 
        	createdAt,
			updatedAt
    	) VALUES ($1, $2, $3, $4, $5)
    	RETURNING id
	`
	var commentId uint64

	err = tx.QueryRow(
		query,
		comment.PostId,
		comment.Username,
		comment.Content,
		comment.CreatedAt,
		comment.UpdatedAt).Scan(&commentId)
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
			createdAt,
			updatedAt
		FROM commentservice.comments 
		WHERE id = $1
	`

	var comment model.Comment
	err := sd.Client.QueryRow(query, id).Scan(
		&comment.Id,
		&comment.PostId,
		&comment.Username,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt)

	comment.CreatedAt = comment.CreatedAt.UTC()
	comment.UpdatedAt = comment.UpdatedAt.UTC()

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
		SELECT nextval('commentservice.comments_id_seq')
	`

	var lastId uint64
	err := sd.Client.QueryRow(query).Scan(&lastId)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to get next comment id")
		return 0
	}

	return lastId + uint64(1)
}

func (sd *SqlDatabase) UpdateComment(comment *model.Comment) error {
	tx, err := sd.Client.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `
		UPDATE commentservice.comments 
		SET 
			content = $1,
			updatedAt = $2
		WHERE id = $3
	`

	result, err := tx.Exec(
		query,
		comment.Content,
		comment.UpdatedAt,
		comment.Id)

	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to update comment with id %d", comment.Id)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to get rows affected when updating comment with id %d", comment.Id)
		return err
	}

	if rowsAffected == 0 {
		log.Warn().Msgf("No comment found with id %d to update", comment.Id)
		return nil
	}

	log.Info().Msgf("Comment with id %d updated successfully", comment.Id)
	return nil
}

func (sd *SqlDatabase) DeleteComment(id uint64) error {
	tx, err := sd.Client.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `DELETE FROM commentservice.comments WHERE id = $1`
	result, err := tx.Exec(query, id)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to delete comment with id %d", id)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to get rows affected when deleting comment with id %d", id)
		return err
	}

	if rowsAffected == 0 {
		log.Warn().Msgf("No comment found with id %d to delete", id)
		return nil
	}

	log.Info().Msgf("Comment with id %d deleted successfully", id)
	return nil
}
