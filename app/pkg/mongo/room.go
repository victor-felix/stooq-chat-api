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

type RoomStorage struct {
	collection *mongo.Collection
	log zerolog.Logger
}

func NewRoomStorage(db *mongo.Database, log zerolog.Logger) models.RoomStorage {
	return &RoomStorage{
		collection: db.Collection("rooms"),
		log: log,
	}
}

func (rs *RoomStorage) FindAll(ctx context.Context) ([]*models.Room, error) {
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := rs.collection.Find(ctx, bson.M{}, &queryOptions)
	if err != nil {
		rs.log.Error().Msg(err.Error())
		return nil, err
	}

	var rooms []*models.Room
	if err := cursor.All(ctx, &rooms); err != nil {
		rs.log.Error().Msg(err.Error())
		return nil, err
	}

	return rooms, nil
}

func (rs *RoomStorage) Insert(ctx context.Context, room *models.Room) (*models.Room, error) {
	room.CreatedAt = time.Now()
	room.ID = primitive.NewObjectID().Hex()

	if _, err := rs.collection.InsertOne(ctx, room); err != nil {
		rs.log.Error().Msg(err.Error())
		return nil, err
	}

	return room, nil
}

func (rs *RoomStorage) FindOneByID(ctx context.Context, id string) (*models.Room, error) {
	var room models.Room
	if err := rs.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&room); err != nil {
		rs.log.Error().Msg(err.Error())
		return nil, err
	}

	return &room, nil
}
