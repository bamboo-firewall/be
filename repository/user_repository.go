package repository

import (
	"context"

	"github.com/bamboo-firewall/be/domain"
	"github.com/bamboo-firewall/be/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func (ur *userRepository) Update(c context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)
	_, err := collection.UpdateOne(c, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (ur *userRepository) DeleteById(c context.Context, id string) error {
	collection := ur.database.Collection(ur.collection)
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(c, bson.M{"_id": idHex})
	return err
}

func (ur *userRepository) GetTotal(c context.Context) (int64, error) {
	collection := ur.database.Collection(ur.collection)
	total, err := collection.CountDocuments(c, bson.D{})
	return total, err
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)

	_, err := collection.InsertOne(c, user)

	return err
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []domain.User

	err = cursor.All(c, &users)
	if users == nil {
		return []domain.User{}, err
	}

	return users, err
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)
	var user domain.User
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	return user, err
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	var user domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return user, err
}

func NewUserRepository(db mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}
