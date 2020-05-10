package repositories

import (
	"errors"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

//ParseTest is used to test the ParseQuestionFunction
type ParseTest struct {
	Name            string
	Type            string
	JSON            string
	ExpectedText    string
	ExpectedOptions []string
}

var parseTests = []ParseTest{
	ParseTest{
		Name: "Valid",
		Type: "standard",
		JSON: `{"text": "Dennis Bergkamp turned 51 today. Who was Arsenal manager when he signed?",
		"correct": "Bruce Rioch",
		"options":[ "Arsene Wenger","George Graham"]
		}`,
		ExpectedText:    "Dennis Bergkamp turned 51 today. Who was Arsenal manager when he signed?",
		ExpectedOptions: []string{"Arsene Wenger", "George Graham", "Bruce Rioch"},
	},
	ParseTest{
		Name: "Case Insensitive",
		Type: "STandard",
		JSON: `{"text": "Dennis Bergkamp turned 51 today. Who was Arsenal manager when he signed?",
		"correct": "Bruce Rioch",
		"options":[ "Arsene Wenger","George Graham"]
		}`,
		ExpectedText:    "Dennis Bergkamp turned 51 today. Who was Arsenal manager when he signed?",
		ExpectedOptions: []string{"Arsene Wenger", "George Graham", "Bruce Rioch"},
	},
	ParseTest{
		Name:            "Invalid JSON",
		Type:            "STandard",
		JSON:            `{"text}`,
		ExpectedText:    "",
		ExpectedOptions: []string{},
	},
	ParseTest{
		Name:            "Invalid Type",
		Type:            "Wrong",
		JSON:            `{"text}`,
		ExpectedText:    "",
		ExpectedOptions: []string{},
	},
}

func TestParseQuestion(t *testing.T) {
	for _, pt := range parseTests {
		text, opt := parseQuestion(pt.Type, pt.JSON)
		assert.Equal(t, pt.ExpectedText, text, fmt.Sprintf("Unexpected result for test: %v result %v\n expected: %v", pt.Name, text, pt.ExpectedText))
		assert.Equal(t, pt.ExpectedOptions, opt, fmt.Sprintf("Unexpected result for test: %v result %v\n expected: %v", pt.Name, opt, pt.ExpectedOptions))

	}
}

func TestViewQuestionNilDB(t *testing.T) {
	//Arrange
	assert := assert.New(t)
	qr := QuestionsRepo{}

	//Act
	_, err := qr.View(1)

	//Assert
	assert.Equal(errors.New("No DB Connection found cannot find question"), err, fmt.Sprintf("error did not match expected: %s", err))
}

func Test404Question(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	reqID := 404
	rows := sqlmock.NewRows([]string{"id", "question_type", "question", "created_at", "updated_at"})
	mock.ExpectQuery(`SELECT id, question_type, question, created_at,updated_at from questions where id = ?`).WillReturnRows(rows)
	qr := QuestionsRepo{Db: db}

	//Act
	_, err = qr.View(int32(reqID))

	//Assert
	if err == nil {
		t.Errorf("error was expected while viewing Question: %s", err)
	}

	if err, ok := err.(*MissingError); !ok {
		t.Errorf("MissingError was expected while viewing Question got : %s", err)

	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOtherError(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	reqID := 500

	mock.ExpectQuery(`SELECT id, question_type, question, created_at,updated_at from questions where id = ?`).WillReturnError(fmt.Errorf("some error"))
	qr := QuestionsRepo{Db: db}

	//Act
	_, err = qr.View(int32(reqID))

	//Assert
	if err == nil {
		t.Errorf("error was expected while viewing Question: %s", err)
	}

	if err, ok := err.(*MissingError); ok {
		t.Errorf("MissingError was not expected while viewing Question got : %s", err)

	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidQuestion(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	reqID := 1
	t1 := time.Now()
	rows := sqlmock.NewRows([]string{"id", "question_type", "question", "created_at", "updated_at"}).
		AddRow(1, "invalid", "{}", &t1, &t1)
	mock.ExpectQuery(`SELECT id, question_type, question, created_at,updated_at from questions where id = ?`).WillReturnRows(rows)
	qr := QuestionsRepo{Db: db}

	//Act
	question, err := qr.View(int32(reqID))

	//Assert
	if err != nil {
		t.Errorf("error was not expected while viewing Question: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	expectQuest := Question{
		ID:      1,
		Text:    "",
		Options: []string{},
		Created: &t1,
		Updated: &t1,
	}
	assert.Equal(t, &expectQuest, question, fmt.Sprintf("Unexpected Question"))

}

func TestMissingError(t *testing.T) {
	//Arrange
	me := &MissingError{Details: "test"}
	expected := "test"
	//Assert

	// we make sure that all expectations were met

	assert.Equal(t, me.Error(), expected, fmt.Sprintf("Unexpected error"))

}
