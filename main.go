package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

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

func main() {
	path := "files/bjørn.fch" // bjørn
	charname := "Bjørn"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	file_stats, err := file.Stat()
	if err != nil {
		log.Fatal("could not get file size", err)
	}
	file_size := file_stats.Size()
	fmt.Println("file size: ", file_size)

	// fmt.Printf("%s opened\n", path)

	header := Header{}
	// data := readNextBytes(file, 96) //  24 * uint32 (4) = 96 bytes

	full_data := readNextBytes(file, file_size)

	file.Close()

	new_file, err := os.Create("bjørn.fch")
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	// for i := 0; i < 100; i++ {
	// 	fmt.Printf("%x", full_data[i])
	// }
	// fmt.Println("-")
	fmt.Println("Character: ", charname)
	full_string := string(full_data)
	i := strings.Index(full_string, charname)
	// fmt.Println("Index: ", i)

	// fmt.Println(data)

	buffer := bytes.NewBuffer(full_data[:96])
	err = binary.Read(buffer, binary.LittleEndian, &header)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	// fmt.Printf("Parsed data:\n%+v\n", header)

	player_data_string := full_string[i:]

	// fmt.Println("Full data: \n", player_data_string)

	byte_pattern := []byte{1, 0, 0, 0, 0, 0, 0, 0}
	string_pattern := string(byte_pattern)

	// string_pattern_i := strings.Index(player_data_string, string_pattern)
	// fmt.Println("index: ", string_pattern_i)
	// fmt.Println(player_data_string[string_pattern_i-17 : string_pattern_i+17])

	arr := make([]string, 1)
	arr[0] = string_pattern

	res := findAllOccurrences([]byte(player_data_string), arr)
	// fmt.Println(res)
	string_pattern_matches := res[string_pattern]

	fmt.Println("Items found: ", len(string_pattern_matches))

	string_pattern_matches = reverseIntSlice(string_pattern_matches)

	for _, match := range string_pattern_matches {
		hasExtraByte := checkIfItemPayloadHasExtraByte(player_data_string, match)
		start_byte_i := match - 17
		item_payload_size := 17
		if hasExtraByte {
			item_payload_size += 1
		}

		item_payload_start_byte := (i + start_byte_i)

		item_count_index := full_string[(i + start_byte_i)]
		item_name := getItemName(player_data_string, start_byte_i)

		if item_name == "CookedMeat" {
			fmt.Println("index of player data: ", i)
			fmt.Println("index of item payload in player data:", start_byte_i)

			// full_string = replaceAtIndex(full_string, rune(20), count_byte_index)
			full_data[item_payload_start_byte] = 20
			// fmt.Printf("new count: %s \n", full_data[count_byte_index-10:count_byte_index])
			// fmt.Println("new count index: ", count_byte_index)
			// item_count_index = full_string[count_byte_index]
		}
		item_payload := []byte(full_data[item_payload_start_byte : item_payload_start_byte+item_payload_size])
		fmt.Println("Count:", item_count_index)
		fmt.Printf("name: %s | Count: %d |\t % 20x\n", item_name, item_payload[0], item_payload)
		// fmt.Printf("name: %s :\t % 20x\n", getItemName(player_data_string, start_byte_i), item_payload)
		// fmt.Printf("name: %s :\t % 20x \t| len: %d, extra byte: %v\n", getItemName(player_data_string, start_byte_i), item_payload, len(item_payload), hasExtraByte)
		// fmt.Printf("%+q", arr)
	}

	new_file.Write(full_data)
	new_file.Close()

}

func readNextBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func getItemName(player_data string, start_byte_i int) string {
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
	item_name = ReverseString(item_name)
	// fmt.Println([]byte(item_name))
	return strings.Replace(strings.Trim(strings.TrimSpace(item_name[1:]), "\n"), string(0xc), "", -1)
}

func checkIfItemPayloadHasExtraByte(player_data string, index int) bool {
	extra_byte_payload_marker := string([]byte{0x6d, 0x1e, 0xf7, 0xd1})

	if extra_byte_payload_marker == player_data[index+8:index+12] {
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

func reverseIntSlice(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
