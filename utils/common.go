package utils

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GenerateAuthKey() string {
	rand.Seed(time.Now().UnixNano())
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

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
