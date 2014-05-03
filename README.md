kvlite3
========
simple key-value store backed by sqlite3 for the Go language

[![Build Status](https://travis-ci.org/walkr/kvlite3.svg?branch=master)](https://travis-ci.org/walkr/kvlite3)

## "Features":

 * basic put/get/getall operations
 * encodes values as json


## Usage:

```go
client, _ := kvlite3.NewClient("./todos.db")
bucket, _ := client.NewBucket("todos")

todo := Todo{kvlite3.NewKey(), "Take out trash"}
bucket.Put(todo.Key, todo) // Insert
```

[*] See examples/simple.go for for more details


## TODO:

 * Implement `putMany`
 * Implement `getMany`
 * Implement limit/offset in `getAll`
