package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rlarkin212/url-shortener/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client  *mongo.Client
	databse string
	timeout time.Duration
}

//newMongoClient generates new mongo client
func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	//conntct to db
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	//check if connected
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

//NewMongoRepository create new version of repository
func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortener.RedirectRespository, error) {
	repo := &mongoRepository{
		timeout: time.Duration(mongoTimeout) * time.Second,
		databse: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	}
	repo.client = client

	return repo, nil
}

func (r *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := generateContext(r.timeout)
	defer cancel()

	redirect := &shortener.Redirect{}

	collection := getCollection(r.client, r.databse, "redirects")
	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(&redirect)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}

		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil
}

func (r *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := generateContext(r.timeout)
	defer cancel()

	collection := getCollection(r.client, r.databse, "redirects")

	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)

	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}

func generateContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func getCollection(client *mongo.Client, database string, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}
