package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect returns a mongo client already connected.
func Connect(ctx context.Context, mongoURL string, log zerolog.Logger) (*mongo.Client, error) {
	if strings.TrimSpace(mongoURL) == "" {
		return nil, errors.New("mongoURL is required")
	}

	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, fmt.Errorf("error on MongoDB connection: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, fmt.Errorf("error after MongoDB connection: %w", err)
	}

	return client, nil
}

// NewDatabase returns a new mongo database.
func NewDatabase(client *mongo.Client, databaseName string) *mongo.Database {
	return client.Database(databaseName)
}

// IsDuplicatedError check if error is deplicated mongo type.
func IsDuplicatedError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

// IsNotFoundError check if error is not found mongo type.
func IsNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}
