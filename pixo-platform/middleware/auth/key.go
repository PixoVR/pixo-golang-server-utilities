package auth

import "os"

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}

func IsValidSecretKey(input string) bool {
	key := GetSecretKey()
	return input != "" && input == key
}
