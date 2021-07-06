package handlers

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DATE_FORMAT = "%a, %d %b %Y %H:%M:%S GMT"

	STATUS_OK               = "OK"
	STATUS_ERR              = "ERR"
	LAST_UPDATED            = "_updated"
	DATE_CREATED            = "_created"
	ISSUES                  = "_issues"
	STATUS                  = "_status"
	ERROR                   = "_error"
	ITEMS                   = "_items"
	LINKS                   = "_links"
	ETAG                    = "_etag"
	VERSION                 = "_version" // field that stores the version number
	DELETED                 = "_deleted" // field to store soft delete status
	META                    = "_meta"
	INFO                    = ""
	VALIDATION_ERROR_STATUS = 422
	NORMALIZE_DOTTED_FIELDS = true

	// return a single field validation error as a list (by default a single error
	// is retuned as string, while multiple errors are returned as a list).
	VALIDATION_ERROR_AS_LIST = false

	// codes for which we want to return a standard response which includes
	// a JSON body with the status, code, and description.
	STANDARD_ERRORS = []int{400, 401, 403, 404, 405, 406, 409, 410, 412, 422, 428, 429}

	// field returned on GET requests so we know if we have the latest copy even if
	// we access a specific version
	LATEST_VERSION = "_latest_version"

	// appended to ID_FIELD, holds the original document id in parallel collection
	VERSION_ID_SUFFIX    = "_document"
	VERSION_DIFF_INCLUDE = []string{} // always include these fields when diffing

	API_VERSION         = ""
	URL_PREFIX          = ""
	ID_FIELD            = "_id"
	CACHE_CONTROL       = ""
	CACHE_EXPIRES       = 0
	ITEM_CACHE_CONTROL  = ""
	X_DOMAINS           = ""    // CORS disabled by default.
	X_DOMAINS_RE        = ""    // CORS disabled by default.
	X_HEADERS           = ""    // CORS disabled by default.
	X_EXPOSE_HEADERS    = ""    // CORS disabled by default.
	X_ALLOW_CREDENTIALS = ""    // CORS disabled by default.
	X_MAX_AGE           = 21600 // Access-Control-Max-Age when CORS is enabled
	HATEOAS             = true  // HATEOAS enabled by default.
	IF_MATCH            = true  // IF_MATCH (ETag match) enabled by default.
	ENFORCE_IF_MATCH    = true  // ENFORCE_IF_MATCH enabled by default.

	ALLOWED_FILTERS    = []string{"*"} // filtering enabled by default
	VALIDATE_FILTERS   = false
	SORTING            = true  // sorting enabled by default.
	JSON_SORT_KEYS     = false // json key sorting
	RENDERERS          = []string{"eve.render.JSONRenderer", "eve.render.XMLRenderer"}
	EMBEDDING          = true // embedding enabled by default
	PROJECTION         = true // projection enabled by default
	PAGINATION         = true // pagination enabled by default.
	PAGINATION_LIMIT   = 50
	PAGINATION_DEFAULT = 25
	VERSIONING         = false       // turn document versioning on or off.
	VERSIONS           = "_versions" // suffix for parallel collection w/old versions
	VERSION_PARAM      = "version"   // URL param for specific version of a document.
	INTERNAL_RESOURCE  = false       // resources are public by default.
	JSONP_ARGUMENT     = ""          // JSONP disabled by default.
	SOFT_DELETE        = false       // soft delete disabled by default.
	SHOW_DELETED_PARAM = "show_deleted"
	BULK_ENABLED       = true

	OPLOG          = false   // oplog is disabled by default.
	OPLOG_NAME     = "oplog" // default oplog resource name.
	OPLOG_ENDPOINT = ""      // oplog endpoint is disabled by default.
	OPLOG_AUDIT    = true    // oplog audit enabled by default.
	OPLOG_METHODS  = []string{
		"DELETE",
		"POST",
		"PATCH",
		"PUT",
	} // oplog logs all operations by default.
	OPLOG_CHANGE_METHODS = []string{
		"DELETE",
		"PATCH",
		"PUT",
	} // methods which write changes to the oplog
	OPLOG_RETURN_EXTRA_FIELD = false // oplog does not return the 'extra' field.

	RESOURCE_METHODS         = []string{"GET"}
	ITEM_METHODS             = []string{"GET"}
	PUBLIC_METHODS           = []string{}
	ALLOWED_ROLES            = []string{}
	ALLOWED_READ_ROLES       = []string{}
	ALLOWED_WRITE_ROLES      = []string{}
	PUBLIC_ITEM_METHODS      = []string{}
	ALLOWED_ITEM_ROLES       = []string{}
	ALLOWED_ITEM_READ_ROLES  = []string{}
	ALLOWED_ITEM_WRITE_ROLES = []string{}
	// globally enables / disables HTTP method overriding
	ALLOW_OVERRIDE_HTTP_METHOD = true
	ITEM_LOOKUP                = true
	ITEM_LOOKUP_FIELD          = ID_FIELD
	ITEM_URL                   = "[a-f0-9]{24}"
	UPSERT_ON_PUT              = true // insert unexisting documents on PUT.
	MERGE_NESTED_DOCUMENTS     = true

	// use a simple file response format by default
	EXTENDED_MEDIA_INFO           = []string{}
	RETURN_MEDIA_AS_BASE64_STRING = true
	RETURN_MEDIA_AS_URL           = false
	MEDIA_ENDPOINT                = "media"
	MEDIA_URL                     = "[a-f0-9]{24}"
	MEDIA_BASE_URL                = ""

	MULTIPART_FORM_FIELDS_AS_JSON = false
	AUTO_COLLAPSE_MULTI_KEYS      = false
	AUTO_CREATE_LISTS             = false
	JSON_REQUEST_CONTENT_TYPES    = []string{"application/json"}

	SCHEMA_ENDPOINT = ""

	// list of extra fields to be included with every POST response. This list
	// should not include the 'standard' fields (ID_FIELD, LAST_UPDATED,
	// DATE_CREATED, and ETAG). Only relevant when bandwidth saving mode is on.
	EXTRA_RESPONSE_FIELDS = []string{}
	BANDWIDTH_SAVER       = true

	// default query parameters
	QUERY_WHERE       = "where"
	QUERY_PROJECTION  = "projection"
	QUERY_SORT        = "sort"
	QUERY_PAGE        = "page"
	QUERY_MAX_RESULTS = "max_results"
	QUERY_EMBEDDED    = "embedded"
	QUERY_AGGREGATION = "aggregate"

	HEADER_TOTAL_COUNT            = "X-Total-Count"
	OPTIMIZE_PAGINATION_FOR_SPEED = false

	// user-restricted resource access is disabled by default.
	AUTH_FIELD = ""

	// don't allow unknown key/value pairs for POST/PATCH payloads.
	ALLOW_UNKNOWN = false

	// GeoJSON specs allows any number of key/value pairs
	// http://geojson.org/geojson-spec.html#geojson-objects
	ALLOW_CUSTOM_FIELDS_IN_GEOJSON = false

	// Rate limits are disabled by default. Needs a running redis-server.
	RATE_LIMIT_GET    = ""
	RATE_LIMIT_POST   = ""
	RATE_LIMIT_PATCH  = ""
	RATE_LIMIT_DELETE = ""

	// disallow Mongo's javascript queries as they might be vulnerable to injection
	// attacks ('ReDoS' especially), are probably too complex for the average API
	// end-user and finally can  seriously impact overall performance.
	MONGO_QUERY_BLACKLIST = []string{"$where", "$regex"}
	MONGO_QUERY_WHITELIST = []string{}
	// Explicitly set default write_concern to 'safe' (do regular
	// aknowledged writes). This is also the current PyMongo/Mongo default setting.
	MONGO_WRITE_CONCERN = map[string]int{"w": 1}
	MONGO_OPTIONS       = map[string]bool{"connect": true, "tz_aware": true}

	// if true, the document will be normalized according to the schema during patch
	// this means fields will be reset their the default value, if any, unless
	// contained in the patch body.
	NORMALIZE_ON_PATCH = true
)

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
	// meta.Total = 20
	meta.Page = 1
	return &meta
}

