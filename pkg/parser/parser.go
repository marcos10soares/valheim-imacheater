package parser

import (
	"bytes"
	"strings"
	utils "vimacheater/pkg/utils"
)

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
	// fmt.Println([]byte(item_name))
	return strings.Replace(strings.Trim(strings.TrimSpace(item_name[1:]), "\n"), string(0xc), "", -1)
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
