package database

import (
	"context"
	"errors"
	"time"

	"github.com/cfegela/azure-aca-go-templ-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InviteRepository struct {
	collection *mongo.Collection
}

func NewInviteRepository(client *mongo.Client, dbName string) *InviteRepository {
	collection := client.Database(dbName).Collection("invites")
	return &InviteRepository{
		collection: collection,
	}
}

func (r *InviteRepository) Create(ctx context.Context, invite *models.Invite) error {
	invite.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, invite)
	if err != nil {
		return err
	}

	invite.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *InviteRepository) FindByToken(ctx context.Context, token string) (*models.Invite, error) {
	var invite models.Invite
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&invite)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}
	return &invite, nil
}

func (r *InviteRepository) MarkUsed(ctx context.Context, token string) error {
	now := time.Now()
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"used_at": now}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("invite not found")
	}
	return nil
}

func (r *InviteRepository) FindAll(ctx context.Context) ([]models.Invite, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invites []models.Invite
	if err = cursor.All(ctx, &invites); err != nil {
		return nil, err
	}

	if invites == nil {
		invites = []models.Invite{}
	}

	return invites, nil
}

func (r *InviteRepository) CreateIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "token", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}
