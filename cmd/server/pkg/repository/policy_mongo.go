package repository

import "github.com/bamboo-firewall/be/cmd/server/pkg/storage"

type PolicyMongo struct {
	mongo *storage.PolicyMongo
}

func NewPolicyMongo(mongo *storage.PolicyMongo) *PolicyMongo {
	return &PolicyMongo{mongo: mongo}
}
