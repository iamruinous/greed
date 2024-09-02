package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
)

type model struct {
	list   list.Model
	client *feedbin.Client
	cache  *cache.Cache
}

func NewTUI(client *feedbin.Client, cache *cache.Cache) *tea.Program {
	m := model{
		list:   list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		client: client,
		cache:  cache,
	}
	return tea.NewProgram(m)
}

func (m model) Init() tea.Cmd {
	return m.fetchEntries
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 1)
	case []feedbin.Entry:
		items := make([]list.Item, len(msg))
		for i, entry := range msg {
			items[i] = entryItem(entry)
		}
		m.list.SetItems(items)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("Feedbin Latest Feeds\n\n%s\nPress q to quit", m.list.View())
}

func (m model) fetchEntries() tea.Msg {
	if cachedEntries, ok := m.cache.Get(); ok {
		return cachedEntries
	}

	entries, err := m.client.GetLatestFeeds()
	if err != nil {
		return nil
	}

	m.cache.Set(entries)
	return entries
}

type entryItem feedbin.Entry

// func (e entryItem) Title() string       { return e.Title }
func (e entryItem) Description() string { return e.URL }
func (e entryItem) FilterValue() string { return e.Title }
