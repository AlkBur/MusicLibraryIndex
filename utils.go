package MusicLibraryIndex

import (
	"log"
	"strings"
)

func strictEqual(src, dst interface{}) {
	if src != dst {
		log.Fatalf("Error: %v != %v", src, dst)
	}
}

func findEndStr(arr []string, start, end string, skip ...string) int {
	iL := len(arr)
	for i := 0; i < iL; i++ {
		field := arr[i]
		bSkip := false
		for _, sk := range skip{
			if strings.HasSuffix(field, sk) || strings.HasPrefix(field, sk) {
				bSkip = true
				break
			}
		}
		if bSkip {
			continue
		}
		if strings.HasSuffix(field, end) {
			return i
		}else if strings.HasSuffix(field, start) && i != 0 {
			i = findEndStr(arr[i:], start, end)+i
		}
	}
	return iL - 1
}
