package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNextTaskID(t *testing.T) {
	tests := []struct {
		name string
		in   []*Task
		out  int
	}{
		{
			name: "Empty tasks slice should return 1",
			in:   []*Task{},
			out:  1,
		},
		{
			name: "Single task item should return 2",
			in: []*Task{
				{
					ID:          1,
					Description: "single tasks slice",
					Status:      "todo",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			out: 2,
		},
		{
			name: "task with ID 99 should return 100",
			in: []*Task{
				{
					ID:          99,
					Description: "task with ID 99",
					Status:      "todo",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			out: 100,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			taskID := getNextTaskID(tc.in)
			assert.Equal(t, tc.out, taskID)
		})
	}

}

func TestGetTaskByID(t *testing.T) {
	tasks := []*Task{
		{
			ID:          10,
			Description: "Get task by ID",
			Status:      "todo",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	i, task, err := GetTaskByID(tasks, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, task)
	assert.Equal(t, 0, i)
	assert.Equal(t, 10, task.ID)
}
