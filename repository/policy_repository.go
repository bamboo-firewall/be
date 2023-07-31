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

type policyRepository struct {
	database   mongo.Database
	collection string
}

func (r *policyRepository) GetTotal(c context.Context) (int64, error) {
	collection := r.database.Collection(r.collection)
	total, err := collection.CountDocuments(c, bson.D{})
	return total, err
}

func (r *policyRepository) Fetch(c context.Context) ([]models.GlobalNetworkPolicies, error) {
	collection := r.database.Collection(r.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var policies []models.GlobalNetworkPolicies

	err = cursor.All(c, &policies)
	if policies == nil {
		return []models.GlobalNetworkPolicies{}, err
	}

	return policies, err
}

func (r *policyRepository) Search(c context.Context, searchOptions bson.M) ([]models.GlobalNetworkPolicies, error) {
	collection := r.database.Collection(r.collection)

	opts := options.Find()
	cursor, err := collection.Find(c, searchOptions, opts)

	if err != nil {
		return nil, err
	}

	var policies []models.GlobalNetworkPolicies

	err = cursor.All(c, &policies)
	if policies == nil {
		return []models.GlobalNetworkPolicies{}, err
	}

	return policies, err
}

func (r *policyRepository) AggGroupBy(c context.Context, filter bson.M, key string, jsonPath string) ([]domain.Option, error) {
	collection := r.database.Collection(r.collection)

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

func NewPolicyRepository(db mongo.Database, collection string) domain.PolicyRepository {
	return &policyRepository{
		database:   db,
		collection: collection,
	}
}
