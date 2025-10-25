package img_cache

import (
	"fmt"
	"strconv"
	"strings"
)

type node struct {
	key  string
	next *node
}

type imgQueue struct {
	front *node
	rear  *node
	size  int
}

func (q *imgQueue) encodeKey(key string, imgType ImgT) string {
	return fmt.Sprintf("%d_%s", imgType, key)
}

func (q *imgQueue) decodeKey(key string) (string, ImgT, error) {
	arr := strings.Split(key, "_")
	if len(arr) != 2 {
		return "", 0, fmt.Errorf("invalid key")
	}

	imgType, err := strconv.ParseInt(arr[0], 10, 8)
	if err != nil {
		return "", 0, err
	}

	resKey := arr[1]

	return resKey, ImgT(imgType), nil
}

func (q *imgQueue) push(key string) {
	newNode := &node{key: key}

	if q.rear == nil {
		q.front = newNode
		q.rear = newNode
	} else {
		q.rear.next = newNode
		q.rear = newNode
	}
	q.size++
}

func (q *imgQueue) pop() string {
	if q.size == 0 {
		return ""
	}

	value := q.front.key
	q.front = q.front.next

	if q.front == nil {
		q.rear = nil
	}

	q.size--
	return value
}

func (q *imgQueue) clear() {
	q.front = nil
	q.rear = nil
	q.size = 0
}
