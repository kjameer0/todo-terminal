package main

import (
	"fmt"
	"time"
)

// if completion date is zero value then the task is incomplete
type task struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	CompletionDate time.Time `json:"completionDate"`
	BeginDate      time.Time `json:"beginDate"`
}

func (t *task) isComplete() bool {
	return !t.CompletionDate.IsZero()
}

func (t *task) String() string {
	var completed string
	if !t.isComplete() {
		completed = "❌"
	} else {
		completed = "✅"
	}
	var completionDate string
	if t.CompletionDate.IsZero() {
		completionDate = ""
	} else {
		completionDate = monthDayYear(t.CompletionDate)
	}
	return fmt.Sprintf("%s %s %s", t.Name, completed, completionDate)
}

func addTask(a *app, taskText string, beginTime time.Time) {
	addedTask := newTask(taskText, beginTime)
	a.Tasks[addedTask.Id] = addedTask
	a.InsertionOrder = append(a.InsertionOrder, addedTask.Id)
	saveToFile(a)
}

func removeTask(a *app, taskId string) bool {
	if _, ok := a.Tasks[taskId]; !ok {
		fmt.Println("hi")
		return false
	}
	delete(a.Tasks, taskId)
	//remove deleted id from insertion order
	filteredInsertionOrder := []string{}
	for _, id := range a.InsertionOrder {
		if id == taskId {
			continue
		}
		filteredInsertionOrder = append(filteredInsertionOrder, id)
	}
	a.InsertionOrder = filteredInsertionOrder
	saveToFile(a)
	return true
}

func removeAllTasks(a *app) {
	a.InsertionOrder = []string{}
	clear(a.Tasks)
	saveToFile(a)
}

func listTasks(a *app) []string {
	tasks := []string{}
	for _, taskId := range a.InsertionOrder {
		if taskId == "" {
			continue
		}
		curTask := a.Tasks[taskId]
		//show a task if it not complete or if show complete and task
		if !a.config.ShowComplete && curTask.isComplete() {
			continue
		}
		if time.Now().Compare(curTask.BeginDate) == -1 {
			continue
		}
		var completed string
		if !curTask.isComplete() {
			completed = "❌"
		} else {
			completed = "✅"
		}
		t := monthDayYear(curTask.CompletionDate)
		if curTask.CompletionDate.IsZero() {
			t = ""
		}
		tasks = append(tasks, fmt.Sprintf("%s %s %s", curTask.Name, completed, t))
	}
	return tasks
}

func updateTask(a *app, t *task) {
	var zeroTime time.Time
	if !t.isComplete() {
		t.CompletionDate = time.Now()
	} else {
		t.CompletionDate = zeroTime
	}
	saveToFile(a)
}
