package httphand

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/other_project/crockroach/internal/logs"
	db "github.com/other_project/crockroach/internal/storage"
	"github.com/other_project/crockroach/models"
)

// NewHandlerRequest ...
func NewHandlerRequest(store *db.Store) *HandlerRequest {
	return &HandlerRequest{
		store: store,
	}
}

// HandlerRequest ...
type HandlerRequest struct {
	store *db.Store
}

// RequestBody contain the information of body of the request
type RequestBody struct {
	DomainName string
}

// Create a new domain
func (p *HandlerRequest) Create(w http.ResponseWriter, r *http.Request) {
	cmd := parseRequest(r, w)

	domain, err := p.store.StoreDomain(r.Context(), cmd)
	if err != nil {
		respondWithError(w, http.StatusNoContent, "error in create domain")
		return
	}

	respondwithJSON(w, http.StatusCreated, domain)
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logs.Log().Errorf("Error Marshal response ", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		logs.Log().Errorf("Error Write response ", err.Error())
	}
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}

func parseRequest(r *http.Request, w http.ResponseWriter) *models.Domain {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't read body")
		print(err)
	}

	var reqBody RequestBody

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		log.Println(err)
	}

	domain, err := ProcessData(reqBody.DomainName)
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "can't create the domain")
	}

	return domain
}
