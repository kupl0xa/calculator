package services

import (
	"testing"
)

func TestCreateTask(t *testing.T) {
	service := NewTaskService()
	id := service.CreateTask(10, 5, "+")
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestGetTaskResult(t *testing.T) {
	service := NewTaskService()
	id := service.CreateTask(10, 5, "+")
	service.ExecuteTask(id)
	result, err := service.GetTaskResult(id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}

	result, err = service.GetTaskResult(123)
	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
	if err == nil {
		t.Error("expected error, got no error")
	}
}

func TestListTasks(t *testing.T) {
	service := NewTaskService()
	service.CreateTask(10, 5, "+")
	service.CreateTask(10, 5, "*")
	tasks := service.ListTasks()
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestExecuteTask(t *testing.T) {
	service := NewTaskService()

	id := service.CreateTask(10, 5, "+")
	service.ExecuteTask(id)
	task := service.tasks[id]
	if task.Status != "Completed" || task.Result != 15 {
		t.Errorf("expected Status: Completed, Result: 15, got %+v", task)
	}

	id = service.CreateTask(10, 5, "-")
	service.ExecuteTask(id)
	task = service.tasks[id]
	if task.Status != "Completed" || task.Result != 5 {
		t.Errorf("expected Status: Completed, Result: 5, got %+v", task)
	}

	id = service.CreateTask(10, 5, "*")
	service.ExecuteTask(id)
	task = service.tasks[id]
	if task.Status != "Completed" || task.Result != 50 {
		t.Errorf("expected Status: Completed, Result: 50, got %+v", task)
	}

	id = service.CreateTask(10, 5, "/")
	service.ExecuteTask(id)
	task = service.tasks[id]
	if task.Status != "Completed" || task.Result != 2 {
		t.Errorf("expected Status: Completed, Result: 2, got %+v", task)
	}
}
