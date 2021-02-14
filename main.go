package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	parser "vimacheater/pkg/parser"
	utils "vimacheater/pkg/utils"
)

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

	res := parser.FindAllOccurrences([]byte(player_data_string), arr)
	// fmt.Println(res)
	string_pattern_matches := res[string_pattern]

	fmt.Println("Items found: ", len(string_pattern_matches))

	string_pattern_matches = utils.ReverseIntSlice(string_pattern_matches)

	for _, match := range string_pattern_matches {
		hasExtraByte := parser.CheckIfItemPayloadHasExtraByte(player_data_string, match)
		start_byte_i := match - 17
		item_payload_size := 17
		if hasExtraByte {
			item_payload_size += 1
		}

		item_payload_start_byte := (i + start_byte_i)

		item_count_index := full_string[(i + start_byte_i)]
		item_name := parser.GetItemName(player_data_string, start_byte_i)

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
