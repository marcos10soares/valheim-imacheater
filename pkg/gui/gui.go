//go:generate go run -tags generate gen.go

package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"vimacheater/pkg/parser"
	"vimacheater/pkg/utils"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type UiItems struct {
	sync.Mutex
	Items     []parser.Item
	FileData  []byte
	Character string
}

// type UiChars struct {
// 	sync.Mutex
// 	Chars []string
// }

func (u *UiItems) GetItems(character string) string {
	u.Lock()
	defer u.Unlock()

	u.Items, u.FileData = parser.LoadItems(character)
	u.Character = character

	b, err := json.Marshal(u.Items)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// fmt.Println("debug getitems: ", string(b))
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

	full_data := parser.ModifyItemData(u.FileData, u.Items)

	utils.WriteOutputFile(full_data, u.Character)
	log.Println("Items saved.")
}

func (u *UiItems) GetChars() []string {
	u.Lock()
	defer u.Unlock()

	var path string

	user, err := utils.GetCurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	if runtime.GOOS == "windows" { // production
		path = user.HomeDir + utils.WinPath
		// fmt.Println(path)

	} else { // mac - for debugging
		path = "files/" // bjørn
	}
	utils.CharactersFolder = path

	return utils.GetAllAvailableCharacters(path)
}
