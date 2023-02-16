package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/SeungyeonHwang/tool-goaal/login"
	"github.com/SeungyeonHwang/tool-goaal/model"
	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	login.GetSessionId = func(r *http.Request) string {
		return "testsessionId"
	}

	os.Remove("../db/test_todo.db")
	assert := assert.New(t)
	ah := MakeHandler("../db/test_todo.db")
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	//addTodoListHandler
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data1"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Id, 1)
	assert.Equal(todo.Name, "Test Data1")
	assert.Equal(reflect.ValueOf(todo.CreatedAt).String(), "<time.Time Value>")
	id1 := todo.Id

	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Data2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Id, 2)
	assert.Equal(todo.Name, "Test Data2")
	assert.Equal(reflect.ValueOf(todo.CreatedAt).String(), "<time.Time Value>")
	id2 := todo.Id

	//getTodoListHandlerTest
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
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

	//completeTodoListHandler
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.Id == id1 {
			assert.True(t.Completed)
		}
	}

	//removeTodoListHandler
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, t := range todos {
		assert.Equal(t.Id, id2)
	}
}
