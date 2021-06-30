package handlers

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gutorc92/eve-go/collections"
	"github.com/gutorc92/eve-go/config"
	"github.com/gutorc92/eve-go/dao"
	"github.com/gutorc92/eve-go/metrics"
)

type API interface {
	InitConfig(w *config.WebConfig, dt *dao.DataMongo) API
	GetUrl() string
	GETHandler() http.Handler
	POSTHandler() http.Handler
}

func GETHandler(dt *dao.DataMongo, metrics *metrics.Metrics, domain collections.Domain) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := NewRequestParameters(r.URL.Query())
		result, err := GetItems(dt, domain, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics.CountApiCall("http", "200", "false", "", "GET", domain.GetUrl(), time.Since(start).Seconds())
		WriteJSONResponse(result, 200, w)
	})
}

func OPTIONSHandler(dt *dao.DataMongo, metrics *metrics.Metrics, domain collections.Domain) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")
		metrics.CountApiCall("http", "200", "false", "", "OPTIONS", domain.GetUrl(), time.Since(start).Seconds())
		w.WriteHeader(http.StatusOK)
		return
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
		hasher := sha1.New()
		fields := reflect.Indirect(x).Field(0).String()
		// TODO: add all fields to etag creation
		hasher.Write([]byte(fields))
		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		reflect.Indirect(x).FieldByName("Etag").SetString(sha)
		// TODO: set real datetime
		reflect.Indirect(x).FieldByName("CreatedAt").SetInt(time.Now().Unix())
		reflect.Indirect(x).FieldByName("UpdateAt").SetInt(time.Now().Unix())
		fmt.Println("struct to insert", x)
		id, err := dt.Insert(collectionName, x.Interface())
		if err != nil {
			fmt.Println("Error to error")
			w.WriteHeader(http.StatusBadRequest)
			return
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
		WriteJSONResponse(x.Interface(), 200, w)
	})
}

func WriteJSONResponse(payload interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
