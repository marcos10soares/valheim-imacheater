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

const WinPath = "\\AppData\\LocalLow\\IronGate\\Valheim\\characters\\"
const MacPath = "files/" // for debugging
var CharactersFolder string

const Bckp_folder = "./bckp/"

func Init() {
	// check if bckp folder exists otherwise creates it
	CreateBckpFolder(Bckp_folder)
}

func ReadNextBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("Error reading bytes: ", file.Name(), err)
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

func WriteOutputFile(full_data []byte, character string) {
	path := CharactersFolder

	character_filename := character + ".fch"

	// backup original file
	_, err := FileCopy(path+character_filename, Bckp_folder+character+"__"+GetTimestampString()+".fch")
	if err != nil {
		log.Fatal("Could not backup original file.")
	}

	// update file
	char_file, err := os.OpenFile(path+character_filename, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	char_file.Truncate(0)
	char_file.Write(full_data)
	char_file.Close()
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

func CreateBckpFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}

func GetTimestampString() string {
	now := time.Now()
	unixTimeUTC := time.Unix(now.Unix(), 0)
	timestamp := unixTimeUTC.Format(time.RFC3339)
	return strings.Replace(timestamp, ":", "_", -1)
}

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

// func GetCharNameFromPath(path string) string {
// 	s := GetFileNameFromPath(path)
// 	return
// }

// func WriteItemLogFile(charname string, data *tabwriter.Writer) {
// 	f, err := os.Create(Bckp_folder + charname + "_itemslog_" + GetTimestampString() + ".txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	_, err2 := f.WriteString(data)

// 	if err2 != nil {
// 		log.Fatal(err2)
// 	}

// 	fmt.Println("done")
// }

func ReadFileFromPath(path string) []byte {
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
	// fmt.Println("file size: ", file_size)

	// read all data and close file
	full_data := ReadNextBytes(file, file_size)
	file.Close()

	return full_data
}

func MakeTitle(s string) string {
	s = string(unicode.ToUpper(rune(s[0]))) + s[1:]

	var space_index = []int{}
	for i := 0; i < len(s); i++ {
		if unicode.IsSpace(rune(s[i])) {
			space_index = append(space_index, i)
		}
	}

	for _, v := range space_index {
		s = s[:v+1] + string(unicode.ToUpper(rune(s[v+1]))) + s[v+2:]
	}

	return s
}

func StringSliceCheckIfContains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			// return i, true
			return true
		}
	}
	// return -1, false
	return false
}
