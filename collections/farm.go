package collections


import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Farm struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
}