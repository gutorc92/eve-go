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
	FARM = "farm"
)

// DefaultDummyAPI holds the default implementation of the Dummy API interface
type DefaultFarmAPI struct {
	*config.WebConfig
	dt *dao.DataMongo
	*metrics.Metrics
	collection string
}

func (dapi *DefaultFarmAPI) InitConfig(w *config.WebConfig, dt *dao.DataMongo) API {
	dapi.WebConfig = w
	dapi.dt = dt
	dapi.Metrics = w.Metrics
	dapi.collection = FARM
	return dapi
}

func (dapi *DefaultFarmAPI) GetUrl() string {
	return fmt.Sprintf("/%s", dapi.collection)
}

func (dapi *DefaultFarmAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []collections.Farm
		req := NewRequestParameters(r.URL.Query())
		err := dapi.dt.FindAll(dapi.collection, &data, req.RequestParameters2MongOptions())
		if err != nil {
			fmt.Println("Error to error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		WriteJSONResponse(data, 200, w)
	})
}

func (dapi *DefaultFarmAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var batch collections.Farm
		err := json.NewDecoder(r.Body).Decode(&batch)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := dapi.dt.Insert(dapi.collection, batch)
		if err != nil {
			fmt.Println("Error to error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		batch.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jEncoder := json.NewEncoder(w)
		jEncoder.SetEscapeHTML(false)
		err = jEncoder.Encode(batch)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}
