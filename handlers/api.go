package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gutorc92/api-farm/config"
	"github.com/gutorc92/api-farm/dao"
)

type API interface {
	InitConfig(w *config.WebConfig, dt *dao.DataMongo) API
	GetUrl() string
	GETHandler() http.Handler
	POSTHandler() http.Handler
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

type ResultPage struct {
	Items interface{} `json:"_items"`
	Meta  *MetaPage   `json:"_meta"`
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
