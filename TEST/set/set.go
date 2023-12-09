package set

import (
	"bufio"
	"os"
	"strings"
)

// size определяет размер множества
const size = 20

// SetItem представляет элемент множества
type SetItem struct {
	Key    string // Ключ элемента
	IsUsed bool   // Флаг, указывающий, используется ли элемент
}

// Set представляет структуру множества
type Set struct {
	Items [size]*SetItem // Массив элементов множества
	Count int            // Количество элементов в множестве
}

// NewSet создает и возвращает новый экземпляр множества
func NewSet() *Set {
	s := &Set{Count: 0}
	for i := 0; i < size; i++ {
		s.Items[i] = &SetItem{}
	}
	return s
}

// HashFun реализует хеш-функцию для ключа
func (s *Set) HashFun(key string) int {
	const prime = 53
	hash := 0
	for _, letter := range key {
		hash = (hash*prime + int(letter)) % size
	}
	return hash
}

// CreateItem создает и возвращает новый элемент множества с указанным ключом
func (s *Set) CreateItem(key string) *SetItem {
	return &SetItem{Key: key, IsUsed: true}
}

// Push добавляет элемент в множество
func (s *Set) Push(key string) string {
	index := s.HashFun(key)
	item := s.CreateItem(key)

	if s.Count >= size {
		return "Множество полное"
	}

	if !s.Items[index].IsUsed {
		s.Items[index] = item
		s.Count++
	} else if s.Items[index].Key == key {
		return "Не оригинальный элемент"
	} else {
		for s.Items[index].IsUsed {
			if s.Items[index].Key == key {
				return "Не оригинальный элемент"
			}
			index = (index + 1) % size
		}
		s.Items[index] = item
		s.Count++
	}
	return "Элемент успешно добавлен"
}

// Pop удаляет элемент из множества по указанному ключу
func (s *Set) Pop(key string) string {
	if s.Count == 0 {
		return "Множество пустое"
	}

	index := s.HashFun(key)

	for s.Items[index].Key != key {
		if !s.Items[index].IsUsed {
			return "Нет такого элемента :("
		}
		index = index + 1
		if index >= size {
			return "Нет такого элемента :("
		}
	}

	s.Items[index].IsUsed = false

	for i := index + 1; i < size; i++ {
		if s.Items[i] != nil && s.HashFun(s.Items[i].Key) == index {
			s.Items[index] = s.Items[i]
			s.Items[i] = nil
			break
		}
	}
	s.Count--
	return "Элемент удален"
}

// Search ищет элемент в множестве по указанному ключу
func (s *Set) Search(key string) string {
	i := 0
	for i < size {
		if s.Items[i].IsUsed && s.Items[i].Key == key {
			return "Найден!"
		}
		i += 1
	}

	return "Нет такого элемента :("
}

// WriteToFile записывает содержимое множества в файл
func (s *Set) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, item := range s.Items {
		if item != nil && item.IsUsed {
			_, err := writer.WriteString(item.Key + "\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ReadFromFile читает содержимое файла и добавляет элементы в множество
func (s *Set) ReadFromFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := strings.TrimSpace(scanner.Text())
		if key != "" {
			s.Push(key)
		}
	}

	return scanner.Err()
}
