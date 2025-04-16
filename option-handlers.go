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

func (a *app) createTaskTableWithCells() (*tview.Table, []*tview.TableCell) {
	table := tview.NewTable().
		SetBorders(true)
	word := 0
	tasks := a.listIncompleteInsertionOrder()
	if len(tasks) == 0 {
		table.SetCell(0, 0,
			tview.NewTableCell("No tasks in list").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	cells := []*tview.TableCell{}
	cols, rows := 1, len(tasks)
	for r := 0; r < int(rows); r++ {
		for c := 0; c < int(cols); c++ {
			color := tcell.ColorWhite
			text := ""
			if word < len(tasks) {
				text = tasks[word].String()
			}
			cell := tview.NewTableCell(wrapText(text, 10))
			table.SetCell(r, c,
				cell.SetExpansion(1).
					SetMaxWidth(20).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			cells = append(cells, cell)
			word = (word + 1)
		}
	}
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})
	return table, cells
}

func deleteTaskHandler(ui *ui, app *app) {
	ui.output.Clear()
	//generate task menu for deletion
	deleteMenu, cells := app.createTaskTableWithCells()
	taskList := app.listIncompleteInsertionOrder()
	taskMap := make(map[rune]*task)
	r := 'a'
	for idx, t := range taskList {
		cell := cells[idx]
		taskMap[r] = t
		cell.SetText(string(r) + ") " + cell.Text)
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
	updateMenu, cells := app.createTaskTableWithCells()
	taskList := app.listIncompleteInsertionOrder()
	taskMap := make(map[rune]*task)
	r := 'a'
	for idx, t := range taskList {
		cell := cells[idx]
		taskMap[r] = t
		cell.SetText(string(r) + ") " + cell.Text)
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
