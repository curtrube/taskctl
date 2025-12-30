package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	file = "tasks.json"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	tasks := readTasks()

	cmd := os.Args[1]
	switch cmd {
	case "add":
		addCmd(os.Args[2:], tasks)
	case "list":
		listCmd(os.Args[2:])
	case "delete":
		deleteCmd(os.Args[2:])
	}
}

// GetNextTaskID returns the next available task ID.
func GetNextTaskID(tasks []*Task) int {
	taskID := 1
	// If tasks is not zero
	if len(tasks) > 0 {
		// loop and check each task id
		for _, t := range tasks {
			if t.ID > taskID {
				taskID = t.ID
			}
		}
		taskID++
	}
	return taskID
}

func addCmd(args []string, tasks []*Task) {
	if len(args) < 1 {
		log.Println("add cmd missing args TODO: Add Add help")
		os.Exit(0)
	}
	log.Printf("adding a new task...")
	// Get the next available task ID
	taskID := GetNextTaskID(tasks)
	task := NewTask(taskID, args[0])
	fmt.Printf("adding new task: %+v", task)
	tasks = append(tasks, task)
	writeTasks(tasks)
}

func listCmd(args []string) {
	fmt.Println("Listing tasks..")
	fmt.Println(args)
}

func deleteCmd(args []string) {
	fmt.Println("Deleting task with id..")
	fmt.Println(args)
}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` //TODO: make status an enum
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Status is todo, in-progress, done

func NewTask(taskID int, description string) *Task {
	// NewTask should be able to determine the next ID
	return &Task{
		ID:          taskID,
		Description: description,
		Status:      "todo", // todo/new, in-progress, done
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false // file does not exist
	}
	if err != nil {
		// some other error occured e.g. permissions issue
		return false
	}
	return !info.IsDir() // return true if it's a file, false if directory
}

func readFile(filename string) []byte {
	if !fileExists(filename) {
		log.Printf("file not found: %s", filename)
		dir, err := os.Getwd()
		path := filepath.Join(dir, file)
		check(err)

		log.Println("creating file:", path)
		f, err := os.Create(path)
		check(err)

		f.Close()
	}
	data, err := os.ReadFile(filename)
	check(err)
	return data
}

func readTasks() []*Task {
	taskBytes := readFile(file)
	var tasks []*Task
	json.Unmarshal([]byte(taskBytes), &tasks)
	return tasks
}

func writeFile(data []byte) {
	dir, err := os.Getwd()
	check(err)
	path := filepath.Join(dir, file)
	err = os.WriteFile(path, data, 0644)
	log.Println("writing file: ", path)
	check(err)
}

func writeTasks(tasks []*Task) {
	bytes, err := json.Marshal(tasks)
	check(err)
	writeFile(bytes)
}

func printHelp() {
	fmt.Println("taskctl - A simple CLI for managing tasks")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  taskctl  [commands] [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add      Add a new task")
	fmt.Println("  list     List all tasks")
	fmt.Println("  delete   Delete a task")
}
