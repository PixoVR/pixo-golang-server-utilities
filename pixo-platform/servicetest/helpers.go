package servicetest

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-faker/faker/v4"
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GenerateRandomString(length int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateRandomUUID() string {
	id := uuid.New()
	return id.String()
}

func generateRandomText() string {
	return faker.Sentence()
}

func generateRandomID() string {
	return fmt.Sprint(rand.Intn(1000))
}

func EncodeItem(item map[string]interface{}) string {
	itemBytes, err := json.Marshal(item)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to encode data due to marshaling data")
		return ""
	}

	return base64.StdEncoding.EncodeToString(itemBytes)
}

func EncodeData(data map[string]interface{}, encodedPath string) (encodedData map[string]interface{}, err error) {
	pathList := strings.Split(encodedPath, ".")
	if len(pathList) > 2 {
		return data, errors.New("invalid path: can only encode up to two levels")
	}

	if len(pathList) == 1 {
		encodedData = data[pathList[0]].(map[string]interface{})
	} else if len(pathList) == 2 {
		data[pathList[0]].(map[string]interface{})[pathList[1]] = EncodeItem(data[pathList[0]].(map[string]interface{})[pathList[1]].(map[string]interface{}))
		encodedData = data
	}

	return encodedData, nil
}

func TrimString(input string) string {
	return strings.TrimSpace(strings.ReplaceAll(input, "\n", ""))
}

func RandomString(length int, letters []rune) string {
	stringBytes := make([]rune, length)

	for i := range stringBytes {
		stringBytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(stringBytes)
}

func GenerateRandomPassword() string {
	uppercaseLetters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowercaseLetters := []rune("abcdefghijklmnopqrstuvwxyz")
	numbers := []rune("0123456789")
	specialCharacters := []rune("!@#$%^&*()")

	password := ""

	for i := 0; i <= 3; i++ {
		password += RandomString(1, uppercaseLetters)
		password += RandomString(3, lowercaseLetters)
		password += RandomString(1, numbers)
		password += RandomString(1, specialCharacters)
	}

	return password
}
