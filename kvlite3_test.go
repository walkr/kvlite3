package kvlite3

import (
	"os"
	"strconv"
	"testing"
)

type Model struct {
	Key   string
	Value string
}

func setUp() *Bucket {
	client, _ := NewClient("./tests.db")
	bucket, _ := client.NewBucket("tests")
	return bucket
}
func tearDown() {
	os.Remove("./tests.db")
}

func TestNewClient(t *testing.T) {
	_, err := NewClient("./tests.db")
	if err != nil {
		t.Error(err)
	}

}

func TestNewBucket(t *testing.T) {
	client, _ := NewClient("./tests.db")
	_, err := client.NewBucket("tests")
	if err != nil {
		t.Error(err)
	}
}

func generateItems(n int) []*Model {
	items := make([]*Model, n)

	for i := 0; i < n; i++ {
		m := &Model{strconv.Itoa(i), strconv.Itoa(i * i)}
		items[i] = m
	}
	return items
}

func insertItems(bucket *Bucket, t *testing.T) []*Model {
	// Insert items
	items := generateItems(10)
	for _, item := range items {
		err := bucket.Put(item.Key, item)
		if err != nil {
			t.Error(err)
		}
	}
	return items
}

func TestPut(t *testing.T) {
	bucket := setUp()
	insertItems(bucket, t)
	tearDown()
}

func TestUpdate(t *testing.T) {
	bucket := setUp()
	items := insertItems(bucket, t)

	// Update items
	for _, item := range items {
		item.Value = item.Value + item.Value
		err := bucket.Put(item.Key, item)
		if err != nil {
			t.Error(err)
		}
	}

	// Fetch and check
	fetched := Model{}
	for _, item := range items {
		err := bucket.Get(item.Key, &fetched)
		if err != nil {
			t.Error(err)
		}
		if fetched != *item {
			t.Error("Retrieved item doesn't match the original")
		}
	}

	tearDown()
}
func TestGet(t *testing.T) {
	bucket := setUp()
	items := insertItems(bucket, t)

	// Fetch and check
	fetched := Model{}
	for _, item := range items {
		err := bucket.Get(item.Key, &fetched)
		if err != nil {
			t.Error(err)
		}
		if fetched != *item {
			t.Error("Retrieved item doesn't match the original")
		}
	}
	tearDown()
}

func TestGetAll(t *testing.T) {
	bucket := setUp()
	items := insertItems(bucket, t)
	N := len(items)

	found := []*Model{}
	err := bucket.GetAll(&found)
	if err != nil {
		t.Error(err)
	}

	if len(found) != N {
		t.Error("Count fetched != count inserted")
	}

	tearDown()
}

// --------------
// Benchmarks
// --------------

func BenchmarkPut(b *testing.B) {
	bucket := setUp()
	// Insert items
	items := generateItems(300)
	b.ResetTimer()
	for _, item := range items {
		bucket.Put(item.Key, item)
	}
	tearDown()
}

func BenchmarkPutGet(b *testing.B) {
	bucket := setUp()
	// Insert items
	items := generateItems(300)
	for _, item := range items {
		bucket.Put(item.Key, item)
	}
	b.ResetTimer()

	for _, item := range items {
		bucket.Get(item.Key, &Model{})
	}
	tearDown()
}
