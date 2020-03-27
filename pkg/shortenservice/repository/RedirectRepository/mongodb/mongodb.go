package mongo

import (
	"context"
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/repository"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

const MONGODB_COLLECTION string = "redirects"

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (repository.IRedirectRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout),
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	} else {
		logger.Log.Info("Connect Mongodb Successfully")
	}
	repo.client = client
	return repo, nil
}

func (r *mongoRepository) Find(code string) (*model.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout*time.Second)
	defer cancel()
	redirect := &model.Redirect{}
	collection := r.client.Database(r.database).Collection(MONGODB_COLLECTION)
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(config.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

func (r *mongoRepository) Store(redirect *model.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout*time.Second)
	defer cancel()
	collection := r.client.Database(r.database).Collection(MONGODB_COLLECTION)
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":      redirect.Code,
			"url":       redirect.URL,
			"create_at": redirect.CreateAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
