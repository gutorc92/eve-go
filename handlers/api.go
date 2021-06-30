package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gutorc92/eve-go/collections"
	"github.com/gutorc92/eve-go/config"
	"github.com/gutorc92/eve-go/dao"
)

type API interface {
	InitConfig(w *config.WebConfig, dt *dao.DataMongo) API
	GetUrl() string
	GETHandler() http.Handler
	POSTHandler() http.Handler
}

func GETHandler(dt *dao.DataMongo, domain collections.Domain) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequestParameters(r.URL.Query())
		result, err := GetItems(dt, domain, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		WriteJSONResponse(result, 200, w)
	})
}

func OPTIONSHandler(dt *dao.DataMongo, domain collections.Domain) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func POSTHandler(dt *dao.DataMongo, domain collections.Domain) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x, err := SaveItem(dt, domain, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		WriteJSONResponse(x, 200, w)
	})
}

func WriteJSONResponse(payload reflect.Value, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	meta := newMeta()
	meta.Total = reflect.Indirect(payload).Len()
	result := ResultPage{Items: payload.Interface(), Meta: meta}
	jEncoder := json.NewEncoder(w)
	jEncoder.SetEscapeHTML(false)
	err := jEncoder.Encode(result)
	if err != nil {
		fmt.Println("Error")
	}
}
