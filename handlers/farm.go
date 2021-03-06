package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gutorc92/api-farm/dao"
	"github.com/gutorc92/api-farm/config"
	"github.com/gutorc92/api-farm/metrics"
	"github.com/gutorc92/api-farm/collections"
)

// DummyAPI defines the available dummy apis
type FarmAPI interface {
	InitConfig(w *config.WebConfig, dt *dao.DataMongo) FarmAPI
	GETHandler() http.Handler
	POSTHandler() http.Handler
}

// DefaultDummyAPI holds the default implementation of the Dummy API interface
type DefaultFarmAPI struct {
	*config.WebConfig
	dt *dao.DataMongo
	*metrics.Metrics
}

func (dapi *DefaultFarmAPI) InitConfig(w *config.WebConfig, dt *dao.DataMongo) FarmAPI {
	dapi.WebConfig = w
	dapi.dt = dt
	dapi.Metrics = w.Metrics
	return dapi
}

func (dapi *DefaultFarmAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JSON body and decode into credentials
		farms, err := dapi.dt.FindFarm()
		if err != nil {
			fmt.Println("Error to error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jEncoder := json.NewEncoder(w)
		jEncoder.SetEscapeHTML(false)
		err = jEncoder.Encode(farms)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}

func (dapi *DefaultFarmAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var farm collections.Farm
		err := json.NewDecoder(r.Body).Decode(&farm)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := dapi.dt.InsertFarm(farm)
		if err != nil {
			fmt.Println("Error to error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		farm.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jEncoder := json.NewEncoder(w)
		jEncoder.SetEscapeHTML(false)
		err = jEncoder.Encode(farm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}