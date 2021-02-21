package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"vimacheater/pkg/utils"
)

type Item struct {
	Name          string
	PayloadIndex  int
	Payload       []byte
	OriginalCount int
	ModifiedCount int
	MaxCount      int
	ToModify      bool
}

// still not sure what to make of it
type Header struct {
	Field1  uint32
	Field2  uint32
	Field3  uint32
	Field4  uint32
	Field5  uint32
	Field6  uint32
	Field7  uint32
	Field8  uint32
	Field9  uint32
	Field10 uint32
	Field11 uint32
	Field12 uint32
	Field13 uint32
	Field14 uint32
	Field15 uint32
	Field16 uint32
	Field17 uint32
	Field18 uint32
	Field19 uint32
	Field20 uint32
	Field21 uint32
	Field22 uint32
	Field23 uint32
	Field24 uint32
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
	EquipedPower             string        // variablen len
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
	Unknown4          []byte // 16 bytes
	HasOwner          bool   // optional
	OwnerName         string // optional
}

func GetItemName(player_data string, start_byte_i int) string {
	var item_name string
	for {
		if player_data[start_byte_i] == 0x0 {
			break
		}
		item_name += string(player_data[start_byte_i])
		start_byte_i--
	}
	if len(item_name) < 1 {
		return ""
	}
	item_name = utils.ReverseString(item_name)
	return utils.CleanString(strings.Replace(strings.Trim(strings.TrimSpace(item_name[1:]), "\n"), string(0xc), "", -1))
}

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

func CleanItemMatches(full_data []byte, player_data_string string, i int, matches []int) []int {
	var clean_matches []int
	for _, match := range matches {
		// items payload length is variable, this checks the item payload size
		hasExtraByte := CheckIfItemPayloadHasExtraByte(player_data_string, match)
		start_byte_i := match - 17
		item_payload_size := 34
		if hasExtraByte {
			item_payload_size += 1
		}

		// get payload limits index
		item_payload_start_byte := (i + start_byte_i)
		end_byte_i := item_payload_start_byte + item_payload_size

		// get payload
		item_payload := []byte(full_data[item_payload_start_byte:end_byte_i])

		//verify payload
		if (string(item_payload[1:4]) != string([]byte{0, 0, 0})) || (string(item_payload[29:33]) != string([]byte{0, 0, 0, 0})) {
			continue
		}

		clean_matches = append(clean_matches, match)
	}
	return clean_matches
}

func GetItems(matches []int, full_data []byte, player_data_string string, i int, character string) []Item {
	totalItems := []Item{}

	// file for debug only FOR DEBUG
	// items_log_file_debug, err := os.Create(utils.Bckp_folder + character + "_itemslog_" + utils.GetTimestampString() + ".txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer items_log_file_debug.Close()

	// write as a table
	// w_debug := tabwriter.NewWriter(items_log_file_debug, 10, 2, 1, ' ', 0)
	//// END OF DEBUG

	// write as a table
	// w := tabwriter.NewWriter(os.Stdout, 10, 2, 1, ' ', 0)
	for _, match := range matches {
		// items payload length is variable, this checks the item payload size
		hasExtraByte := CheckIfItemPayloadHasExtraByte(player_data_string, match)
		start_byte_i := match - 17
		item_payload_size := 34
		if hasExtraByte {
			item_payload_size += 1
		}

		// get payload limits index
		item_payload_start_byte := (i + start_byte_i)
		end_byte_i := item_payload_start_byte + item_payload_size

		// get item name (this string is modified for printing reasons, do not use for changing output file)
		item_name := GetItemName(player_data_string, start_byte_i)

		// get payload
		item_payload := []byte(full_data[item_payload_start_byte:end_byte_i])

		// save to items
		totalItems = append(totalItems, Item{
			Name:          item_name,
			PayloadIndex:  item_payload_start_byte,
			Payload:       item_payload,
			OriginalCount: int(item_payload[0]),
		})

		// format string
		// s_out := fmt.Sprintf("| %s\t| Count: %d\t| % 20x \t|", item_name, item_payload[0], item_payload)

		// // add to table
		// fmt.Fprintln(w, s_out)

		// // add TO DEBUG table
		// fmt.Fprintln(w_debug, s_out)

		// fmt.Printf("name: %s :\t % 20x \t| len: %d, extra byte: %v\n", getItemName(player_data_string, start_byte_i), item_payload, len(item_payload), hasExtraByte)
		// fmt.Printf("%+q", patterns)
	}
	// print table
	// w.Flush()

	// print to DEBUG file
	// w_debug.Flush()

	return totalItems
}

