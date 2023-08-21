package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bamboo-firewall/be/bootstrap"
	"github.com/bamboo-firewall/be/domain"
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const adminRole = "Admin"

func main() {
	app := bootstrap.App()
	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	// Check if admin user already exists
	var admin domain.User
	res := db.Collection(domain.CollectionUser).FindOne(context.Background(), bson.M{"email": "admin@example.com"}).Decode(&admin)
	if res == nil {
		log.Println("Admin user already exists")
		return
	}

	// Create admin user
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(env.AdminPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Println(err)
		return
	}

	user := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     env.AdminAccount,
		Username: env.AdminAccount,
		Email:    fmt.Sprintf("%s@%s", env.AdminAccount, env.EmailDomain),
		Password: string(encryptedPassword),
		Role:     adminRole,
	}

	_, err = db.Collection(domain.CollectionUser).InsertOne(context.Background(), &user)
	if err != nil {
		log.Println(err)
		return
	}

	adapterConfig := mongodbadapter.AdapterConfig{
		DatabaseName: env.DBName,
	}

	adapter, err := mongodbadapter.NewAdapterByDB(db.Client().MongoClient(), &adapterConfig)

	if err != nil {
		log.Println(err)
		return
	}

	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		log.Println(err)
		return
	}
	// Add policy for admin
	enforcer.AddGroupingPolicy(user.ID.Hex(), user.Role)

	log.Println("Admin user created successfully")
}
