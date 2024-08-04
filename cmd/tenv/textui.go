/*
 *
 * Copyright 2024 tofuutils authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"fmt"
	"io"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/config/cmdconst"
	"github.com/tofuutils/tenv/v3/pkg/loghelper"
	"github.com/tofuutils/tenv/v3/versionmanager"
	"github.com/tofuutils/tenv/v3/versionmanager/builder"
	"github.com/tofuutils/tenv/v3/versionmanager/semantic"
)

const (
	defaultWidth      = 20
	listHeight        = 14
	selectedColorCode = "170"
)

var tools = []list.Item{item(cmdconst.TofuName), item(cmdconst.TerraformName), item(cmdconst.TerragruntName), item(cmdconst.AtmosName)}

var (
	helpStyle         = list.DefaultStyles().HelpStyle
	paginationStyle   = list.DefaultStyles().PaginationStyle
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(selectedColorCode))
	titleStyle        = lipgloss.NewStyle()
)

type item string

func (i item) FilterValue() string {
	return string(i)
}

type itemDelegate struct {
	choices map[string]struct{}
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	version, selected := listItem.FilterValue(), " "
	if _, ok := d.choices[version]; ok {
		selected = "X"
	}
	line := loghelper.Concat("[", selected, "] ", version)

	if index == m.Index() {
		line = selectedItemStyle.Render(line)
	}

	fmt.Fprint(w, line)
}

type manageItemDelegate struct {
	choices   map[string]struct{}
	installed map[string]struct{}
}

func (d manageItemDelegate) Height() int                             { return 1 }
func (d manageItemDelegate) Spacing() int                            { return 0 }
func (d manageItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d manageItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	version, selectedStr := listItem.FilterValue(), " "
	_, selected := d.choices[version]
	_, installed := d.installed[version]
	if selected {
		// display what will be done
		if installed {
			selectedStr = "U"
		} else {
			selectedStr = "I"
		}
	} else {
		if installed {
			selectedStr = "X"
		}
	}

	line := loghelper.Concat("[", selectedStr, "] ", version)

	if index == m.Index() {
		line = selectedItemStyle.Render(line)
	}

	fmt.Fprint(w, line)
}

type itemModel struct {
	choices  map[string]struct{}
	list     list.Model
	quitting bool
}

func (m itemModel) Init() tea.Cmd {
	return nil
}

func (m itemModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)

		return m, nil
	case tea.KeyMsg:
		version := m.list.SelectedItem().FilterValue()

		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc", "q":
			m.quitting = true

			clear(m.choices)

			return m, tea.Quit
		case "enter":
			m.quitting = true

			if len(m.choices) == 0 {
				m.choices[version] = struct{}{}
			}

			return m, tea.Quit
		case " ":
			if _, ok := m.choices[version]; ok {
				delete(m.choices, version)
			} else {
				m.choices[version] = struct{}{}
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m itemModel) View() string {
	if m.quitting {
		return ""
	}

	return "\n" + m.list.View()
}

func toolUI(conf *config.Config, hclParser *hclparse.Parser) error {
	conf.InitDisplayer(false)

	// shared object
	selection := map[string]struct{}{}

	delegate := itemDelegate{
		choices: selection,
	}

	l := list.New(tools, delegate, defaultWidth, listHeight)
	l.Title = "Which tool do you want to manage ?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	l.AdditionalFullHelpKeys = additionalFullHelpKeys
	l.AdditionalShortHelpKeys = additionalShortHelpKeys

	m := itemModel{
		choices: selection,
		list:    l,
	}

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		return err
	}

	if len(m.choices) == 0 {
		loghelper.StdDisplay("No selected tool")

		return nil
	}

	for _, toolItem := range tools {
		tool := toolItem.FilterValue()
		if _, selected := selection[tool]; selected {
			if err = manageUI(builder.Builders[tool](conf, hclParser)); err != nil {
				return err
			}
		}
	}

	return nil
}

func manageUI(versionManager versionmanager.VersionManager) error {
	installed := versionManager.LocalSet()

	remoteVersions, err := versionManager.ListRemote(true)
	if err != nil {
		return err
	}

	items := make([]list.Item, 0, len(remoteVersions))
	for _, remoteVersion := range remoteVersions {
		items = append(items, item(remoteVersion))
	}

	// shared object
	selection := map[string]struct{}{}

	delegate := manageItemDelegate{
		choices:   selection,
		installed: installed,
	}

	l := list.New(items, delegate, defaultWidth, listHeight)
	l.Title = loghelper.Concat("Which ", versionManager.FolderName, " version(s) do you want to install(I) or uninstall(U) ? (X mark already installed)")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	l.AdditionalFullHelpKeys = additionalFullHelpKeys
	l.AdditionalShortHelpKeys = additionalShortHelpKeys

	m := itemModel{
		choices: selection,
		list:    l,
	}

	_, err = tea.NewProgram(m).Run()
	if err != nil {
		return err
	}

	if len(m.choices) == 0 {
		loghelper.StdDisplay(loghelper.Concat("No selected ", versionManager.FolderName, " versions"))

		return nil
	}

	toInstall := make([]string, 0, len(m.choices))
	toUninstall := make([]string, 0, len(m.choices))
	for version := range m.choices {
		if _, installed := installed[version]; installed {
			toUninstall = append(toUninstall, version)
		} else {
			toInstall = append(toInstall, version)
		}
	}
	slices.SortFunc(toInstall, semantic.CmpVersion)
	slices.SortFunc(toUninstall, semantic.CmpVersion)

	err = versionManager.UninstallMultiple(toUninstall)
	if err != nil {
		return nil
	}

	return versionManager.InstallMultiple(toInstall)
}

func uninstallUI(versionManager versionmanager.VersionManager) error {
	datedVersions, err := versionManager.ListLocal(false)
	if err != nil {
		return err
	}

	items := make([]list.Item, 0, len(datedVersions))
	for _, datedVersion := range datedVersions {
		items = append(items, item(datedVersion.Version))
	}

	// shared object
	selection := map[string]struct{}{}

	delegate := itemDelegate{
		choices: selection,
	}

	l := list.New(items, delegate, defaultWidth, listHeight)
	l.Title = loghelper.Concat("Which ", versionManager.FolderName, " version(s) do you want to uninstall ?")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	l.AdditionalFullHelpKeys = additionalFullHelpKeys
	l.AdditionalShortHelpKeys = additionalShortHelpKeys

	m := itemModel{
		choices: selection,
		list:    l,
	}

	_, err = tea.NewProgram(m).Run()
	if err != nil {
		return err
	}

	if len(m.choices) == 0 {
		loghelper.StdDisplay(loghelper.Concat("No selected ", versionManager.FolderName, " versions"))

		return nil
	}

	selected := make([]string, 0, len(m.choices))
	for version := range m.choices {
		selected = append(selected, version)
	}
	slices.SortFunc(selected, semantic.CmpVersion)

	return versionManager.UninstallMultiple(selected)
}

func additionalFullHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "select item"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "validate uninstallation"),
		),
	}
}

func additionalShortHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "select"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "validate"),
		),
	}
}
