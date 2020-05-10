package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	repo "github.com/PatrickWalker/buffStreams/repositories"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

//MockSRepoError so we can return errors and test controller
type MockSRepo struct {
	Resp *repo.StreamListResponse
	Err  error
}

func (mr MockSRepo) List(req repo.ListRequest) (*repo.StreamListResponse, error) {
	return mr.Resp, mr.Err
}

//ListTest is used to test the List function on the controller
type ListTest struct {
	Name         string
	Resp         *repo.StreamListResponse
	Err          error
	ExpectedBody string
	ExpectedCode int
}

var listTests = []ListTest{
	ListTest{
		Name:         "Throws Error",
		Resp:         nil,
		Err:          errors.New("An error"),
		ExpectedBody: `{"error": "An error"}`,
		ExpectedCode: 500,
	},
}

func TestList(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	for _, lt := range listTests {
		sc := StreamsController{
			MockSRepo{
				Resp: lt.Resp,
				Err:  lt.Err,
			},
		}
		sc.List(rr, req, httprouter.Params{})
		assert.Equal(t, lt.ExpectedBody, rr.Body.String(), "Unexpected Body")
		assert.Equal(t, lt.ExpectedCode, rr.Code, "Unexpected Code")
	}
}
