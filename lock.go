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
	// default name of the lock
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
	// continually try to lock
	for {
		// try to open+create file and error if it already exists
		l.f, err = os.OpenFile(l.name, os.O_CREATE|os.O_EXCL, 0755)
		if err != nil {
			// if not using a timeout, return an error
			if l.timeout.Nanoseconds() == 0 {
				return errors.New("could not obtain lock")
			} else {
				// we are using a timeout, wait for lock
				time.Sleep(l.interval)
			}
		} else {
			// no error, close file
			// and break out of loop
			l.f.Close()
			break
		}
		// check if we are using a timeout
		if l.timeout.Nanoseconds() > 0 {
			// check if timeout is exceeded and return error
			if time.Since(start).Nanoseconds() > l.timeout.Nanoseconds() {
				return errors.New("could not obtain lock, timeout")
			}
		}
	}
	return
}

// Unlock will remove the file that it used for locking
func (l *Lock) Unlock() (err error) {
	err = os.Remove(l.name)
	return
}
