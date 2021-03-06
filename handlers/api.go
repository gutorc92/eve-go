package handlers

import (
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
	Total int `json:"total"`
	Page int `json:"page"`
	MaxResults int `json:"max_results"`
}

func newMeta(data []interface{}) *MetaPage {
	var meta MetaPage
	meta.Total = 20
	return &meta
}

type ResultPage struct {
	Items []interface{}  `json:"_items"`
	Meta 	*MetaPage			 `json:"_meta"`
}

