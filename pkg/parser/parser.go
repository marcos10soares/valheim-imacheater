package parser

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	utils "vimacheater/pkg/utils"
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

func GetItems(matches []int, full_data []byte, player_data_string string, i int) []Item {
	totalItems := []Item{}

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

		// change the value for the item CookedMeat
		// if item_name == "CookedMeat" {
		// 	full_data[item_payload_start_byte] = 20
		// }
		// item_count := full_string[item_payload_start_byte]

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

		// fmt.Printf("name: %s :\t % 20x \t| len: %d, extra byte: %v\n", getItemName(player_data_string, start_byte_i), item_payload, len(item_payload), hasExtraByte)
		// fmt.Printf("%+q", patterns)
	}
	// print table
	w.Flush()

	return totalItems
}

func ModifyItemData(full_data []byte, items []Item) []byte {
	for _, item := range items {
		full_data[int(item.PayloadIndex)] = byte(item.ModifiedCount)
	}
	return full_data
}
