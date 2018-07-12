package golock

import (
	"errors"
	"os"
	"time"
)

// Lock is the structure for containing the lock attributes.
type Lock struct {
	name     string
	f        *os.File
	haveLock bool
	timeout  time.Duration
	interval time.Duration
}

// Option is the type all options need to adhere to
type Option func(l *Lock)

// OptionSetName sets the name
func OptionSetName(name string) Option {
	return func(l *Lock) {
		l.name = name
	}
}

// OptionSetTimeout sets the timeout (default: none)
func OptionSetTimeout(timeout time.Duration) Option {
	return func(l *Lock) {
		l.timeout = timeout
	}
}

// OptionSetInterval sets the interval to check (default: none)
func OptionSetInterval(interval time.Duration) Option {
	return func(l *Lock) {
		l.interval = interval
	}
}

// New creates a new lock with the specified options.
// If no name is specified with the options, the default name
// is "golock.lock".
func New(options ...Option) *Lock {
	l := new(Lock)
	l.name = "golock.lock"
	for _, o := range options {
		o(l)
	}
	return l
}

// Lock will try to lock by creating a new file. If the file exists, it will
// throw an error, unless you are using timeouts in which case it will poll until
// it can create the file. If it is still unable to create the file it will throw
// an error.
func (l *Lock) Lock() (err error) {
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

// Unlock will remove the file that it used for locking
// and throw an error if anything goes wrong.
func (l *Lock) Unlock() (err error) {
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
