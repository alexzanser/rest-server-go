package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id int			`json:"id"`
	Text string 	`json:"text"`
	Tags []string 	`json:"tags"`
	Due time.Time 	`json:"due"`
}

type TaskStore struct {
	sync.Mutex

	tasks	map[int]Task
	nextId	int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0

	return ts 
}

func (ts *TaskStore) CreateTask (text string, tags[]string, due time.Time)  int {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id : ts.nextId,
		Text : text,
		Tags : tags,
		Due : due,
	}

	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	
	return task.Id
}

func (ts *TaskStore) GetTask (id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, ok := ts.tasks[id]
	if ok {
		return task, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

func (ts *TaskStore) DeleteTask (id int) error {
	ts.Lock()
	defer ts.Unlock()

	_, ok := ts.tasks[id]
	if ok {
		delete(ts.tasks, id)
	} else {
		return fmt.Errorf("cant't delete task with id=%d (not found)", id)
	}

	return nil
}

func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]Task, 0)
	for _, el := range ts.tasks {
		allTasks = append(allTasks, el)
	}

	return allTasks
}

func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	allTasksTag := make([]Task, 0)

taskloop:
		for _, task := range ts.tasks {
			for _, t := range task.Tags {
				if t == tag {
					allTasksTag = append(allTasksTag, task)
					continue taskloop
				}
			}
		}

	return allTasksTag
}

func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task {
	ts.Lock()
	defer ts.Unlock()

	tasks := make([]Task, 0)

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
