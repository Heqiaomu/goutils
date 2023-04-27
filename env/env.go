// Package env
/**
 * @Author: sunyang
 * @Email: sunyang@hyperchain.cn
 * @Date: 2023/4/23
 */
package env

import (
	"os"
	"strings"
)

func Get(key string, def interface{}) interface{} {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return v
}

func GetString(key, def string) string {
	val := Get(key, def)
	s, ok := val.(string)
	if !ok {
		return def
	}
	return s
}

func GetStringSlice(key, split string, def []string) []string {
	val := Get(key, def)
	s, ok := val.(string)
	if !ok {
		return def
	}
	return strings.Split(s, split)
}
