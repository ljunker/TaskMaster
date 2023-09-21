package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var taskFile = "tasks.json"

func main() {
	var tasks = loadTasks()
	reader := bufio.NewReader(os.Stdin)
	shouldRun := true
	for shouldRun {
		for _, task := range tasks {
			fmt.Println(task)
		}
		fmt.Println("------------------------------------")
		fmt.Println("(s) to save\n(n) for a new task\n(c) to complete\n(d) to delete\n(q) to quit")
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = cleanString(text)
		switch text {
		case "s":
			fmt.Println("Saving...")
			saveTasks(tasks)
		case "n":
			text, tasks = createTask(text, reader, tasks)
		case "c":
			text = toggleCompletion(text, reader, tasks)
		case "d":
			deleteTask(text, reader, tasks)
		case "q":
			shouldRun = false

		}
	}
}

func createTask(text string, reader *bufio.Reader, tasks []Task) (string, []Task) {
	fmt.Println("Enter name of task:")
	fmt.Print("-> ")
	text, _ = reader.ReadString('\n')
	text = cleanString(text)
	t := Task{
		Id:          getHighestId(tasks) + 1,
		Content:     text,
		DateCreated: DateTime{time.Now()},
		Completed:   false,
	}
	tasks = append(tasks, t)
	fmt.Println("Task created")
	return text, tasks
}

func toggleCompletion(text string, reader *bufio.Reader, tasks []Task) string {
	fmt.Println("Which task is complete? Id:")
	fmt.Print("-> ")
	text, _ = reader.ReadString('\n')
	text = cleanString(text)
	id, err := strconv.ParseUint(text, 10, 64)
	check(err)
	for i, task := range tasks {
		if task.Id == id {
			task.Completed = !task.Completed
			tasks[i] = task
			fmt.Println("toggled status for id " + strconv.FormatUint(id, 10))
			break
		}
	}
	return text
}

func deleteTask(text string, reader *bufio.Reader, tasks []Task) {
	fmt.Println("Which task to delete? Id:")
	fmt.Print("-> ")
	text, _ = reader.ReadString('\n')
	text = cleanString(text)
	id, err := strconv.ParseUint(text, 10, 64)
	check(err)
	for i, task := range tasks {
		if task.Id == id {
			tasks = remove(tasks, i)
			fmt.Println("deleted task with id " + strconv.FormatUint(id, 10))
			break
		}
	}
}

// helper functions
func remove(tasks []Task, i int) []Task {
	tasks[i] = tasks[len(tasks)-1]
	return tasks[:len(tasks)-1]
}

func cleanString(text string) string {
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	return text
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getHighestId(tasks []Task) uint64 {
	var highest uint64 = 0
	for _, task := range tasks {
		if task.Id > highest {
			highest = task.Id
		}
	}
	return highest
}

func saveTasks(tasks []Task) {
	js, _ := json.Marshal(tasks)
	err := os.WriteFile(taskFile, js, 0644)
	check(err)
}

func loadTasks() []Task {
	jsonFile, err := os.Open(taskFile)
	if err != nil {
		_, err := os.Create(taskFile)
		check(err)
		fmt.Println(err)
		return []Task{}
	}
	fmt.Println("Successfully Opened tasks")
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		check(err)
	}(jsonFile)

	byteValue, _ := io.ReadAll(jsonFile)
	var tasks []Task
	err = json.Unmarshal(byteValue, &tasks)
	if err != nil {
		fmt.Println(err)
	}
	return tasks
}
