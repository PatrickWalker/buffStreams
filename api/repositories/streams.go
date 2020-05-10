package repositories

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"strconv"
	"strings"

	helpers "github.com/PatrickWalker/buffStreams/helpers"
)

//SRInterface allows us to unit test the controller
type SRInterface interface {
	List(req ListRequest) (*StreamListResponse, error)
}

// StreamsRepo type handles access to the underlying data store for Stream Info
type StreamsRepo struct {
	Db       *sql.DB
	PageSize int32
}

//NewStreamsRepo creates a new instance of StreamsRepo
func NewStreamsRepo(cfg helpers.Config) SRInterface {
	//this should get a config object in to build this or just another function to do it because its not testable this way
	db, err := sql.Open("mysql", cfg.DB.ConnectionString)
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Fatal(err.Error())

	}
	return StreamsRepo{
		Db:       db,
		PageSize: cfg.PageSize,
	}

}

//Stream contains information pertaining to a single stream
type Stream struct {
	ID        int32      `json:"ID"`
	Title     string     `json:"title"`
	Questions []int32    `json:"questions"`
	Created   *time.Time `json:"created"`
	Updated   *time.Time `json:"updated"`
}

//StreamListResponse is the top level response object back from the List function
type StreamListResponse struct {
	Streams    []Stream `json:"streams"`
	PageNumber int32    `json:"pageNum"`
	PageSize   int32    `json:"pageSize"`
	TotalPages int32    `json:"totalPages"`
}

//ListRequest is used to specify any query parameters. For example page Number
type ListRequest struct {
	PageNumber int32
}

//List returns a page of stream results
func (sr StreamsRepo) List(req ListRequest) (*StreamListResponse, error) {
	resp := StreamListResponse{}
	if sr.Db == nil {
		return nil, errors.New("No DB Connection found cannot list streams")
	}
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	if sr.PageSize == 0 {
		sr.PageSize = 10
	}
	//get count of results so we can return that too
	countStmt := `SELECT count(id) as count from streams;`
	var count int32
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	countRow := sr.Db.QueryRow(countStmt)
	switch err := countRow.Scan(&count); err {
	case nil:
		break
	case sql.ErrNoRows:
		//We will encapsultate the error to not bleed context
		log.Println("Unable to fetch result count for streams. No Rows returned")
		return nil, errors.New("Unable to fetch streams")
	default:
		log.Printf("Error fetching stream count : %v \n", err)
		return nil, errors.New("Unable to fetch streams")
	}

	//Setting the meta data for the response
	resp.TotalPages = int32((count / sr.PageSize) + 1)
	resp.PageSize = sr.PageSize
	resp.PageNumber = req.PageNumber

	rows, err := sr.Db.Query(
		`SELECT s.id as id,title,created_at, updated_at, GROUP_CONCAT(question_id) as questions
		FROM buffup.streams s
		join question_stream qs
		ON s.id = qs.stream_id
		GROUP BY
			s.id
		limit ? offset ?;`, sr.PageSize, sr.PageSize*(req.PageNumber-1))
	if err != nil {
		//Returning this sets success to
		log.Printf("Error fetching stream list : %v \n", err)
		return nil, errors.New("Unable to fetch streams")
	}

	//loop through the rows, scan and then add them to the result object
	defer rows.Close()
	for rows.Next() {
		var questionString string

		s := Stream{}

		err := rows.Scan(&s.ID, &s.Title, &s.Created, &s.Updated, &questionString)
		if err != nil {
			log.Printf("Error fetching stream list : %v \n", err)
			return nil, errors.New("Unable to fetch streams")
		}
		s.Questions = convertQuestionList(questionString)
		resp.Streams = append(resp.Streams, s)
	}
	//row.Scan(&resp.InstallationID, &pwd)
	// check if the password is right

	return &resp, nil
}

//convertQuestionList takes a csv string of question IDs and returns a slice of Question Links
func convertQuestionList(csvQuestion string) []int32 {
	stArr := strings.Split(csvQuestion, ",")
	resp := make([]int32, len(stArr))
	for i, quest := range stArr {
		//these are auto geenrated sql ids so should always scan to ints
		qid, err := strconv.Atoi(strings.TrimSpace(quest))
		//if not something weird has happened so add a 0
		if err != nil {
			log.Printf("Error converting ID. Received : %v \n", quest)
			return []int32{}
		}
		resp[i] = int32(qid)
	}
	return resp
}
