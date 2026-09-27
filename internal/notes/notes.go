package notes

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrBody = errors.New("write a note with 1 to 2000 characters")
var ErrNotFound = errors.New("note not found")

type Note struct {
	ID   int64  `json:"id,string"`
	Body string `json:"body"`
}

type Store struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Store { return &Store{pool: pool} }

func (s *Store) List(ctx context.Context, userID int64) ([]Note, error) {
	rows, err := s.pool.Query(ctx, "SELECT id,body FROM notes WHERE user_id=$1 ORDER BY id DESC LIMIT 50", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	notes := make([]Note, 0)
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Body); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func (s *Store) Create(ctx context.Context, userID int64, body string) error {
	body = strings.TrimSpace(body)
	if !utf8.ValidString(body) || utf8.RuneCountInString(body) < 1 || utf8.RuneCountInString(body) > 2000 {
		return ErrBody
	}
	_, err := s.pool.Exec(ctx, "INSERT INTO notes(user_id,body) VALUES($1,$2)", userID, body)
	return err
}

func (s *Store) Delete(ctx context.Context, userID, id int64) error {
	result, err := s.pool.Exec(ctx, "DELETE FROM notes WHERE user_id=$1 AND id=$2", userID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
