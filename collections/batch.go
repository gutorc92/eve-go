package collections

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Batch struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
	Farm string             `bson:"farm,omitempty" json:"farm,omitempty"`
}
