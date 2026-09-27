package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	pool      *pgxpool.Pool
	dummyHash []byte
}

func New(pool *pgxpool.Pool) (*Store, error) {
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
	var user User
	err = s.pool.QueryRow(ctx, "INSERT INTO users(email,password_hash,admin) VALUES($1,$2,$3) RETURNING id,email,admin", email, hash, admin).Scan(&user.ID, &user.Email, &user.Admin)
	return user, err
}

func (s *Store) Authenticate(ctx context.Context, email, password string) (User, error) {
	var user User
	var hash []byte
	err := s.pool.QueryRow(ctx, "SELECT id,email,admin,password_hash FROM users WHERE email=$1", strings.ToLower(strings.TrimSpace(email))).Scan(&user.ID, &user.Email, &user.Admin, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		bcrypt.CompareHashAndPassword(s.dummyHash, []byte(password))
		return User{}, ErrCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("find account: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return User{}, ErrCredentials
		}
		return User{}, fmt.Errorf("verify password hash: %w", err)
	}
	return user, nil
}

func (s *Store) Find(ctx context.Context, id int64) (User, error) {
	var user User
	err := s.pool.QueryRow(ctx, "SELECT id,email,admin FROM users WHERE id=$1", id).Scan(&user.ID, &user.Email, &user.Admin)
	return user, err
}

func (s *Store) Count(ctx context.Context) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&count)
	return count, err
}
