package user

import (
	"github.com/gomscourse/auth/internal/repository"
	"github.com/gomscourse/auth/internal/service"
	"github.com/gomscourse/common/pkg/db"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{userRepository: userRepository, txManager: txManager}
}
