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

var dbItems []DbItem

// Item represents a character Item from the character file
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

// DbItem is an item from the database (wiki) of items
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

// Header - still not sure what to make of it
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

type itemControlData struct {
	PayloadIndex              int
	PayloadIndexStartingOnQty int
	Payload                   []byte
	OriginalCount             int
	ModifiedCount             int
	MaxCount                  int
	ToModify                  bool
}

// CharData contains the character data
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

// ItemPayload holds item information
type ItemPayload struct {
	ControlData       itemControlData
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

func checkIfItemPayloadHasExtraByte(playerData string, index int) bool {
	extraBytePayloadMarker := string([]byte{0x6d, 0x1e, 0xf7, 0xd1})

	if extraBytePayloadMarker == playerData[index+8:index+12] {
		return true
	}
	return false
}

func findAllOccurrences(data []byte, searches []string) map[string][]int {
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

// ModifyItemData modifies file data, used to update items
func ModifyItemData(fullData []byte, items []Item) []byte {
	for _, item := range items {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(item.ModifiedCount))
		for i, b := range buf {
			fullData[int(item.PayloadIndex)+i] = b
		}
		// log.Printf("Name: %s | current_lvl: %20x | new_lvl:%20x \n", item.Name, fullData[int(item.PayloadIndex)+item.LvlIndex], item.Lvl)

		// update item lvl
		fullData[int(item.PayloadIndex)+item.LvlIndex] = byte(item.ModifiedLvl)
	}
	return fullData
}

// ModifyPowerCooldownData modifies file data, used to update power cooldown
func ModifyPowerCooldownData(fullData []byte, index uint32) []byte {
	buf := []byte{0, 0, 0, 0}
	for i, b := range buf {
		fullData[int(index)+i] = b
	}
	return fullData
}

// LoadItems loads items from character file
func LoadItems(character string, path string, filename string) (charData CharData, FileData []byte) {
	// selected character
	characterPath := path + filename
	// log.Println("characterPath:", characterPath)

	// read all data from selected character
	fullData := utils.ReadFileFromPath(characterPath)

	// log.Println("filesize:", len(fullData))

	// parse file
	charData = ParseFile(character, fullData)

	return charData, fullData
}

// ParseFile parses the character file and returns the character data object
func ParseFile(charname string, data []byte) CharData {
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
	startIndex := 0
	for {
		tmpIndex := strings.Index(string(data[startIndex:]), charname)

		startIndex += tmpIndex
		// log.Println("start_index:", start_index)

		// log.Printf("name: %s | % 20x", charname, charname)
		// log.Println("len:", data[start_index-1])
		// log.Printf("% 20x", data[start_index:start_index+40])
		// log.Printf("% 20x", data[start_index-9:start_index-1])
		// fmt.Println("nbr: ", (int(data[start_index-1]) + 16))
		// log.Printf("% 20x \n", data[start_index+int(data[start_index-1])+14])

		// it seems that all characters have this byte with the value of 0x18, so I'm using it to make sure i'm in the character data region
		byte0x18 := data[startIndex+int(data[startIndex-1])+14]

		if startIndex <= 0 {
			return CharData{}
		}

		// if byte_0x18 doesn't work, try to use the previous 4 bytes before the name
		// && string(data[start_index-5:start_index-1]) == string([]byte{0, 0, 0, 0})
		// checks if the leading byte of the name is the name size and if 15th byte has the value of 0x18
		if int(data[startIndex-1]) == len(charname) && byte0x18 == 0x18 {
			startIndex--
			break
		}
		startIndex++

		time.Sleep(time.Second * 5)
	}

	charSectionRegion := data[startIndex:]

	charnameLen := charSectionRegion[0]
	charnameFromFile := string(charSectionRegion[1 : charnameLen+1])
	unknown1 := charSectionRegion[charnameLen+1 : charnameLen+1+35]
	equipedPowerLen := charSectionRegion[charnameLen+1+35]
	equipedPower := string(charSectionRegion[charnameLen+1+35+1 : charnameLen+1+35+1+equipedPowerLen])

	nextIndex := charnameLen + 1 + 35 + 1 + equipedPowerLen
	i := nextIndex
	powerCooldownIndex := uint32(startIndex) + uint32(i)
	nextIndex = i + 4
	powerCooldown := binary.LittleEndian.Uint32(charSectionRegion[i:nextIndex])

	unknown2 := charSectionRegion[nextIndex : nextIndex+4]
	numberOfItemsInInventory := binary.LittleEndian.Uint32(charSectionRegion[nextIndex+4 : nextIndex+4+4])

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

	charData := &CharData{
		CharnameLen:              uint8(charnameLen),
		Charname:                 string(charnameFromFile),
		Unknown1:                 unknown1,
		EquipedPowerLenIndex:     uint32(startIndex) + uint32(charnameLen+1+35),
		EquipedPowerLen:          uint8(equipedPowerLen),
		EquipedPower:             equipedPower,
		PowerCooldownIndex:       powerCooldownIndex,
		PowerCooldown:            powerCooldown,
		Unknown2:                 binary.LittleEndian.Uint32(unknown2),
		NumberOfItemsInInventory: numberOfItemsInInventory,
		NewPowerToEquip:          "",
		ModifyPowerEquiped:       false,
		// ItemSection:,
	}

	startOfItemsSection := startIndex + int(nextIndex+4+4)
	items, _ := GetItems(int(numberOfItemsInInventory), startOfItemsSection, data)
	charData.ItemSection = items

	return *charData
}

// GetItems parses items from file raw data
func GetItems(numberOfItemsInInventory int, startOfItemsSection int, data []byte) (items []ItemPayload, endOfItemSectionIndex int) {
	itemData := data[startOfItemsSection:]
	byteOffset := 0

	// write as a table
	w := tabwriter.NewWriter(os.Stdout, 10, 2, 1, ' ', 0)
	for i := 0; i < int(numberOfItemsInInventory); i++ {
		offset := byteOffset
		hasOwner := false
		owner := ""

		itemNameLen := uint8(itemData[offset])
		if string(itemData[offset+int(itemNameLen)+1+2:offset+int(itemNameLen)+1+2+2]) != string([]byte{0, 0}) {
			// using zero bytes in quantity to check if its quantity or string, if it's string it has owner
			hasOwner = true
		}

		itemName := string(itemData[offset+1 : offset+int(itemNameLen)+1])
		if hasOwner {
			hasOwner = true
			secondNameLen := uint8(itemData[offset+int(itemNameLen)+1])
			nameBytes := itemData[offset+int(itemNameLen)+2 : offset+int(itemNameLen+secondNameLen)+2]
			secondItemName := string(nameBytes)
			itemNameLen += secondNameLen + 1

			owner = itemName
			itemName = secondItemName
		}

		itemCount := binary.LittleEndian.Uint32(itemData[offset+int(itemNameLen)+1 : offset+int(itemNameLen)+1+4])
		unknown3 := itemData[offset+int(itemNameLen)+1+4 : offset+int(itemNameLen)+1+4+4]
		xCoordInInventory := binary.LittleEndian.Uint32(itemData[offset+int(itemNameLen)+1+4+4 : offset+int(itemNameLen)+1+4+4+4])
		yCoordInInventory := binary.LittleEndian.Uint32(itemData[offset+int(itemNameLen)+1+4+4+4 : offset+int(itemNameLen)+1+4+4+4+4])
		lvlIndex := 4 + 4 + 4 + 5
		lvl := itemData[offset+int(itemNameLen)+1+lvlIndex]
		unknown4 := itemData[offset+int(itemNameLen)+1+4+4+4+4+1 : offset+int(itemNameLen)+1+4+4+4+4+17]

		// for debugging
		// fmt.Println("item name len:", item_name_len)
		// fmt.Println("item name:", item_name)
		// fmt.Println("item count:", item_count)
		// fmt.Println("unknown3:", unknown3)
		// fmt.Println("x coord:", x_coord_in_inventory)
		// fmt.Println("y coord::", y_coord_in_inventory)
		// fmt.Println("unknown4:", unknown4)

		fullPayload := itemData[byteOffset : byteOffset+34+int(itemNameLen)]

		// debugging
		// format string
		sOut := fmt.Sprintf("| %s\t| Count: %d\t| % 20x \t|", itemName, itemCount, itemData[offset+int(itemNameLen)+1:offset+int(itemNameLen)+1+4+4+4+4+17])
		// add to table
		fmt.Fprintln(w, sOut)

		payload := &ItemPayload{
			ItemNameLen:       itemNameLen,
			ItemName:          itemName,
			ItemCount:         itemCount,
			Unknown3:          binary.LittleEndian.Uint32(unknown3),
			XcoordInInventory: xCoordInInventory,
			YcoordInInventory: yCoordInInventory,
			Lvl:               lvl,
			LvlIndex:          lvlIndex,
			Unknown4:          unknown4,
			HasOwner:          hasOwner,
			OwnerName:         owner,
			ControlData: itemControlData{
				PayloadIndex:              startOfItemsSection + byteOffset,
				PayloadIndexStartingOnQty: startOfItemsSection + byteOffset + int(itemNameLen) + 1,
				Payload:                   fullPayload,
				OriginalCount:             int(itemCount),
				// ModifiedCount int
				// MaxCount      int
				// ToModify      bool
			},
		}

		byteOffset = byteOffset + int(itemNameLen) + 1 + 4 + 4 + 4 + 4 + 17
		items = append(items, *payload)
	}
	w.Flush()

	remainingData := data[startOfItemsSection+byteOffset:]
	fmt.Println("offset:", startOfItemsSection+byteOffset)
	fmt.Printf("data: % 20x \n", remainingData[:30])
	fmt.Printf("data: %s \n", remainingData[:30])

	nextSectionCount := binary.LittleEndian.Uint32(remainingData[1:5])
	fmt.Println("next_section_count: ", nextSectionCount)

	offset := 5
	i := 0
	// for i := 0; i < int(next_section_count); i++ {
	for j, char := range remainingData[offset:] {
		if string(char) == "$" {
			fmt.Printf("i:%d | str: %s\n", i, remainingData[offset+j:offset+j+8])
			i++
		}
		if i == int(nextSectionCount) {
			offset += int(remainingData[offset+j])
			break
		}
	}

	nextRemainingData := remainingData[offset:]
	fmt.Printf("d: % 20x\n", nextRemainingData[:30])
	fmt.Printf("d: %c", nextRemainingData[:30])

	return items, byteOffset
}

// LoadDbItems loads items database from json file
func LoadDbItems(bytes []byte) {
	if bytes == nil {
		log.Panic("could not load json file")
	}
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'db_items' which we defined on the beginning of the file
	json.Unmarshal(bytes, &dbItems)
}

// GetItemFromDbItemWithName gets an item from items database by name
func GetItemFromDbItemWithName(name string) DbItem {
	for _, item := range dbItems {
		if item.InternalID == name {
			// fmt.Println(item.Name)
			return item
		}
	}
	return DbItem{}
}
