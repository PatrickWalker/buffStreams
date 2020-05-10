package controllers

import (
	"encoding/json"
	"fmt"

	helpers "github.com/PatrickWalker/buffStreams/helpers"
	repo "github.com/PatrickWalker/buffStreams/repositories"

	"github.com/julienschmidt/httprouter"

	"net/http"

	"strconv"
)

// StreamsController struct handles all routes/handlers for Stream access
type StreamsController struct {
	SRepo repo.SRInterface
}

// NewStreamsController returns a StreamsController for use by main handler
func NewStreamsController(conf helpers.Config) StreamsController {
	return StreamsController{
		SRepo: repo.NewStreamsRepo(conf),
	}
}

// List will list all streams
func (sc StreamsController) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//We will get our pagination info from the request. They will be query parameters
	pageParam := r.URL.Query().Get("page")
	var err error
	pageNum := 1 //default to page 1
	if pageParam != "" {
		pageNum, err = strconv.Atoi(pageParam)
		if err != nil {
			//HTTP Status Code 400 Bad Request
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"Invalid Page Number Supplied. Received : %v"}`, pageParam)
			return
		}
	}
	result, err := sc.SRepo.List(repo.ListRequest{
		PageNumber: int32(pageNum),
	})
	if err != nil {
		//HTTP Status Code 500 Internal Server Error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	//return a 200
}
