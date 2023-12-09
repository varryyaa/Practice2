package main

import (
	"TEST/QuerProc"
	"fmt"
	"net"
)

func main() {
	// Запуск сервера для ожидания соединений по протоколу tcp на порту 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	defer listener.Close() // Закрытие слушателя при завершении работы

	fmt.Println("Сервер ожидает соединений на порту 6379...")

	for {
		conn, err := listener.Accept() // Принятие входящего соединения
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}
		go QuerProc.HandleConnection(conn) // Запуск обработки соединения в отдельной горутине
	}
}
