package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountModel_ToEntity(t *testing.T) {
	model := &model{
		Category:   "wechat",
		Name:       "test_account",
		LastURL:    "http://test.com",
		Cookies:    []byte("test_cookies"),
		InstanceID: 1,
		State:      int(stateNew),
	}

	entity := model.toEntity()

	assert.Equal(t, model.Category, entity.Category)
	assert.Equal(t, model.Name, entity.Name)
	assert.Equal(t, model.LastURL, entity.lastURL)
	assert.Equal(t, model.Cookies, entity.cookies)
	assert.Equal(t, model.InstanceID, entity.instanceID)
	assert.Equal(t, state(model.State), entity.State)
}

func TestAccountModel_FromEntity(t *testing.T) {
	entity := &account{
		Category:   "wechat",
		Name:       "test_account",
		lastURL:    "http://test.com",
		cookies:    []byte("test_cookies"),
		instanceID: 1,
		State:      stateNew,
	}

	model := &model{}
	model.fromEntity(entity)

	assert.Equal(t, entity.Category, model.Category)
	assert.Equal(t, entity.Name, model.Name)
	assert.Equal(t, entity.lastURL, model.LastURL)
	assert.Equal(t, entity.cookies, model.Cookies)
	assert.Equal(t, entity.instanceID, model.InstanceID)
	assert.Equal(t, int(entity.State), model.State)
}
