package main

import (
	"encoding/json"
	"fmt"
	"kvstore/db"
	"strconv"
	"sync"
	"time"
)

//Person is sample test object
type Person struct {
	ID   int
	Name string
}

func main() {

	var k = "ADSHSDF DKJHFJHGBDSF JHDFGJ HJDFG JHBDFHGDBF HJDSBFGJ HDFB HJBDFGHD"
	person := Person{ID: 12568, Name: "some test user"}

	db := db.NewKVStore("E:\\test2")
	defer db.Close()

	//do concurrent Creates upto 1GB
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
	t1 := time.Now()

	//get something from last...
	fmt.Println(db.Get("14995ADSHSDF DKJHFJHGBDSF JHDFGJ"))

	t2 := time.Now()

	//output from console : time to traverse  999.4µs
	fmt.Println("time to traverse ", t2.Sub(t1))

	fmt.Println("completed")

}
