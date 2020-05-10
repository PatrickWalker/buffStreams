package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	helpers "github.com/buffup/api/helpers"
	repo "github.com/buffup/api/repositories"

	"github.com/julienschmidt/httprouter"

	"net/http"
)

// QuestionsController struct handles all routes/handlers for Stream access
type QuestionsController struct {
	QRepo repo.QuestionsRepo
}

// NewQuestionsController returns a QuestionsController for use by main handler
func NewQuestionsController(conf helpers.Config) QuestionsController {
	return QuestionsController{
		QRepo: repo.NewQuestionsRepo(conf),
	}
}

// View will return details fo a single question
func (qc QuestionsController) View(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qString := ps.ByName("questionID")
	qID, err := strconv.Atoi(qString)
	if err != nil {
		//HTTP Status Code 400 Bad Request
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"Invalid Question ID Supplied. Received : %v"}`, qString)
		return
	}
	result, err := qc.QRepo.View(int32(qID))
	if err != nil {
		if err, ok := err.(*repo.MissingError); ok {
			//we will 404 here
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
			return
		}
		//HTTP Status Code 500 Internal Server Error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	return
}
