package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const fileName = "todos.json"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func loadList() ([]Task, error) {
	var tasks []Task
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return tasks, err
	}
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveList(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, data, 0644)
	return err
}

func listTask(tasks []Task) string {
	if len(tasks) == 0 {
		return "Лист задач пуст"
	}
	var sb strings.Builder
	for _, task := range tasks {
		status := " "
		if task.Done {
			status = "x"
		}
		fmt.Fprintf(&sb, "[%s] %d %s\n", status, task.ID, task.Text)
	}
	return sb.String()

}

func addTask(tasks []Task, text string) []Task {
	var task Task
	if len(tasks) == 0 {
		task = Task{1, text, false}
	} else {
		task = Task{tasks[len(tasks)-1].ID + 1, text, false}
	}
	tasks = append(tasks, task)
	return tasks
}

func doneTask(tasks []Task, index int) ([]Task, error) {

	i := findTaskIndex(tasks, index)
	if i == -1 {
		return tasks, fmt.Errorf("такой задачи не существует")
	}

	tasks[i].Done = !tasks[i].Done
	return tasks, nil

}

func deleteTask(tasks []Task, index int) ([]Task, error) {
	i := findTaskIndex(tasks, index)
	if i == -1 {
		return tasks, fmt.Errorf("такой задачи не существует")
	}
	tasks = append(tasks[:i], tasks[i+1:]...)
	return tasks, nil

}

func main() {
	arrayTasks, err := loadList()
	if err != nil {
		fmt.Println("Ошибка загрузки файла: ", err)
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("Использование:\nlist - увидеть список задач\nadd <text> - добавить новую задачу\ndone <num> - изменить статус задачи\ndelete <num> - удалить задачу")
		return
	}
	switch os.Args[1] {
	case "list":
		fmt.Println(listTask(arrayTasks))
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("необходимо указать название задачи")
			return
		}
		arrayTasks = addTask(arrayTasks, os.Args[2])
		fmt.Println("Задача успешно добавлена!")
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("необходимо указать номер задачи")
			return
		}
		index, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Тут необходимо ввести число")
			return
		}
		if arrayTasks, err = doneTask(arrayTasks, index); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Состояние задачи #%s успешно изменено\n", os.Args[2])

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("необходимо указать номер задачи")
			return
		}
		index, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Тут необходимо ввести число")
			return
		}
		if arrayTasks, err = deleteTask(arrayTasks, index); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Задача успешно удалена")

	}
	if err := saveList(arrayTasks); err != nil {
		fmt.Println("Ошибка сохранения файла: ", err)
		return
	}
}

func findTaskIndex(tasks []Task, id int) int {
	if id < 0 {
		return -1
	}
	for i, value := range tasks {
		if value.ID == id {
			return i
		}
	}
	return -1
}
