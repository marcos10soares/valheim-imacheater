package utils

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
)

func ReadNextBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func ReverseIntSlice(s []int) []int {
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

func ReplaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func CleanString(s string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(s, "")
	return processedString
}

func WriteOutputFile(full_data []byte) {
	// create file copy with modified data - for debugging
	new_file, err := os.Create("output/bjørn.fch")
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	new_file.Write(full_data)
	new_file.Close()
}

func GetCurrentUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	return usr, err
}

func GetAllAvailableCharacters(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	var characters []string

	for _, f := range files {
		filename := f.Name()
		if filename[len(filename)-4:] == ".fch" {
			characters = append(characters, GetFileNameFromPath(dir+f.Name()))
		}
	}
	return characters
}

func GetFileNameFromPath(path string) string {
	var filename string
	var s []string
	if runtime.GOOS == "windows" {
		s = strings.Split(path, "\\")

	} else {
		s = strings.Split(path, "/")
	}
	str := s[len(s)-1]
	filename = str[:len(str)-4]
	return filename
}

// func GetCharNameFromPath(path string) string {
// 	s := GetFileNameFromPath(path)
// 	return
// }
