package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"
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
	items_log_file_debug, err := os.Create(utils.Bckp_folder + character + "_itemslog_" + utils.GetTimestampString() + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer items_log_file_debug.Close()

	// write as a table
	w_debug := tabwriter.NewWriter(items_log_file_debug, 10, 2, 1, ' ', 0)
	//// END OF DEBUG

	// write as a table
	w := tabwriter.NewWriter(os.Stdout, 10, 2, 1, ' ', 0)
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
		s_out := fmt.Sprintf("| %s\t| Count: %d\t| % 20x \t|", item_name, item_payload[0], item_payload)

		// add to table
		fmt.Fprintln(w, s_out)

		// add TO DEBUG table
		fmt.Fprintln(w_debug, s_out)

		// fmt.Printf("name: %s :\t % 20x \t| len: %d, extra byte: %v\n", getItemName(player_data_string, start_byte_i), item_payload, len(item_payload), hasExtraByte)
		// fmt.Printf("%+q", patterns)
	}
	// print table
	w.Flush()

	// print to DEBUG file
	w_debug.Flush()

	return totalItems
}

func ModifyItemData(full_data []byte, items []Item) []byte {
	for _, item := range items {
		full_data[int(item.PayloadIndex)] = byte(item.ModifiedCount)
	}
	return full_data
}

func LoadItems(character string) (loadedItems []Item, FileData []byte) {
	path := utils.CharactersFolder

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

	character_path := path + character + ".fch"

	// open selected character
	file, err := os.Open(character_path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	// get total amount of bytes
	file_stats, err := file.Stat()
	if err != nil {
		log.Fatal("could not get file size", err)
	}
	file_size := file_stats.Size()
	fmt.Println("file size: ", file_size)

	// read all data and close file
	full_data := utils.ReadNextBytes(file, file_size)
	file.Close()

	charname := string(unicode.ToUpper(rune(character[0]))) + character[1:]
	fmt.Println("Character: ", charname)
	full_string := string(full_data)
	i := strings.Index(full_string, charname)

	// parse header, still not sure of structure format and meaning, probably date and time somewhere?
	header := Header{}
	buffer := bytes.NewBuffer(full_data[:96])
	err = binary.Read(buffer, binary.LittleEndian, &header)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	// fmt.Printf("Header data:\n%+v\n", header)

	player_data_string := full_string[i:]

	// pattern to look for in items
	byte_pattern := []byte{1, 0, 0, 0, 0, 0, 0, 0}
	string_pattern := string(byte_pattern)

	// create a slice of patterns to look for
	patterns := make([]string, 1)
	// inserts the pattern to look for
	patterns[0] = string_pattern

	// finds all indexes where the pattern occurs, result is a map because there can be more patterns to look for
	result := FindAllOccurrences([]byte(player_data_string), patterns)

	// get the match indexes
	matches := result[string_pattern]

	// reverse order of matches
	matches = utils.ReverseIntSlice(matches)

	// clean matches by verifying some extra patterns on each item
	matches = CleanItemMatches(full_data, player_data_string, i, matches)

	fmt.Println("Items found: ", len(matches))

	totalItems := GetItems(matches, full_data, player_data_string, i, character)
	fmt.Printf("Items: %v\n", totalItems)

	return totalItems, full_data
}
