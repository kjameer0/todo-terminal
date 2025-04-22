package main

import (
	"log"
	"os"
	"time"

	"github.com/aidarkhanov/nanoid"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func wrapText(s string, limit int) string {
	var result string
	for i, char := range s {
		if i > 0 && i%limit == 0 {
			result += "\n"
		}
		result += string(char)
	}
	return result
}

type taskDate struct {
	t time.Time
}

func (t taskDate) String() string {
	return monthDayYear(t.t)
}

type stringWrapper string

func (s stringWrapper) String() string {
	return string(s)
}

const CHECK_TASKS stringWrapper = "Check tasks"
const UPDATE_TASK stringWrapper = "Update tasks"
const ADD_A_TASK stringWrapper = "Add a task"
const DELETE_A_TASK stringWrapper = "Delete a specific task"
const DELETE_ALL_TASKS stringWrapper = "Delete -all- tasks"
const QUIT stringWrapper = "Quit"

var options = []stringWrapper{CHECK_TASKS, UPDATE_TASK, ADD_A_TASK, DELETE_A_TASK, DELETE_ALL_TASKS, QUIT}

type app struct {
	Tasks          map[string]*task `json:"tasks"`
	InsertionOrder []string         `json:"insertionOrder"`
	saveLocation   string
	configPath     string
	config         *config
}

func newApp() *app {
	tasks := make(map[string]*task, 100)
	return &app{Tasks: tasks}
}

func newTask(name string, beginDate time.Time) *task {
	if name == "" {
		log.Fatal("a task must have a name")
	}
	var taskId string
	taskId, err := nanoid.Generate(nanoid.DefaultAlphabet, 20)
	if err != nil {
		log.Fatal(err)
	}
	t := &task{Id: taskId, Name: name, BeginDate: beginDate}
	return t
}

type messageText struct {
	text  string
	color tcell.Color
}

func exitCleanup(a *app) {
	os.Exit(0)
}

func (a *app) listInsertionOrder(showComplete bool, showFutureTasks bool) []*task {
	tasks := make([]*task, 0, len(a.InsertionOrder))
	for _, t := range a.InsertionOrder {
		curTask := a.Tasks[t]
		if !showComplete && curTask.isComplete() {
			continue
		}
		if time.Now().Compare(curTask.BeginDate) == -1 && showFutureTasks {
			continue
		}
		tasks = append(tasks, curTask)
	}
	return tasks
}

type ui struct {
	app              *tview.Application
	optionsMenu      *tview.List
	output           *tview.Flex
	messageContainer *tview.TextView
}

func (a *app) createTaskTable() *tview.Table {
	table := tview.NewTable().
		SetBorders(true)
	word := 0
	showComplete := false
	showFutureTasks := false
	tasks := a.listInsertionOrder(showComplete, showFutureTasks)
	if len(tasks) == 0 {
		table.SetCell(0, 0,
			tview.NewTableCell("No tasks in list").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	cols, rows := 1, len(tasks)
	for r := 0; r < int(rows); r++ {
		for c := 0; c < int(cols); c++ {
			color := tcell.ColorWhite
			text := ""
			if word < len(tasks) {
				text = tasks[word].String()
			}
			table.SetCell(r, c,
				tview.NewTableCell(wrapText(text, 10)).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			word = (word + 1)
		}
	}
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})
	return table
}

func generateDaysList(t time.Time, numDays int) []string {
	dateList := []string{}
	for i := 0; i < numDays; i++ {
		addedDay := addDayToDate(t, i)
		dateList = append(dateList, monthDayYear(addedDay))
	}
	return dateList
}

type handler struct {
	Label    string
	Shortcut rune
	Action   func()
}

func (u *ui) ResetUi(a *app) {
	u.output.Clear()
	u.output.AddItem(a.createTaskTable(), 0, 2, false)
	u.app.SetFocus(u.optionsMenu)
	u.optionsMenu.SetCurrentItem(0)
}

func generateOptionsHandlers(ui *ui, app *app) []handler {
	output := ui.output
	handlers := []handler{
		{"List Tasks", 'a', func() {
			cells, _ := app.createTaskTableWithCells(false, false)
			output.Clear().AddItem(cells, 0, 1, false)
		},
		},
		{"Add Task", 'b', func() { addtaskHandler(ui, app) }},
		{"Delete Task", 'c', func() { deleteTaskHandler(ui, app) }},
		{"Update Task", 'd', func() { updateTaskHandler(ui, app) }},
	}
	return handlers
}

func main() {
	a := newApp()
	a.saveLocation = "./tasks.json"
	a.configPath = "./config.json"
	c, err := a.loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	a.config = c
	readTasksFromFile(a)
	ui := &ui{app: tview.NewApplication()}
	output := tview.NewFlex()
	optionsMenu := tview.NewList()
	tApp := ui.app
	ui.output = output
	ui.optionsMenu = optionsMenu

	handlers := generateOptionsHandlers(ui, a)
	for _, opt := range handlers {
		action := opt.Action
		optionsMenu.AddItem(opt.Label, "", opt.Shortcut, action)
	}
	table, _ := a.createTaskTableWithCells(false, false)
	output.AddItem(table, 0, 1, false)
	message := tview.NewTextView().SetText("Message")
	message.SetBorder(false)
	ui.messageContainer = message
	layout := tview.NewFlex().AddItem(optionsMenu, 0, 1, true).
		AddItem(output, 0, 4, false)
	grid := tview.NewGrid().
		SetColumns(10).
		AddItem(message, 0, 0, 1, 3, 0, 0, false).
		AddItem(layout, 1, 0, 4, 60, 0, 0, true)
	message.SetChangedFunc(func() {
		if message.GetText(true) == "Message" {
			return
		}
		<-time.After(time.Second * 2)
		message.SetText("Message").SetTextColor(tcell.ColorWhite)
		ui.app.Draw()
	})
	if err := tApp.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
