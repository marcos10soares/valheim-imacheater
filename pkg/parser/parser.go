package parser

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"vimacheater/pkg/utils"
)

var db_items []DbItem

type Item struct {
	Name          string
	PayloadIndex  int
	Payload       []byte
	OriginalCount int
	ModifiedCount int
	MaxCount      int
	Lvl           int
	ModifiedLvl   int
	LvlIndex      int
	ToModify      bool
}

type DbItem struct {
	Name         string
	Type         string
	InternalID   string
	Usage        string
	URL          string
	ImageURL     string
	Weight       string
	Stack        string
	Teleportable string
}

// still not sure what to make of it
// header has 72 bytes, where the first 4 is the filesize in little endian, minus the header size (72 bytes)
type Header struct {
	Filesize uint32 // filesize in little endian, minus some offset (minus the header size (72 bytes)?)
	Field2   uint32 // unknown
	Field3   uint32 // unknown
	Field4   uint32 // unknown
	Field5   uint32 // unknown
	Field6   uint32 // unknown
	Field7   uint32 // unknown
	Field8   uint32 // unknown
	Field9   uint32 // unknown
}

type ItemControlData struct {
	PayloadIndex              int
	PayloadIndexStartingOnQty int
	Payload                   []byte
	OriginalCount             int
	ModifiedCount             int
	MaxCount                  int
	ToModify                  bool
}

type CharData struct {
	CharnameLen              uint8         // 1 byte
	Charname                 string        // variable len
	Unknown1                 []byte        // 35 bytes
	EquipedPowerLen          uint8         // 1 byte
	EquipedPowerLenIndex     uint32        // reference for accessing power name
	EquipedPower             string        // variablen len
	NewPowerToEquip          string        // string to save new power to equip
	ModifyPowerEquiped       bool          // flag to check if have to modify power
	PowerCooldownIndex       uint32        // reference for accessing cooldown
	PowerCooldown            uint32        // 4 bytes
	Unknown2                 uint32        // 4 bytes
	NumberOfItemsInInventory uint32        // 4 bytes
	ItemSection              []ItemPayload // variable
}

type ItemPayload struct {
	ControlData       ItemControlData
	ItemNameLen       uint8  // 1 byte
	ItemName          string // variable len
	ItemCount         uint32 // 4 bytes
	Unknown3          uint32 // 4 bytes
	XcoordInInventory uint32 // 4 bytes
	YcoordInInventory uint32 // 4 bytes
	Lvl               uint8  // 1 byte
	Unknown4          []byte // 15 bytes
	HasOwner          bool   // optional
	OwnerName         string // optional
	LvlIndex          int    // NOT PART OF PAYLOAD, just a position of the lvl byte
}

// func GetItemName(player_data string, start_byte_i int) string {
// 	var item_name string
// 	for {
// 		if player_data[start_byte_i] == 0x0 {
// 			break
// 		}
// 		item_name += string(player_data[start_byte_i])
// 		start_byte_i--
// 	}
// 	if len(item_name) < 1 {
// 		return ""
// 	}
// 	item_name = utils.ReverseString(item_name)
// 	return utils.CleanString(strings.Replace(strings.Trim(strings.TrimSpace(item_name[1:]), "\n"), string(0xc), "", -1))
// }

func CheckIfItemPayloadHasExtraByte(player_data string, index int) bool {
	extra_byte_payload_marker := string([]byte{0x6d, 0x1e, 0xf7, 0xd1})

	if extra_byte_payload_marker == player_data[index+8:index+12] {
		return true
	}
	return false
}

func FindAllOccurrences(data []byte, searches []string) map[string][]int {
	results := make(map[string][]int, 0)

	for _, search := range searches {
		index := len(data)
		tmp := data
		for true {
			match := bytes.LastIndex(tmp[0:index], []byte(search))
			if match == -1 {
				break
			} else {
				index = match
				results[search] = append(results[search], match)
			}
		}
	}

	return results
}

func ModifyItemData(full_data []byte, items []Item) []byte {
	for _, item := range items {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(item.ModifiedCount))
		for i, b := range buf {
			full_data[int(item.PayloadIndex)+i] = b
		}
		// fmt.Printf("Name: %s | current_lvl: %20x | new_lvl:%20x \n", item.Name, full_data[int(item.PayloadIndex)+item.LvlIndex], item.Lvl)

		// update item lvl
		full_data[int(item.PayloadIndex)+item.LvlIndex] = byte(item.ModifiedLvl)
	}
	return full_data
}

func ModifyPowerCooldownData(full_data []byte, index uint32) []byte {
	buf := []byte{0, 0, 0, 0}
	for i, b := range buf {
		full_data[int(index)+i] = b
	}
	return full_data
}

