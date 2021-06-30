package handlers

import (
	"fmt"
	"reflect"

	"github.com/gutorc92/eve-go/collections"
	"github.com/gutorc92/eve-go/dao"
	log "github.com/sirupsen/logrus"
)

func GetItems(dt *dao.DataMongo, domain collections.Domain, req *RequestParameters) (interface{}, error) {
	typ, err := domain.Schema.CreateStruct(true)
	if err != nil {
		fmt.Println("Error to create struct", err)
		return nil, err
	}
	// v := reflect.New(typ).Elem()
	// req := NewRequestParameters(r.URL.Query())
	doc := req.WhereClause()
	// fmt.Println("Collection name:", domain.GetCollectionName(), req.MaxResults)
	slice := reflect.MakeSlice(reflect.SliceOf(typ), 5, req.MaxResults)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	err = dt.FindAll(domain.GetCollectionName(), doc, x.Interface(), req.RequestParameters2MongOptions())
	if err != nil {
		log.Error("Error to find all", err)
		return nil, err
	}
	return x.Interface(), nil
}
