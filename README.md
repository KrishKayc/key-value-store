# key-value-store
A file based key-value store.

Supports string keys and json values, concurrent read and writes from the store and 'eviction' of keys with expiration time.

Functions to Create, Get, check if key Exists and Delete a key from the store

### Usage

```
  //Object in the application
  person := Person{ID: 12568, Name: "some test user"}
  v, _ := json.Marshal(person)
  
  //create new KV store in specified path
  db := NewKVStore("E:\\test2")
  defer db.Close()
  
  //Add to store
  db.Create("key1", v, 0)
  
  //Get from store
  value, err := db.Get(k)
  var p1 Person
  json.Unmarshal(value, &p1)
  
  ```

### TODO :

- Currently supports 1GB of storage, must be made configurable.
- Supports values upto 16KB in size and keys lenght is truncated to 32 chars, this also need to be configurable

