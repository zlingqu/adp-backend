package tools

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"strconv"
	"time"
	"reflect"
)

func GetEnvDefault(key string, defVal interface{}) interface{} {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

func StringToInt(string string) int {
	//if s, err := strconv.Atoi(string); err == nil {
	//	fmt.Printf("%T, %v", s, s)
	//}
	s, _ := strconv.Atoi(string)
	return s
}

func StringToUint(str string) uint {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return uint(i)
}

func GetMysqlLimitOffset(page string, size string) (int, int) {
	intPage := StringToInt(page)
	limit := StringToInt(size)
	offset := (intPage - 1) * limit
	return limit, offset
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(8)
}

// string to bool
func StringToBool(string string) bool {
	b, err := strconv.ParseBool(string)
	if err != nil {
		return false
	}
	return b
}

//生成随机字符串
func GetRandomString(lenght int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	bytesLen := len(bytes)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenght; i++ {
		result = append(result, bytes[r.Intn(bytesLen)])
	}
	return string(result)
}

func Float64ToInt64(num float64) int64 {
	string := strconv.FormatFloat(num, 'f', -1, 64)
	int64, _ := strconv.ParseInt(string, 10, 64)
	return int64
}


//结构体转成map
func Struct2Map(obj interface{}) map[string]interface{} {
    t := reflect.TypeOf(obj)
    v := reflect.ValueOf(obj)

    var data = make(map[string]interface{})
    for i := 0; i < t.NumField(); i++ {
        data[t.Field(i).Name] = v.Field(i).Interface()
    }
    return data
}