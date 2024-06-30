package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID          string   `json:"id"`          // ID задачи
	Description string   `json:"description"` // Заголовок
	Note        string   `json:"note"`        // Описание задачи
	Application []string `json:"application"` // Приложения, которыми будете пользоваться
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func main() {
	r := chi.NewRouter()
	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `createTask`
	r.Post("/tasks", createTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", getTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTask`
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	_, errMessage := w.Write(resp)
	if errMessage != nil {
		fmt.Printf("%v", errMessage)
		return
	}

}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task        // создаем новый экземпляр структуры Task
	var buf bytes.Buffer // создаем буфер для хранения полученных данных

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, ok := tasks[task.ID] // ищем задачу в слайсе, проверяя ID
	if ok {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}
	// добавляем задачу в слайс
	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // получаем ID задачи из URL-параметра
	task, ok := tasks[id]       // ищем задачу в слайсе
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task) // сериализуем данные из слайса tasks
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, errMessage := w.Write(resp)
	if errMessage != nil {
		fmt.Printf("%v", errMessage)
		return
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // получаем ID задачи из URL-параметра
	_, ok := tasks[id]          //	ищем задачу в слайсе
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
