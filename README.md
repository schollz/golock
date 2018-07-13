# golock

[![travis](https://travis-ci.org/schollz/golock.svg?branch=master)](https://travis-ci.org/schollz/golock) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/golock)](https://goreportcard.com/report/github.com/schollz/golock) 
[![coverage](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)](https://gocover.io/github.com/schollz/golock)
[![godocs](https://godoc.org/github.com/schollz/golock?status.svg)](https://godoc.org/github.com/schollz/golock) 

Very simple file locking with optional waiting/timeouts.


## Install

```
go get -u github.com/schollz/golock
```

## Usage 

```golang
// first initiate lockfile
l := golock.New(
    "mylockfile", 
    golock.OptionSetInterval(1*time.Millisecond), 
    golock.OptionSetTimeout(60*time.Second),
)

// use it wherever

// lock it
err := l.Lock()
if err != nil {
    panic(err)
}

// be sure to unlock it
err = l.Unlock()
if err != nil {
    panic(err)
}
```

## Benchmarks

```
goos: windows
goarch: amd64
pkg: github.com/schollz/golock
BenchmarkLocking-8    5000    302266 ns/op    720 B/op    5 allocs/op
```

## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT