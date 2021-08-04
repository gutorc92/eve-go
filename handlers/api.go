package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

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
		meta := newMeta()
		meta.Total = reflect.Indirect(result).Len()
		page := ResultPage{Items: result.Interface(), Meta: meta}
		WriteJSONResponse(page, 200, w)
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
		page := CreatePage{}
		etag := reflect.Indirect(x).FieldByName("Etag").String()
		id := reflect.Indirect(x).FieldByName("ID")
		created := reflect.Indirect(x).FieldByName("CreatedAt")
		updated := reflect.Indirect(x).FieldByName("UpdateAt")
		createdTime := created.Addr().MethodByName("Time").Call([]reflect.Value{})[0]
		updateTime := updated.Addr().MethodByName("Time").Call([]reflect.Value{})[0]
		createdString := createdTime.Interface().(time.Time)
		updateString := updateTime.Interface().(time.Time)
		fmt.Println("Time formated", createdString.Format(time.RFC1123))
		idString := id.Addr().MethodByName("Hex").Call([]reflect.Value{})[0].String()
		links := SelfResult{Title: domain.GetUrl(), Href: domain.GetUrlSelfItem(idString)}
		page.Links = LinksResult{Self: links}
		page.Etag = etag
		page.Status = STATUS_OK
		page.ID = idString
		page.Created = createdString.Format(time.RFC1123)
		page.Updated = updateString.Format(time.RFC1123)
		WriteJSONResponse(page, 200, w)
	})
}

func WriteJSONResponse(page interface{}, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	jEncoder := json.NewEncoder(w)
	var err error
	err = jEncoder.Encode(page)
	jEncoder.SetEscapeHTML(false)
	if err != nil {
		fmt.Println("Error")
	}
}