func LoadItems(character string, path string, filename string) (charData CharData, FileData []byte) {
	// selected character
	character_path := path + filename
	log.Println("character_path:", character_path)

	// read all data from selected character
	full_data := utils.ReadFileFromPath(character_path)

	log.Println("filesize:", len(full_data))

	// parse file
	charData = ParseFileNewMethod(character, full_data)

	return charData, full_data
}

func ParseFileNewMethod(charname string, data []byte) CharData {
	// start_index := 0
	// size := len(data)
	// fmt.Println("data size: ", size)
	// if size > 8000000 {
	// 	start_index = 8389395
	// } else if size > 4000000 {
	// 	start_index = 4194408
	// }

	// this could be removed if there is a way to always know where the character data starts
	// instead of parsing the file looking for the name
	// like the commented code above
	charname = utils.MakeTitle(charname)

	// look for where character region starts by searching for the character name and some flags
	start_index := 0
	for {
		tmp_index := strings.Index(string(data[start_index:]), charname)

		start_index += tmp_index
		// log.Println("start_index:", start_index)

		// log.Printf("name: %s | % 20x", charname, charname)
		// log.Println("len:", data[start_index-1])
		// log.Printf("% 20x", data[start_index:start_index+40])
		// log.Printf("% 20x", data[start_index-9:start_index-1])
		// fmt.Println("nbr: ", (int(data[start_index-1]) + 16))
		// log.Printf("% 20x \n", data[start_index+int(data[start_index-1])+14])

		// it seems that all characters have this byte with the value of 0x18, so I'm using it to make sure i'm in the character data region
		byte_0x18 := data[start_index+int(data[start_index-1])+14]

		if start_index <= 0 {
			return CharData{}
		}

		// if byte_0x18 doesn't work, try to use the previous 4 bytes before the name
		// && string(data[start_index-5:start_index-1]) == string([]byte{0, 0, 0, 0})
		// checks if the leading byte of the name is the name size and if 15th byte has the value of 0x18
		if int(data[start_index-1]) == len(charname) && byte_0x18 == 0x18 {
			start_index--
			break
		}
		start_index += 1

		time.Sleep(time.Second * 5)
	}

	char_section_region := data[start_index:]

	charname_len := char_section_region[0]
	charname_from_file := string(char_section_region[1 : charname_len+1])
	unknown1 := char_section_region[charname_len+1 : charname_len+1+35]
	equiped_power_len := char_section_region[charname_len+1+35]
	equiped_power := string(char_section_region[charname_len+1+35+1 : charname_len+1+35+1+equiped_power_len])

	next_i := charname_len + 1 + 35 + 1 + equiped_power_len
	i := next_i
	power_cooldown_index := uint32(start_index) + uint32(i)
	next_i = i + 4
	power_cooldown := binary.LittleEndian.Uint32(char_section_region[i:next_i])

	unknown2 := char_section_region[next_i : next_i+4]
	number_of_items_in_inventory := binary.LittleEndian.Uint32(char_section_region[next_i+4 : next_i+4+4])

	// for debugging
	// fmt.Printf("charname_len: %d, % 20x\n", charname_len, char_section_region[0])
	// fmt.Printf("charname: %s, % 20x\n", charname_from_file, char_section_region[1:charname_len+1])
	// fmt.Printf("unknown1: % 20x\n", unknown1)
	// fmt.Printf("equiped_power_len: %d\n", char_section_region[charname_len+1+35])
	// fmt.Printf("equiped_power: %s, % 20x\n", equiped_power, char_section_region[charname_len+1+35+1:charname_len+1+35+1+equiped_power_len])
	// power_cooldown_bytes := char_section_region[i:next_i]
	// fmt.Printf("power_cooldown: %d, % 20x\n", power_cooldown, power_cooldown_bytes)
	// // Max cooldown 1200s
	// fmt.Println("power_cooldown: microseconds", power_cooldown/1000000)
	// fmt.Printf("unknown2: %d, % 20x\n", unknown2, char_section_region[next_i:next_i+4])
	// fmt.Printf("number_of_items_in_inventory: %d, % 20x\n", number_of_items_in_inventory, char_section_region[next_i+4:next_i+4+4])

	char_data := &CharData{
		CharnameLen:              uint8(charname_len),
		Charname:                 string(charname_from_file),
		Unknown1:                 unknown1,
		EquipedPowerLenIndex:     uint32(start_index) + uint32(charname_len+1+35),
		EquipedPowerLen:          uint8(equiped_power_len),
		EquipedPower:             equiped_power,
		PowerCooldownIndex:       power_cooldown_index,
		PowerCooldown:            power_cooldown,
		Unknown2:                 binary.LittleEndian.Uint32(unknown2),
		NumberOfItemsInInventory: number_of_items_in_inventory,
		NewPowerToEquip:          "",
		ModifyPowerEquiped:       false,
		// ItemSection:,
	}

	start_of_items_section := start_index + int(next_i+4+4)
	items, _ := GetItemsNewMethod(int(number_of_items_in_inventory), start_of_items_section, data)
	char_data.ItemSection = items

	return *char_data

}

