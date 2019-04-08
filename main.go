package main

// This program read text files from given directory and prints the number of uniq symbols from them

import (
	"fileScanner/letterStorage"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const maxWorkingGoroutines = 100

func main() {
	if len(os.Args) < 2 {
		panic("err: no path")
	}

	path := os.Args[1]

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		panic("err: path not exist")
	}

	if !fileInfo.IsDir() {
		panic("err: is not dir")
	}

	lettersChannel := make(chan letterStorage.LetterStorage)
	semaphore := make(chan int, maxWorkingGoroutines)
	wgForWalkFunc := sync.WaitGroup{}

	walkFunc := getWalkFunc(lettersChannel, &wgForWalkFunc, semaphore)

	wgForPrintResult := sync.WaitGroup{}
	wgForPrintResult.Add(1)

	go printResult(lettersChannel, &wgForPrintResult)

	filepath.Walk(path, walkFunc)

	wgForWalkFunc.Wait()
	close(lettersChannel)
	wgForPrintResult.Wait()
}

func printResult(lettersChannel chan letterStorage.LetterStorage, wg *sync.WaitGroup) {
	defer wg.Done()

	resultMap := letterStorage.New()
	for e := range lettersChannel {
		resultMap.Join(e)
	}

	fmt.Println(resultMap.ToString())
}

func getWalkFunc(lettersChannel chan letterStorage.LetterStorage, wg *sync.WaitGroup, semaphore chan int) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		semaphore <- 0
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()

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
