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

package tenvlib

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/pkg/loghelper"
	"github.com/tofuutils/tenv/v3/versionmanager"
	"github.com/tofuutils/tenv/v3/versionmanager/builder"
	"github.com/tofuutils/tenv/v3/versionmanager/lastuse"
)

var (
	errNoBuilder = errors.New("no builder for this tool")
	errNoConfig  = errors.New("need config or configInitFunc")
)

type tenvConfig struct {
	builders       map[string]builder.BuilderFunc
	conf           *config.Config
	displayer      loghelper.Displayer
	hclParser      *hclparse.Parser
	initConfigFunc func() (config.Config, error)
}

type TenvOption func(*tenvConfig)

// add builder or override default builder (see builder.Builders).
func AddTool(toolName string, builderFunc builder.BuilderFunc) TenvOption {
	return func(tc *tenvConfig) {
		tc.builders[toolName] = builderFunc
	}
}

func WithConfig(conf *config.Config) TenvOption {
	return func(tc *tenvConfig) {
		tc.conf = conf
	}
}

func WithDisplayer(displayer loghelper.Displayer) TenvOption {
	return func(tc *tenvConfig) {
		tc.displayer = displayer
	}
}

func WithHCLParser(hclParser *hclparse.Parser) TenvOption {
	return func(tc *tenvConfig) {
		tc.hclParser = hclParser
	}
}

func WithInertDisplayer(tc *tenvConfig) {
	tc.displayer = loghelper.InertDisplayer
}

func WithInitConfig(initConfigFunc func() (config.Config, error)) TenvOption {
	return func(tc *tenvConfig) {
		tc.initConfigFunc = initConfigFunc
	}
}

type Tenv struct {
	builders  map[string]builder.BuilderFunc
	conf      *config.Config
	hclParser *hclparse.Parser
	managers  map[string]versionmanager.VersionManager
}

func Make(options ...TenvOption) (Tenv, error) {
	builders := map[string]builder.BuilderFunc{}
	for toolName, builderFunc := range builder.Builders {
		builders[toolName] = builderFunc
	}

	tc := tenvConfig{
		builders:       builders,
		initConfigFunc: config.InitConfigFromEnv,
	}

	for _, option := range options {
		option(&tc)
	}

	if tc.conf == nil {
		if tc.initConfigFunc == nil {
			return Tenv{}, errNoConfig
		}

		conf, err := tc.initConfigFunc()
		if err != nil {
			return Tenv{}, err
		}

		tc.conf = &conf
	}

	if tc.displayer == nil {
		tc.conf.InitDisplayer(false)
	} else {
		tc.conf.Displayer = tc.displayer
	}

	if tc.hclParser == nil {
		tc.hclParser = hclparse.NewParser()
	}

	return Tenv{
		builders:  builders,
		conf:      tc.conf,
		hclParser: tc.hclParser,
		managers:  map[string]versionmanager.VersionManager{},
	}, nil
}

func (t Tenv) Command(toolName string, requestedVersion string, cmdArgs ...string) (*exec.Cmd, error) {
	if err := t.init(toolName); err != nil {
		return nil, err
	}

	installPath, err := t.managers[toolName].InstallPath()
	if err != nil {
		return nil, err
	}

	versionPath := filepath.Join(installPath, requestedVersion)
	lastuse.WriteNow(versionPath, t.conf.Displayer)

	return exec.Command(filepath.Join(versionPath, toolName), cmdArgs...), nil
}

// Evaluate version resolution strategy or version constraint (can install depending on configuration).
func (t Tenv) Evaluate(toolName string, requestedVersion string) (string, error) {
	if err := t.init(toolName); err != nil {
		return "", err
	}

	return t.managers[toolName].Evaluate(requestedVersion, false)
}

func (t Tenv) Install(toolName string, requestedVersion string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].Install(requestedVersion)
}

func (t Tenv) InstallMultiple(toolName string, versions []string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].InstallMultiple(versions)
}

func (t Tenv) ListLocal(toolName string, reverseOrder bool) ([]versionmanager.DatedVersion, error) {
	if err := t.init(toolName); err != nil {
		return nil, err
	}

	return t.managers[toolName].ListLocal(reverseOrder)
}

func (t Tenv) ListRemote(toolName string, reverseOrder bool) ([]string, error) {
	if err := t.init(toolName); err != nil {
		return nil, err
	}

	return t.managers[toolName].ListRemote(reverseOrder)
}

func (t Tenv) LocallyInstalled(toolName string) (map[string]struct{}, error) {
	if err := t.init(toolName); err != nil {
		return nil, err
	}

	return t.managers[toolName].LocalSet(), nil
}

func (t Tenv) ResolveWithVersionFiles(toolName string) (string, error) {
	if err := t.init(toolName); err != nil {
		return "", err
	}

	return t.managers[toolName].ResolveWithVersionFiles()
}

// does not handle special behavior.
func (t Tenv) Uninstall(toolName string, requestedVersion string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].UninstallMultiple([]string{requestedVersion})
}

func (t Tenv) UninstallMultiple(toolName string, versions []string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].UninstallMultiple(versions)
}

func (t Tenv) init(toolName string) error {
	if _, ok := t.managers[toolName]; ok {
		return nil
	}

	builderFunc := t.builders[toolName]
	if builderFunc == nil {
		return errNoBuilder
	}

	t.managers[toolName] = builderFunc(t.conf, t.hclParser)

	return nil
}
