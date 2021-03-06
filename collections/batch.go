package collections

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Batch struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
	Farm string `bson:"name,omitempty"`
}