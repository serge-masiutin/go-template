package accounts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateLoginSession grants access only after authentication. SCS writes can
// refresh session data, but cannot recreate a grant revoked by another request.
func (s *Store) CreateLoginSession(ctx context.Context, userID int64, token, previousToken string, deadline time.Time) error {
	if token == "" {
		return errors.New("login requires a renewed session token")
	}
	hash, previousHash := sha256.Sum256([]byte(token)), sha256.Sum256([]byte(previousToken))
	return s.pool.ORM.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Expired grants are never accepted; reclaim them on subsequent sign-ins.
		if err := tx.Exec("DELETE FROM login_sessions WHERE expires_at <= now() OR token_hash = ?", hex.EncodeToString(previousHash[:])).Error; err != nil {
			return err
		}
		return tx.Exec("INSERT INTO login_sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)", hex.EncodeToString(hash[:]), userID, deadline).Error
	})
}

func (s *Store) FindBySession(ctx context.Context, token string) (User, error) {
	hash := sha256.Sum256([]byte(token))
	record, err := gorm.G[accountRecord](s.pool.ORM).
		Select("users.id", "users.email", "users.admin").
		Where("id IN (SELECT user_id FROM login_sessions WHERE token_hash = ? AND expires_at > now())", hex.EncodeToString(hash[:])).First(ctx)
	return record.public(), err
}

func (s *Store) RevokeSession(ctx context.Context, token string) error {
	hash := sha256.Sum256([]byte(token))
	return s.pool.ORM.WithContext(ctx).Exec("DELETE FROM login_sessions WHERE token_hash = ?", hex.EncodeToString(hash[:])).Error
}
