# Taskctl

`taskctl` is a simple cli based task tracker

## Requirements

- [ ] Add, Update, and Delete tasks
- [ ] Mark a task as in progress or completed
- [ ] List all tasks
- [ ] List all completed tasks
- [ ] List all in progress tasks
- [ ] List all todo tasks

## Implementation

- Use a JSON file to store the tasks
- The JSON file should be created if not exists
- Use native std lib os/file packages

```
# Adding a new task

$ taskctl add "mow the yard"

# Output: Task added successfully (ID: 1)


# Updating and deleting tasks

$ taskctl update 1 "buy groceries"
$ taskctl delete 1

# Marking a task as in progress or done

$ taskctl mark-in-progress 1
$ taskctl mark-done 1

# Listing all tasks

$ taskctl list

# List tasks by status

$ taskctl list done
$ taskctl list in-progress
$ taskctl list todo
```

## Task Properties

Each task should have the following properties

- `id`
- `description`
- `status`
- `createdAt`
- `updatedAt`
