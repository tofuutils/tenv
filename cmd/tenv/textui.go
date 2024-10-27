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
	"context"
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

var tools = []list.Item{item(cmdconst.TofuName), item(cmdconst.TerraformName), item(cmdconst.TerragruntName), item(cmdconst.AtmosName)} //nolint

var (
	helpStyle         = list.DefaultStyles().HelpStyle                                    //nolint
	paginationStyle   = list.DefaultStyles().PaginationStyle                              //nolint
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(selectedColorCode)) //nolint
	titleStyle        = lipgloss.NewStyle()                                               //nolint
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
func (d itemDelegate) Render(writer io.Writer, displayList list.Model, index int, listItem list.Item) {
	version, selected := listItem.FilterValue(), " "
	if _, ok := d.choices[version]; ok {
		selected = "X"
	}
	line := loghelper.Concat("[", selected, "] ", version)

	if index == displayList.Index() {
		line = selectedItemStyle.Render(line)
	}

	fmt.Fprint(writer, line)
}

type manageItemDelegate struct {
	choices   map[string]struct{}
	installed map[string]struct{}
}

func (d manageItemDelegate) Height() int                             { return 1 }
func (d manageItemDelegate) Spacing() int                            { return 0 }
func (d manageItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d manageItemDelegate) Render(writer io.Writer, displayList list.Model, index int, listItem list.Item) {
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

	if index == displayList.Index() {
		line = selectedItemStyle.Render(line)
	}

	fmt.Fprint(writer, line)
}

type itemSelector struct {
	choices  map[string]struct{}
	list     list.Model
	quitting bool
}

func (m itemSelector) Init() tea.Cmd {
	return nil
}

func (m itemSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint
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

func (m itemSelector) View() string {
	if m.quitting {
		return ""
	}

	return "\n" + m.list.View()
}

func toolUI(ctx context.Context, conf *config.Config, hclParser *hclparse.Parser) error {
	conf.InitDisplayer(false)

	// shared object
	selection := map[string]struct{}{}

	delegate := itemDelegate{
		choices: selection,
	}

	displayList := list.New(tools, delegate, defaultWidth, listHeight)
	displayList.Title = "Which tool do you want to manage ?"
	displayList.SetShowStatusBar(false)
	displayList.SetFilteringEnabled(false)
	displayList.Styles.Title = titleStyle
	displayList.Styles.PaginationStyle = paginationStyle
	displayList.Styles.HelpStyle = helpStyle

	displayList.AdditionalFullHelpKeys = additionalFullHelpKeys
	displayList.AdditionalShortHelpKeys = additionalShortHelpKeys

	selector := itemSelector{
		choices: selection,
		list:    displayList,
	}

	_, err := tea.NewProgram(selector).Run()
	if err != nil {
		return err
	}

	if len(selector.choices) == 0 {
		loghelper.StdDisplay("No selected tool")

		return nil
	}

	for _, toolItem := range tools {
		tool := toolItem.FilterValue()
		if _, selected := selection[tool]; selected {
			if err = manageUI(ctx, builder.Builders[tool](conf, hclParser)); err != nil {
				return err
			}
		}
	}

	return nil
}

func manageUI(ctx context.Context, versionManager versionmanager.VersionManager) error {
	installed := versionManager.LocalSet()

	remoteVersions, err := versionManager.ListRemote(ctx, true)
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

	displayList := list.New(items, delegate, defaultWidth, listHeight)
	displayList.Title = loghelper.Concat("Which ", versionManager.FolderName, " version(s) do you want to install(I) or uninstall(U) ? (X mark already installed)")
	displayList.SetShowStatusBar(false)
	displayList.SetFilteringEnabled(false)
	displayList.Styles.Title = titleStyle
	displayList.Styles.PaginationStyle = paginationStyle
	displayList.Styles.HelpStyle = helpStyle

	displayList.AdditionalFullHelpKeys = additionalFullHelpKeys
	displayList.AdditionalShortHelpKeys = additionalShortHelpKeys

	selector := itemSelector{
		choices: selection,
		list:    displayList,
	}

	_, err = tea.NewProgram(selector).Run()
	if err != nil {
		return err
	}

	if len(selector.choices) == 0 {
		loghelper.StdDisplay(loghelper.Concat("No selected ", versionManager.FolderName, " versions"))

		return nil
	}

	toInstall := make([]string, 0, len(selector.choices))
	toUninstall := make([]string, 0, len(selector.choices))
	for version := range selector.choices {
		if _, installed := installed[version]; installed {
			toUninstall = append(toUninstall, version)
		} else {
			toInstall = append(toInstall, version)
		}
	}
	slices.SortFunc(toInstall, semantic.CmpVersion)
	slices.SortFunc(toUninstall, semantic.CmpVersion)

	if err = versionManager.UninstallMultiple(toUninstall); err != nil {
		return err
	}

	return versionManager.InstallMultiple(ctx, toInstall)
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

	displayList := list.New(items, delegate, defaultWidth, listHeight)
	displayList.Title = loghelper.Concat("Which ", versionManager.FolderName, " version(s) do you want to uninstall ?")
	displayList.SetShowStatusBar(false)
	displayList.SetFilteringEnabled(false)
	displayList.Styles.Title = titleStyle
	displayList.Styles.PaginationStyle = paginationStyle
	displayList.Styles.HelpStyle = helpStyle

	displayList.AdditionalFullHelpKeys = additionalFullHelpKeys
	displayList.AdditionalShortHelpKeys = additionalShortHelpKeys

	selector := itemSelector{
		choices: selection,
		list:    displayList,
	}

	_, err = tea.NewProgram(selector).Run()
	if err != nil {
		return err
	}

	if len(selector.choices) == 0 {
		loghelper.StdDisplay(loghelper.Concat("No selected ", versionManager.FolderName, " versions"))

		return nil
	}

	selected := make([]string, 0, len(selector.choices))
	for version := range selector.choices {
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
