package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/gutorc92/api-farm/collections"
	"github.com/gutorc92/api-farm/config"
	"github.com/gutorc92/api-farm/dao"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type API interface {
	InitConfig(w *config.WebConfig, dt *dao.DataMongo) API
	GetUrl() string
	GETHandler() http.Handler
	POSTHandler() http.Handler
}

type ResultPage struct {
	Items interface{} `json:"_items"`
	Meta  *MetaPage   `json:"_meta"`
}

type MetaPage struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	MaxResults int `json:"max_results"`
}

func newMeta() *MetaPage {
	var meta MetaPage
	meta.Total = 20
	return &meta
}

type RequestParameters struct {
	MaxResults int
	Where      string
}

func NewRequestParameters(values url.Values) RequestParameters {
	req := RequestParameters{0, ""}
	max_results := values.Get("max_results")
	if max_results != "" {
		i, err := strconv.Atoi(max_results)
		if err != nil {
			log.Error("Cannot convert max_values to int")
			req.MaxResults = 50
		} else {
			req.MaxResults = i
		}
	} else {
		req.MaxResults = 50
	}
	where := values.Get("where")
	if where != "" {
		req.Where = where
	}
	return req
}

func (req *RequestParameters) WhereClause() interface{} {
	fmt.Println("where", req.Where)
	var doc interface{}
	if req.Where != "" {
		err := bson.UnmarshalExtJSON([]byte(req.Where), true, &doc)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		doc = bson.M{}
	}
	fmt.Println("where compiled", doc)
	return doc
}

func (req *RequestParameters) RequestParameters2MongOptions() *options.FindOptions {
	findOptions := options.FindOptions{}
	m := int64(req.MaxResults)
	findOptions.Limit = &m
	return &findOptions
}

func GETHandler(dt *dao.DataMongo, domain collections.Domain) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		typ, err := domain.Schema.CreateStruct(true)
		if err != nil {
			fmt.Println("Error to create struct", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// v := reflect.New(typ).Elem()
		req := NewRequestParameters(r.URL.Query())
		doc := req.WhereClause()
		fmt.Println("Collection name:", domain.GetCollectionName(), req.MaxResults)
		slice := reflect.MakeSlice(reflect.SliceOf(typ), 5, req.MaxResults)
		x := reflect.New(slice.Type())
		x.Elem().Set(slice)
		err = dt.FindAll(domain.GetCollectionName(), doc, x.Interface(), req.RequestParameters2MongOptions())
		if err != nil {
			fmt.Println("Error to find all", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("x", x)
		WriteJSONResponse(x.Interface(), 200, w)
	})
}

func POSTHandler(dt *dao.DataMongo, collectionName string, schema collections.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		typ, err := schema.CreateStruct(true)
		if err != nil {
			fmt.Println("Error to create struct", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		x := reflect.New(typ)
		err = json.NewDecoder(r.Body).Decode(x.Interface())
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = dt.Insert(collectionName, x.Interface())
		if err != nil {
			fmt.Println("Error to error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("x", x)
		WriteJSONResponse(x.Interface(), 200, w)
	})
}

func WriteJSONResponse(payload interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	meta := newMeta()
	result := ResultPage{Items: payload, Meta: meta}
	jEncoder := json.NewEncoder(w)
	jEncoder.SetEscapeHTML(false)
	err := jEncoder.Encode(result)
	if err != nil {
		fmt.Println("Error")
	}
}
