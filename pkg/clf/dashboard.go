package clf

import (
	"fmt"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type point struct {
	Timestamp int64
	Value     float64
}

// Dashboard defines a dashboard to show common log format metrics
type Dashboard struct {
	sync.Mutex

	// config
	width    int
	height   int
	interval int64

	// panels
	totalRequestsPanel *widgets.Plot
	totalBytesPanel    *widgets.Plot
	pathRequestsPanel  *widgets.Table
	pathBytesPanel     *widgets.Table
	messagesPanel      *widgets.List
	metrics            map[string][]point

	//data
	messages []string
}

// Resize resizes the dashboard
func (d *Dashboard) Resize(width, height int) {
	d.width = width
	d.height = height
	d.Render()
}

// Render renders the dashboard
func (d *Dashboard) Render() {

	// refresh panels data
	d.updateTotalRequestsPanel()
	d.updateTotalBytesPanel()
	d.updatePathRequestsPanel()
	d.updatePathBytesPanel()
	d.updateMessagesPanel()

	// draw
	ui.Render(
		d.totalRequestsPanel,
		d.totalBytesPanel,
		d.pathRequestsPanel,
		d.pathBytesPanel,
		d.messagesPanel,
	)
}

// Message adds a new message to the messages panel
func (d *Dashboard) Message(t time.Time, message string) {
	d.messages = append(d.messages, fmt.Sprintf("[%v] %s", t, message))
}

// AddPoint adds a new point of the correspoinding metric
func (d *Dashboard) AddPoint(metric string, timestamp int64, value float64) {
	d.Lock()
	defer d.Unlock()

	if _, ok := d.metrics[metric]; !ok {
		d.metrics[metric] = make([]point, 0)
	}
	d.metrics[metric] = append(d.metrics[metric], point{timestamp, value})
}

func (d *Dashboard) updateTotalRequestsPanel() {
	d.Lock()
	defer d.Unlock()

	limit := time.Now().Unix() - d.interval

	points := make([][]float64, 1)
	points[0] = make([]float64, 0)

	for _, v := range d.metrics["requests.total"] {
		if v.Timestamp >= limit {
			points[0] = append(points[0], v.Value)
		}
	}

	if len(points[0]) < 2 {
		points[0] = []float64{0, 0}
	} else if len(points[0]) >= d.width/2 {
		points[0] = points[0][len(points[0])-d.width/2 : len(points[0])]
	}

	// requests.total panel at top left
	d.totalRequestsPanel.Title = "Total Requests"
	d.totalRequestsPanel.Data = points
	d.totalRequestsPanel.SetRect(0, 0, d.width/2, d.height/3)
	d.totalRequestsPanel.ShowAxes = true
	d.totalRequestsPanel.AxesColor = ui.ColorRed
	d.totalRequestsPanel.LineColors = []ui.Color{ui.ColorYellow}
}

func (d *Dashboard) updateTotalBytesPanel() {
	d.Lock()
	defer d.Unlock()

	limit := time.Now().Unix() - d.interval

	points := make([][]float64, 1)
	points[0] = make([]float64, 0)

	for _, v := range d.metrics["bytes.total"] {
		if v.Timestamp > limit {
			points[0] = append(points[0], v.Value)
		}
	}

	if len(points[0]) < 2 {
		points[0] = []float64{0, 0}
	} else if len(points[0]) >= d.width/2 {
		points[0] = points[0][len(points[0])-d.width/2 : len(points[0])]
	}

	// bytes.total panel at top right
	d.totalBytesPanel.Title = "Total Bytes"
	d.totalBytesPanel.Data = points
	d.totalBytesPanel.SetRect(d.width/2, 0, d.width, d.height/3)
	d.totalBytesPanel.AxesColor = ui.ColorWhite
	d.totalBytesPanel.LineColors = []ui.Color{ui.ColorBlue}
}

func (d *Dashboard) updatePathRequestsPanel() {
	// path.requests panel at middle left
	d.pathRequestsPanel.Rows = [][]string{
		[]string{"Path", "Requests", "2xx", "4xx", "5xx"},
		[]string{"AAA", "BBB", "CCC", "DDD", "EEE"},
		[]string{"2AA", "BBB", "CCC", "DDD", "EEE"},
		[]string{"3AA", "BBB", "CCC", "DDD", "EEE"},
		[]string{"4AA", "BBB", "CCC", "DDD", "EEE"},
		[]string{"5AA", "BBB", "CCC", "DDD", "EEE"},
		[]string{"5AA", "BBB", "CCC", "DDD", "EEE"},
		[]string{"5AA", "BBB", "CCC", "DDD", "EEE"},
		[]string{"5AA", "BBB", "CCC", "DDD", "EEE"},
	}
	d.pathRequestsPanel.TextStyle = ui.NewStyle(ui.ColorWhite)
	d.pathRequestsPanel.RowSeparator = false
	d.pathRequestsPanel.BorderStyle = ui.NewStyle(ui.ColorWhite)
	d.pathRequestsPanel.SetRect(0, d.height/3, d.width/2, d.height/3+d.height/2)
	d.pathRequestsPanel.FillRow = true
	d.pathRequestsPanel.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
}

func (d *Dashboard) updatePathBytesPanel() {
	// path.bytes panel at middle right
	d.pathBytesPanel.Rows = [][]string{
		[]string{"Path", "Bytes", "GET", "POST"},
		[]string{"AAA", "BBB", "CCC", "DDD"},
		[]string{"2AA", "BBB", "CCC", "DDD"},
		[]string{"3AA", "BBB", "CCC", "DDD"},
		[]string{"4AA", "BBB", "CCC", "DDD"},
		[]string{"5AA", "BBB", "CCC", "DDD"},
		[]string{"5AA", "BBB", "CCC", "DDD"},
		[]string{"5AA", "BBB", "CCC", "DDD"},
		[]string{"5AA", "BBB", "CCC", "DDD"},
	}
	d.pathBytesPanel.TextStyle = ui.NewStyle(ui.ColorWhite)
	d.pathBytesPanel.RowSeparator = false
	d.pathBytesPanel.BorderStyle = ui.NewStyle(ui.ColorWhite)
	d.pathBytesPanel.SetRect(d.width/2, d.height/3, d.width, d.height/3+d.height/2)
	d.pathBytesPanel.FillRow = true
	d.pathBytesPanel.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
}

func (d *Dashboard) updateMessagesPanel() {
	d.messagesPanel.Title = "Alerts"
	d.messagesPanel.SetRect(0, d.height/3+d.height/2, d.width, d.height)
	d.messagesPanel.Rows = d.messages
	d.messagesPanel.WrapText = false

}

// NewDashboard creates a new dashboard
func NewDashboard(width, height int, interval int64) *Dashboard {

	d := Dashboard{
		width:              width,
		height:             height,
		interval:           interval,
		totalRequestsPanel: widgets.NewPlot(),
		totalBytesPanel:    widgets.NewPlot(),
		pathRequestsPanel:  widgets.NewTable(),
		pathBytesPanel:     widgets.NewTable(),
		messagesPanel:      widgets.NewList(),
		messages:           make([]string, 0),
		metrics:            make(map[string][]point, 0),
	}

	return &d
}
