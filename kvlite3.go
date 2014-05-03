package kvlite3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	Db *sql.DB
}

type Bucket struct {
	Name   string
	Client *Client
}

type Result struct {
	Key   string
	Value string
}

// Create a unique 16 bytes key
func NewKey() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := make([]byte, 16)
	f.Read(b)
	return fmt.Sprintf("%x", b)
}

// Create a new client
func NewClient(path string) (*Client, error) {
	db, err := sql.Open("sqlite3", path)
	return &Client{db}, err
}

// Create a new bucket
func (c *Client) NewBucket(name string) (*Bucket, error) {
	bucket := &Bucket{}
	sql := "create table if not exists " + name + " (key text primary key unique, value text)"
	_, err := c.Db.Exec(sql)
	if err == nil {
		bucket.Name, bucket.Client = name, c
	}
	return bucket, err
}

// Insert or replace an item of a given key
func (b *Bucket) Put(key string, value interface{}) error {
	err := insert(b, key, value)
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return update(b, key, value)
		} else {
			panic(err)
		}
	}

	return err
}

// Gets an item from the bucket
func (b *Bucket) Get(key string, result interface{}) error {
	sql := "select key, value from " + b.Name + " where key=?"
	stmt, err := b.Client.Db.Prepare(sql)

	if err != nil {
		return err
	}
	defer stmt.Close()

	var k, v string
	err = stmt.QueryRow(key).Scan(&k, &v)
	json.Unmarshal([]byte(v), result)
	return err
}

// Delete an item
func (b *Bucket) Del(key string) error {
	sql := "delete from " + b.Name + " where key=?"
	return execute(b.Client.Db, sql, key)
}

// insert a new sqlite3 row
func insert(b *Bucket, key string, value interface{}) error {
	sql := "insert into " + b.Name + " (key, value) values (?, ?)"
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return execute(b.Client.Db, sql, key, string(encoded))
}

// update existing sqlite3 row
func update(b *Bucket, key string, value interface{}) error {
	sql := "update " + b.Name + " set key=?, value=? where key=?"
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return execute(b.Client.Db, sql, key, string(encoded), key)
}

// Execute sql command
func execute(db *sql.DB, sql string, args ...interface{}) error {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	return err
}

func query(db *sql.DB, sql string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	return rows, err
}

// Get all items in
func (b *Bucket) GetAll(result interface{}) error {
	sql := "select value from " + b.Name // + " limit " + "10" + " offset " + "10"
	rows, err := b.Client.Db.Query(sql)

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("result must be slice address")
	}

	slicev := resultv.Elem()
	elemt := slicev.Type().Elem()

	count := 0
	var value string

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
		elemp := reflect.New(elemt)
		json.Unmarshal([]byte(value), elemp.Interface())
		slicev = reflect.Append(slicev, elemp.Elem())
		count += 1
	}
	resultv.Elem().Set(slicev.Slice(0, count))
	return err
}
