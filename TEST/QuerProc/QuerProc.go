package QuerProc

import (
	"TEST/hashtable"
	"TEST/queue"
	"TEST/set"
	"TEST/stack"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Структура DBLock для управления блокировками на структуре базы данных
type DBLock struct {
	setLock       sync.Mutex // Мьютекс для операций с множеством
	queueLock     sync.Mutex // Мьютекс для операций с очередью
	stackLock     sync.Mutex // Мьютекс для операций со стеком
	hashTableLock sync.Mutex // Мьютекс для операций с хэш-таблицей
}

// Глобальная переменная для управления блокировками
var dbLock DBLock

// Функция HandleConnection обрабатывает входящие подключения
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Буфер для чтения данных из соединения
	buffer := make([]byte, 1024)

	for {
		// Чтение команды от клиента
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Ошибка чтения данных из подключения:", err)
			return
		}

		// Парсинг строки
		params := strings.Fields(string(buffer[:n]))

		// Вывод полученной строки
		fmt.Println("Получен запрос:", string(buffer[:n]))

		// Обработка запроса
		switch params[0] {
		// Обработка множества
		case "SADD":
			if len(params) != 3 { //при неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue
			}
			// Получаем блокировку для операций с множеством для обеспечения взаимного исключения
			dbLock.setLock.Lock()
			// Создаем новый экземпляр множества
			set := set.NewSet()
			// Читаем содержимое множества из файла, указанного клиентом
			set.ReadFromFile(params[1])

			// Выполняем операцию "Push" на множестве с предоставленным параметром
			// и отправляем результат клиенту в виде среза байтов
			_, err := conn.Write([]byte(set.Push(params[2])))
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			// Записываем обновленное множество обратно в файл
			set.WriteToFile(params[1])
			// Освобождаем блокировку для операций с множеством
			dbLock.setLock.Unlock()

		case "SREM":
			if len(params) != 3 {
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue
			}
			dbLock.setLock.Lock()       // Получаем блокировку для операций с множеством для обеспечения взаимного исключения
			set := set.NewSet()         // Создаем новый экземпляр множества
			set.ReadFromFile(params[1]) // Читаем содержимое множества из файла, указанного клиентом
			// Выполняем операцию "pop" на множестве с предоставленным параметром
			// и отправляем результат клиенту в виде среза байтов
			_, err := conn.Write([]byte(set.Pop(params[2])))
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			set.WriteToFile(params[1]) // Записываем обновленное множество обратно в файл
			dbLock.setLock.Unlock()    // Освобождаем блокировку для операций с множеством

		case "SISMEMBER":
			if len(params) != 3 {
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue
			}
			dbLock.setLock.Lock()       // Получаем блокировку для операций с множеством для обеспечения взаимного исключения
			set := set.NewSet()         // Создаем новый экземпляр множества
			set.ReadFromFile(params[1]) // Читаем содержимое множества из файла, указанного клиентом
			// Выполняем операцию "search" на множестве с предоставленным параметром
			// и отправляем результат клиенту в виде среза байтов
			_, err := conn.Write([]byte(set.Search(params[2])))
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			dbLock.setLock.Unlock() // Освобождаем блокировку для операций с множеством

		// Обработка очереди
		case "QPUSH":
			if len(params) != 3 { //При неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue

			}
			dbLock.queueLock.Lock()                                // Получаем блокировку для операций с очередью для обеспечения взаимного исключения
			qManager := queue.NewQueue()                           //Создаем новый экземпляр очереди
			qManager.ReadFromFile(params[1])                       //Читаем содержимое очереди из файла, указанного клиентом
			_, err := conn.Write([]byte(qManager.Push(params[2]))) //Выполняем операцию Push и отправляем результат клиенту
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			qManager.WriteToFile(params[1]) //Записываем обновленную очередь в файл
			dbLock.queueLock.Unlock()       // Освобождаем блокировку для операций с очередью

		case "QPOP":
			if len(params) != 2 { //При неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue

			}
			dbLock.queueLock.Lock()                      // Получаем блокировку для операций с очередью для обеспечения взаимного исключения
			qManager := queue.NewQueue()                 //Создаем новый экземпляр очереди
			qManager.ReadFromFile(params[1])             //Читаем содержимое очереди из файла, указанного клиентом
			_, err := conn.Write([]byte(qManager.Pop())) //Выполняем операцию Pop и отправляем результат клиенту
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			qManager.WriteToFile(params[1]) //Записываем обновленную очередь в файл
			dbLock.queueLock.Unlock()       // Освобождаем блокировку для операций с очередью

		// Обработка стека
		case "SPUSH":
			if len(params) != 3 { //При неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue

			}
			dbLock.stackLock.Lock()                             // Получаем блокировку для операций со стеком для обеспечения взаимного исключения
			stack := stack.NewStack()                           //Создаем новый экземпляр стека
			stack.ReadFromFile(params[1])                       //Читаем содержимое стека из файла, указанного клиентом
			_, err := conn.Write([]byte(stack.Push(params[2]))) //Выполняем операцию Push и отправляем результат клиенту
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			stack.WriteToFile(params[1]) //Записываем обновленный стек в файл
			dbLock.stackLock.Unlock()    // Освобождаем блокировку для операций со стеком

		case "SPOP":
			if len(params) != 2 { //При неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue
			}
			dbLock.stackLock.Lock()                   // Получаем блокировку для операций со стеком для обеспечения взаимного исключения
			stack := stack.NewStack()                 //Создаем новый экземпляр стека
			stack.ReadFromFile(params[1])             //Читаем содержимое стека из файла, указанного клиентом
			_, err := conn.Write([]byte(stack.Pop())) //Выполняем операцию Pop и отправляем результат клиенту
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			stack.WriteToFile(params[1]) //Записываем обновленный стек в файл
			dbLock.stackLock.Unlock()    // Освобождаем блокировку для операций со стеком

		// Обработка хэш-таблицы
		case "HSET":
			if len(params) != 4 { //При неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue
			}
			dbLock.hashTableLock.Lock()                                        // Получаем блокировку для операций с таблицей для обеспечения взаимного исключения
			hashTable := hashtable.NewHashTable()                              //Создаем новый экземпляр таблицы
			hashTable.ReadFromFile(params[1])                                  //Читаем содержимое таблицы из файла, указанного клиентом
			_, err := conn.Write([]byte(hashTable.Push(params[2], params[3]))) //Выполняем операцию Push и отправляем результат клиенту
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			hashTable.WriteToFile(params[1]) //Записываем обновленную таблицу в файл
			dbLock.hashTableLock.Unlock()    // Освобождаем блокировку для операций с таблицей

		case "HDEL":
			if len(params) != 3 { //При неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue
			}
			dbLock.hashTableLock.Lock()                            // Получаем блокировку для операций с таблицей для обеспечения взаимного исключения
			hashTable := hashtable.NewHashTable()                  //Создаем новый экземпляр таблицы
			hashTable.ReadFromFile(params[1])                      //Читаем содержимое таблицы из файла, указанного клиентом
			_, err := conn.Write([]byte(hashTable.Pop(params[2]))) //Выполняем операцию Pop и отправляем результат клиенту
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			hashTable.WriteToFile(params[1]) //Записываем обновленную таблицу в файл
			dbLock.hashTableLock.Unlock()    // Освобождаем блокировку для операций с таблицей

		case "HGET":
			if len(params) != 3 { //При неверно введенном количестве параметров
				fmt.Println("Некорректное количество параметров.")
				sendEndResponse(conn)
				continue

			}
			dbLock.hashTableLock.Lock()                               // Получаем блокировку для операций с таблицей для обеспечения взаимного исключения
			hashTable := hashtable.NewHashTable()                     //Создаем новый экземпляр таблицы
			hashTable.ReadFromFile(params[1])                         //Читаем содержимое таблицы из файла, указанного клиентом
			_, err := conn.Write([]byte(hashTable.Search(params[2]))) //Выполняем операцию Search и отправляем результат клиенту
			if err != nil {
				fmt.Println("Ошибка отправки ответа клиенту:", err)
			}
			dbLock.hashTableLock.Unlock() // Освобождаем блокировку для операций с таблицей

		default:
			fmt.Println("Неподдерживаемая операция.")
			sendEndResponse(conn)
		}
		buffer = make([]byte, 1024) //новый буфер
	}
}

func sendEndResponse(conn net.Conn) {
	_, err := conn.Write([]byte("END\n"))
	if err != nil {
		fmt.Println("Ошибка отправки ответа клиенту:", err)
	}
}
