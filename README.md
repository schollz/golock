# golock

Very simple file locking.

```golang
l := golock.New("mylockfile",golock.OptionSetInterval(1*time.Millisecond), golock.OptionSetTimeout(60*time.Second))
err := l.Lock()
if err != nil {
    panic(err)
}
err = l.Unlock()
if err != nil {
    panic(err)
}
```

# License

MIT