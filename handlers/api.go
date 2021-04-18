package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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
	}
	return req
}

func (req *RequestParameters) RequestParameters2MongOptions() *options.FindOptions {
	findOptions := options.FindOptions{}
	m := int64(req.MaxResults)
	findOptions.Limit = &m
	return &findOptions
}

// func GETHandler(wc *config.WebConfig, collectionName string, data interface{}) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Get the JSON body and decode into credentials
// 		var dt *dao.DataMongo
// 		dt, err := dao.NewDataMongo(wc.Uri, wc.Database)
// 		if err != nil {
// 			panic(err)
// 		}
// 		err = dt.FindAll(collectionName, &data)
// 		if err != nil {
// 			fmt.Println("Error to error", err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		q := r.URL.Query()

// 		WriteJSONResponse(data, 200, w)
// 	})
// }

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
