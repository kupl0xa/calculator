package services

import (
	"sync"

	"github.com/kupl0xa/calculator/internal/models"
)

type TaskService struct {
	tasks  map[int]*models.Task
	mu     sync.RWMutex
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]*models.Task),
		nextID: 1,
	}
}

func (s *TaskService) CreateTask(x, y int, operator string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := &models.Task{
		X:        x,
		Y:        y,
		Operator: operator,
		Status:   "Pending",
	}
	s.tasks[s.nextID] = task
	s.nextID++
	return s.nextID - 1
}

func (s *TaskService) ListTasks() []models.TaskStatusResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]models.TaskStatusResponse, 0, len(s.tasks))
	for id, task := range s.tasks {
		taskStatus := models.TaskStatusResponse{
			ID:     id,
			Status: task.Status,
		}
		tasks = append(tasks, taskStatus)
	}
	return tasks
}

func (s *TaskService) GetTaskResult(id int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return 0, models.ErrNotFound
	}
	return task.Result, nil
}

func (s *TaskService) ExecuteTask(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := s.tasks[id]

	var result int
	switch task.Operator {
	case "+":
		result = task.X + task.Y
	case "-":
		result = task.X - task.Y
	case "*":
		result = task.X * task.Y
	case "/":
		result = task.X / task.Y
	}
	task.Status = "Completed"
	task.Result = result
}
