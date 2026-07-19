package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

func listTask(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("Лист задач пуст")
		return
	}
	for _, task := range tasks {
		status := " "
		if task.Done {
			status = "x"
		}
		fmt.Printf("[%s] %d %s\n", status, task.ID, task.Text)
	}
}

func addTask(tasks []Task, text string) []Task {
	var task Task
	if len(tasks) == 0 {
		task = Task{1, text, false}
	} else {
		task = Task{tasks[len(tasks)-1].ID + 1, text, false}
	}
	tasks = append(tasks, task)
	fmt.Println("Задача успешно добавлена!")
	return tasks
}

func doneTask(tasks []Task, indexStr string) []Task {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		fmt.Println("Тут необходимо ввести число")
		return tasks
	}

	i := findTaskIndex(tasks, index)
	if i == -1 {
		fmt.Println("Такой задачи не существует")
		return tasks
	}

	tasks[i].Done = !tasks[i].Done
	fmt.Printf("состояние задачи №%d изменено на %t\n", index, tasks[i].Done)
	return tasks

}

func deleteTask(tasks []Task, indexStr string) []Task {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		fmt.Println("Тут необходимо ввести число")
		return tasks
	}
	i := findTaskIndex(tasks, index)
	if i == -1 {
		fmt.Println("Такой задачи не существует")
		return tasks
	}
	tasks = append(tasks[:i], tasks[i+1:]...)
	fmt.Println("Задача успешно удалена!")
	return tasks

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
		listTask(arrayTasks)
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("необходимо указать название задачи")
			return
		}
		arrayTasks = addTask(arrayTasks, os.Args[2])
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("необходимо указать номер задачи")
			return
		}
		arrayTasks = doneTask(arrayTasks, os.Args[2])
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("необходимо указать номер задачи")
			return
		}
		arrayTasks = deleteTask(arrayTasks, os.Args[2])
	}
	if err := saveList(arrayTasks); err != nil {
		fmt.Println("Ошибка сохранения файла: ", err)
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
