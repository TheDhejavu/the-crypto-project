package blockchain

import (
	log "github.com/sirupsen/logrus"
)

func Handle(err error) {

	if err != nil {
		log.Panic(err)
	}
}
