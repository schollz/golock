# golock

[![travis](https://travis-ci.org/schollz/golock.svg?branch=master)](https://travis-ci.org/schollz/golock) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/golock)](https://goreportcard.com/report/github.com/schollz/golock) 
[![coverage](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)](https://gocover.io/github.com/schollz/golock)
[![godocs](https://godoc.org/github.com/schollz/golock?status.svg)](https://godoc.org/github.com/schollz/golock) 

Very simple (<100 LOC) file locking with optional timeouts. 


## Install

```
go get github.com/schollz/golock
```

## Usage 

Initialize the lock and then obtain it. If you specify the timeout and interval, it will poll at the given interval for the specified time until it successfully gets a lock, otherwise throw an error. If you don't specify the timeout, then it will throw an error immediately if it does not obtain the lock.

If you get no errors from locking, then you are good to go. Make sure to unlock it when you are done.

```golang
// first initiate lockfile
l := golock.New(
    "mylockfile", 
    golock.OptionSetInterval(1*time.Millisecond), 
    golock.OptionSetTimeout(60*time.Second),
)

// lock it
err := l.Lock()
if err != nil {
    // error means we didn't get the lock
    // handle it
}

// do stuff

// unlock it
err = l.Unlock()
if err != nil {
    panic(err)
}
```

## Benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/schollz/golock
BenchmarkLocking-4     200000     12224 ns/op    128 B/op    5 allocs/op
```

## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT