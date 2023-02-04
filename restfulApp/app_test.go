package restfulapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	//mocked WebServer
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	//mocked WebServer
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(string(data), "No Users")
}

func TestUsers_WithUsersData(t *testing.T) {
	assert := assert.New(t)

	//mocked WebServer
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(
			`{"first_name": "first_name1","last_name": "last_name1","email":"email1@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(
			`{"first_name": "first_name2","last_name": "last_name2","email":"email2@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(
			`{"first_name": "first_name3","last_name": "last_name3","email":"email3@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(3, len(users))
}

func TestGetNotExsistUserInfo(t *testing.T) {
	assert := assert.New(t)

	//mocked WebServer
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/12")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:12")
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	//mocked WebServer
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(
			`{"first_name": "seungyeon","last_name": "hwang","email":"hwang@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	//User1(Created)
	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	//User2
	id := user.ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	user2 := new(User)
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)
	assert.Equal(user.FirstName, user2.FirstName)
}

// Delete User
func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	//mocked WebServer
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	//Usermap[Empty]
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:1")

	//Usermap[1]
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(
			`{"first_name": "seungyeon","last_name": "hwang","email":"hwang@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(string(data), "Deleted User ID:1")
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	//mocked WebServer
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	//Usermap[Empty]
	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(
			`{"id":1, "first_name": "updated","last_name": "updated","email":"updated@gmail.com"}`))
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:1")

	//Usermap[1]
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(
			`{"first_name": "hwang","last_name": "seungyeon","email":"hwang@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	updateStr := fmt.Sprintf(`{"id":%d, "first_name": "json"}`, 1)

	req, _ = http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)
	assert.Equal(updateUser.ID, user.ID)
	assert.Equal("json", updateUser.FirstName)
	assert.Equal(updateUser.LastName, updateUser.LastName)
	assert.Equal(user.Email, updateUser.Email)
}
