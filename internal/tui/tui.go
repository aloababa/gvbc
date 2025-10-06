package tui

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/aloababa/gvbc/internal/video"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Model struct {
	d video.Display

	screenWidth   int
	screenHeight  int
	screenRefresh bool

	flexbox         *flexbox.HorizontalFlexBox
	displayCell     *flexbox.Cell
	displayPortCell *flexbox.Cell
	hdmiCell        *flexbox.Cell

	displayPortTable *table.Table
	hdmiTable        *table.Table

	inputs []textinput.Model

	colorDepthItems    []list.Item
	colorDepthList     list.Model
	showColorDepthList bool

	timingItems    []list.Item
	timingList     list.Model
	showTimingList bool

	displayPortItems    []list.Item
	displayPortList     list.Model
	showDisplayPortList bool

	hdmiItems    []list.Item
	hdmiList     list.Model
	showHdmiList bool

	presetItems    []list.Item
	presetList     list.Model
	showPresetList bool

	focusIndex int
}

func NewModel() *Model {
	fb := flexbox.NewHorizontal(0, 0)
	displayCell := flexbox.NewCell(1, 2).SetStyle(flexCell)
	displayPortCell := flexbox.NewCell(1, 1).SetStyle(flexCell)
	hdmiCell := flexbox.NewCell(1, 1).SetStyle(flexCell)
	fb.AddColumns([]*flexbox.Column{
		fb.NewColumn().AddCells(
			displayCell,
		),
		fb.NewColumn().AddCells(
			displayPortCell,
			hdmiCell,
		),
	})

	colorDepths := video.ColorDepths()
	timings := video.Timings()
	displayPorts := video.DisplayPortVersions()
	hdmis := video.HDMIVersions()
	presets := video.Presets()

	m := &Model{
		screenRefresh:    true,
		flexbox:          fb,
		displayCell:      displayCell,
		displayPortCell:  displayPortCell,
		hdmiCell:         hdmiCell,
		inputs:           make([]textinput.Model, 3),
		colorDepthItems:  make([]list.Item, len(colorDepths)),
		timingItems:      make([]list.Item, len(timings)),
		displayPortItems: make([]list.Item, len(displayPorts)+1),
		hdmiItems:        make([]list.Item, len(hdmis)+1),
		presetItems:      make([]list.Item, len(presets)),
	}

	for i, c := range colorDepths {
		switch c {
		case video.ColorDepth8bit():
			m.colorDepthItems[i] = colorDepthListItem{
				colorDepth: c,
				title:      "8 bit",
				desc:       c.String(),
			}
		case video.ColorDepth10bit():
			m.colorDepthItems[i] = colorDepthListItem{
				colorDepth: c,
				title:      "10 bit",
				desc:       c.String(),
			}
		case video.ColorDepth12bit():
			m.colorDepthItems[i] = colorDepthListItem{
				colorDepth: c,
				title:      "12 bit",
				desc:       c.String(),
			}
		case video.ColorDepth16bit():
			m.colorDepthItems[i] = colorDepthListItem{
				colorDepth: c,
				title:      "16 bit",
				desc:       c.String(),
			}
		}
	}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = focus
	delegate.Styles.SelectedDesc = focus
	m.colorDepthList = list.New(m.colorDepthItems, delegate, 0, 0)
	m.colorDepthList.Styles.FilterCursor = focus
	m.colorDepthList.SetShowPagination(false)
	m.colorDepthList.SetShowFilter(false)
	m.colorDepthList.SetShowHelp(false)
	m.colorDepthList.SetShowStatusBar(false)
	m.colorDepthList.SetShowTitle(false)
	m.colorDepthList.Select(1)

	for i, t := range timings {
		m.timingItems[i] = timingListItem{
			timing: t,
		}
	}
	delegate = list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = focus
	delegate.Styles.SelectedDesc = focus
	m.timingList = list.New(m.timingItems, delegate, 0, 0)
	m.timingList.Styles.FilterCursor = focus
	m.timingList.SetShowPagination(false)
	m.timingList.SetShowFilter(false)
	m.timingList.SetShowHelp(false)
	m.timingList.SetShowStatusBar(false)
	m.timingList.SetShowTitle(false)
	m.timingList.Select(1)

	m.displayPortItems[0] = displayPortListItem{
		dp: video.DisplayPort{
			Version: "All",
		},
	}
	i := 1
	for _, dp := range displayPorts {
		m.displayPortItems[i] = displayPortListItem{
			dp: dp,
		}
		i++
	}
	delegate = list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = focus
	delegate.Styles.SelectedDesc = focus
	m.displayPortList = list.New(m.displayPortItems, delegate, 0, 0)
	m.displayPortList.Styles.FilterCursor = focus
	m.displayPortList.SetShowPagination(false)
	m.displayPortList.SetShowFilter(false)
	m.displayPortList.SetShowHelp(false)
	m.displayPortList.SetShowStatusBar(false)
	m.displayPortList.SetShowTitle(false)
	m.displayPortList.Select(0)

	m.hdmiItems[0] = hdmiListItem{
		hdmi: video.HDMI{
			Version: "All",
		},
	}
	i = 1
	for _, hdmi := range hdmis {
		m.hdmiItems[i] = hdmiListItem{
			hdmi: hdmi,
		}
		i++
	}
	delegate = list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = focus
	delegate.Styles.SelectedDesc = focus
	m.hdmiList = list.New(m.hdmiItems, delegate, 0, 0)
	m.hdmiList.Styles.FilterCursor = focus
	m.hdmiList.SetShowPagination(false)
	m.hdmiList.SetShowFilter(false)
	m.hdmiList.SetShowHelp(false)
	m.hdmiList.SetShowStatusBar(false)
	m.hdmiList.SetShowTitle(false)
	m.hdmiList.Select(0)

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Prompt = ""
		t.Cursor.Style = focus
		switch i {
		case 0:
			t.Focus()
			t.CharLimit = 4
			t.Width = 4
			t.PromptStyle = focus
			t.TextStyle = focus
			t.SetValue("3840")
		case 1:
			t.Blur()
			t.CharLimit = 4
			t.Width = 4
			t.SetValue("2160")
		case 2:
			t.Blur()
			t.CharLimit = 3
			t.Width = 3
			t.SetValue("144")
		}
		m.inputs[i] = t
	}

	for i, p := range presets {
		m.presetItems[i] = presetListItem{
			preset: p,
		}
	}
	delegate = list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = focus
	delegate.Styles.SelectedDesc = focus
	m.presetList = list.New(m.presetItems, delegate, 0, 0)
	m.presetList.Styles.FilterCursor = focus
	m.presetList.SetShowPagination(false)
	m.presetList.SetShowFilter(false)
	m.presetList.SetShowHelp(false)
	m.presetList.SetShowStatusBar(false)
	m.presetList.SetShowTitle(false)
	m.presetList.Select(0)

	m.updateDisplay()

	slog.Debug("updated display", slog.Any("display", m.d))

	m.displayPortTable = table.New().
		Headers([]string{"VERSION", "MODE", "MAX", "EFFECTIVE", "USAGE", "HDR", "STATUS"}...).
		Rows(m.displayPortTableData()...).
		BorderStyle(focus)

	m.hdmiTable = table.New().
		Headers([]string{"VERSION", "MODE", "MAX", "EFFECTIVE", "USAGE", "HDR", "STATUS"}...).
		Rows(m.hdmiTableData()...).
		BorderStyle(focus)

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type refreshScreen struct {
	Width  int
	Height int
}

