//go:generate go run -tags generate gen.go

package gui

import (
	"encoding/binary"
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

// UIData is the data structure used to passa data between frontend and backend
type UIData struct {
	sync.Mutex
	Items     []uiItem
	FileData  []byte
	Character string
	CharData  parser.CharData
}

type uiItem struct {
	DbItem   parser.DbItem
	FileItem parser.Item
}

// NewUIData returns a UIData initialized
func NewUIData() *UIData {
	return &UIData{
		Items: []uiItem{},
	}
}

// GetItems fetches items from character
func (u *UIData) GetItems(character string) string {
	u.Lock()
	defer u.Unlock()

	// var loadedItems []parser.Item
	var loadedItems []uiItem
	charData, fileData := parser.LoadItems(character, utils.CharactersFolderPath, character+".fch")

	// create []Item because it's the struct already in use, refactor this later
	for _, item := range charData.ItemSection {
		loadedItems = append(loadedItems, uiItem{
			DbItem: parser.GetItemFromDbItemWithName(item.ItemName),
			FileItem: parser.Item{
				Name:          item.ItemName,
				PayloadIndex:  item.ControlData.PayloadIndexStartingOnQty,
				Payload:       item.ControlData.Payload,
				OriginalCount: item.ControlData.OriginalCount,
				ModifiedCount: item.ControlData.ModifiedCount,
				MaxCount:      item.ControlData.MaxCount,
				Lvl:           int(item.Lvl),
				LvlIndex:      item.LvlIndex,
				ToModify:      item.ControlData.ToModify,
			},
		})
	}
	u.FileData = fileData
	u.Items = loadedItems
	u.CharData = charData
	u.Character = character

	b, err := json.Marshal(u.Items)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(b)
}

// GetPowers gets powers from character
func (u *UIData) GetPowers() string {
	u.Lock()
	defer u.Unlock()

	var powers []string
	availablePowers := []string{"Eikthyr", "TheElder", "Bonemass", "Moder", "Yagluth"}

	equipedPower := u.CharData.EquipedPower
	if len(u.CharData.EquipedPower) > 3 {
		equipedPower = equipedPower[3:]
	}

	powers = append(powers, equipedPower)
	for _, power := range availablePowers {
		if !utils.StringSliceCheckIfContains(powers, power) {
			powers = append(powers, power)
		}
	}

	b, err := json.Marshal(powers)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}

// SetInfiniteLife gives player infinite life
func (u *UIData) SetInfiniteLife() {
	u.Lock()
	defer u.Unlock()

	buf := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for i, b := range buf {
		u.FileData[int(u.CharData.EquipedPowerLenIndex)-17+i] = b
	}

	log.Println("Set infinite life.")
}

// RemoveInfiniteLife removes infinite life ability from the player
func (u *UIData) RemoveInfiniteLife() {
	u.Lock()
	defer u.Unlock()

	buf := []byte{0x00, 0x00, 0xc8, 0x41}
	for i, b := range buf {
		u.FileData[int(u.CharData.EquipedPowerLenIndex)-13+i] = b
	}

	log.Println("Remove infinite life.")
}

// ResetPowerCooldown resets player power cooldown
func (u *UIData) ResetPowerCooldown() {
	u.Lock()
	defer u.Unlock()

	u.FileData = parser.ModifyPowerCooldownData(u.FileData, u.CharData.PowerCooldownIndex)
	log.Println("Power Cooldown Reset.")
}

// UpdatePower replaces the player equipped power with the given one
func (u *UIData) UpdatePower(power string) {
	u.Lock()
	defer u.Unlock()

	power = "GP_" + power

	// check if already had a power, otherwise cannot do the change
	if len(u.CharData.EquipedPower) != 0 {
		bufFileData := u.FileData

		// get difference in power length
		fmt.Println("equiped power: ", int(u.CharData.EquipedPowerLen))
		fmt.Println("new power len: ", int(len(power)))
		byteDifference := int(u.CharData.EquipedPowerLen) - int(len(power))
		fmt.Println("byte_difference:", byteDifference)

		// get current file size
		currentFileSizeBytes := binary.LittleEndian.Uint32(u.FileData[:4])

		// add or subtract difference and create new file size
		newFileSize := int(currentFileSizeBytes) - byteDifference

		u.CharData.PowerCooldownIndex = uint32(int(u.CharData.PowerCooldownIndex) - byteDifference)

		// update fileData with new filesize in little endian
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(newFileSize))
		for i, b := range buf {
			u.FileData[i] = b
		}

		bufFileData[int(u.CharData.EquipedPowerLenIndex)] = byte(len(power))

		// this byte might be a uint32 - be careful - to check in the future
		weirdAloneByteNew := byte(int(bufFileData[int(u.CharData.EquipedPowerLenIndex)-25]) - byteDifference)
		bufFileData[int(u.CharData.EquipedPowerLenIndex)-25] = weirdAloneByteNew

		var newBuf []byte
		newBuf = append(newBuf, bufFileData[:int(u.CharData.EquipedPowerLenIndex)+1]...)
		newBuf = append(newBuf, []byte(power)...)
		newBuf = append(newBuf, bufFileData[int(u.CharData.EquipedPowerLenIndex)+1+int(u.CharData.EquipedPowerLen):]...)
		u.FileData = newBuf
	}
}

// UpdateItems updates character items quantity and level
// receives two json arrays as strings
func (u *UIData) UpdateItems(qtyString string, lvlString string) {
	u.Lock()
	defer u.Unlock()

	qty := parseIntSlice(qtyString)
	lvl := parseIntSlice(lvlString)

	for i, v := range qty {
		u.Items[i].FileItem.ModifiedCount = v
		u.Items[i].FileItem.ModifiedLvl = lvl[i]
	}
	log.Println("Items updated.")

	var items []parser.Item
	for _, item := range u.Items {
		items = append(items, item.FileItem)
	}

	u.FileData = parser.ModifyItemData(u.FileData, items)
}

// SaveData commits current changes to character file
func (u *UIData) SaveData() {
	utils.WriteOutputFile(u.FileData, u.Character)
	log.Println("Items saved.")
}

// GetChars fetches available character files
func (u *UIData) GetChars() []string {
	u.Lock()
	defer u.Unlock()

	var path string

	user, err := utils.GetCurrentUser()
	if err != nil {
		log.Fatal("could not get current user: ", err)
	}

	if runtime.GOOS == "windows" {
		path = user.HomeDir + utils.WinPath
	} else if runtime.GOOS == "darwin" {
		// mac - for debugging
		path = utils.MacPath
	} else if runtime.GOOS == "linux" {
		path = user.HomeDir + utils.LinuxPath
	}
	utils.CharactersFolderPath = path

	return utils.GetAllAvailableCharacters(path)
}

func parseIntSlice(str string) []int {
	var ints []int
	err := json.Unmarshal([]byte(str), &ints)
	if err != nil {
		log.Fatal("Error updating items:", err)
	}
	return ints
}
