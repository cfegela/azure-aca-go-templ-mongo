package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invite struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token     string             `json:"token" bson:"token"`
	Email     string             `json:"email" bson:"email"`
	InvitedBy primitive.ObjectID `json:"invited_by" bson:"invited_by"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	UsedAt    *time.Time         `json:"used_at,omitempty" bson:"used_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func GenerateInviteToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (i *Invite) IsValid() bool {
	return i.UsedAt == nil && time.Now().Before(i.ExpiresAt)
}
