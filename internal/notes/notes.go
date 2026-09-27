package notes

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/serge-masiutin/go-template/internal/database"
	"gorm.io/gorm"
)

var ErrBody = errors.New("write a note with 1 to 2000 characters")
var ErrNotFound = errors.New("note not found")

type Note struct {
	ID   int64  `json:"id,string"`
	Body string `json:"body"`
}

type Store struct{ pool *database.DB }

func New(pool *database.DB) *Store { return &Store{pool: pool} }

func (s *Store) List(ctx context.Context, userID int64) ([]Note, error) {
	records, err := gorm.G[noteRecord](s.pool.ORM).Select("id", "body").Where("user_id = ?", userID).Order("id DESC").Limit(50).Find(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]Note, 0, len(records))
	for _, record := range records {
		items = append(items, Note{ID: record.ID, Body: record.Body})
	}
	return items, nil
}

func (s *Store) Create(ctx context.Context, userID int64, body string) error {
	body = strings.TrimSpace(body)
	if !utf8.ValidString(body) || utf8.RuneCountInString(body) < 1 || utf8.RuneCountInString(body) > 2000 {
		return ErrBody
	}
	return gorm.G[noteRecord](s.pool.ORM).Create(ctx, &noteRecord{UserID: userID, Body: body})
}

func (s *Store) Delete(ctx context.Context, userID, id int64) error {
	affected, err := gorm.G[noteRecord](s.pool.ORM).Where("user_id = ? AND id = ?", userID, id).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

type noteRecord struct {
	ID     int64 `gorm:"primaryKey"`
	UserID int64
	Body   string
}

func (noteRecord) TableName() string { return "notes" }