func ModifyItemData(full_data []byte, items []Item) []byte {
	for _, item := range items {
		full_data[int(item.PayloadIndex)] = byte(item.ModifiedCount)
	}
	return full_data
}

func LoadItems(character string, path string, filename string) (loadedItems []Item, FileData []byte) {
	// path := utils.CharactersFolder

	// user, err := utils.GetCurrentUser()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if runtime.GOOS == "windows" { // production
	// 	path = user.HomeDir + WinPath
	// 	fmt.Println(path)
	// } else { // mac - for debugging
	// 	path = MacPath
	// }

	// selected character
	character_path := path + filename

	// read all data from selected character
	full_data := utils.ReadFileFromPath(character_path)

	fmt.Println(len(full_data))

	// parse file
	charData := ParseFileNewMethod(character, full_data)

	// charname := utils.MakeTitle(character)

	// full_string := string(full_data)
	// i := strings.Index(full_string, charname)

	// fmt.Println(i)

	// // parse header, still not sure of structure format and meaning, probably date and time somewhere?
	// header := Header{}
	// buffer := bytes.NewBuffer(full_data[:96])
	// err := binary.Read(buffer, binary.LittleEndian, &header)
	// if err != nil {
	// 	log.Fatal("binary.Read failed", err)
	// }
	// fmt.Printf("Header data:\n%+v\n", header)

	// player_data_string := full_string[i:]

	// // pattern to look for in items
	// byte_pattern := []byte{1, 0, 0, 0, 0, 0, 0, 0}
	// string_pattern := string(byte_pattern)

	// // create a slice of patterns to look for
	// patterns := make([]string, 1)
	// // inserts the pattern to look for
	// patterns[0] = string_pattern

	// // finds all indexes where the pattern occurs, result is a map because there can be more patterns to look for
	// result := FindAllOccurrences([]byte(player_data_string), patterns)

	// // get the match indexes
	// matches := result[string_pattern]

	// // reverse order of matches
	// matches = utils.ReverseIntSlice(matches)

	// // clean matches by verifying some extra patterns on each item
	// matches = CleanItemMatches(full_data, player_data_string, i, matches)

	// // fmt.Println("Items found: ", len(matches))

	// totalItems := GetItems(matches, full_data, player_data_string, i, character)
	// // fmt.Printf("Items: %v\n", totalItems)

	for _, item := range charData.ItemSection {
		// fmt.Println(item)
		loadedItems = append(loadedItems, Item{
			Name:          item.ItemName,
			PayloadIndex:  item.ControlData.PayloadIndexStartingOnQty,
			Payload:       item.ControlData.Payload,
			OriginalCount: item.ControlData.OriginalCount,
			ModifiedCount: item.ControlData.ModifiedCount,
			MaxCount:      item.ControlData.MaxCount,
			ToModify:      item.ControlData.ToModify,
		})
		// fmt.Println("payload_index:", item.ControlData.PayloadIndexStartingOnQty)
		// fmt.Println("item_count:", item.ControlData.OriginalCount)
		// fmt.Println("payload:", full_data[item.ControlData.PayloadIndexStartingOnQty:item.ControlData.PayloadIndexStartingOnQty+40])
	}

	// panic("fuck")

	return loadedItems, full_data
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

	charname = utils.MakeTitle(charname)
	fmt.Println(charname)
	start_index := strings.Index(string(data), charname) - 1

	char_section_region := data[start_index:]

	charname_len := char_section_region[0]
	charname_from_file := string(char_section_region[1 : charname_len+1])
	unknown1 := char_section_region[charname_len+1 : charname_len+1+35]
	equiped_power_len := char_section_region[charname_len+1+35]
	equiped_power := string(char_section_region[charname_len+1+35+1 : charname_len+1+35+1+equiped_power_len])

	next_i := charname_len + 1 + 35 + 1 + equiped_power_len
	i := next_i
	next_i = i + 4
	power_cooldown := binary.LittleEndian.Uint32(char_section_region[i:next_i])

	unknown2 := char_section_region[next_i : next_i+4]
	number_of_items_in_inventory := binary.LittleEndian.Uint32(char_section_region[next_i+4 : next_i+4+4])

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
		EquipedPowerLen:          uint8(equiped_power_len),
		EquipedPower:             equiped_power,
		PowerCooldown:            power_cooldown,
		Unknown2:                 binary.LittleEndian.Uint32(unknown2),
		NumberOfItemsInInventory: number_of_items_in_inventory,
		// ItemSection:,
	}

	// fmt.Printf("% 20x\n", char_section_region[:100])
	// fmt.Println(char_data)

	start_of_items_section := start_index + int(next_i+4+4)
	items, _ := GetItemsNewMethod(int(number_of_items_in_inventory), start_of_items_section, data)
	// fmt.Println(items)

	char_data.ItemSection = items

	return *char_data

}

