package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ScanStatus represents the current status of a scan
// (copied from internal/models/scan.go for simplicity)
type ScanStatus string

const (
	StatusQueued    ScanStatus = "queued"
	StatusRunning   ScanStatus = "running"
	StatusCompleted ScanStatus = "completed"
	StatusFailed    ScanStatus = "failed"
)

type ScanProgress struct {
	Percentage     float64   `json:"percentage"`
	TemplatesRun   int       `json:"templates_run"`
	TotalTemplates int       `json:"total_templates"`
	VulnsFound     int       `json:"vulns_found"`
	CurrentTarget  string    `json:"current_target,omitempty"`
	LastUpdate     time.Time `json:"last_update"`
}

type ScanRequest struct {
	ID        string       `json:"id"`
	Target    string       `json:"target"`
	Status    ScanStatus   `json:"status"`
	Progress  ScanProgress `json:"progress"`
	CreatedAt time.Time    `json:"created_at"`
	StartedAt *time.Time   `json:"started_at,omitempty"`
	EndedAt   *time.Time   `json:"ended_at,omitempty"`
	Error     string       `json:"error,omitempty"`
	Results   string       `json:"results,omitempty"`
}

type model struct {
	scans   []ScanRequest
	err     error
	loading bool
}

func initialModel() model {
	return model{loading: true}
}

func (m model) Init() tea.Cmd {
	return fetchScansCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case scansMsg:
		m.scans = msg.scans
		m.err = msg.err
		m.loading = false
		return m, tea.Tick(5*time.Second, func(time.Time) tea.Msg { return fetchScans() })
	case tickMsg:
		return m, fetchScansCmd
	}
	return m, nil
}

type scansMsg struct {
	scans []ScanRequest
	err   error
}

type tickMsg struct{}

func fetchScans() tea.Msg {
	resp, err := http.Get("http://localhost:8080/scans")
	if err != nil {
		return scansMsg{nil, err}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return scansMsg{nil, err}
	}
	var scans []ScanRequest
	err = json.Unmarshal(body, &scans)
	return scansMsg{scans, err}
}

var fetchScansCmd = func() tea.Msg {
	return fetchScans()
}

func (m model) View() string {
	if m.loading {
		return "Loading scan jobs..."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}
	if len(m.scans) == 0 {
		return "No scan jobs found."
	}
	headers := lipgloss.NewStyle().Bold(true).Render("ID        Target      Status      Progress   Vulns   Last Update")
	rows := ""
	for _, scan := range m.scans {
		progress := fmt.Sprintf("%.1f%%", scan.Progress.Percentage)
		lastUpdate := scan.Progress.LastUpdate.Format("15:04:05")
		rows += fmt.Sprintf("%s  %s  %s  %s  %d  %s\n",
			scan.ID[:8], scan.Target, scan.Status, progress, scan.Progress.VulnsFound, lastUpdate)
	}
	return headers + "\n" + rows + "\nPress q to quit."
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}
