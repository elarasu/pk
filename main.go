package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/dustin/go-humanize"
)

func hashFileMd5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func main() {
	// fmt.Println(os.Args)
	// commands.Execute()
	dirPath := "."
	var blackListExp = regexp.MustCompile(`^\..*$`)
	fileSizes := map[string]int64{}
	// blackList := `^\..*`
	// files, err := ioutil.ReadDir(dirPath)
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			matched := blackListExp.MatchString(path)
			if matched {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			fileSizes[path] = info.Size()
			// fmt.Println(path, info.Name(), info.Size())
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	// lets first run thru filters
	// for _, f := range files {
	// 	matched, _ := regexp.MatchString(blackList, f.Name())
	// 	if matched {
	// 		continue
	// 	}
	// 	if f.IsDir() {
	// 		continue
	// 	}
	// 	fileSizes[f.Name()] = f.Size()
	// }

	// lets sort the list by size
	type kv struct {
		Key   string
		Value int64
	}
	var ss []kv
	for k, v := range fileSizes {
		ss = append(ss, kv{k, v})
	}
	// ascending <, descending >
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	// lets loop thru and print - only find hashes if the sizes are same
	var prev kv
	var prevHash string
	var needHash bool = false

	for _, kv := range ss {
		if prev.Value == kv.Value {
			prevHash, _ = hashFileMd5(prev.Key)
			needHash = true
		}
		if needHash && prevHash == "" {
			prevHash, _ = hashFileMd5(prev.Key)
			needHash = false
		}
		if prev.Value != 0 {
			fmt.Printf("%-32s %s\t%s\n", prev.Key, humanize.Bytes(uint64(prev.Value)), prevHash)
		}
		prev = kv
		prevHash = ""
	}
	if needHash {
		prevHash, _ = hashFileMd5(prev.Key)
	}
	fmt.Printf("%-32s %s\t%s\n", prev.Key, humanize.Bytes(uint64(prev.Value)), prevHash)
}
