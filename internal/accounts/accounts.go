package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/serge-masiutin/go-template/internal/database"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

const MaxEmailBytes = 254

var ErrCredentials = errors.New("invalid email or password")

type User struct {
	ID    int64  `json:"id,string"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

type Store struct {
	pool      *database.DB
	dummyHash []byte
}

func New(pool *database.DB) (*Store, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("unused-constant-time-comparison"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Store{pool: pool, dummyHash: hash}, nil
}

func (s *Store) Create(ctx context.Context, email, password string, admin bool) (User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || len(email) > MaxEmailBytes {
		return User{}, fmt.Errorf("email must be a valid address with at most 254 bytes")
	}
	if len(password) < 12 || len(password) > 72 {
		return User{}, fmt.Errorf("password must contain 12 to 72 bytes")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	record := accountRecord{Email: email, PasswordHash: hash, Admin: admin}
	err = gorm.G[accountRecord](s.pool.ORM).Create(ctx, &record)
	return record.public(), err
}

func (s *Store) Authenticate(ctx context.Context, email, password string) (User, error) {
	record, err := gorm.G[accountRecord](s.pool.ORM).Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		bcrypt.CompareHashAndPassword(s.dummyHash, []byte(password))
		return User{}, ErrCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("find account: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(record.PasswordHash, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return User{}, ErrCredentials
		}
		return User{}, fmt.Errorf("verify password hash: %w", err)
	}
	return record.public(), nil
}

func (s *Store) Find(ctx context.Context, id int64) (User, error) {
	record, err := gorm.G[accountRecord](s.pool.ORM).Select("id", "email", "admin").Where("id = ?", id).First(ctx)
	return record.public(), err
}

func (s *Store) Count(ctx context.Context) (int, error) {
	count, err := gorm.G[accountRecord](s.pool.ORM).Count(ctx, "*")
	return int(count), err
}

// Persistence records never cross an HTTP or AI boundary.
type accountRecord struct {
	ID           int64 `gorm:"primaryKey"`
	Email        string
	PasswordHash []byte
	Admin        bool
}

func (accountRecord) TableName() string { return "users" }
func (r accountRecord) public() User    { return User{ID: r.ID, Email: r.Email, Admin: r.Admin} }
