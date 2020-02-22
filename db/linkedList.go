package db

import (
	"os"
)

type linkedListStore struct {
	root       *Node // root of the linked list
	currentPos int64 //current position in file
	sm         storageManager
	nodeCap    int64 //determines the max capacity a single node will occupy
}

func newLinkedListStore(filepath string) (*linkedListStore, map[string]int64) {

	var cache map[string]int64
	sm, mode := newFileStorageManager(filepath)

	lls := &linkedListStore{sm: sm, nodeCap: capsize}

	//if file is created newly, create root node
	if mode == os.O_CREATE {
		lls.createRootNode()
	}

	//If data file already exits, read existing file and load cache and currentPos
	if mode == os.O_RDWR {
		currentPos, c := lls.readFile()

		//set currentPos, cache and root from read file
		lls.currentPos = currentPos
		cache = c
		lls.root = lls.getRootNode()
	}

	return lls, cache
}

func (lls *linkedListStore) addNode(n *Node) {

	n.Pos = lls.currentPos
	nextPos := lls.getNextPos()
	n.NextPos = nextPos
	lls.sm.write(encode(n), n.Pos)
	lls.currentPos = nextPos

}

//this is costly, get position from memcache and use getNodeByPos
func (lls *linkedListStore) getNode(key string) (*Node, error) {
	n, err := lls.traverse(key)

	if err != nil {
		return &Node{}, err
	}

	return lls.getNodeByPos(key, n.Pos)
}

func (lls *linkedListStore) getNodeByPos(key string, pos int64) (*Node, error) {
	buf, err := lls.sm.read(lls.nodeCap, pos)
	n := decode(buf)
	if key == key {
		return n, nil
	}
	return &Node{}, err
}

//this is costly, get position from memcache and use deleteNodeByPos
func (lls *linkedListStore) deleteNode(key string) error {

	node, err := lls.traverse(key)

	if err != nil {
		return err
	}

	lls.deleteNodeByPos(key, node.Pos)

	return nil

}

func (lls *linkedListStore) deleteNodeByPos(key string, pos int64) {
	empty := make([]byte, capsize)
	prevPos := pos - lls.nodeCap
	nextPos := pos + lls.nodeCap

	lls.sm.write(empty, pos)

	prevBuf, _ := lls.sm.read(capsize, prevPos)
	nextBuf, _ := lls.sm.read(capsize, nextPos)

	prevNode := decode(prevBuf)
	nextNode := decode(nextBuf)

	if prevNode.Key != "" && nextNode.Key != "" {
		prevNode.NextPos = nextNode.Pos

	}

	lls.sm.write(encode(prevNode), prevNode.Pos)
}

//this is costly, use memcache instead, traversing 1GB file takes 6 secs
func (lls *linkedListStore) traverse(key string) (*Node, error) {

	return lls.traverseRecur(key, lls.root.Pos)
}

func (lls *linkedListStore) traverseRecur(key string, offset int64) (*Node, error) {

	for {

		res, err := lls.sm.read(lls.nodeCap, offset)

		if err != nil && err.Error() != "EOF" {
			return &Node{}, err
		}

		n := decode(res)

		if offset > 0 && n.Key == "" {
			return &Node{}, &keyNotFoundError{}
		}

		if n.Key == key {
			return n, nil
		}
		return lls.traverseRecur(key, n.NextPos)
	}
}
func (lls *linkedListStore) createRootNode() {
	root := lls.getRootNode()
	lls.addNode(root)
	lls.root = root

}

func (lls *linkedListStore) getRootNode() *Node {
	root := &Node{Pos: rootPosition, Key: "", Value: ""}
	root.NextPos = root.Pos + lls.nodeCap
	return root

}

func (lls *linkedListStore) readFile() (int64, map[string]int64) {
	return lls.sm.loadFile()
}

func (lls *linkedListStore) getNextPos() int64 {

	return lls.currentPos + lls.nodeCap
}

func (lls *linkedListStore) isfull() bool {
	return lls.sm.isfull()
}

func (lls *linkedListStore) close() {
	lls.sm.close()
}
