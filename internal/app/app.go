package app

import (
	"wblayerzero/internal/config"
)

func Run() error {
	_, err := config.Load(config.CfgFilePath)
	if err != nil {
		return err
	}

	return nil
}
