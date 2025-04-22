package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func addtaskHandler(ui *ui, app *app) {
	ui.output.Clear()
	nameVal := ""
	dateIdx := 0

	taskForm := tview.NewForm()
	taskForm.AddInputField("task name", "", 20, nil, func(text string) { nameVal = text })

	taskForm.AddDropDown("date to appear on calendar", generateDaysList(time.Now(), 14), 0, func(option string, optionIndex int) {
		dateIdx = optionIndex
	})

	taskForm.AddButton("Submit", func() {
		addTask(app, nameVal, addDayToDate(time.Now(), dateIdx))
		ui.messageContainer.SetText("Task Added").SetTextColor(tcell.ColorDarkGreen)
		ui.ResetUi(app)
	})

	handler := taskForm.GetFormItem(1).InputHandler()
	if handler != nil {
		enterKey := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
		setFocus := func(p tview.Primitive) {
			ui.app.SetFocus(taskForm.GetButton(0))
		}
		handler(enterKey, setFocus)
	}
	ui.output.AddItem(taskForm, 0, 1, true)
	ui.app.SetFocus(taskForm)
}

func (a *app) createTaskTableWithCells(showComplete bool, showFutureTasks bool) (*tview.Table, []*tview.TableCell) {
	table := tview.NewTable().
		SetBorders(true).SetSelectable(false, false)
	tasks := a.listInsertionOrder(showComplete, showFutureTasks)
	if len(tasks) == 0 {
		table.SetCell(0, 0,
			tview.NewTableCell("No tasks in list").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	cells := []*tview.TableCell{}
	headers := []string{"ID", "Task", "Status"}
	for i, h := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(h).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	for r := 1; r < len(tasks)+1; r++ {
		t := tasks[r-1]
		complete := ""
		if t.isComplete() {
			complete = "yes"
		} else {
			complete = "no"
		}
		fields := []string{"key", t.Name, complete}
		for c := 0; c < len(headers); c++ {
			color := tcell.ColorWhite
			cell := tview.NewTableCell(fields[c])
			table.SetCell(r, c,
				cell.SetExpansion(1).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			cells = append(cells, cell)
		}
	}
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(true, false)
	})
	return table, cells
}

func deleteTaskHandler(ui *ui, app *app) {
	ui.output.Clear()
	//generate task menu for deletion
	deleteMenu, _ := app.createTaskTableWithCells(true, true)
	showComplete := true
	showFutureTasks := false
	taskList := app.listInsertionOrder(showComplete, showFutureTasks)
	taskMap := make(map[rune]*task)
	r := 'a'
	for idx, t := range taskList {
		cell := deleteMenu.GetCell(idx+1, 0)
		taskMap[r] = t
		cell.SetText(string(r) + ") ")
		r += 1
		if r == 'z'+1 {
			r = 'A'
		}
	}
	var selectedTask rune
	//run delete task function
	deleteMenu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune(tcell.KeyEnter) {
			if selectedTask == 0 {
				return event
			}
			t, ok := taskMap[selectedTask]
			if ok {
				removeTask(app, t.Id)
				ui.messageContainer.SetText("Task Deleted").SetTextColor(tcell.ColorYellow)
			}
			ui.ResetUi(app)
		} else {
			selectedTask = event.Rune()
		}
		return event
	})
	//reset
	ui.output.AddItem(deleteMenu, 0, 2, true)
	ui.app.SetFocus(deleteMenu)
}
func updateTaskHandler(ui *ui, app *app) {
	ui.output.Clear()
	//generate task menu for deletion
	updateMenu, _ := app.createTaskTableWithCells(false, true)
	showComplete := true
	showFutureTasks := false
	taskList := app.listInsertionOrder(showComplete, showFutureTasks)
	taskMap := make(map[rune]*task)
	r := 'a'
	for idx, t := range taskList {
		cell := updateMenu.GetCell(idx+1, 0)
		taskMap[r] = t
		cell.SetText(string(r) + ") ")
		r += 1
		if r == 'z'+1 {
			r = 'A'
		}
	}
	var selectedTask rune
	//run delete task function
	updateMenu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune(tcell.KeyEnter) {
			if selectedTask == 0 {
				return event
			}
			t, ok := taskMap[selectedTask]
			if ok {
				updateTask(app, t)
				ui.messageContainer.SetText("Task Updated").SetTextColor(tcell.ColorYellow)
			}
			ui.ResetUi(app)
		} else {
			selectedTask = event.Rune()
		}
		return event
	})
	//reset
	ui.output.AddItem(updateMenu, 0, 2, true)
	ui.app.SetFocus(updateMenu)
}
