package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"vimacheater/pkg/gui"
	"vimacheater/pkg/parser"
	"vimacheater/pkg/utils"

	"github.com/zserge/lorca"
)

func main() {
	// setup log output file -- uncomment this for debugging
	// f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	// load items from json
	// parser.LoadDbItems(loadJsonFile())
	parser.LoadDbItems(assets["/items_list.json"])

	// render app
	renderApp()
}

func renderApp() {
	utils.Init()
	customArgs := []string{}
	ui, err := lorca.New("", "", 750, 800, customArgs...)
	if err != nil {
		log.Fatal("lorca error:", err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	u := &gui.UiData{
		Items: []gui.UiItem{},
	}

	ui.Bind("updatePower", u.UpdatePower)
	ui.Bind("updateItems", u.UpdateItems)
	ui.Bind("saveData", u.SaveData)
	ui.Bind("resetPowerCooldown", u.ResetPowerCooldown)
	ui.Bind("getPowers", u.GetPowers)
	ui.Bind("getItems", u.GetItems)
	ui.Bind("getChars", u.GetChars)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	// ui.Eval(`
	// 	console.log("Hello, world!");
	// 	console.log('Multiple values:', [1, false, {"x":5}]);
	// `)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}
