# go-chat
chat server in go from [this](https://www.youtube.com/watch?v=5buaPyJ0XeQ) great talk by Dave Cheney. Using channels for concurrency instead of mutexes.

# usage
```bash
$ cd single-channel
# server
$ go run *.go
# client
$ nc localhost 13990
```

# license
MIT
