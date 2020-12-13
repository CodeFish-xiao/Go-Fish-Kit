package thread

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestRoutineId(t *testing.T) {
	assert.True(t, RoutineId() > 0)
}

func TestGo(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan int)
	go Go(func() {
		defer func() {
			ch <- 1
		}()

		panic("panic")
	})

	<-ch
	i++
}
