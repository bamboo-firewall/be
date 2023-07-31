package domain

import "time"

type Metadata struct {
	Name              string            `bson:"name"`
	UID               string            `bson:"uid"`
	CreationTimestamp time.Time         `bson:"creationTimestamp"`
	Labels            map[string]string `bson:"labels"`
}
