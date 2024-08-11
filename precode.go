package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик эндпоинта /tasks, метод GET, возвращаем все задачки
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks) // Формируем JSON из нашей мапы tasks
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем статус 500 при наличии ошибки
		return
	}

	w.Header().Set("Content-type", "application/json") // Добавим в хедер тип контента = JSON
	w.WriteHeader(http.StatusOK)                       // Добавим в хедер 200 статус
	_, err = w.Write(resp)                             // Запишем в ответ наш JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Вернём 400 при ошибке записи
		return
	}
}

// Обработчик эндпоинта /tasks, метод POST, добавляем задачку
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body) // Читаем данные из тела и запишем в буфер
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем 400 при наличии ошибки
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil { // Переводим JSON обратно в структуру
		http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем 400 при наличии ошибки
		return
	}

	// Проверим, есть ли таска в мапе
	if _, innit := tasks[task.ID]; innit {
		_, err = w.Write([]byte("Задача с id = " + string(task.ID) + " уже существует!")) // Напишем юзеру о наличии задачки с таким ID
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Вернём 400 при ошибке записи
			return
		}
		return
	}

	tasks[task.ID] = task // Добавим в мапу новый таск

	w.Header().Set("Content-Type", "application/json")                                   // Добавим в хедер тип контента = JSON
	w.WriteHeader(http.StatusCreated)                                                    // Добавим статус 201, Created
	_, err = w.Write([]byte("Задача с id = " + string(task.ID) + " успешно добавлена!")) // Решил добавить сообщение при успешном добавлении таски
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Вернём 400 при ошибке записи
		return
	}
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Вытянем ID переданную в URL

	task, statusOK := tasks[id] // Поищем таску с нужной нам ID
	if !statusOK {
		http.Error(w, "Нет задачи с id = "+string(id), http.StatusBadRequest) // Вернём 400 при отсутствии таски
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Вернём 400 при ошибке сериализации
		return
	}

	w.Header().Set("Content-type", "application/json") // Добавим в хедер тип контента = JSON
	w.WriteHeader(http.StatusOK)                       // Добавим в хедер 200 статус
	_, err = w.Write(resp)                             // Запишем в ответ наш JSON

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Вернём 400 при ошибке записи
		return
	}
}

func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Вытянем ID переданную в URL

	task, statusOK := tasks[id] // Поищем таску с нужной нам ID
	if !statusOK {
		http.Error(w, "Нет задачи с id = "+string(id), http.StatusBadRequest) // Вернём 400 при отсутствии таски
		return
	}

	delete(tasks, task.ID) // Удаляем из мапы переданную нам таску по ID

	w.Header().Set("Content-Type", "application/json")                                  // Добавим в хедер тип контента = JSON
	w.WriteHeader(http.StatusOK)                                                        // Добавим статус 200, OK
	_, err := w.Write([]byte("Задача с id = " + string(task.ID) + " успешно удалена!")) // Решил добавить сообщение при успешном добавлении таски
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Вернём 400 при ошибке записи
		return
	}

}

func main() {
	r := chi.NewRouter()

	// Обработчик для вывода всех тасок
	r.Get("/tasks", getTasks)

	// Обработчик для добавления таски
	r.Post("/tasks", postTask)

	// Обработчик для вывода таски по ID
	r.Get("/tasks/{id}", getTaskById)

	// Обработчик для удаления таски по ID
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