type RequestParameters struct {
	MaxResults int
	Where      string
	Page       int
}

func NewRequestParameters(values url.Values) RequestParameters {
	req := RequestParameters{}
	max_results := values.Get(QUERY_MAX_RESULTS)
	if max_results != "" {
		i, err := strconv.Atoi(max_results)
		if err != nil {
			log.Debug("Cannot convert %s to int", QUERY_MAX_RESULTS)
			req.MaxResults = PAGINATION_DEFAULT
		} else {
			req.MaxResults = i
		}
	} else {
		req.MaxResults = PAGINATION_DEFAULT
	}
	where := values.Get(QUERY_WHERE)
	if where != "" {
		req.Where = where
	}
	page := values.Get(QUERY_PAGE)
	if page != "" {
		i, err := strconv.Atoi(page)
		if err != nil {
			log.Debug("Cannot convert %s to int", QUERY_PAGE)
			req.Page = 1
		} else {
			req.Page = i
		}
	} else {
		req.Page = 1
	}
	return req
}

func (req *RequestParameters) WhereClause() interface{} {
	fmt.Println("where", req.Where)
	var doc interface{}
	if req.Where != "" {
		err := bson.UnmarshalExtJSON([]byte(req.Where), true, &doc)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		doc = bson.M{}
	}
	fmt.Println("where compiled", doc)
	return doc
}

func (req *RequestParameters) RequestParameters2MongOptions() *options.FindOptions {
	findOptions := options.FindOptions{}
	m := int64(req.MaxResults)
	findOptions.Limit = &m
	return &findOptions
}
