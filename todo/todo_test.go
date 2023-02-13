package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_addTodoListHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()

	//addTodoListHandler
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data1"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	var todo Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Id, 1)
	assert.Equal(todo.Name, "Test Data1")
	assert.Equal(reflect.ValueOf(todo.CreatedAt).String(), "<time.Time Value>")

	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Id, 2)
	assert.Equal(todo.Name, "Test Data2")
	assert.Equal(reflect.ValueOf(todo.CreatedAt).String(), "<time.Time Value>")
}

func Test_getTodoListHandlerTest(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()

	var todo Todo
	resp, _ := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data1"}})
	json.NewDecoder(resp.Body).Decode(&todo)
	id1 := todo.Id
	resp, _ = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data2"}})
	json.NewDecoder(resp.Body).Decode(&todo)
	id2 := todo.Id

	resp, err := http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.Id == id1 {
			assert.Equal("Test Data1", t.Name)
		} else if t.Id == id2 {
			assert.Equal("Test Data2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testId should be id1 or id2"))
		}
	}
}

func Test_completeTodoListHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()

	var todo Todo
	resp, _ := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data1"}})
	json.NewDecoder(resp.Body).Decode(&todo)
	id1 := todo.Id
	resp, _ = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data2"}})
	json.NewDecoder(resp.Body).Decode(&todo)

	resp, err := http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)

	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.Id == id1 {
			assert.True(t.Completed)
		}
	}
}
func Test_removeTodoListHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()

	var todo Todo
	resp, _ := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data1"}})
	json.NewDecoder(resp.Body).Decode(&todo)
	id1 := todo.Id
	resp, _ = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data2"}})
	json.NewDecoder(resp.Body).Decode(&todo)
	id2 := todo.Id

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, t := range todos {
		assert.Equal(t.Id, id2)
	}
}
