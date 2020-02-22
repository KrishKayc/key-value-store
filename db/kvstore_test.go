package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Person struct {
	ID   int
	Name string
}

const k = "ADSHSDF DKJHFJHGBDSF JHDFGJ HJDFG JHBDFHGDBF HJDSBFGJ HDFB HJBDFGHD"

var p = Person{ID: 12568, Name: "some test user"}

func TestSimpleSetGet(t *testing.T) {
	db := NewKVStore("E:\\test2")
	defer db.Close()

	v, _ := json.Marshal(p)

	err := db.Create(k, v, 0)

	if err != nil {
		t.Errorf("set not working")
	}

	got, err := db.Get(k)
	var p1 Person
	json.Unmarshal(got, &p1)
	if err != nil || p1 != p {
		t.Errorf("get not working")
	}

}

func TestGetFromExistingFile(t *testing.T) {
	//file already has data written to it through db.Set(k)
	db := NewKVStore("E:\\test2")
	defer db.Close()

	got, err := db.Get(k)
	var p1 Person
	json.Unmarshal(got, &p1)
	if err != nil || p1 != p {
		t.Errorf("Get set not working")
	}

}

func TestSimpleSetGet_WithNoPath(t *testing.T) {
	db := NewKVStore("")
	defer db.Close()

	v, _ := json.Marshal(p)

	err := db.Create(k, v, 0)

	if err != nil {
		t.Errorf("set not working")
	}

	got, err := db.Get(k)
	var p1 Person
	json.Unmarshal(got, &p1)
	if err != nil || p1 != p {
		t.Errorf("get not working")
	}

}

func TestConcurrentSet_Upto_1GB(t *testing.T) {
	person := Person{ID: 12568, Name: "some test user"}

	db := NewKVStore("E:\\test2")
	defer db.Close()

	var wg sync.WaitGroup
	for i := 1; i <= 15000; i++ {
		v, _ := json.Marshal(person)
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup, v []byte) {
			defer wg.Done()
			err := db.Create(strconv.Itoa(i)+k, v, 0)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(i, &wg, v)

	}
	wg.Wait()
}

func Test_Get_After_ConcurrentSet(t *testing.T) {
	db := NewKVStore("E:\\test2")
	defer db.Close()

	//get something from last
	//output from console : time to traverse  999.4µs
	got, err := db.Get("14995ADSHSDF DKJHFJHGBDSF JHDFGJ")
	var p1 Person
	json.Unmarshal(got, &p1)
	if err != nil || p != p1 {
		t.Errorf("Get not working after restarting the app.")
	}
}

func Test_Key_Exists(t *testing.T) {
	db := NewKVStore("E:\\test2")
	defer db.Close()
	v, _ := json.Marshal(p)
	db.Create(k, v, 0)

	err := db.Create(k, v, 0)

	if err == nil || err.Error() != "Key already exists !!" {
		t.Errorf("Duplicate key validation failed")
	}
}

func TestExpiration(t *testing.T) {
	db := NewKVStore("E:\\test2")
	defer db.Close()

	v, _ := json.Marshal(p)

	err := db.Create(k, v, 5)

	if err != nil {
		t.Errorf("set not working")
	}

	time.Sleep(7 * time.Second)

	_, err = db.Get(k)

	if err == nil || err.Error() != "Key not found" {
		t.Errorf("expiration not working")
	}

}
func Benchmark_Set(b *testing.B) {
	db := NewKVStore("E:\\test2")
	defer db.Close()
	v, _ := json.Marshal(p)
	for i := 0; i < b.N; i++ {
		db.Create(strconv.Itoa(i)+k, v, 0)
	}
}
