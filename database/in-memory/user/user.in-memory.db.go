package user

import (
	"FaRyuk/config"
	"FaRyuk/internal/types"
)

type UserInMemoryRepository struct {
	config *config.AppConfig
}

func NewInMemoryRepository(config *config.AppConfig) *UserInMemoryRepository {
	return &UserInMemoryRepository{}
}

// CloseConnection : closes connection with mongo db
func (repo *UserInMemoryRepository) CloseUserDBConnection() {

}

var data []*types.User = []*types.User{}
