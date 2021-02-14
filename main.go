package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	gui "vimacheater/pkg/gui"
	parser "vimacheater/pkg/parser"
	utils "vimacheater/pkg/utils"

	"github.com/zserge/lorca"
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
	full_data := utils.ReadNextBytes(file, file_size)
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

	// reverse order of matches
	matches = utils.ReverseIntSlice(matches)

	// clean matches by verifying some extra patterns on each item
	matches = parser.CleanItemMatches(full_data, player_data_string, i, matches)

	fmt.Println("Items found: ", len(matches))

	totalItems := parser.GetItems(matches, full_data, player_data_string, i)
	fmt.Printf("Items: %v\n", totalItems)

	renderApp(totalItems, full_data)
}

func renderApp(totalItems []parser.Item, full_data []byte) {
	customArgs := []string{}
	ui, err := lorca.New("", "", 640, 480, customArgs...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	u := &gui.UiItems{
		Items: totalItems,
	}
	ui.Bind("updateItems", u.UpdateItems)
	ui.Bind("getItems", u.GetItems)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	ui.Eval(`
		console.log("Hello, world!");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	full_data = parser.ModifyItemData(full_data, u.Items)

	// fmt.Printf("stone: %d\n", u.Items[2].ModifiedCount)
	utils.WriteOutputFile(full_data)

	log.Println("exiting...")
}
