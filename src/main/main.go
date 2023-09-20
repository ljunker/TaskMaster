package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	var tasks = loadTasks()
	fmt.Println(tasks[0])
}

func loadTasks() []Task {
	jsonFile, err := os.Open("tasks.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened tasks")
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var tasks []Task
	err = json.Unmarshal(byteValue, &tasks)
	if err != nil {
		fmt.Println(err)
	}
	return tasks
}
