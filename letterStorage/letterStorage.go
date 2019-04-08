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

	v, ok := ls.letters[letter]
	if ok {
		ls.letters[letter] = v + 1
	} else {
		ls.letters[letter] = 1
	}

	ls.Unlock()
}

func (ls LetterStorage) Join(letterStorage LetterStorage) {
	ls.Lock()

	for k, v1 := range letterStorage.letters {
		v2, ok := ls.letters[k]
		if ok {
			ls.letters[k] = v1 + v2
		} else {
			ls.letters[k] = 1
		}
	}

	ls.Unlock()
}

func (ls LetterStorage) ToString() string {
	letters := []string{}

	for k := range ls.letters {
		letters = append(letters, k)
	}

	sort.Strings(letters)

	sb := strings.Builder{}

	for _, letter := range letters {
		sb.WriteString(fmt.Sprintf("%s: %d\r\n", letter, ls.letters[letter]))
	}

	return sb.String()
}
