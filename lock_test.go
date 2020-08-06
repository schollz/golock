package golock

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkLocking(b *testing.B) {
	l := New(OptionSetName("testlock"))
	for i := 0; i < b.N; i++ {
		l.Lock()
		l.Unlock()
	}
}

func TestLocking(t *testing.T) {
	l := New(OptionSetName("testlock"))
	err := l.Lock()
	assert.Nil(t, err)

	l2 := New(OptionSetName("testlock"), OptionSetInterval(1*time.Millisecond), OptionSetTimeout(100*time.Millisecond))
	err = l2.Lock()
	assert.NotNil(t, err)

	l3 := New(OptionSetName("testlock3"))
	err = l3.Lock()
	assert.Nil(t, err)

	err = l.Unlock()
	assert.Nil(t, err)

	err = l2.Lock()
	assert.Nil(t, err)
	err = l2.Unlock()
	assert.Nil(t, err)
}

func TestTimeout(t *testing.T) {
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
	l1 := New()
	l2 := New(OptionSetInterval(1*time.Millisecond), OptionSetTimeout(100*time.Millisecond))
	l1.Lock()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		time.Sleep(50 * time.Millisecond)
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

func TestMultiprocessing(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			l := New(OptionSetName("lock11"), OptionSetInterval(500*time.Microsecond), OptionSetTimeout(100*time.Second))
			err := l.Lock()
			assert.Nil(t, err)
			time.Sleep(100 * time.Millisecond)
			err = l.Unlock()
			if err != nil {
				fmt.Println(i, err.Error())
			}
			assert.Nil(t, err)
		}(i)
	}
	wg.Wait()
}
