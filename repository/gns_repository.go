package repository

import (
	"context"

	models "github.com/bamboo-firewall/watcher/model"

	"github.com/bamboo-firewall/be/domain"
	"github.com/bamboo-firewall/be/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type gnsRepository struct {
	database   mongo.Database
	collection string
}

func (r *gnsRepository) GetTotal(c context.Context) (int64, error) {
	collection := r.database.Collection(r.collection)
	total, err := collection.CountDocuments(c, bson.D{})
	return total, err
}

func (gr *gnsRepository) Fetch(c context.Context) ([]models.GlobalNetworkSet, error) {
	collection := gr.database.Collection(gr.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var gns []models.GlobalNetworkSet

	err = cursor.All(c, &gns)
	if gns == nil {
		return []models.GlobalNetworkSet{}, err
	}

	return gns, err
}

func (gr *gnsRepository) Search(c context.Context, searchOptions bson.M) ([]models.GlobalNetworkSet, error) {
	collection := gr.database.Collection(gr.collection)

	opts := options.Find()
	cursor, err := collection.Find(c, searchOptions, opts)

	if err != nil {
		return nil, err
	}

	var networkset []models.GlobalNetworkSet

	err = cursor.All(c, &networkset)
	if networkset == nil {
		return []models.GlobalNetworkSet{}, err
	}

	return networkset, err
}

func (gr *gnsRepository) AggGroupBy(c context.Context, filter bson.M, key string, jsonPath string) ([]domain.Option, error) {
	collection := gr.database.Collection(gr.collection)

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$group": bson.M{
				"_id": bson.TypeNull,
				key: bson.M{
					"$addToSet": jsonPath,
				},
			},
		},
		{
			"$project": bson.M{
				"_id": 0,
				key:   1,
			},
		},
	}

	cursor, err := collection.Aggregate(c, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := cursor.All(c, &results); err != nil {
		return nil, err
	}
	var options []domain.Option
	if len(results) == 0 {
		return []domain.Option{}, nil
	}
	for _, item := range results[0][key].(primitive.A) {
		option := domain.Option{
			Key:   key,
			Value: item.(string),
		}
		options = append(options, option)
	}

	if options == nil {
		return []domain.Option{}, err
	}

	return options, err
}

func NewGNSRepository(db mongo.Database, collection string) domain.GNSRepository {
	return &gnsRepository{
		database:   db,
		collection: collection,
	}
}
