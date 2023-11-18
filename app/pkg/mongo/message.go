package mongo

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-service/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageStorage struct {
	collection *mongo.Collection
	log zerolog.Logger
}

func NewMessageStorage(db *mongo.Database, log zerolog.Logger) models.MessageStorage {
	return &MessageStorage{
		collection: db.Collection("messages"),
		log: log,
	}
}

func (ms *MessageStorage) FindAllByRoomId(ctx context.Context, roomID string) ([]*models.Message, error) {
	query := bson.M{"room_id": roomID}
	
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{Key: "created_at", Value: 1}})
	queryOptions.SetLimit(50)

	cursor, err := ms.collection.Find(ctx, query, &queryOptions)
	if err != nil {
		ms.log.Error().Msg(err.Error())
		return nil, err
	}

	var messages []*models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		ms.log.Error().Msg(err.Error())
		return nil, err
	}

	return messages, nil
}

func (ms *MessageStorage) Insert(ctx context.Context, message *models.Message) (*models.Message, error) {
	message.CreatedAt = time.Now()
	message.ID = primitive.NewObjectID().Hex()

	if _, err := ms.collection.InsertOne(ctx, message); err != nil {
		ms.log.Error().Msg(err.Error())
		return nil, err
	}

	return message, nil
}