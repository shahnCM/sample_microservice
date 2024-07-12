package safeasync

import (
	"log"
)

func Run(callback func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("async.RunAsync recovered from panic: ", r)
			}
		}()
		callback()
	}()
}
