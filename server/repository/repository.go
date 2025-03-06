package repository

import "github.com/muhamadrizkiariffadillah/bookshop-vue-go-monggo/config"

type Repository struct {
	UserRepository MongoRepositoryInterface
}

func GetRepository() *Repository {
	return &Repository{UserRepository: GetMongoRepository(config.GetEnvProperties("database_name"), "user")}
}
