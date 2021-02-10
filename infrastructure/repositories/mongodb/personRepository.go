package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"persons.com/api/domain/person"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func NewMongoRepository(mongoURL string, mongoDB string, mongoTimeout int) (person.PersonRepository, error) {
	repository := &mongoRepository{
		timeout:  time.Duration(mongoTimeout),
		database: mongoDB,
	}

	client, err := NewMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.NewMongoRepository")
	}

	repository.client = client

	return repository, nil

}

func (m *mongoRepository) FindById(id string) (*person.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	personFound := &person.Person{}

	collection := m.client.Database(m.database).Collection("persons")

	filter := bson.M{"ID": id}

	err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"_id": 0})).Decode(&personFound)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(person.ErrPersonNotFound, "repository.Person.FindById")
		}
		return nil, errors.Wrap(err, "repository.Person.FindById")
	}

	return personFound, nil
}

func (m *mongoRepository) Create(person *person.Person) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("persons")

	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"ID":       person.ID,
			"Name":     person.Name,
			"LastName": person.LastName,
			"Age":      person.Age,
		},
	)

	if err != nil {
		return errors.Wrap(err, "repository.Person.Create")
	}

	return nil
}

func (m *mongoRepository) GetAll() ([]*person.Person, error) {

	recordsCollection := make([]*person.Person, 0)

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("persons")

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"_id": 0}))
	if err != nil {
		return nil, errors.Wrap(err, "repository.Person.Create")
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		personFound := new(person.Person)

		if err = cursor.Decode(&personFound); err != nil {
			return nil, errors.Wrap(err, "repository.Person.Create")
		}

		recordsCollection = append(recordsCollection, personFound)
	}

	return recordsCollection, nil

}
