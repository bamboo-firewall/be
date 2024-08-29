package repository

import "github.com/bamboo-firewall/be/cmd/server/pkg/storage"

type PolicyDB struct {
	mongo *storage.PolicyDB
}

func NewPolicy(mongo *storage.PolicyDB) *PolicyDB {
	return &PolicyDB{mongo: mongo}
}
