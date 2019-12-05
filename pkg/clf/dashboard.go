package clf

import (
	"fmt"
	"strings"
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
	pathRequests       map[string]float64
	pathBytes          map[string]float64
	pathStatus         map[string]float64
	pathMethods        map[string]float64

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

	// used for top panels
	if metric == "requests.total" || metric == "bytes.total" {
		d.metrics[metric] = append(d.metrics[metric], point{timestamp, value})
	} else {
		// used for middle panels
		metricWords := strings.Split(metric, ".")
		if len(metricWords) == 3 {
			rootPath := metricWords[1]
			switch t := metricWords[2]; t {
			case "requests":
				d.pathRequests[rootPath] = value
			case "bytes":
				d.pathBytes[rootPath] = value
			}
		} else if len(metricWords) == 5 {
			rootPath := metricWords[1]
			subtype := metricWords[3]
			switch t := metricWords[4]; t {
			case "requests":
				d.pathStatus[strings.Join([]string{rootPath, string(subtype[0]) + "xx"}, ".")] += value
			case "bytes":
				d.pathMethods[strings.Join([]string{rootPath, subtype}, ".")] = value
			}
		}
	}
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
	rows := make([][]string, 1)
	rows[0] = []string{"Path", "Requests", "2xx", "3xx", "4xx", "5xx"}

	for k, v := range d.pathRequests {
		row := make([]string, 6)
		row[0] = k
		row[1] = fmt.Sprintf("%f", v)
		row[2] = fmt.Sprintf("%f", d.pathStatus[k+".2xx"])
		row[3] = fmt.Sprintf("%f", d.pathStatus[k+".3xx"])
		row[4] = fmt.Sprintf("%f", d.pathStatus[k+".4xx"])
		row[5] = fmt.Sprintf("%f", d.pathStatus[k+".5xx"])

		rows = append(rows, row)
	}

	// reset after rendering as we aggregate points here
	d.pathStatus = make(map[string]float64)

	d.pathRequestsPanel.Rows = rows
	d.pathRequestsPanel.TextStyle = ui.NewStyle(ui.ColorWhite)
	d.pathRequestsPanel.RowSeparator = false
	d.pathRequestsPanel.BorderStyle = ui.NewStyle(ui.ColorWhite)
	d.pathRequestsPanel.SetRect(0, d.height/3, d.width/2, d.height/3+d.height/2)
	d.pathRequestsPanel.FillRow = true
	d.pathRequestsPanel.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
}

func (d *Dashboard) updatePathBytesPanel() {
	rows := make([][]string, 1)
	rows[0] = []string{"Path", "Bytes", "GET", "POST", "PUT", "DELETE"}

	for k, v := range d.pathBytes {
		row := make([]string, 6)
		row[0] = k
		row[1] = fmt.Sprintf("%f", v)
		row[2] = fmt.Sprintf("%f", d.pathMethods[k+".GET"])
		row[3] = fmt.Sprintf("%f", d.pathMethods[k+".POST"])
		row[4] = fmt.Sprintf("%f", d.pathMethods[k+".PUT"])
		row[5] = fmt.Sprintf("%f", d.pathMethods[k+".DELETE"])

		rows = append(rows, row)
	}

	// path.bytes panel at middle right
	d.pathBytesPanel.Rows = rows
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
		pathRequests:       make(map[string]float64),
		pathBytes:          make(map[string]float64),
		pathStatus:         make(map[string]float64),
		pathMethods:        make(map[string]float64),
	}

	return &d
}
