// +build e2e

package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/dionaditya/go-production-ready-api/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/dionaditya/go-production-ready-api/internal/user"
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

type ErrorMessage struct {
	Error   string
	Message string
}

func TestSignUp(t *testing.T) {
	fmt.Println("Running E2E for sign up user")

	client := resty.New()

	resp, err := client.R().SetBody(&models.User{
		Username: "test1234",
		Email:    "test122@gmail.com",
		Password: "semarang11",
	}).Post(BASE_URL + "/api/register")

	if err != nil {
		t.Fail()
	}

	var errorMessage ErrorMessage

	if err := json.Unmarshal([]byte(resp.Body()), &errorMessage); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestLogin(t *testing.T) {
	fmt.Println("Running E2E for login user")

	client := resty.New()

	resp, err := client.R().SetBody(&models.User{
		Email:    "test12@gmail.com",
		Password: "semarang11",
	}).Post(BASE_URL + "/api/login")

	if err != nil {
		fmt.Println("Error test login user")
	}

	var user user.Payload

	if err = json.Unmarshal([]byte(resp.Body()), &user); err != nil {
		fmt.Println(err)
	}

	os.Setenv("access_token", user.Access_Token)
	os.Setenv("refresh_token", user.Refresh_Token)

	assert.Equal(t, 200, resp.StatusCode())
}

func TestGetComments(t *testing.T) {
	fmt.Println("Running E2E test for get comments")

	client := resty.New()
	client.SetAuthToken(os.Getenv("access_token"))
	resp, err := client.R().Get(BASE_URL + "/api/comment")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestPostComment(t *testing.T) {
	fmt.Println("Running E2E for post new comments")

	client := resty.New()
	client.SetAuthToken(os.Getenv("access_token"))
	resp, err := client.R().SetBody(&Comment{
		Slug:   "/testing",
		Author: "john doe",
		Body:   "fix",
	}).Post(BASE_URL + "/api/comment")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())

}
