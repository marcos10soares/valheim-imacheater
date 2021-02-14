package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
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

const path = "\\AppData\\LocalLow\\IronGate\\Valheim\\Characters\\"

func main() {
	// config, to be removed when GUI is created
	path := "files/bjørn.fch" // bjørn
	charname := "Bjørn"

	// open selected character
	file, err := os.Open(path)
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
	full_data := readNextBytes(file, file_size)
	file.Close()

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
	result := parser.FindAllOccurrences([]byte(player_data_string), patterns)

	// get the match indexes
	matches := result[string_pattern]
	fmt.Println("Items found: ", len(matches))

	// reverse order of matches
	matches = utils.ReverseIntSlice(matches)

	// write as a table
	w := tabwriter.NewWriter(os.Stdout, 10, 2, 1, ' ', 0)
	for _, match := range matches {
		// items payload length is variable, this checks the item payload size
		hasExtraByte := parser.CheckIfItemPayloadHasExtraByte(player_data_string, match)
		start_byte_i := match - 17
		item_payload_size := 34
		if hasExtraByte {
			item_payload_size += 1
		}

		// get payload limits index
		item_payload_start_byte := (i + start_byte_i)
		end_byte_i := item_payload_start_byte + item_payload_size

		// get item name (this string is modified for printing reasons, do not use for changing output file)
		item_name := parser.GetItemName(player_data_string, start_byte_i)

		// change the value for the item CookedMeat
		if item_name == "CookedMeat" {
			full_data[item_payload_start_byte] = 20
		}
		// item_count := full_string[item_payload_start_byte]
		item_payload := []byte(full_data[item_payload_start_byte:end_byte_i])

		// format string
		s_out := fmt.Sprintf("| %s\t| Count: %d\t| % 20x \t|", item_name, item_payload[0], item_payload)

		// add to table
		fmt.Fprintln(w, s_out)

		// fmt.Printf("name: %s :\t % 20x \t| len: %d, extra byte: %v\n", getItemName(player_data_string, start_byte_i), item_payload, len(item_payload), hasExtraByte)
		// fmt.Printf("%+q", patterns)
	}
	// print table
	w.Flush()

	// create file copy with modified data - for debugging
	new_file, err := os.Create("bjørn.fch")
	if err != nil {
		log.Fatal("Error while opening file", err)
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