func GetItemsNewMethod(number_of_items_in_inventory int, start_of_items_section int, data []byte) (items []ItemPayload, end_of_item_section_index int) {
	item_data := data[start_of_items_section:]
	byte_offset := 0

	// write as a table
	w := tabwriter.NewWriter(os.Stdout, 10, 2, 1, ' ', 0)
	for i := 0; i < int(number_of_items_in_inventory); i++ {
		offset := byte_offset
		has_owner := false
		owner := ""

		item_name_len := uint8(item_data[offset])
		if string(item_data[offset+int(item_name_len)+1+2:offset+int(item_name_len)+1+2+2]) != string([]byte{0, 0}) {
			// using zero bytes in quantity to check if its quantity or string, if it's string it has owner
			has_owner = true
		}

		item_name := string(item_data[offset+1 : offset+int(item_name_len)+1])
		if has_owner {
			has_owner = true
			second_name_len := uint8(item_data[offset+int(item_name_len)+1])
			name_bytes := item_data[offset+int(item_name_len)+2 : offset+int(item_name_len+second_name_len)+2]
			second_item_name := string(name_bytes)
			item_name_len += second_name_len + 1

			owner = item_name
			item_name = second_item_name
		}

		item_count := binary.LittleEndian.Uint32(item_data[offset+int(item_name_len)+1 : offset+int(item_name_len)+1+4])
		unknown3 := item_data[offset+int(item_name_len)+1+4 : offset+int(item_name_len)+1+4+4]
		x_coord_in_inventory := binary.LittleEndian.Uint32(item_data[offset+int(item_name_len)+1+4+4 : offset+int(item_name_len)+1+4+4+4])
		y_coord_in_inventory := binary.LittleEndian.Uint32(item_data[offset+int(item_name_len)+1+4+4+4 : offset+int(item_name_len)+1+4+4+4+4])
		lvl_index := 4 + 4 + 4 + 5
		lvl := item_data[offset+int(item_name_len)+1+lvl_index]
		unknown4 := item_data[offset+int(item_name_len)+1+4+4+4+4+1 : offset+int(item_name_len)+1+4+4+4+4+17]

		// for debugging
		// fmt.Println("item name len:", item_name_len)
		// fmt.Println("item name:", item_name)
		// fmt.Println("item count:", item_count)
		// fmt.Println("unknown3:", unknown3)
		// fmt.Println("x coord:", x_coord_in_inventory)
		// fmt.Println("y coord::", y_coord_in_inventory)
		// fmt.Println("unknown4:", unknown4)

		full_payload := item_data[byte_offset : byte_offset+34+int(item_name_len)]

		// debugging
		// format string
		s_out := fmt.Sprintf("| %s\t| Count: %d\t| lvl: %d\t| % 20x \t|", item_name, item_count, lvl, full_payload)
		// add to table
		fmt.Fprintln(w, s_out)

		item_payload := &ItemPayload{
			ItemNameLen:       item_name_len,
			ItemName:          item_name,
			ItemCount:         item_count,
			Unknown3:          binary.LittleEndian.Uint32(unknown3),
			XcoordInInventory: x_coord_in_inventory,
			YcoordInInventory: y_coord_in_inventory,
			Lvl:               lvl,
			LvlIndex:          lvl_index,
			Unknown4:          unknown4,
			HasOwner:          has_owner,
			OwnerName:         owner,
			ControlData: ItemControlData{
				PayloadIndex:              start_of_items_section + byte_offset,
				PayloadIndexStartingOnQty: start_of_items_section + byte_offset + int(item_name_len) + 1,
				Payload:                   full_payload,
				OriginalCount:             int(item_count),
				// ModifiedCount int
				// MaxCount      int
				// ToModify      bool
			},
		}

		byte_offset = byte_offset + int(item_name_len) + 1 + 4 + 4 + 4 + 4 + 17
		items = append(items, *item_payload)
	}
	w.Flush()

	return items, byte_offset
}

func LoadDbItems(bytes []byte) {
	if bytes == nil {
		log.Panic("could not load json file")
	}
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'db_items' which we defined on the beginning of the file
	json.Unmarshal(bytes, &db_items)
}

func GetItemFromDbItemWithName(name string) DbItem {
	for _, item := range db_items {
		if item.InternalID == name {
			// fmt.Println(item.Name)
			return item
		}
	}
	return DbItem{}
}
