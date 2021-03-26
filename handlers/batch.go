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
	BATCH = "batch"
)

type DefaultBatchAPI struct {
	*config.WebConfig
	dt *dao.DataMongo
	*metrics.Metrics
	collection string
}

func (dapi *DefaultBatchAPI) InitConfig(w *config.WebConfig, dt *dao.DataMongo) API {
	dapi.WebConfig = w
	dapi.dt = dt
	dapi.Metrics = w.Metrics
	dapi.collection = BATCH
	return dapi
}

func (dapi *DefaultBatchAPI) GetUrl() string {
	return fmt.Sprintf("/%s", dapi.collection)
}

func (dapi *DefaultBatchAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JSON body and decode into credentials
		var data []collections.Batch
		err := dapi.dt.FindAll(dapi.collection, &data)
		if err != nil {
			fmt.Println("Error to error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		WriteJSONResponse(data, 200, w)
	})
}

func (dapi *DefaultBatchAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var batch collections.Batch
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