func GetItemsNewMethod(number_of_items_in_inventory int, start_of_items_section int, data []byte) (items []ItemPayload, end_of_item_section_index int) {
	item_data := data[start_of_items_section:]
	byte_offset := 0

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
			// fmt.Println("HAS OWNER")
			has_owner = true
			second_name_len := uint8(item_data[offset+int(item_name_len)+1])
			name_bytes := item_data[offset+int(item_name_len)+2 : offset+int(item_name_len+second_name_len)+2]
			second_item_name := string(name_bytes)
			item_name_len += second_name_len + 1

			owner = item_name
			item_name = second_item_name
		}

		// last_byte := item_data[byte_offset+int(item_name_len)+33]
		// fmt.Println("last_byte", last_byte) // just for debugging

		item_count := binary.LittleEndian.Uint32(item_data[offset+int(item_name_len)+1 : offset+int(item_name_len)+1+4])
		unknown3 := item_data[offset+int(item_name_len)+1+4 : offset+int(item_name_len)+1+4+4]
		x_coord_in_inventory := binary.LittleEndian.Uint32(item_data[offset+int(item_name_len)+1+4+4 : offset+int(item_name_len)+1+4+4+4])
		y_coord_in_inventory := binary.LittleEndian.Uint32(item_data[offset+int(item_name_len)+1+4+4+4 : offset+int(item_name_len)+1+4+4+4+4])
		unknown4 := item_data[offset+int(item_name_len)+1+4+4+4+4 : offset+int(item_name_len)+1+4+4+4+4+17]

		// fmt.Println("item name len:", item_name_len)
		// fmt.Println("item name:", item_name)
		// fmt.Println("item count:", item_count)
		// fmt.Println("unknown3:", unknown3)
		// fmt.Println("x coord:", x_coord_in_inventory)
		// fmt.Println("y coord::", y_coord_in_inventory)
		// fmt.Println("unknown4:", unknown4)

		full_payload := item_data[byte_offset : byte_offset+34+int(item_name_len)]

		// fmt.Println(item_data[byte_offset : byte_offset+34+int(item_name_len)]) // full payload
		// fmt.Println(item_data[byte_offset : byte_offset+44+int(item_name_len)]) // shows a bit of next item
		// fmt.Println("-----")

		item_payload := &ItemPayload{
			ItemNameLen:       item_name_len,
			ItemName:          item_name,
			ItemCount:         item_count,
			Unknown3:          binary.LittleEndian.Uint32(unknown3),
			XcoordInInventory: x_coord_in_inventory,
			YcoordInInventory: y_coord_in_inventory,
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

	return items, byte_offset
}
