package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gutorc92/api-farm/collections"
	"github.com/gutorc92/api-farm/config"
	"github.com/gutorc92/api-farm/dao"
	"github.com/gutorc92/api-farm/metrics"
)

const (
	COW = "cow"
)

type DefaultCowAPI struct {
	*config.WebConfig
	dt *dao.DataMongo
	*metrics.Metrics
	collection string
}

func (dapi *DefaultCowAPI) InitConfig(w *config.WebConfig, dt *dao.DataMongo) API {
	dapi.WebConfig = w
	dapi.dt = dt
	dapi.Metrics = w.Metrics
	dapi.collection = COW
	return dapi
}

func (dapi *DefaultCowAPI) GetUrl() string {
	return fmt.Sprintf("/%s", dapi.collection)
}

func (dapi *DefaultCowAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JSON body and decode into credentials
		var data []collections.Cow
		err := dapi.dt.FindAll(dapi.collection, &data)
		if err != nil {
			fmt.Println("Error to error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		WriteJSONResponse(data, 200, w)
	})
}

func (dapi *DefaultCowAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cow collections.Cow
		err := json.NewDecoder(r.Body).Decode(&cow)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			fmt.Println("Error to error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := dapi.dt.Insert(dapi.collection, cow)
		if err != nil {
			fmt.Println("Error to error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cow.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jEncoder := json.NewEncoder(w)
		jEncoder.SetEscapeHTML(false)
		err = jEncoder.Encode(cow)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}
