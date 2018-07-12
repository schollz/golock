package golock

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkLocking(b *testing.B) {
	os.Remove("testlock")
	l := New(OptionSetName("testlock"))
	for i := 0; i < b.N; i++ {
		l.Lock()
		l.Unlock()
	}
}

func TestLocking(t *testing.T) {
	os.Remove("testlock")

	l := New(OptionSetName("testlock"))
	err := l.Lock()
	assert.Nil(t, err)

	l2 := New(OptionSetName("testlock"))
	err = l2.Lock()
	assert.NotNil(t, err)

	err = l.Unlock()
	assert.Nil(t, err)

	err = l2.Unlock()
	assert.NotNil(t, err)

	err = l2.Lock()
	assert.Nil(t, err)

}

func TestTimeout(t *testing.T) {
	os.Remove("golock.lock")

	l1 := New()
	l2 := New(OptionSetInterval(1*time.Millisecond), OptionSetTimeout(100*time.Millisecond))
	l1.Lock()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		time.Sleep(200 * time.Millisecond)
		l1.Unlock()
		wg.Done()
	}()
	go func() {
		err := l2.Lock()
		assert.NotNil(t, err)
		wg.Done()
	}()
	wg.Wait()
}

func TestNoTimeout(t *testing.T) {
	os.Remove("golock.lock")
	l1 := New()
	l2 := New(OptionSetInterval(1*time.Millisecond), OptionSetTimeout(100*time.Millisecond))
	l1.Lock()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		time.Sleep(10 * time.Millisecond)
		l1.Unlock()
		wg.Done()
	}()
	go func() {
		err := l2.Lock()
		assert.Nil(t, err)
		wg.Done()
	}()
	wg.Wait()
}
