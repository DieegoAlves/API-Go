package database

import (
	"fmt"
	"github.com/DieegoAlves/API/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
	"testing"
)

func ConexaoCompletaDB() *gorm.DB {
	//Criar conexão com o banco de dados
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil
	}

	//Criar tabela Product no DB
	db.AutoMigrate(&entity.Product{})

	return db
}

func Test_CreateNewProduct(t *testing.T) {
	db := ConexaoCompletaDB()

	//Criando a Estrutura Product no Banco de Dados
	productDB := NewProduct(db)

	//Criar um produto
	product, err := entity.NewProduct("Produto 1", 10.5)
	assert.NoError(t, err)

	//Criar produto no banco de dados
	err = productDB.CreateProduct(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func Test_FindAllProducts(t *testing.T) {
	db := ConexaoCompletaDB()

	//passando estrutura Product para o banco de dados
	productDB := NewProduct(db)

	//Criando vários produtos
	for i := 1; i < 25; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		err = productDB.CreateProduct(product)
		assert.NoError(t, err)
		assert.NotEmpty(t, product.ID)

	}

	//Buscar todos os produtos da página 1, com limite de 10 produtos e ordenados de forma ascendente
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.NotEmpty(t, products)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	//Buscar todos os produtos da página 2, com limite de 10 produtos e ordenados de forma ascendente
	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.NotEmpty(t, products)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	//Buscar todos os produtos da página 3, com limite de 10 produtos e ordenados de forma ascendente
	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.NotEmpty(t, products)
	assert.Len(t, products, 4)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 24", products[3].Name)
}

func Test_FindProductByID(t *testing.T) {
	db := ConexaoCompletaDB()

	//passando estrutura Product para o banco de dados
	productDB := NewProduct(db)

	//Criar um produto
	product, err := entity.NewProduct("Produto 1", 10.5)
	assert.NoError(t, err)

	//Passando o produto para o banco de dados
	err = productDB.CreateProduct(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)

	//Buscar produto pelo ID
	productByID, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.NotNil(t, productByID)
	assert.Equal(t, product.ID, productByID.ID)
}

func Test_UpdateProduct(t *testing.T) {
	db := ConexaoCompletaDB()

	//passando estrutura Product para o banco de dados
	productDB := NewProduct(db)

	//Criar um produto
	product, err := entity.NewProduct("Produto 1", 10.5)
	assert.NoError(t, err)

	//Passando o produto para o banco de dados
	err = productDB.CreateProduct(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)

	//Alterar o nome do produto
	product.Name = "Produto 2"
	err = productDB.UpdateProduct(product)
	assert.NoError(t, err)

	//Buscar produto pelo ID
	productByID, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.NotNil(t, productByID)
	assert.Equal(t, product.Name, productByID.Name)
}

func Test_DeleteProduct(t *testing.T) {
	db := ConexaoCompletaDB()

	//passando estrutura Product para o banco de dados
	productDB := NewProduct(db)

	//Criar um produto
	product, err := entity.NewProduct("Produto 1", 10.5)
	assert.NoError(t, err)

	//Passando o produto para o banco de dados
	err = productDB.CreateProduct(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)

	//Deletar produto pelo ID
	err = productDB.DeleteProduct(product.ID.String())
	assert.NoError(t, err)

	//Buscar produto pelo ID
	ProductByID, err := productDB.FindByID(product.ID.String())
	assert.Error(t, err)
	assert.Empty(t, ProductByID)
}
