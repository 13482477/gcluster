package utils

import (
	"github.com/go-xorm/core"
	"regexp"
)

func CheckStringInSlice(search string, list []string) bool {
	for _, key := range list {
		if key == search {
			return true
		}
	}

	return false
}

func CheckInt32InSlice(search int32, list []int32) bool {
	for _, key := range list {
		if key == search {
			return true
		}
	}

	return false
}

func CamelToDbName(list []string) []string {
	mapper := core.SnakeMapper{}
	newList := make([]string, 0)
	for _, item := range list {
		newList = append(newList, mapper.Obj2Table(item))
	}
	return newList
}
func CheckEmailValid(email string) bool {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return false
	}

	return true
}

func CheckMobileValid(mobile string) bool {
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{8})$`, mobile); !m {
		return false
	}

	return true
}
