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
}

func NewRequestParameters(values url.Values) RequestParameters {
	req := RequestParameters{0}
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
	return req
}

func (req *RequestParameters) RequestParameters2MongOptions() *options.FindOptions {
	findOptions := options.FindOptions{}
	m := int64(req.MaxResults)
	findOptions.Limit = &m
	return &findOptions
}

func GETHandler(dt *dao.DataMongo, collectionName string, schema collections.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		typ, err := schema.CreateStruct()
		if err != nil {
			fmt.Println("Error to create struct", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// v := reflect.New(typ).Elem()
		req := NewRequestParameters(r.URL.Query())
		fmt.Println("Collection name:", collectionName, req.MaxResults)
		slice := reflect.MakeSlice(reflect.SliceOf(typ), 5, req.MaxResults)
		x := reflect.New(slice.Type())
		x.Elem().Set(slice)
		err = dt.FindAll(collectionName, x.Interface(), req.RequestParameters2MongOptions())
		if err != nil {
			fmt.Println("Error to find all", err)
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
