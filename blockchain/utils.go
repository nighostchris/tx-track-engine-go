package blockchain

import "strings"

func Contains(elements []string, target string, checkCasing bool) bool {
	for _, element := range elements {
		if checkCasing {
			if strings.ToLower(target) == strings.ToLower(element) {
				return true
			}

			return false
		} else {
			if strings.ToLower(target) == strings.ToLower(element) {
				return true
			}

			return false
		}
	}

	return false
}
