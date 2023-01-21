package shortener

import (
	"fmt"
	"net/http"

	clarkezoneLog "github.com/clarkezone/pocketshorten/pkg/log"
)

type storeLoader interface {
	Init(urlLookupService)
}

type urlLookupService interface {
	Store(string, string) error
	Lookup(string) (string, error)
	Count() int
}

// NewDictLookupHandler creates a new instance of type
//
//lint:ignore U1000 reason backend not selected
func NewDictLookupHandler() *ShortenHandler {
	vl := &viperLoader{}
	ds := newDictStore(vl)
	lh := &ShortenHandler{storage: ds}
	vl.Init(ds)
	return lh
}

// NewGrpcLookupHandler returns a new lookuphandler instance
func NewGrpcLookupHandler(s string) (*ShortenHandler, error) {
	// dictstore
	// grpcloader
	ds, err := newGrpcStore(s)
	if err != nil {
		return nil, err
	}
	lh := &ShortenHandler{storage: ds}
	return lh, nil
}

// ShortenHandler core logic
type ShortenHandler struct {
	storage urlLookupService
}

// RegisterHandlers attaches handlers to Mux that is passed in
func (lh *ShortenHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", lh.redirectHandler)
}

func (lh *ShortenHandler) redirectHandler(w http.ResponseWriter, r *http.Request) {
	//TODO telemetry
	requested := r.URL.Query().Get("shortlink")

	if requested == "" {
		clarkezoneLog.Errorf("redirecthandler called with  missingshortlink")
		writeOutputError(w, "please supply shortlink query parameter")
		return
	}
	uri, err := lh.storage.Lookup(requested)
	if err != nil {
		writeOutputError(w, fmt.Sprintf("shortlink %v notfound", requested))
		return
	}
	clarkezoneLog.Debugf("redirecting to %v", uri)

	http.Redirect(w, r, uri, http.StatusMovedPermanently)

}

func writeOutputError(w http.ResponseWriter, message string) {
	clarkezoneLog.Debugf("Error reported to user %v", message)
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(message))
	if err != nil {
		clarkezoneLog.Debugf("writeOutputError: Failed to write bytes %v\n", err)
	}
}
