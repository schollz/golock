package golock

import (
	"errors"
	"os"
	"time"
)

type lock struct {
	name     string
	f        *os.File
	haveLock bool
	timeout  time.Duration
	interval time.Duration
}

// Option is the type all options need to adhere to
type Option func(l *lock)

// OptionSetName sets the name
func OptionSetName(name string) Option {
	return func(l *lock) {
		l.name = name
	}
}

// OptionSetTimeout sets the timeout (default: none)
func OptionSetTimeout(timeout time.Duration) Option {
	return func(l *lock) {
		l.timeout = timeout
	}
}

// OptionSetInterval sets the interval to check (default: none)
func OptionSetInterval(interval time.Duration) Option {
	return func(l *lock) {
		l.interval = interval
	}
}

func New(options ...Option) *lock {
	l := new(lock)
	l.name = "golock.lock"
	for _, o := range options {
		o(l)
	}
	return l
}

func (l *lock) Lock() (err error) {
	start := time.Now()
	for {
		if l.timeout.Nanoseconds() > 0 {
			if time.Since(start).Nanoseconds() > l.timeout.Nanoseconds() {
				return errors.New("could not obtain lock, timeout")
			}
		}
		l.f, err = os.OpenFile(l.name, os.O_CREATE|os.O_EXCL, 0755)
		if err != nil {
			if l.timeout.Nanoseconds() == 0 {
				return errors.New("could not obtain lock")
			} else {
				// wait for lock
				time.Sleep(l.interval)
			}
		} else {
			break
		}
	}
	l.haveLock = true
	return
}

func (l *lock) Unlock() (err error) {
	if !l.haveLock {
		return errors.New("no lock")
	}
	err = l.f.Close()
	if err != nil {
		return
	}
	err = os.Remove(l.name)
	return
}
