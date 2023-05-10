package main

import (
	"fmt"
	"github.com/jumaniyozov/gdo/models"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected")

	us := models.UserService{
		DB: db,
	}

	user, err := us.Create("bbb@mail.com", "123456qwe")
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

}
