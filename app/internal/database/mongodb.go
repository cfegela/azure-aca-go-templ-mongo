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

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(client *mongo.Client, dbName string) *TaskRepository {
	collection := client.Database(dbName).Collection("tasks")
	return &TaskRepository{
		collection: collection,
	}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	task.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *TaskRepository) FindAll(ctx context.Context) ([]models.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func (r *TaskRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}

	var task models.Task
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) FindByIDAndUserID(ctx context.Context, id string, userID primitive.ObjectID) (*models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}

	var task models.Task
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "user_id": userID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, id string, task *models.Task) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	task.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"due_date":    task.DueDate,
			"updated_at":  task.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *TaskRepository) UpdateByUserID(ctx context.Context, id string, userID primitive.ObjectID, task *models.Task) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	task.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"due_date":    task.DueDate,
			"updated_at":  task.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID, "user_id": userID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *TaskRepository) DeleteByUserID(ctx context.Context, id string, userID primitive.ObjectID) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID, "user_id": userID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
