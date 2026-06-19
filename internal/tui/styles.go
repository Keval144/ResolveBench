package tui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	SubtitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	SuccessStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	DimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	TableHeader    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	PrimaryStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	SecondaryStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	BorderStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).MarginTop(1)
)
