package handlers

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/gutorc92/eve-go/collections"
	"github.com/gutorc92/eve-go/dao"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetItems(dt *dao.DataMongo, domain collections.Domain, req *RequestParameters) (reflect.Value, error) {
	typ, err := domain.Schema.CreateStruct(true)
	if err != nil {
		fmt.Println("Error to create struct", err)
		return reflect.Value{}, err
	}
	// v := reflect.New(typ).Elem()
	// req := NewRequestParameters(r.URL.Query())
	doc := req.WhereClause()
	// fmt.Println("Collection name:", domain.GetCollectionName(), req.MaxResults)
	slice := reflect.MakeSlice(reflect.SliceOf(typ), 5, req.MaxResults)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	err = dt.FindAll(domain.GetCollectionName(), doc, x.Interface(), req.RequestParameters2MongOptions())
	fmt.Println("Kind on find", reflect.Indirect(x).Len())
	if err != nil {
		log.Error("Error to find all", err)
		return reflect.Value{}, err
	}
	return x, nil
}

func SaveItem(dt *dao.DataMongo, domain collections.Domain, body io.ReadCloser) (reflect.Value, error) {
	typ, err := domain.Schema.CreateStruct(true)
	if err != nil {
		log.Debug("Error to create struct: %s", err)
		return reflect.Value{}, err
	}
	x := reflect.New(typ)
	err = json.NewDecoder(body).Decode(x.Interface())
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		log.Debug("Cannot unmarshal json: %s", err)
		return reflect.Value{}, err
	}
	// TODO: set real datetime
	date := primitive.NewDateTimeFromTime(time.Now())
	reflect.Indirect(x).FieldByName("CreatedAt").Set(reflect.ValueOf(date))
	reflect.Indirect(x).FieldByName("UpdateAt").Set(reflect.ValueOf(date))
	hasher := sha1.New()
	fields := reflect.Indirect(x).Field(0).String()
	// TODO: add all fields to etag creation
	hasher.Write([]byte(fields))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	reflect.Indirect(x).FieldByName("Etag").SetString(sha)
	fmt.Println("struct to insert", x)
	id, err := dt.Insert(domain.GetCollectionName(), x.Interface())
	if err != nil {
		log.Debug("Error to insert data: %s", err)
		return reflect.Value{}, err
	}
	structFieldValue := reflect.Indirect(x).FieldByName("ID")
	if !structFieldValue.CanSet() {
		fmt.Errorf("Cannot set %s field value")
	}

	if !structFieldValue.IsValid() {
		fmt.Errorf("No such field: in obj")
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(id)
	if structFieldType != val.Type() {
		fmt.Errorf("Provided value %v type %v didn't match obj field type %v", val, val.Type(), structFieldType)
	}
	structFieldValue.Set(val)
	fmt.Println("x", x)
	return x, nil
}
