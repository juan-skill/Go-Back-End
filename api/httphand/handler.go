package httphand

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/internal/storage"
)

// NewHandlerRequest ...
func NewHandlerRequest(store *storage.Store) *HandlerRequest {
	return &HandlerRequest{
		store: store,
	}
}

// HandlerRequest ...
type HandlerRequest struct {
	store *storage.Store
}

// RequestBody contain the information of body of the request
type RequestBody struct {
	DomainName string
}

// Create a new domain
func (p *HandlerRequest) Create(w http.ResponseWriter, r *http.Request) {
	domainName := parseRequest(r, w)

	ctx, cancelfunc := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancelfunc()

	_, err := p.store.ReloadRecord(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't reload the last domains")
		return
	}

	domain, err := ProcessData(ctx, domainName)
	if err != nil {
		respondWithError(w, http.StatusConflict, "can't create the domain")
		return
	}

	// reasignar el attributo Servers
	argPre := storage.TransferTxParamsServers{
		FromDomain: domain,
	}

	result1, err := p.store.TransferTxServers(r.Context(), argPre)
	if err != nil {
		respondWithError(w, http.StatusConflict, "error in create a server of the domain")
		return
	}

	nDomain := result1.FromDomain

	_, err = p.store.NewRecord(nDomain)
	if err != nil {
		respondWithError(w, http.StatusConflict, "error saving log domain")
		return
	}

	// reasignar el attributo previoGradeSSL
	argIni := storage.TransferTxParamsInitialize{
		FromDomain: nDomain,
	}

	result2, err := p.store.TransferTxInitialize(ctx, argIni)
	if err != nil {
		respondWithError(w, http.StatusConflict, "error in create a server of the domain")
		return
	}

	_, err = p.store.NewRecord(nDomain)
	if err != nil {
		respondWithError(w, http.StatusConflict, "error saving log domain")
		return
	}

	parseResponse := parseJSON(result2.ToDomain)

	respondwithJSON(w, http.StatusCreated, parseResponse)
}

// RequestLastDomains get the last records
func (p *HandlerRequest) RequestLastDomains(w http.ResponseWriter, r *http.Request) {
	ctx, cancelfunc := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancelfunc()

	_, err := p.store.ReloadRecord(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't reload the last domains")
		return
	}

	mapl := p.store.GetLastDomain()
	parseResponse := parseListJSON(mapl)

	respondwithJSON(w, http.StatusOK, parseResponse)
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logs.Log().Errorf("Error Marshal response ", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		logs.Log().Errorf("Error Write response %s", err.Error())
		return
	}
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}

// parseRequest extract the body of the request
func parseRequest(r *http.Request, w http.ResponseWriter) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't read body")
		return ""
	}

	var reqBody RequestBody

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		logs.Log().Errorf("Error Unmarshal request body: %s", err.Error())
		return ""
	}

	return reqBody.DomainName
}
