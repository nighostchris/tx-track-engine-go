package blockchain

import (
	"strconv"
	"strings"
)

func HexToDec(input string) (decStr string, err error) {
	postProcessed := strings.Replace(input, "0x", "", -1)
	decimal, decParseError := strconv.ParseUint(postProcessed, 16, 64)

	if decParseError != nil {
		return "", decParseError
	}

	return strconv.FormatUint(decimal, 10), nil
}

func Contains(elements []string, target string, checkCasing bool) bool {
	for _, element := range elements {
		if checkCasing {
			return target == element
		} else {
			return strings.EqualFold(target, element)
		}
	}

	return false
}
