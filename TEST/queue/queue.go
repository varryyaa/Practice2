package queue

import (
	"fmt"
	"os"
)

// Queue представляет узел очереди
type Queue struct {
	Data string // Данные узла
	Next *Queue // Следующий узел в очереди
}

// QueueManager представляет структуру управления очередью
type QueueManager struct {
	Head *Queue // Начало очереди
	Tail *Queue // Конец очереди
}

// NewQueue создает и возвращает новый экземпляр управляющей структуры для очереди
func NewQueue() *QueueManager {
	q := &QueueManager{}
	return q
}

// Push добавляет элемент в конец очереди
func (q *QueueManager) Push(value string) string {
	node := &Queue{Data: value, Next: nil}
	if q.Head == nil {
		q.Head = node
		q.Tail = node
		return "Элемент успешно добавлен"
	} else {
		q.Tail.Next = node
		q.Tail = node
		return "Элемент успешно добавлен"
	}
}

// Pop удаляет и возвращает элемент из начала очереди
func (q *QueueManager) Pop() string {
	if q.Head == nil {
		return "Очередь пуста"
	} else {
		temp := q.Head
		q.Head = q.Head.Next
		if q.Head == nil {
			q.Tail = nil
		}
		return "Удалённый элемент:  " + temp.Data
	}
}

// WriteToFile записывает содержимое очереди в файл
func (q *QueueManager) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	current := q.Head
	for current != nil {
		file.WriteString(current.Data + "\n")
		current = current.Next
	}

	return nil
}

// ReadFromFile читает содержимое файла и добавляет элементы в очередь
func (q *QueueManager) ReadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var value string
	for {
		_, err := fmt.Fscanln(file, &value)
		if err != nil {
			break
		}
		q.Push(value)
	}

	return nil
}
