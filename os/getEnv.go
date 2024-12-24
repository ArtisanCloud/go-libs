package os

import (
	"errors"
	"os"
	"strconv"
)

const LOCALE_EN = "en_US"
const LOCALE_CN = "zh_CN"
const LOCALE_TW = "zh_TW"

var ErrEnvVarEmpty = errors.New("getEnv: environment variable empty")

func GetEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func GetEnvInt(key string) (int, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetEnvBool(key string) (bool, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
