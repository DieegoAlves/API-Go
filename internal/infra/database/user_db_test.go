package database

import (
	"github.com/DieegoAlves/API/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func Test_CreateNewUser(t *testing.T) {
	//Criar conexão com o banco de dados
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	//Criar tabela User no DB
	db.AutoMigrate(&entity.User{})

	//Criar um usuário
	user, _ := entity.NewUser("Diego", "diegoaf@ucl.br", "Revolution22#")

	//passando usuário para o banco de dados
	userDB := NewUser(db)

	//Criar usuário no banco de dados
	err = userDB.CreateUser(user)
	assert.Nil(t, err)

	//Buscar usuário no banco de dados
	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)
}
