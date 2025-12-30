package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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
		listCmd(os.Args[2:], tasks)
	case "delete":
		deleteCmd(os.Args[2:], tasks)
	}
}

// GetNextTaskID returns the next available task ID.
func getNextTaskID(tasks []*Task) int {
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

// GetTaskByID retrieves a task by the provided ID.
func GetTaskByID(tasks []*Task, ID int) (int, *Task, error) {
	for i, t := range tasks {
		if t.ID == ID {
			return i, t, nil
		}
	}
	return 0, &Task{}, fmt.Errorf("No matching task found with ID: %d", ID)
}

func addCmd(args []string, tasks []*Task) {
	if len(args) < 1 {
		log.Println("add cmd missing args TODO: Add Add help")
		os.Exit(0)
	}
	log.Printf("adding a new task...")
	// Get the next available task ID
	taskID := getNextTaskID(tasks)
	task := NewTask(taskID, args[0])

	tasks = append(tasks, task)
	writeTasks(tasks) // TODO: interesting this doesn't return an error

	//fmt.Printf("adding new task: %+v", task) // TODO convert this to debug log
	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
}

func convertStringToInt(s string) (int, error) {
	int, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("convert string to int: %v", err)
	}
	return int, nil
}

func listCmd(args []string, tasks []*Task) {
	log.Print("Listing tasks...")
	// TODO: if length args > 0 match taskID
	if len(args) > 0 {
		// get the task id from args
		ID, err := convertStringToInt(args[0])
		if err != nil {
			log.Fatal(err)
		}
		_, task, err := GetTaskByID(tasks, ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", task)

	} else {
		fmt.Println("| ID | Description | Status | Created | LastUpdated |")
		for _, t := range tasks {
			fmt.Printf("| %d  | %s | %s | %s | %s |\n", t.ID, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
		}
	}
}

func deleteCmd(args []string, tasks []*Task) {
	log.Print("Deleteing task...")
	fmt.Println(args)

	// Delete requires a param of ID
	ID, err := convertStringToInt(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// Should we check to make sure the task exists first?
	// Not necessarily, but we will need the index for the task.
	i, task, err := GetTaskByID(tasks, ID)
	if err != nil {
		log.Fatal(err)
	}

	// If so, go ahead and delete by calling slice.Delete()
	log.Printf("Deleting task with ID: %d at index: %d", task.ID, i)
	updatedTasks := slices.Delete(tasks, i, i+1)

	// Write the tasks back to the file
	writeTasks(updatedTasks)
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
