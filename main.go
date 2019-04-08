package main

// This program read text files from given directory and prints the number of uniq symbols in them

import (
	"fileScanner/letterStorage"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		panic("err: no path")
	}

	path := os.Args[1]

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		panic("err: path no exist")
	}

	if !fileInfo.IsDir() {
		panic("err: is not dir")
	}

	wg := sync.WaitGroup{}
	lettersChannel := make(chan letterStorage.LetterStorage)
	walkFunc := getWalkFunc(lettersChannel, &wg)

	err = filepath.Walk(path, walkFunc)
	checkErr(err)

	go closeChannel(lettersChannel, &wg)

	resultMap := letterStorage.New()
	for e := range lettersChannel {
		resultMap.Join(e)
	}

	fmt.Println(resultMap.ToString())
}

func closeChannel(lettersChannel chan letterStorage.LetterStorage, wg *sync.WaitGroup) {
	wg.Wait()
	close(lettersChannel)
}

func getWalkFunc(lettersChannel chan letterStorage.LetterStorage, wg *sync.WaitGroup) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if info.IsDir() {
				return
			}

			fileBytes, err := ioutil.ReadFile(path)
			checkErr(err)

			letters := letterStorage.New()

			fileString := string(fileBytes)
			for _, ch := range fileString {
				letter := string(ch)
				if letter == " " || letter == "\n" || letter == "\r" {
					continue
				}

				letters.Add(letter)
			}

			lettersChannel <- letters
		}()
		return nil
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
