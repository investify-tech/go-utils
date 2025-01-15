package str

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gookit/goutil/maputil"
	"sort"
)

// Other nice functions regarding string manipulations available at 'github.com/stoewer/go-strcase'

func DeduplicateSliceOfStrings(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func SortMapOfStrings(inputMap map[string]string) map[string]string {
	mapKeys := maputil.Keys(inputMap)
	sort.Strings(mapKeys)
	sortedMap := make(map[string]string)
	for _, key := range mapKeys {
		sortedMap[key] = inputMap[key]
	}
	return sortedMap
}

// HashAndTrim hashes the given input string and trims the result to the given length (i.e. get X chars from 0)
func HashAndTrim(dataToHash string, resultLength int) string {
	hash := md5.Sum([]byte(dataToHash))
	hexHash := hex.EncodeToString(hash[:])
	return hexHash[:resultLength]
}
