package letterStorage

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type LetterStorage struct {
	letters map[string]int
	sync.Mutex
}

func New() LetterStorage {
	return LetterStorage{letters: make(map[string]int)}
}

func (ls LetterStorage) Add(letter string) {
	ls.Lock()

	numOfLetters, ok := ls.letters[letter]
	if ok {
		ls.letters[letter] = numOfLetters + 1
	} else {
		ls.letters[letter] = 1
	}

	ls.Unlock()
}

func (ls LetterStorage) Join(outerLs LetterStorage) {
	ls.Lock()

	for innerLetter, outerNumOfLetters := range outerLs.letters {
		innerNumOfLetters, ok := ls.letters[innerLetter]
		if ok {
			ls.letters[innerLetter] = outerNumOfLetters + innerNumOfLetters
		} else {
			ls.letters[innerLetter] = outerNumOfLetters
		}
	}

	ls.Unlock()
}

func (ls LetterStorage) ToString() string {
	var letters []string

	for letter := range ls.letters {
		letters = append(letters, letter)
	}

	sort.Strings(letters)

	sb := strings.Builder{}

	for _, letter := range letters {
		sb.WriteString(fmt.Sprintf("%s: %d\r\n", letter, ls.letters[letter]))
	}

	return sb.String()
}
