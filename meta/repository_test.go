package meta

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestLoadRepository(t *testing.T) {
	r := &reader_test{}
	r.content = `{"products":[{"code":"ecs"}]}`
	repository := LoadRepository(r)
	assert.NotNil(t, repository)
	assert.Contains(t, repository.Names, "ecs")

	defer func(){
		err := recover()
		assert.NotNil(t, err)
	}()
	r.content = ""
	repository = LoadRepository(r)
	assert.Nil(t, repository)
}

func TestLoadRepository1(t *testing.T) {
	r := &reader_test{}
	defer func(){
		err := recover()
		assert.NotNil(t, err)
	}()
	r.content = `{"products":[{"code":"ecs"},{"code":"ecs"}]}`
	repository := LoadRepository(r)
	assert.Nil(t, repository)
}

func TestRepository_GetApi(t *testing.T) {
	repository := &Repository{
		index: map[string]Product{
			"ecs": {
				Code: "ecs",
			},
		},
		reader: &reader_test{
			content: `{"name":"ecs""protocol":"http"}`,
		},
	}
	_, ok := repository.GetApi("ros", "", "")
	assert.False(t, ok)

	_, ok = repository.GetApi("ecs", "", "")
	assert.False(t, ok)

	repository.reader = &reader_test{
		content: `{"name":"ecs","protocol":"http"}`,
	}
	_, ok = repository.GetApi("ecs", "", "")
	assert.True(t, ok)
}