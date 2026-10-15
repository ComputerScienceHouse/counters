package database

import (
	"context"
	"fmt"
	"time"

	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/util"
	"github.com/jackc/pgx/v5"
)

const MAX_RESETS_IN_QUERY = 10

func CreateReset(ctx context.Context, userID string, username string, req *models.CreateResetInput) (string, error) {
	resetID, err := GenerateUUID()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO resets (reset_id, counter_id, user_id, username, description, occured_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, resetID, req.CounterID, userID, username, req.Description, time.Unix(util.ValidateTime(req.ResetTime), 0))

	return resetID, err
}

func CreateResetWithTransaction(ctx context.Context, tx pgx.Tx, userID string, username string, req *models.CreateResetInput) (string, error) {
	resetID, err := GenerateUUID()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO resets (reset_id, counter_id, user_id, username, description, occured_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, resetID, req.CounterID, userID, username, req.Description, time.Unix(util.ValidateTime(req.ResetTime), 0))

	return resetID, err
}

func DeleteReset(ctx context.Context, counterid string, resetid string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = db.Exec(ctx, `
		DELETE FROM resets WHERE reset_id = $1
	`, resetid)

	return err
}

func GetResetOwner(ctx context.Context, resetID string) (string, error) {
	var ownerId string

	err := db.QueryRow(ctx, `
		SELECT user_id FROM resets WHERE reset_id = $1
	`, resetID).Scan(&ownerId)

	return ownerId, err
}

func GetCounterFromReset(ctx context.Context, resetID string) (string, error) {
	var counterID string

	if err := db.QueryRow(ctx, `
		SELECT r.counter_id FROM resets r WHERE r.reset_id = $1
	`, resetID).Scan(&counterID); err != nil {
		return "", err
	}

	return counterID, nil
}

func GetResetsFromCounterId(ctx context.Context, counterID string, filter *models.GetResetsFilter) ([]*models.ResetListPart, error) {
	args := make([]any, 0)
	args = append(args, counterID)

	postQuery := `
		SELECT r.reset_id, r.description, r.user_id, r.username, r.occured_at FROM resets r
		WHERE r.counter_id = $1
	`
	if filter.Timestamp != nil {
		args = append(args, time.Unix(*filter.Timestamp, 0))

		if filter.SortDescending {
			postQuery = fmt.Sprintf("%s AND r.occured_at < $%d", postQuery, len(args))
		} else {
			postQuery = fmt.Sprintf("%s AND r.occured_at > $%d", postQuery, len(args))
		}
	}

	if filter.SortDescending {
		postQuery = fmt.Sprintf("%s ORDER BY r.occured_at DESC LIMIT %d", postQuery, MAX_RESETS_IN_QUERY)
	} else {
		postQuery = fmt.Sprintf("%s ORDER BY r.occured_at ASC LIMIT %d", postQuery, MAX_RESETS_IN_QUERY)
	}

	rows, err := db.Query(ctx, postQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resets := make([]*models.ResetListPart, 0)

	for rows.Next() {
		var (
			resetId          string
			resetDescription string
			resetOwner       string
			resetUsername    string
			resetOccurredAt  time.Time
		)

		if err := rows.Scan(&resetId, &resetDescription, &resetOwner, &resetUsername, &resetOccurredAt); err != nil {
			return nil, err
		}

		resets = append(resets, &models.ResetListPart{
			ResetID:          resetId,
			ResetDescription: resetDescription,
			ResetOwner:       resetOwner,
			ResetUsername:    resetUsername,
			ResetOccurredAt:  resetOccurredAt.Unix(),
		})
	}

	return resets, nil
}

func InitReset(ctx context.Context) error {
	if _, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS resets (
			reset_id 	UUID PRIMARY KEY,
			counter_id 	UUID REFERENCES counters(counter_id) ON DELETE CASCADE,
			user_id 	UUID NOT NULL,
			username 	TEXT NOT NULL,
			description TEXT NOT NULL,
			occured_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return err
	}

	return nil
}
