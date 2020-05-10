package repositories

import (
	"errors"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

//ConvertTest is used to test the ConvertQuestionList
type ConvertTest struct {
	Name              string
	CSVQuestions      string
	ExpectedQuestions []int32
}

var convertTests = []ConvertTest{
	ConvertTest{
		Name:              "Empty string",
		CSVQuestions:      "",
		ExpectedQuestions: []int32{},
	},
	ConvertTest{
		Name:              "Valid single",
		CSVQuestions:      "1",
		ExpectedQuestions: []int32{1},
	},
	ConvertTest{
		Name:              "Valid multiple",
		CSVQuestions:      "1  ,2  ,   3",
		ExpectedQuestions: []int32{1, 2, 3},
	},
	ConvertTest{
		Name:              "Invalid string",
		CSVQuestions:      "anc",
		ExpectedQuestions: []int32{},
	},
}

func TestConvertQuestionList(t *testing.T) {
	for _, ct := range convertTests {
		quest := convertQuestionList(ct.CSVQuestions)
		assert.Equal(t, ct.ExpectedQuestions, quest, fmt.Sprintf("Unexpected result for test: %v result %v\n expected: %v", ct.Name, quest, ct.ExpectedQuestions))

	}
}

func TestListStreamsNilDB(t *testing.T) {
	//Arrange
	assert := assert.New(t)
	sr := StreamsRepo{}

	//Act
	_, err := sr.List(ListRequest{PageNumber: 1})

	//Assert
	assert.Equal(errors.New("No DB Connection found cannot list streams"), err, fmt.Sprintf("error did not match expected: %s", err))
}

func TestStreamCountOtherError(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	lr := ListRequest{PageNumber: 1}

	mock.ExpectQuery("SELECT count").WillReturnError(fmt.Errorf("some error"))
	sr := StreamsRepo{Db: db}

	//Act
	_, err = sr.List(lr)

	//Assert
	if err == nil {
		t.Errorf("error was expected while viewing Streams: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStreamCountNoRows(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	lr := ListRequest{PageNumber: 1}
	cRows := sqlmock.NewRows([]string{"count"})
	mock.ExpectQuery("SELECT count").WillReturnRows(cRows)
	sr := StreamsRepo{Db: db}

	//Act
	_, err = sr.List(lr)

	//Assert
	if err == nil {
		t.Errorf("error was expected while viewing Streams: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStreamSelectError(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	lr := ListRequest{PageNumber: 1}
	cRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT count").WillReturnRows(cRows)
	mock.ExpectQuery("SELECT s.id as id,title,created_at, updated_at, GROUP_CONCAT").WillReturnError(fmt.Errorf("stream error"))
	sr := StreamsRepo{Db: db}

	//Act
	_, err = sr.List(lr)

	//Assert
	if err == nil {
		t.Errorf("error was expected while viewing Streams: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStreamScanError(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	lr := ListRequest{PageNumber: 1}
	cRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT count").WillReturnRows(cRows)
	t1 := time.Now()
	sRows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at", "questions"}).
		AddRow("abc", "title", &t1, &t1, "")
	mock.ExpectQuery("SELECT s.id as id,title,created_at, updated_at, GROUP_CONCAT").WillReturnRows(sRows)
	sr := StreamsRepo{Db: db}

	//Act
	_, err = sr.List(lr)

	//Assert
	if err == nil {
		t.Errorf("error was expected while viewing Streams: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStreamValid(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	lr := ListRequest{PageNumber: 1}
	cRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT count").WillReturnRows(cRows)
	t1 := time.Now()
	sRows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at", "questions"}).
		AddRow(1, "title", &t1, &t1, "")
	mock.ExpectQuery("SELECT s.id as id,title,created_at, updated_at, GROUP_CONCAT").WillReturnRows(sRows)
	sr := StreamsRepo{Db: db}

	//Act
	res, err := sr.List(lr)

	//Assert
	if err != nil {
		t.Errorf("error was not expected while viewing Streams: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	expected := StreamListResponse{
		Streams: []Stream{
			Stream{
				ID:        1,
				Title:     "title",
				Questions: []int32{},
				Created:   &t1,
				Updated:   &t1,
			},
		},
		PageNumber: 1,
		PageSize:   10,
		TotalPages: 1,
	}
	assert.Equal(t, &expected, res, "Unexpected Result")
}

func TestStreamDefaultPage(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	lr := ListRequest{PageNumber: 0}
	cRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT count").WillReturnRows(cRows)
	t1 := time.Now()
	sRows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at", "questions"}).
		AddRow(1, "title", &t1, &t1, "")
	mock.ExpectQuery("SELECT s.id as id,title,created_at, updated_at, GROUP_CONCAT").WillReturnRows(sRows)
	sr := StreamsRepo{Db: db}

	//Act
	res, err := sr.List(lr)

	//Assert
	if err != nil {
		t.Errorf("error was not expected while viewing Streams: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equal(t, int32(1), res.PageNumber, "Unexpected Page Number")
}
