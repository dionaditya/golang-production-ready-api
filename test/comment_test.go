// +build e2e

package test

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type Comments struct {
	Result []Comment `json:"Result"`
}

// User struct which contains a name
// a type and a list of social links
type Comment struct {
	Slug    string `json:"Slug"`
	Body    string `json:"Body"`
	Author  string `json:"Author"`
	Created string `json:"Created"`
}

func TestGetComments(t *testing.T) {
	fmt.Println("Running E2E test for get comments")

	client := resty.New()

	resp, err := client.R().Get(BASE_URL + "/api/comment")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestPostComment(t *testing.T) {
	fmt.Println("Running E2E for post new comments")

	client := resty.New()

	resp, err := client.R().SetBody(&Comment{
		Slug:   "/testing",
		Author: "john doe",
		Body:   "Fix",
	}).
		Post(BASE_URL + "/api/comment")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())

}
