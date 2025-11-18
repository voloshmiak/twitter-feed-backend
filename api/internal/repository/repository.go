package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CockroachRepo struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *CockroachRepo {
	return &CockroachRepo{
		conn: conn,
	}
}

func (r *CockroachRepo) SaveMessage(ctx context.Context, msg *Message) error {
	query := `
		INSERT INTO messages (id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.conn.Exec(ctx, query, msg.id, msg.userID, msg.content, msg.createdAt)
	return err
}

func (r *CockroachRepo) GetAllMessages(ctx context.Context) ([]*Message, error) {
	query := `SELECT id, user_id, content, created_at FROM messages ORDER BY created_at ASC`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	messages, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*Message, error) {
		var msg Message
		err := row.Scan(&msg.id, &msg.userID, &msg.content, &msg.createdAt)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	})
	if err != nil {
		return nil, err
	}

	return messages, nil
}
