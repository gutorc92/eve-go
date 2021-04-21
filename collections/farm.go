package collections

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Farm struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
}

func (f *Farm) GetUrl() string {
	return fmt.Sprintf("/%s", "farm")
}

func (f *Farm) GetCollectionName() string {
	return "farm"
}
