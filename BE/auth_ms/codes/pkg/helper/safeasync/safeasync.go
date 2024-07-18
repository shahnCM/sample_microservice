package safeasync

import (
	"auth_ms/pkg/provider/database/mariadb10"
	"log"
)

func Run(callback func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				mariadb10.TransactionRollback()
				log.Println("safeasync.Run recovered from panic: ", r)
			}
		}()
		callback()
	}()
}
