//go:generate go run -tags generate gen.go

package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	parser "vimacheater/pkg/parser"
)

type UiItems struct {
	sync.Mutex
	Items []parser.Item
}

func (u *UiItems) GetItems() string {
	u.Lock()
	defer u.Unlock()

	b, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(b))
	return string(b)
}

func (u *UiItems) UpdateItems(str string) {
	u.Lock()
	defer u.Unlock()

	var ints []int
	err := json.Unmarshal([]byte(str), &ints)
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range ints {
		u.Items[i].ModifiedCount = v
	}
	log.Println("Items updated.")
}

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type Counter struct {
	sync.Mutex
	count int
}

func (c *Counter) Add(n int) {
	c.Lock()
	defer c.Unlock()
	c.count = c.count + n
	fmt.Println("add +1")
}

func (c *Counter) Value() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}
