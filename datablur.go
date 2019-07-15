package datablur

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

type DataBlur interface {
	blur(dataIn string) (string, bool)
}

type Substitute struct {
	lookupTable map[string]string
	lookupFile  string
}

func (s *Substitute) blur(d string) (string, bool) {
	// Lets lookup the file if we have it
	subStr, found := s.lookupFromFile(d)

	// We didn't find it in a file, let look up the supplied lookupTable
	if !found {
		subStr, found = s.lookupTable[d]
	}

	// We didn't find it in the lookup table, we'll return the original string
	// unchanged
	if !found {
		return d, found
	}

	return subStr, found
}

func (s *Substitute) lookupFromFile(d string) (string, bool) {
	var err error

	f, err := os.Open(s.lookupFile)
	if err != nil {
		return d, false
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		t := scanner.Text()
		tSplit := strings.Split(t, ",")
		if tSplit[0] == d {
			return tSplit[1], true
		}
	}

	return d, false
}

type Rot13 struct{}

func (r *Rot13) blur(d string) (string, bool) {
	rotated := make([]byte, len(d))
	for i := 0; i < len(d); i++ {
		s := unicode.ToLower(rune(d[i]))
		switch {
		case s >= 'a' && s <= 'm':
			rotated[i] = d[i] + 13
		case s >= 'n' && s <= 'z':
			rotated[i] = d[i] - 13
		default:
			rotated[i] = d[i]
		}
	}

	return string(rotated), true
}
