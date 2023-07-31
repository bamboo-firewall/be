package optionutil

import (
	"strings"

	"github.com/bamboo-firewall/be/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func ConvertToBsonM(options []domain.Option, mapping map[string]string) bson.M {
	query := bson.M{}
	for _, item := range options {
		key := strings.Split(mapping[item.Key], "$")
		if len(key) > 1 {
			query[key[1]] = item.Value
		}
	}
	return query
}
