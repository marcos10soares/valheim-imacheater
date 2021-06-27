package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"
)

const (
	// WinPath is the Windows relative path for the character files folder
	WinPath = "\\AppData\\LocalLow\\IronGate\\Valheim\\characters\\"
	// MacPath is the MacOS relative path for the character files folder
	MacPath = "files/" // for debugging
	// LinuxPath is the Linux relative path for the character files folder
	LinuxPath = "/.config/unity3d/IronGate/Valheim/characters/"
	// BckpFolder is the backup folder
	BckpFolder = "./bckp/"
)

// CharactersFolderPath is the absolute path for the characters folder
var CharactersFolderPath string

// Init - check if bckp folder exists otherwise creates it
func Init() {
	// check if bckp folder exists otherwise creates it
	CreateBckpFolder(BckpFolder)
}

// ReadNextBytes reads the next bytes of file
func ReadNextBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("Error reading bytes: ", file.Name(), err)
	}

	return bytes
}

// ReverseIntSlice reverts the order of an int slice
func ReverseIntSlice(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseString reverses a string
func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

// ReplaceAtIndex replaces a character in a string at an index
func ReplaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

// CleanString cleans string from unwanted chars
func CleanString(s string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(s, "")
	return processedString
}

// WriteOutputFile writes data to character file
func WriteOutputFile(fullData []byte, character string) {
	path := CharactersFolderPath

	characterFilename := character + ".fch"

	// backup original file
	_, err := FileCopy(path+characterFilename, BckpFolder+character+"__"+GetTimestampString()+".fch")
	if err != nil {
		log.Fatal("Could not backup original file.")
	}

	// update file
	charFile, err := os.OpenFile(path+characterFilename, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	charFile.Truncate(0)
	charFile.Write(fullData)
	charFile.Close()
}

// GetCurrentUser fetches current OS user
func GetCurrentUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	return usr, err
}

// GetAllAvailableCharacters gets a list of character files
func GetAllAvailableCharacters(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Error accessing characters dir: ", err)
	}

	var characters []string

	for _, f := range files {
		filename := f.Name()
		if len(filename) > 4 {
			if filename[len(filename)-4:] == ".fch" {
				characters = append(characters, GetFileNameFromPath(dir+f.Name()))
			}
		}
	}
	return characters
}

// GetFileNameFromPath extracts a file name from a file path
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

// FileCopy copy file
func FileCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		// fmt.Println(err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		// fmt.Println(err)
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		// fmt.Println(err)
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	// fmt.Println(err)
	return nBytes, err
}

// CreateBckpFolder creates a backup folder if it does not exist
func CreateBckpFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}

// GetTimestampString gets a timestamp in string format from time now
func GetTimestampString() string {
	now := time.Now()
	unixTimeUTC := time.Unix(now.Unix(), 0)
	timestamp := unixTimeUTC.Format(time.RFC3339)
	return strings.Replace(timestamp, ":", "_", -1)
}

// ListDirRecursively gets all files from a dir recursively
func ListDirRecursively(root string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false && info.Name() != ".DS_Store" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// ReadFileFromPath gets the data from a file
func ReadFileFromPath(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	// get total amount of bytes
	fileStats, err := file.Stat()
	if err != nil {
		log.Fatal("could not get file size", err)
	}
	fileSize := fileStats.Size()
	// fmt.Println("file size: ", file_size)

	// read all data and close file
	fullData := ReadNextBytes(file, fileSize)
	file.Close()

	return fullData
}

// MakeTitle converts string to title case
func MakeTitle(s string) string {
	s = string(unicode.ToUpper(rune(s[0]))) + s[1:]

	var spaceIndex = []int{}
	for i := 0; i < len(s); i++ {
		if unicode.IsSpace(rune(s[i])) {
			spaceIndex = append(spaceIndex, i)
		}
	}

	for _, v := range spaceIndex {
		s = s[:v+1] + string(unicode.ToUpper(rune(s[v+1]))) + s[v+2:]
	}

	return s
}

// StringSliceCheckIfContains checks if string slice contains string
func StringSliceCheckIfContains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			// return i, true
			return true
		}
	}
	return false
}
