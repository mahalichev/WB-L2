package anagrams

import (
	"slices"
	"strings"
)

// Сортировка символов в строке
func SortLetters(str string) string {
	letters := []rune(str)
	slices.Sort(letters)
	return string(letters)
}

// Удаление множеств из одного элемента
func ClearAnagramsMap(anagramsMap map[string][]string) {
	for key, anagrams := range anagramsMap {
		if len(anagrams) == 1 {
			delete(anagramsMap, key)
		}
	}
}

// Сортировка слов в множествах
func SortAnagramsMap(anagramsMap map[string][]string) {
	for _, anagrams := range anagramsMap {
		slices.Sort(anagrams)
	}
}

// Построение map множества анаграмм
func AnagramsMap(words *[]string) *map[string][]string {
	// map, хранящий результат выполнения функции
	result := make(map[string][]string)

	// map - отображение строки с отсортированными символами и ключом, которому соответствует строка
	sortedWordKeyMap := make(map[string]string)

	// Множество слов
	wordsSet := make(map[string]struct{})

	for _, word := range *words {
		word = strings.ToLower(word)

		// Если слово есть в множестве - пропуск
		if _, ok := wordsSet[word]; ok {
			continue
		}
		wordsSet[word] = struct{}{}

		// Сортировка символов в слове
		sortedLetters := SortLetters(word)

		// Если не найдено ключа для sortedLetters, word - ключ
		if _, ok := sortedWordKeyMap[sortedLetters]; !ok {
			sortedWordKeyMap[sortedLetters] = word
		}
		key := sortedWordKeyMap[sortedLetters]
		// Добавление слова в соответствующее множество
		result[key] = append(result[key], word)
	}
	// Удаление множеств из одного элемента
	ClearAnagramsMap(result)
	// Сортировка множеств
	SortAnagramsMap(result)
	return &result
}
