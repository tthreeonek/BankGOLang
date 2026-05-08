package luhn

import (
	"math/rand"
	"strconv"
	"time"
)

// Генерация номера, начинающегося с prefix
func Generate(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	// общая длина 16 (например, 4xxxxx...)
	length := 16
	pan := prefix
	for i := len(prefix); i < length-1; i++ {
		pan += strconv.Itoa(rand.Intn(10))
	}
	// вычислить контрольную цифру
	checkDigit := calculateLuhnCheckDigit(pan)
	return pan + strconv.Itoa(checkDigit)
}

func calculateLuhnCheckDigit(number string) int {
	sum := 0
	parity := len(number) % 2
	for i := 0; i < len(number); i++ {
		digit, _ := strconv.Atoi(string(number[i]))
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (10 - (sum % 10)) % 10
}
