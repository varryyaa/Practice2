package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Client представляет клиента
type Client struct {
	conn net.Conn
}

// Init инициализирует клиента, устанавливая соединение с сервером
func (c *Client) Init() {
	// Подключение к серверу на локальной машине и порту 6379
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return
	}
	c.conn = conn
}

// SendCommand отправляет команду на сервер
func (c *Client) SendCommand(command string) {
	// Отправка команды на сервер, добавление символа новой строки
	_, err := c.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Ошибка при отправке команды на сервер:", err)
		return
	}
}

// ReceiveResponse читает ответ от сервера и возвращает его
func (c *Client) ReceiveResponse() string {
	// Создание буфера для чтения ответа от сервера
	buffer := make([]byte, 1024)
	// Чтение ответа от сервера
	n, err := c.conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа от сервера:", err)
		return ""
	}
	// Преобразование ответа в строку и удаление лишних пробелов
	response := strings.TrimSpace(string(buffer[:n]))
	return response
}

// Close закрывает соединение клиента с сервером
func (c *Client) Close() {
	c.conn.Close()
}

func main() {
	// Создание экземпляра клиента
	client := &Client{}
	// Инициализация клиента и отложенное закрытие соединения при завершении main()
	client.Init()
	defer client.Close()

	// Создание сканера для считывания ввода пользователя
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Ввод пользователя
		fmt.Print("Введите команду (или 'QUIT' для выхода): ")
		scanner.Scan()
		command := strings.TrimSpace(scanner.Text())

		// Отправка команды на сервер
		client.SendCommand(command)

		// Проверка условия выхода
		if strings.ToUpper(command) == "QUIT" {
			break
		}

		// Чтение и вывод ответа от сервера
		response := client.ReceiveResponse()
		fmt.Println("Ответ сервера:", response)
	}
}
