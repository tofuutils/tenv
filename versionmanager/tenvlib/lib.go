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
	"context"
	"errors"
	"os/exec"

	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/pkg/cmdproxy"
	"github.com/tofuutils/tenv/v3/pkg/loghelper"
	"github.com/tofuutils/tenv/v3/versionmanager"
	"github.com/tofuutils/tenv/v3/versionmanager/builder"
	"github.com/tofuutils/tenv/v3/versionmanager/proxy"
)

var errNoBuilder = errors.New("no builder for this tool")

type tenvConfig struct {
	autoInstall    bool
	builders       map[string]builder.Func
	conf           *config.Config
	displayer      loghelper.Displayer
	hclParser      *hclparse.Parser
	ignoreEnv      bool
	initConfigFunc func() (config.Config, error)
}

type TenvOption func(*tenvConfig)

// add builder or override default builder (see builder.Builders).
func AddTool(toolName string, builderFunc builder.Func) TenvOption {
	return func(tc *tenvConfig) {
		tc.builders[toolName] = builderFunc
	}
}

func AutoInstall(tc *tenvConfig) {
	tc.autoInstall = true
}

func DisableDisplay(tc *tenvConfig) {
	tc.displayer = loghelper.InertDisplayer
}

func IgnoreEnv(tc *tenvConfig) {
	tc.ignoreEnv = true
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

// Not concurrent safe.
type Tenv struct {
	builders  map[string]builder.Func
	conf      *config.Config
	hclParser *hclparse.Parser
	managers  map[string]versionmanager.VersionManager
}

// The returned wrapper is not concurrent safe.
func Make(options ...TenvOption) (Tenv, error) {
	builders := map[string]builder.Func{}
	for toolName, builderFunc := range builder.Builders {
		builders[toolName] = builderFunc
	}

	wrapperConf := tenvConfig{
		builders:       builders,
		initConfigFunc: config.InitConfigFromEnv,
	}

	for _, option := range options {
		option(&wrapperConf)
	}

	if wrapperConf.ignoreEnv {
		wrapperConf.initConfigFunc = config.DefaultConfig
	}

	if wrapperConf.conf == nil {
		innerConf, err := wrapperConf.initConfigFunc()
		if err != nil {
			return Tenv{}, err
		}

		wrapperConf.conf = &innerConf
	}

	if wrapperConf.autoInstall {
		wrapperConf.conf.SkipInstall = false
	}

	if wrapperConf.displayer == nil {
		wrapperConf.conf.InitDisplayer(false)
	} else {
		wrapperConf.conf.Displayer = wrapperConf.displayer
	}

	if wrapperConf.hclParser == nil {
		wrapperConf.hclParser = hclparse.NewParser()
	}

	return Tenv{
		builders:  builders,
		conf:      wrapperConf.conf,
		hclParser: wrapperConf.hclParser,
		managers:  map[string]versionmanager.VersionManager{},
	}, nil
}

// return an exec.Cmd in order to call the specified tool version (need to have it installed for the Cmd call to work).
func (t Tenv) Command(ctx context.Context, toolName string, requestedVersion string, cmdArgs ...string) (*exec.Cmd, error) {
	err := t.init(toolName)
	if err != nil {
		return nil, err
	}

	installPath, err := t.managers[toolName].InstallPath()
	if err != nil {
		return nil, err
	}

	execPath := proxy.ExecPath(installPath, requestedVersion, toolName, t.conf.Displayer)

	return exec.CommandContext(ctx, execPath, cmdArgs...), nil
}

// Use the result of Tenv.Command to call cmdproxy.Run (Always call os.Exit).
func (t Tenv) CommandProxy(ctx context.Context, toolName string, requestedVersion string, cmdArgs ...string) error {
	cmd, err := t.Command(ctx, toolName, requestedVersion, cmdArgs...)
	if err != nil {
		return err
	}

	cmdproxy.Run(cmd, t.conf.GithubActions, t.conf.Getenv)

	return nil
}

// Detect version (resolve and evaluate, can install depending on configuration).
func (t Tenv) Detect(ctx context.Context, toolName string) (string, error) {
	err := t.init(toolName)
	if err != nil {
		return "", err
	}

	return t.managers[toolName].Detect(ctx, false)
}

// Use the result of Tenv.Detect to call Tenv.Command.
func (t Tenv) DetectedCommand(ctx context.Context, toolName string, cmdArgs ...string) (*exec.Cmd, error) {
	detectedVersion, err := t.Detect(ctx, toolName)
	if err != nil {
		return nil, err
	}

	// t.managers[toolName] is initialized by Tenv.Detect
	installPath, err := t.managers[toolName].InstallPath()
	if err != nil {
		return nil, err
	}

	execPath := proxy.ExecPath(installPath, detectedVersion, toolName, t.conf.Displayer)

	return exec.CommandContext(ctx, execPath, cmdArgs...), nil
}

// Use the result of Tenv.DetectedCommand to call cmdproxy.Run (Always call os.Exit).
func (t Tenv) DetectedCommandProxy(ctx context.Context, toolName string, cmdArgs ...string) error {
	cmd, err := t.DetectedCommand(ctx, toolName, cmdArgs...)
	if err != nil {
		return err
	}

	cmdproxy.Run(cmd, t.conf.GithubActions, t.conf.Getenv)

	return nil
}

// Evaluate version resolution strategy or version constraint (can install depending on configuration).
func (t Tenv) Evaluate(ctx context.Context, toolName string, requestedVersion string) (string, error) {
	if err := t.init(toolName); err != nil {
		return "", err
	}

	return t.managers[toolName].Evaluate(ctx, requestedVersion, false)
}

func (t Tenv) Install(ctx context.Context, toolName string, requestedVersion string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].Install(ctx, requestedVersion)
}

func (t Tenv) InstallMultiple(ctx context.Context, toolName string, versions []string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].InstallMultiple(ctx, versions)
}

func (t Tenv) ListLocal(_ context.Context, toolName string, reverseOrder bool) ([]versionmanager.DatedVersion, error) {
	if err := t.init(toolName); err != nil {
		return nil, err
	}

	return t.managers[toolName].ListLocal(reverseOrder)
}

func (t Tenv) ListRemote(ctx context.Context, toolName string, reverseOrder bool) ([]string, error) {
	if err := t.init(toolName); err != nil {
		return nil, err
	}

	return t.managers[toolName].ListRemote(ctx, reverseOrder)
}

func (t Tenv) LocallyInstalled(_ context.Context, toolName string) (map[string]struct{}, error) {
	if err := t.init(toolName); err != nil {
		return nil, err
	}

	return t.managers[toolName].LocalSet(), nil
}

func (t Tenv) ResetDefaultConstraint(_ context.Context, toolName string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].ResetConstraint()
}

func (t Tenv) ResetDefaultVersion(_ context.Context, toolName string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].ResetVersion()
}

func (t Tenv) SetDefaultConstraint(_ context.Context, toolName string, constraint string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].SetConstraint(constraint)
}

func (t Tenv) SetDefaultVersion(ctx context.Context, toolName string, requestedVersion string, workingDir bool) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].Use(ctx, requestedVersion, workingDir)
}

// Does not handle special behavior.
func (t Tenv) Uninstall(_ context.Context, toolName string, requestedVersion string) error {
	if err := t.init(toolName); err != nil {
		return err
	}

	return t.managers[toolName].UninstallMultiple([]string{requestedVersion})
}

func (t Tenv) UninstallMultiple(_ context.Context, toolName string, versions []string) error {
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
