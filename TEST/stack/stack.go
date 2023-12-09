package stack

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Stack представляет структуру стека
type Stack struct {
	Data string // Данные элемента стека
	Next *Stack // Указатель на следующий элемент стека
	Head *Stack // Указатель на вершину стека
}

// NewStack создает и возвращает новый экземпляр стека
func NewStack() *Stack {
	return &Stack{Head: nil}
}

// Push добавляет элемент на вершину стека
func (s *Stack) Push(value string) string {
	node := &Stack{Data: value, Next: nil}

	if s.Head == nil {
		s.Head = node
		return "Элемент успешно добавлен"
	} else {
		node.Next = s.Head
		s.Head = node
		return "Элемент успешно добавлен"
	}
}

// Pop удаляет элемент с вершины стека
func (s *Stack) Pop() string {
	if s.Head == nil {
		return "Стек пуст"
	} else {
		temp := s.Head
		s.Head = s.Head.Next
		return "Удалённый элемент: " + temp.Data
	}
}

// WriteToFile записывает содержимое стека в файл
func (s *Stack) WriteToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := file.WriteString
	current := s.Head
	for current != nil {
		_, err := writer(fmt.Sprintf("%s\n", current.Data))
		if err != nil {
			return err
		}
		current = current.Next
	}

	return nil
}

// ReadFromFile читает содержимое файла и добавляет элементы в стек
func (s *Stack) ReadFromFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var values []string

	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		values = append(values, value)
	}

	// Переворачиваем порядок элементов перед добавлением в стек
	for i := len(values) - 1; i >= 0; i-- {
		s.Push(values[i])
	}

	return scanner.Err()
}
