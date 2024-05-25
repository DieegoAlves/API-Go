package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Diego", "diegoaf@ucl.br", "Revolution22#")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "Diego", user.Name)
	assert.Equal(t, "diegoaf@ucl.br", user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, _ := NewUser("Diego", "diegoaf@ucl.br", "Revolution22#")
	assert.True(t, user.ValidatePassword("Revolution22#"))
	assert.False(t, user.ValidatePassword("Revolution22"))
	assert.NotEqual(t, "Revolution22#", user.Password)
}
