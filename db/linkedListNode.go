package db

//Node represents each node in the data file.
type Node struct {
	Key     string
	Value   string
	Pos     int64 //Current position in the file
	NextPos int64 //Holds the next position
}

func getNewNode(key string, val string) *Node {
	return &Node{Key: key, Value: val}
}
