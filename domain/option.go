package domain

const (
	CollectionOption = "options"
)

type Option struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}
