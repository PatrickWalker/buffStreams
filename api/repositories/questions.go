package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	helpers "github.com/PatrickWalker/buffStreams/helpers"
)

// QuestionsRepo type handles access to the underlying data store for Stream Info
type QuestionsRepo struct {
	Db *sql.DB
}

//NewQuestionsRepo creates a new instance of QuestionsRepo
func NewQuestionsRepo(cfg helpers.Config) QuestionsRepo {
	//this should get a config object in to build this or just another function to do it because its not testable this way
	db, err := sql.Open("mysql", cfg.DB.ConnectionString)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Fatal(err.Error())
	}
	return QuestionsRepo{
		Db: db,
	}

}

//StandardQuestion models a question with 1 correct answer and multiple wrong options
type StandardQuestion struct {
	Text    string   `json:"text"`
	Correct string   `json:"correct"`
	Options []string `json:"options"`
}

//Question models our question object
type Question struct {
	ID      int32      `json:"id"`
	Text    string     `json:"text"`
	Options []string   `json:"options"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
}

//MissingError is a custom error we will use to tell the controller to return a 404
//The repo shouldn't care about the transport so it shouldnt IMO be specifying a http status code here.
//That's not its job
type MissingError struct {
	Details string
}

func (e *MissingError) Error() string {
	return e.Details
}

//View is how we view a question
func (qr QuestionsRepo) View(questionID int32) (*Question, error) {

	resp := Question{}
	if qr.Db == nil {
		return nil, errors.New("No DB Connection found cannot find question")
	}
	//get the question based on the id supplied to us
	questStmt := `SELECT id, question_type, question, created_at,updated_at from questions where id = ?;`
	var questionJSON string
	var qType string
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	countRow := qr.Db.QueryRow(questStmt, questionID)
	switch err := countRow.Scan(&resp.ID, &qType, &questionJSON, &resp.Created, &resp.Updated); err {
	case nil:
		break
	case sql.ErrNoRows:
		//We will return a not found error type which can then be handled by controller to give
		//appropriate http status code
		return nil, &MissingError{Details: fmt.Sprintf("Question: %v not found", questionID)}
	default:
		log.Printf("Error fetching question : %v \n", err)
		return nil, errors.New("Unable to fetch question")
	}
	//go convert the JSON string and type to a struct and then convert that to an Option list
	resp.Text, resp.Options = parseQuestion(qType, questionJSON)

	return &resp, nil
}

//parseQuestion takes a type and json rep of  question. Parses them to a question type and returns options and
//question text
func parseQuestion(qType, questionJSON string) (string, []string) {
	switch strings.ToLower(qType) {
	case "standard":
		stdQues := StandardQuestion{}
		if err := json.Unmarshal([]byte(questionJSON), &stdQues); err != nil {
			log.Printf("Error parsing question %v \n err: %v \n", questionJSON, err)
			return "", []string{}
		}
		return stdQues.Text, append(stdQues.Options, stdQues.Correct)
	default:
		return "", []string{}
	}
}
