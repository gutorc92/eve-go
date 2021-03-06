package collections

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cow struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
	Batch string `bson:"name,omitempty"`
}