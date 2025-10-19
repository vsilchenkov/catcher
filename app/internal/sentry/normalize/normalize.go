package normalize

import (
	"regexp"
	"strings"
)

func Id(id string) string {
	return strings.ReplaceAll(id, "-", "")
}

func Release(id, version string) string {
	return id + "@" + version
}

func Message(input string) string {

	message := removeBraces(input)
	message = strings.TrimSpace(message)
	return message
}

// 192.168.1.2 , 192.168.1.3
func Ip(input string) string {

	var res string

	parts := strings.Split(input, ",")
	if len(parts) > 0 {
		res = parts[0]
	} else {
		res = input
	}

	return strings.TrimSpace(res)

}

// Пример строки
// Infostart Toolkit PROF (2023.4.08)
func SplitString(input, sStart, sEnd string) [2]string {

	var result [2]string

	// Ищем любые символы (кроме скобок) внутри скобок
	// \([^)]*\)
	re := regexp.MustCompile(`\` + sStart + `[^` + sEnd + `]*\` + sEnd)

	first := re.ReplaceAllString(input, "")
	first = strings.TrimSpace(first)

	second := re.FindString(input)
	second = strings.ReplaceAll(second, sStart, "")
	second = strings.ReplaceAll(second, sEnd, "")
	second = strings.TrimSpace(second)

	result[0] = first
	result[1] = second
	return result

}

func RemoveModuleNameSuffix(s string) string {

	suffixs := [2]string{".Модуль", ".Форма"}

	for _, suffix := range suffixs {
		s = strings.TrimSuffix(s, suffix)
	}

	return s
}

func removeBraces(input string) string {

	// Регулярное выражение для нахождения всех вхождений {любые_символы_включая_переносы_строк}
	// re := regexp.MustCompile(`\{[^{}]*\}:|\}`)
	re := regexp.MustCompile(`\{[^{}]*\}:?|\}`)
	// Заменяем все совпадения на пустую строку
	result := re.ReplaceAllString(input, "")
	return result

}

// RemoveFromSecondBrace - удаляет полный стек из строки
// удаляет из строки все символы начиная со второй '{' до первой пустой строки (или до конца текста, если пустой строки нет)
func RemoveFromSecondBrace(input string, addBrace bool) string {

	// Найти все вхождения '{'
	re := regexp.MustCompile(`\{`)
	indices := re.FindAllStringIndex(input, -1)
	if len(indices) < 2 {
		return input // Если меньше двух '{', возвращаем исходный текст
	}
	secondBraceIdx := indices[1][0]

	// Получаем подстроку начиная со второго '{'
	substr := input[secondBraceIdx:]

	// Разбиваем подстроку на строки для поиска пустой строки
	lines := strings.Split(substr, "\n")
	emptyLineIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			emptyLineIdx = i
			break
		}
	}

	result := input[:secondBraceIdx]

	if emptyLineIdx != -1 {
		// Если пустая строка найдена, возвращаем текст до второго '{' + текст после пустой строки
		// Добавляем "\n" между частями, если нужно
		afterEmpty := strings.Join(lines[emptyLineIdx+1:], "\n")
		if len(afterEmpty) > 0 && !strings.HasPrefix(afterEmpty, "\n") {
			afterEmpty = "\n" + afterEmpty
		}
		// Если подстрока заканчивается запятой, удаляем её

		// Добавляем закрывающую фигурную скобку
		if addBrace {
			if strings.HasSuffix(result, "\n") {
				result = strings.TrimSuffix(result, "\n")
				result += "}\n"
			} else {
				result += "}"
			}
		}

		result = result + afterEmpty
	}

	result = strings.TrimSpace(result)
	// Если пустой строки нет, возвращаем текст до второго '{' (удаляем всё от второго '{' до конца)
	return result

}
