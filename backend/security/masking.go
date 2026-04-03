package security

import (
	"strings"
)

// MaskString оставляет видимыми только первые и последние символы.
// Корректно работает с UTF-8 (русскими буквами).
// Например: "Сергей Игоревич" -> "Се***********ич"
func MaskString(s string) string {
	runes := []rune(s)
	l := len(runes)
	if l <= 4 {
		return "****"
	}
	// Оставляем 2 символа в начале и 2 в конце, остальное забиваем звездами
	return string(runes[:2]) + strings.Repeat("*", l-4) + string(runes[l-2:])
}

// MaskDiploma скрывает центральную часть номера диплома.
// Например: "106104-777" -> "106****77"
func MaskDiploma(number string) string {
	runes := []rune(number)
	l := len(runes)
	if l < 6 {
		return "***"
	}
	// Оставляем 3 символа в начале и 2 в конце
	return string(runes[:3]) + "****" + string(runes[l-2:])
}
