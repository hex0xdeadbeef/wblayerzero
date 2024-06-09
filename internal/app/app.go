package app

import (
	"wblayerzero/internal/config"
	"wblayerzero/internal/database/postgre"

	"log"
)

func Run() error {
	err := config.Load(config.CfgFilePath)
	if err != nil {
		return err
	}

	storage, err := postgre.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
