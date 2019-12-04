package clf

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Dashboard defines a dashboard to show common log format metrics
type Dashboard struct {
	// config
	width    int
	height   int
	interval int64

	// panels
	totalRequestsPanel *widgets.Plot
	totalBytesPanel    *widgets.Plot
	pathRequestsPanel  *widgets.Table
	pathBytesPanel     *widgets.Table
	messagesPanel      *widgets.Paragraph

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
	d.updateMessagesPanel()
}

func (d *Dashboard) updateMessagesPanel() {
	d.messagesPanel.Text = strings.Join(d.messages, "\n")
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
		messagesPanel:      widgets.NewParagraph(),
		messages:           make([]string, 0),
	}

	return &d
}
