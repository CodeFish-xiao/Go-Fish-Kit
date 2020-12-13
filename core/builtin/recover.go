package builtin

import (
	"fmt"
	"log"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if painErr := recover(); painErr != nil {
		log.Printf(fmt.Sprint(painErr)) //等后续写了自己的日志部分改为自己的日志
	}
}
