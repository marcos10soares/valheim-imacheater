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
type UiData struct {
	sync.Mutex
	Items     []UiItem
	FileData  []byte
	Character string
	CharData  parser.CharData
}

type UiItem struct {
	DbItem   parser.DbItem
	FileItem parser.Item
}

// type UiChars struct {
// 	sync.Mutex
// 	Chars []string
// }

func (u *UiData) GetItems(character string) string {
	u.Lock()
	defer u.Unlock()

	// var loadedItems []parser.Item
	var loadedItems []UiItem
	charData, fileData := parser.LoadItems(character, utils.CharactersFolder, character+".fch")

	// create []Item because it's the struct already in use, refactor this later
	for _, item := range charData.ItemSection {
		// fmt.Println(item)
		loadedItems = append(loadedItems, UiItem{
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
	// fmt.Println("debug getitems: ", string(b))
	return string(b)
}

func (u *UiData) GetPowers() string {
	u.Lock()
	defer u.Unlock()

	var powers []string

	availablePowers := []string{"Eikthyr", "TheElder", "Bonemass", "Moder", "Yagluth"}

	equiped_power := u.CharData.EquipedPower
	if len(u.CharData.EquipedPower) > 3 {
		equiped_power = equiped_power[3:]
	}

	powers = append(powers, equiped_power)
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
	// fmt.Println("debug getitems: ", string(b))
	return string(b)
}

func (u *UiData) SetInfiniteLife() {
	u.Lock()
	defer u.Unlock()

	buf := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for i, b := range buf {
		u.FileData[int(u.CharData.EquipedPowerLenIndex)-17+i] = b
	}

	log.Println("Set infinite life.")
}

func (u *UiData) RemoveInfiniteLife() {
	u.Lock()
	defer u.Unlock()

	buf := []byte{0x00, 0x00, 0xc8, 0x41}
	for i, b := range buf {
		u.FileData[int(u.CharData.EquipedPowerLenIndex)-13+i] = b
	}

	log.Println("Remove infinite life.")
}

func (u *UiData) ResetPowerCooldown() {
	u.Lock()
	defer u.Unlock()

	u.FileData = parser.ModifyPowerCooldownData(u.FileData, u.CharData.PowerCooldownIndex)
	log.Println("Power Cooldown Reset.")
}

func (u *UiData) UpdatePower(power string) {
	u.Lock()
	defer u.Unlock()

	power = "GP_" + power

	// check if already had a power, otherwise cannot do the change
	if len(u.CharData.EquipedPower) != 0 {
		buf_FileData := u.FileData

		// get difference in power length
		fmt.Println("equiped power: ", int(u.CharData.EquipedPowerLen))
		fmt.Println("new power len: ", int(len(power)))
		byte_difference := int(u.CharData.EquipedPowerLen) - int(len(power))
		fmt.Println("byte_difference:", byte_difference)
		// get current file size
		current_file_size_bytes := binary.LittleEndian.Uint32(u.FileData[:4])
		// add or subtract difference and create new file size
		new_file_size := int(current_file_size_bytes) - byte_difference

		u.CharData.PowerCooldownIndex = uint32(int(u.CharData.PowerCooldownIndex) - byte_difference)

		// update fileData with new filesize in little endian
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(new_file_size))
		for i, b := range buf {
			u.FileData[i] = b
		}

		buf_FileData[int(u.CharData.EquipedPowerLenIndex)] = byte(len(power))

		// this byte might be a uint32 - be careful - to check in the future
		weird_alone_byte_new := byte(int(buf_FileData[int(u.CharData.EquipedPowerLenIndex)-25]) - byte_difference)
		buf_FileData[int(u.CharData.EquipedPowerLenIndex)-25] = weird_alone_byte_new

		var new_buf []byte
		new_buf = append(new_buf, buf_FileData[:int(u.CharData.EquipedPowerLenIndex)+1]...)
		new_buf = append(new_buf, []byte(power)...)
		new_buf = append(new_buf, buf_FileData[int(u.CharData.EquipedPowerLenIndex)+1+int(u.CharData.EquipedPowerLen):]...)
		u.FileData = new_buf
	}

}

func parseIntSlice(str string) []int {
	var ints []int
	err := json.Unmarshal([]byte(str), &ints)
	if err != nil {
		log.Fatal("Error updating items:", err)
	}
	return ints
}

func (u *UiData) UpdateItems(qty_str string, lvl_str string) {
	u.Lock()
	defer u.Unlock()

	qty_ints := parseIntSlice(qty_str)
	lvl_ints := parseIntSlice(lvl_str)

	for i, v := range qty_ints {
		u.Items[i].FileItem.ModifiedCount = v
		u.Items[i].FileItem.ModifiedLvl = lvl_ints[i]
	}
	log.Println("Items updated.")

	var items []parser.Item
	for _, item := range u.Items {
		items = append(items, item.FileItem)
	}

	u.FileData = parser.ModifyItemData(u.FileData, items)
}

func (u *UiData) SaveData() {
	utils.WriteOutputFile(u.FileData, u.Character)
	log.Println("Items saved.")
}

func (u *UiData) GetChars() []string {
	u.Lock()
	defer u.Unlock()

	var path string

	user, err := utils.GetCurrentUser()
	if err != nil {
		log.Fatal("could not get current user: ", err)
	}

	if runtime.GOOS == "windows" { // production
		path = user.HomeDir + utils.WinPath
	} else if runtime.GOOS == "darwin" { // mac - for debugging
		path = utils.MacPath // bjørn
	} else {
		path = user.HomeDir + utils.LinuxPath
	}
	utils.CharactersFolder = path

	return utils.GetAllAvailableCharacters(path)
}
