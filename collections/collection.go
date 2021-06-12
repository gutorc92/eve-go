package collections

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection interface {
	GetUrl() string
	GetCollectionName() string
}

type Field struct {
	Type      string `json:"type"`
	Required  bool   `json:"required,omitempty"`
	MinLength int    `json:"min_length,omitempty"`
}

type Schema struct {
	Fields map[string]Field
}

func typeOf(typeof string) reflect.Type {
	switch typeof {
	case "string":
		return reflect.TypeOf("")
	case "integer":
		return reflect.TypeOf(0)
	case "boolean":
		return reflect.TypeOf(true)
	}

	return reflect.TypeOf("")
}

func createTag(name string) reflect.StructTag {
	return reflect.StructTag(fmt.Sprintf(`bson:"%s,omitempty" json:"%s"`, name, name))
}

func (s *Schema) CreateStruct(list bool) (reflect.Type, error) {
	fields := make([]reflect.StructField, 0, len(s.Fields))
	for key, field := range s.Fields {
		fields = append(fields, reflect.StructField{
			Name: strings.Title(key),
			Type: typeOf(field.Type),
			Tag:  createTag(key),
		})
	}
	if list == true {
		fields = append(fields, reflect.StructField{
			Name: "ID",
			Type: reflect.TypeOf(primitive.ObjectID{}),
			Tag:  `bson:"_id,omitempty" json:"_id"`,
		})
	}
	typ := reflect.StructOf(fields)
	return typ, nil
}

func (s *Schema) UnmarshalJSON(b []byte) error {
	s.Fields = make(map[string]Field)
	var converted map[string]json.RawMessage
	err := json.Unmarshal(b, &converted)
	if err != nil {
		return err
	}
	for key, f := range converted {
		var field Field
		err := json.Unmarshal(f, &field)
		if err != nil {
			return err
		}
		s.Fields[key] = field
	}
	return nil
}

type Datasource struct {
	Source string `json:"source,omitempty"`
}

type Domain struct {
	URL             string     `json:"url,omitempty"`
	Schema          Schema     `json:"schema,omitempty"`
	ResourceMethods []string   `json:"resource_methods,omitempty"`
	Datasource      Datasource `json:"datasource,omitempty"`
}

func (d *Domain) GetUrl() string {
	return fmt.Sprintf("/%s", d.URL)
}

func (d *Domain) GetUrlItem() string {
	return fmt.Sprintf("/%s/{id:[0-9]+}", d.URL)
}

func (d *Domain) GetCollectionName() string {
	if d.Datasource.Source == "" {
		return d.URL
	}
	return d.Datasource.Source
}

// func (dapi *DefaultBatchAPI) GetUrl() string {
// 	return fmt.Sprintf("/%s", dapi.collection)
// }
