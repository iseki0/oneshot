package oneshot

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestOneShot(t *testing.T) {
	source, subject := New()
	wg := sync.WaitGroup{}
	go func() {
		time.Sleep(time.Second * 1)
		subject.Emit("foo")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_t := time.Now()
		assert.False(t, source.Wait(context.TODO(), "foo", "bar"))
		assert.True(t, time.Now().Sub(_t).Milliseconds() >= 1000, "duang~~~")
		fmt.Println("=========", time.Now().Sub(_t).Microseconds())
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeout, cancelFunc := context.WithTimeout(context.TODO(), time.Millisecond*50)
		defer cancelFunc()
		assert.True(t, source.Wait(timeout, "bar"), "duang22")
	}()
	wg.Wait()
}
