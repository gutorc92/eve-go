package collections

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Weight struct {
	Value float64            `bson:"value,omitempty" json:"value,omitempty"`
	Date  primitive.DateTime `bson:"date_measure,omitempty" json:"date_measure,omitempty"`
}

type Cow struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name    string             `bson:"name,omitempty" json:"name,omitempty"`
	Batch   string             `bson:"batch,omitempty" json:"batch,omitempty"`
	Weights []Weight           `bson:"weights,omitempty" json:"weights,omitempty"`
}
