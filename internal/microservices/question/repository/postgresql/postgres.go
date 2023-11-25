package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type QuestionRep struct {
	db postgresql.DbConn
}

func NewRepository(db postgresql.DbConn) *QuestionRep {
	return &QuestionRep{
		db: db,
	}
}

func (r *QuestionRep) CreateAnswer(ctx context.Context, userID uuid.UUID, a models.Answer) error {
	_, err := r.db.Exec(ctx, "INSERT INTO user(user_id, rating, (SELECT question_id FROM question WHERE name = $3)) VALUES ($1, $2)", userID, a.Rating, a.Name)
	if err != nil {
		return fmt.Errorf("error creating answer: %w", err)
	}

	return nil
}

func (r *QuestionRep) CheckUserAnswer(ctx context.Context, userID, questionName string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM answer a JOIN question q ON a.question_id = q.id WHERE a.user_id = $1 AND q.name = $2)", userID, questionName).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r *QuestionRep) CalculateAverageRating(ctx context.Context, questionName string) (int, error) {
	var averageRating sql.NullInt16

	err := r.db.QueryRow(ctx, "SELECT AVG(a.rating) FROM answer a JOIN question q ON a.question_id = q.id WHERE q.name = $1", questionName).Scan(&averageRating)
	if err != nil {
		return 0.0, err
	}

	if !averageRating.Valid {
		return 0.0, nil
	}

	return int(averageRating.Int16), nil
}
