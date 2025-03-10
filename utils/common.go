package utils

import (
	"log"
	"sort"
	"strconv"
	"strings"
)

func FormatVersionStr(v string) int64 {
	vs := strings.Split(v, ".")
	if len(vs) <= 0 {
		log.Panic("Version str error")
	}
	var vNum int64
	ReverseArr(vs)
	for index, v := range vs {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Panic(err.Error())
		}
		for i := 0; i < index; i++ {
			num = num * 100
		}
		vNum += num
	}
	return vNum
}

func ReverseArr(s interface{}) {
	sort.SliceStable(s, func(i, j int) bool {
		return true
	})
}
