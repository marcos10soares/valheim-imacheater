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
type UiData struct {
	sync.Mutex
	Items     []parser.Item
	FileData  []byte
	Character string
	CharData  parser.CharData
}

// type UiChars struct {
// 	sync.Mutex
// 	Chars []string
// }

func (u *UiData) GetItems(character string) string {
	u.Lock()
	defer u.Unlock()

	var loadedItems []parser.Item
	charData, fileData := parser.LoadItems(character, utils.CharactersFolder, character+".fch")

	// create []Item because it's the struct already in use, refactor this later
	for _, item := range charData.ItemSection {
		// fmt.Println(item)
		loadedItems = append(loadedItems, parser.Item{
			Name:          item.ItemName,
			PayloadIndex:  item.ControlData.PayloadIndexStartingOnQty,
			Payload:       item.ControlData.Payload,
			OriginalCount: item.ControlData.OriginalCount,
			ModifiedCount: item.ControlData.ModifiedCount,
			MaxCount:      item.ControlData.MaxCount,
			ToModify:      item.ControlData.ToModify,
		})
	}
	u.FileData = fileData
	u.Items = loadedItems
	u.CharData = charData

	u.Character = character

	b, err := json.Marshal(u.Items)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// fmt.Println("debug getitems: ", string(b))
	return string(b)
}

func (u *UiData) ResetPowerCooldown() {
	u.Lock()
	defer u.Unlock()

	u.FileData = parser.ModifyPowerCooldownData(u.FileData, u.CharData.PowerCooldownIndex)
	log.Println("Power Cooldown Reset.")
}

func (u *UiData) UpdateItems(str string) {
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

func (u *UiData) GetChars() []string {
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
