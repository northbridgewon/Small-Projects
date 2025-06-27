package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v3/process"
)

type processInfo struct {
	pid    int32
	name   string
	cpu    float64
	memory float32
}

type model struct {
	table table.Model
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		procs, err := process.Processes()
		if err != nil {
			return m, nil
		}

		infos := make([]processInfo, len(procs))
		for i, p := range procs {
			name, _ := p.Name()
			cpu, _ := p.CPUPercent()
			mem, _ := p.MemoryPercent()
			infos[i] = processInfo{pid: p.Pid, name: name, cpu: cpu, memory: mem}
		}

		sort.Slice(infos, func(i, j int) bool {
			return infos[i].cpu > infos[j].cpu
		})

		rows := make([]table.Row, len(infos))
		for i, p := range infos {
			rows[i] = table.Row{fmt.Sprintf("%d", p.pid), p.name, fmt.Sprintf("%.2f%%", p.cpu), fmt.Sprintf("%.2f%%", p.memory)}
		}
		m.table.SetRows(rows)
		return m, tick()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.table.View() + "\n"
}

func main() {
	columns := []table.Column{
		{Title: "PID", Width: 10},
		{Title: "Name", Width: 40},
		{Title: "CPU %", Width: 10},
		{Title: "Mem %", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	initialModel := model{table: t}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
