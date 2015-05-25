

$ cd <DIR> 
$ git init
$ git remote add diff https://github.com/bkln/diff.git
$ git pull diff
$ git merge diff/master
$ export GOPATH=<DIR>
$ go get github.com/gorilla/rpc
$ go get github.com/syndtr/goleveldb/leveldb
$ go get github.com/gorilla/rpc
$ go get github.com/cznic/kv 
$ go build src/puzz/server/server.go
$ ./server
$ curl....