func (m Model) refreshScreen(width, height int) func() tea.Msg {
	return func() tea.Msg {
		return refreshScreen{Width: width, Height: height}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshScreen:
		slog.Debug("received refreshScreen message")
		m.updateScreen(msg.Width, msg.Height)
		m.screenRefresh = false
		return m, nil
	case tea.WindowSizeMsg:
		slog.Debug("received tea.WindowSizeMsg message")
		m.updateScreen(msg.Width, msg.Height)
		if m.screenRefresh {
			return m, m.refreshScreen(msg.Width, msg.Height)
		}
		return m, nil
	case tea.KeyMsg:
		switch s := msg.String(); s {
		case "ctrl+c":
			return m, tea.Quit
		case "tab", "up", "down":
			if !m.showColorDepthList && !m.showTimingList &&
				!m.showDisplayPortList && !m.showHdmiList &&
				!m.showPresetList {
				if s == "up" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}
				index := len(m.inputs) + 3
				if m.focusIndex > index {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = index
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focus
						m.inputs[i].TextStyle = focus
						continue
					}
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = normal
					m.inputs[i].TextStyle = normal
				}
				return m, tea.Batch(cmds...)
			}
		case "enter":
			if m.showPresetList {
				m.applyPreset(m.presetItems[m.presetList.GlobalIndex()].(presetListItem).preset)
				m.tooglePresetList()
				m.presetList.Select(0)
			} else {
				switch m.focusIndex {
				case 3:
					m.toogleColorDepthList()
				case 4:
					m.toogleTimingList()
				case 5:
					m.toogleDisplayPortList()
				case 6:
					m.toogleHdmiList()
				}
				return m, nil
			}
		case "p":
			m.tooglePresetList()
			m.presetList.Select(0)
			return m, nil
		case "esc":
			if m.showColorDepthList {
				m.toogleColorDepthList()
			} else if m.showTimingList {
				m.toogleTimingList()
			} else if m.showDisplayPortList {
				m.toogleDisplayPortList()
			} else if m.showHdmiList {
				m.toogleHdmiList()
			} else if m.showPresetList {
				m.tooglePresetList()
				m.presetList.Select(0)
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0, len(m.inputs)+5)
	if !m.showPresetList {
		switch m.focusIndex {
		case 0, 1, 2:
			for i := range m.inputs {
				m.inputs[i], cmd = m.inputs[i].Update(msg)
				cmds = append(cmds, cmd)
			}
		case 3:
			m.colorDepthList, cmd = m.colorDepthList.Update(msg)
			cmds = append(cmds, cmd)
		case 4:
			m.timingList, cmd = m.timingList.Update(msg)
			cmds = append(cmds, cmd)
		case 5:
			m.displayPortList, cmd = m.displayPortList.Update(msg)
			cmds = append(cmds, cmd)
		case 6:
			m.hdmiList, cmd = m.hdmiList.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		m.presetList, cmd = m.presetList.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.updateDisplay()

	m.displayPortTable = m.displayPortTable.ClearRows()
	m.displayPortTable = m.displayPortTable.Rows(m.displayPortTableData()...)
	m.hdmiTable = m.hdmiTable.ClearRows()
	m.hdmiTable = m.hdmiTable.Rows(m.hdmiTableData()...)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateScreen(w, h int) {
	m.screenWidth = w
	m.screenHeight = h
	m.flexbox.SetWidth(w)
	m.flexbox.SetHeight(h)
	m.colorDepthList.SetSize(m.displayCell.GetWidth(), m.displayCell.GetHeight()-10)
	m.timingList.SetSize(m.displayCell.GetWidth(), m.displayCell.GetHeight()-10)
	m.displayPortList.SetSize(m.displayCell.GetWidth(), m.displayCell.GetHeight()-10)
	m.hdmiList.SetSize(m.displayCell.GetWidth(), m.displayCell.GetHeight()-10)
	m.presetList.SetSize(m.displayCell.GetWidth(), m.displayCell.GetHeight()-10)
	slog.Debug("screen updated", slog.Int("width", w), slog.Int("height", h))
}

func (m Model) View() string {
	m.displayCell.SetContent(renderCellContent("Display", m.displayCell, m.renderDisplayContent()))
	m.displayPortCell.SetContent(renderCellContent("Display Port", m.displayPortCell, m.displayPortTable.Render()))
	m.hdmiCell.SetContent(renderCellContent("HDMI", m.hdmiCell, m.hdmiTable.Render()))
	return m.flexbox.Render()
}

func (m Model) renderDisplayContent() string {
	var keyBinds []keyBind
	var displayContent strings.Builder
	if m.showColorDepthList {
		displayContent.WriteString(m.colorDepthList.View())
		keyBinds = listKeyBind
	} else if m.showTimingList {
		displayContent.WriteString(m.timingList.View())
		keyBinds = listKeyBind
	} else if m.showDisplayPortList {
		displayContent.WriteString(m.displayPortList.View())
		keyBinds = listKeyBind
	} else if m.showHdmiList {
		displayContent.WriteString(m.hdmiList.View())
		keyBinds = listKeyBind
	} else if m.showPresetList {
		displayContent.WriteString(m.presetList.View())
		keyBinds = presetKeyBind
	} else {
		displayContent.WriteString(line.Render("Resolution @ Refresh Rate"))
		displayContent.WriteString("\n")
		displayContent.WriteString(m.inputs[0].View())
		displayContent.WriteString("x ")
		displayContent.WriteString(m.inputs[1].View())
		displayContent.WriteString("@ ")
		displayContent.WriteString(m.inputs[2].View())
		displayContent.WriteString("Hz")
		displayContent.WriteString("\n\n\n")
		displayContent.WriteString(line.Render("Color Depth"))
		displayContent.WriteString("\n")
		if m.focusIndex == 3 {
			displayContent.WriteString(focus.Render(m.colorDepthItems[m.colorDepthList.GlobalIndex()].(colorDepthListItem).desc))
		} else {
			displayContent.WriteString(normal.Render(m.colorDepthItems[m.colorDepthList.GlobalIndex()].(colorDepthListItem).desc))
		}
		displayContent.WriteString("\n\n\n")
		displayContent.WriteString(line.Render("Timing"))
		displayContent.WriteString("\n")
		if m.focusIndex == 4 {
			displayContent.WriteString(focus.Render(m.timingItems[m.timingList.GlobalIndex()].(timingListItem).timing.String()))
		} else {
			displayContent.WriteString(normal.Render(m.timingItems[m.timingList.GlobalIndex()].(timingListItem).timing.String()))
		}
		displayContent.WriteString("\n\n\n")
		displayContent.WriteString(line.Render("DisplayPort"))
		displayContent.WriteString("\n")
		if m.focusIndex == 5 {
			displayContent.WriteString(focus.Render(m.displayPortItems[m.displayPortList.GlobalIndex()].(displayPortListItem).dp.Version))
		} else {
			displayContent.WriteString(normal.Render(m.displayPortItems[m.displayPortList.GlobalIndex()].(displayPortListItem).dp.Version))
		}
		displayContent.WriteString("\n\n\n")
		displayContent.WriteString(line.Render("HDMI"))
		displayContent.WriteString("\n")
		if m.focusIndex == 6 {
			displayContent.WriteString(focus.Render(m.hdmiItems[m.hdmiList.GlobalIndex()].(hdmiListItem).hdmi.Version))
		} else {
			displayContent.WriteString(normal.Render(m.hdmiItems[m.hdmiList.GlobalIndex()].(hdmiListItem).hdmi.Version))
		}
		displayContent.WriteString("\n")
		displayContent.WriteString(line.Render(strings.Repeat(" ", 32)))
		displayContent.WriteString("\n\n")
		displayContent.WriteString(normal.Render("Bandwidth: "))
		displayContent.WriteString(highlight.Render(m.d.Bandwidth().String()))
		displayContent.WriteString("\n\n")
		displayContent.WriteString(normal.Render("DSC: "))
		displayContent.WriteString(highlight.Render(m.d.DSC().String()))
		keyBinds = displayKeyBinds
	}
	if len(keyBinds) > 0 {
		displayContent.WriteString("\n\n")
		items := make([]string, 0, len(keyBinds)*4)
		for i, keyBind := range keyBinds {
			items = append(items, highlight.Bold(true).Render(keyBind.Key), " ", keyBind.Value)
			if i < len(keyBinds)-1 {
				items = append(items, subtle.Render(" | "))
			}
		}
		displayContent.WriteString(footer.Render(lipgloss.JoinHorizontal(lipgloss.Bottom, items...)))
	}
	return displayContent.String()
}

func renderCellContent(title string, cell *flexbox.Cell, content string) string {
	var b strings.Builder
	b.WriteString(flexCellHeader.Render(title))
	b.WriteString(flexCellFunc(cell).Render(content))
	return b.String()
}

func (m *Model) toogleColorDepthList() {
	m.showColorDepthList = !m.showColorDepthList
}

func (m *Model) toogleTimingList() {
	m.showTimingList = !m.showTimingList
}

func (m *Model) toogleDisplayPortList() {
	m.showDisplayPortList = !m.showDisplayPortList
}

func (m *Model) toogleHdmiList() {
	m.showHdmiList = !m.showHdmiList
}

func (m *Model) tooglePresetList() {
	m.showPresetList = !m.showPresetList
}

func (m Model) getLowestCompatibleMode(modes []video.TransmissionMode) video.TransmissionMode {
	if len(modes) == 0 {
		return nil
	}
	sort.Sort(byEffectiveBandwidth(modes))
	lastMode := modes[0]
	for _, mode := range modes {
		if mode.EffectiveBandwidth().Bits >= m.d.Bandwidth().Bits {
			lastMode = mode
			continue
		} else if lastMode.EffectiveBandwidth().Bits >= m.d.Bandwidth().Bits {
			return lastMode
		} else if mode.MaxCompressedBandwidth(m.d.ColorDepth).Bits >= m.d.Bandwidth().Bits {
			return mode
		}
	}
	return lastMode
}

func (m Model) displayPortTableData() [][]string {
	var rows [][]string
	if index := m.displayPortList.GlobalIndex(); index > 0 {
		dp := m.displayPortItems[index].(displayPortListItem).dp
		for _, mode := range dp.Modes {
			rows = append(rows, m.displayPortRow(dp, mode))
		}
	} else {
		for _, item := range m.displayPortItems[1:] {
			dp := item.(displayPortListItem).dp
			mode := m.getLowestCompatibleMode(dp.Modes)
			rows = append(rows, m.displayPortRow(dp, mode))
		}
	}
	return rows
}

func (m Model) displayPortRow(dp video.DisplayPort, mode video.TransmissionMode) []string {
	var hdr string
	if dp.CanHDR(m.d.ColorDepth) {
		hdr = "Yes"
	} else {
		hdr = "No"
	}
	var status string
	if mode.EffectiveBandwidth().Bits >= m.d.Bandwidth().Bits {
		status = "✅"
	} else {
		if dp.DSC {
			if mode.MaxCompressedBandwidth(m.d.ColorDepth).Bits >= m.d.Bandwidth().Bits {
				status = "❗ (DSC)"
			} else {
				status = "❌ (Bandwidth)"
				hdr = "No"
			}
		} else {
			status = "❌ (No DSC)"
			hdr = "No"
		}
	}
	return []string{dp.Version, mode.GetName(),
		mode.GetBandwidth().String(), mode.EffectiveBandwidth().String(),
		fmt.Sprintf("%.1f%%", mode.Usage(m.d.Bandwidth())), hdr, status}
}

func (m Model) hdmiTableData() [][]string {
	var rows [][]string
	if index := m.hdmiList.GlobalIndex(); index > 0 {
		hdmi := m.hdmiItems[index].(hdmiListItem).hdmi
		for _, mode := range hdmi.Modes {
			rows = append(rows, m.hdmiRow(hdmi, mode))
		}
	} else {
		for _, item := range m.hdmiItems[1:] {
			hdmi := item.(hdmiListItem).hdmi
			mode := m.getLowestCompatibleMode(hdmi.Modes)
			rows = append(rows, m.hdmiRow(hdmi, mode))
		}
	}
	return rows
}

func (m Model) hdmiRow(hdmi video.HDMI, mode video.TransmissionMode) []string {
	var hdr string
	if hdmi.CanHDR(m.d.ColorDepth) {
		hdr = "Yes"
	} else {
		hdr = "No"
	}
	var status string
	if mode.EffectiveBandwidth().Bits >= m.d.Bandwidth().Bits {
		status = "✅"
	} else {
		if hdmi.DSC {
			if mode.MaxCompressedBandwidth(m.d.ColorDepth).Bits >= m.d.Bandwidth().Bits {
				status = "❗ (DSC)"
			} else {
				status = "❌ (Bandwidth)"
				hdr = "No"
			}
		} else {
			status = "❌ (No DSC)"
			hdr = "No"
		}
	}
	return []string{hdmi.Version, mode.GetName(),
		mode.GetBandwidth().String(), mode.EffectiveBandwidth().String(),
		fmt.Sprintf("%.1f%%", mode.Usage(m.d.Bandwidth())), hdr, status}
}

func (m *Model) updateDisplay() {
	d, err := m.getDisplay()
	if err == nil {
		m.d = d
	}
}

func (m Model) getDisplay() (video.Display, error) {
	width, err := strconv.Atoi(m.inputs[0].Value())
	if err != nil {
		return video.Display{}, err
	}
	height, err := strconv.Atoi(m.inputs[1].Value())
	if err != nil {
		return video.Display{}, err
	}
	refreshRate, err := strconv.Atoi(m.inputs[2].Value())
	if err != nil {
		return video.Display{}, err
	}
	return video.Display{
		Width:       width,
		Height:      height,
		RefreshRate: refreshRate,
		ColorDepth:  m.colorDepthItems[m.colorDepthList.GlobalIndex()].(colorDepthListItem).colorDepth,
		Timing:      m.timingItems[m.timingList.GlobalIndex()].(timingListItem).timing,
	}, nil
}

func (m *Model) applyPreset(p video.Preset) {
	m.d = p.Display
	m.inputs[0].SetValue(strconv.Itoa(m.d.Width))
	m.inputs[1].SetValue(strconv.Itoa(m.d.Height))
	m.inputs[2].SetValue(strconv.Itoa(m.d.RefreshRate))
	m.colorDepthList.Select(m.getColorDepthIndex(m.d.ColorDepth))
	m.timingList.Select(m.getTimingIndex(m.d.Timing))
}

func (m Model) getColorDepthIndex(colorDepth video.ColorDepth) int {
	for i, item := range m.colorDepthItems {
		if item.(colorDepthListItem).colorDepth == colorDepth {
			return i
		}
	}
	return 0
}

func (m Model) getTimingIndex(timing video.Timing) int {
	for i, item := range m.timingItems {
		if item.(timingListItem).timing.String() == timing.String() {
			return i
		}
	}
	return 0
}

type keyBind struct {
	Key   string
	Value string
}

var (
	displayKeyBinds = []keyBind{
		{
			Key:   "↑ / ↓",
			Value: "navigate",
		},
		{
			Key:   "p",
			Value: "presets",
		},
		{
			Key:   "ctrl+c",
			Value: "exit",
		},
	}
	listKeyBind = []keyBind{
		{
			Key:   "↑ / ↓",
			Value: "navigate",
		},
		{
			Key:   "esc / enter",
			Value: "close",
		},
		{
			Key:   "ctrl+c",
			Value: "exit",
		},
	}
	presetKeyBind = []keyBind{
		{
			Key:   "↑ / ↓",
			Value: "navigate",
		},
		{
			Key:   "enter",
			Value: "apply",
		},
		{
			Key:   "esc / p",
			Value: "close",
		},
		{
			Key:   "ctrl+c",
			Value: "exit",
		},
	}
)

type colorDepthListItem struct {
	colorDepth video.ColorDepth
	title      string
	desc       string
}

func (i colorDepthListItem) Title() string       { return i.title }
func (i colorDepthListItem) Description() string { return i.desc }
func (i colorDepthListItem) FilterValue() string { return i.title }

type timingListItem struct {
	timing video.Timing
}

func (i timingListItem) Title() string       { return i.timing.String() }
func (i timingListItem) Description() string { return i.timing.String() }
func (i timingListItem) FilterValue() string { return i.timing.String() }

type displayPortListItem struct {
	dp video.DisplayPort
}

func (i displayPortListItem) Title() string       { return "DisplayPort" }
func (i displayPortListItem) Description() string { return i.dp.Version }
func (i displayPortListItem) FilterValue() string { return i.dp.Version }

type hdmiListItem struct {
	hdmi video.HDMI
}

func (i hdmiListItem) Title() string       { return "HDMI" }
func (i hdmiListItem) Description() string { return i.hdmi.Version }
func (i hdmiListItem) FilterValue() string { return i.hdmi.Version }

type presetListItem struct {
	preset video.Preset
}

func (i presetListItem) Title() string { return i.preset.Name }
func (i presetListItem) Description() string {
	return fmt.Sprintf("%dHz %s", i.preset.Display.RefreshRate, i.preset.Display.ColorDepth.String())
}
func (i presetListItem) FilterValue() string { return i.preset.Name }

type byEffectiveBandwidth []video.TransmissionMode

func (b byEffectiveBandwidth) Len() int {
	return len(b)
}

func (b byEffectiveBandwidth) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byEffectiveBandwidth) Less(i, j int) bool {
	return b[i].EffectiveBandwidth().Bits > b[j].EffectiveBandwidth().Bits
}
