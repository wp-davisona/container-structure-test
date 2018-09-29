// Copyright 2017 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"fmt"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib"

	"github.com/GoogleContainerTools/container-structure-test/pkg/drivers"
	types "github.com/GoogleContainerTools/container-structure-test/pkg/types/unversioned"
	"github.com/GoogleContainerTools/container-structure-test/pkg/utils"
)

type CommandTest struct {
	BaseTest       `yaml:",inline"`
	ExitCode       int      `yaml:"exitCode"`
	Command        string   `yaml:"command"`
	Args           []string `yaml:"args"`
	ExpectedOutput []string `yaml:"expectedOutput"`
	ExcludedOutput []string `yaml:"excludedOutput"`
	ExpectedError  []string `yaml:"expectedError"`
	ExcludedError  []string `yaml:"excludedError"` // excluded error from running command
}

func (ct *CommandTest) Validate(channel chan interface{}) bool {
	res := &types.TestResult{}

	if ct.Command == "" {
		res.Errorf("Please provide a valid command to run for test %s", ct.Name)
	}

	return ct.BaseTest.Validate(channel, res)
}

func (ct *CommandTest) LogName() string {
	return fmt.Sprintf("Command Test: %s", ct.Name)
}

func (ct *CommandTest) Run(driver drivers.Driver) *types.TestResult {
	ctc_lib.Log.Debug(ct.LogName())
	config, err := driver.GetConfig()
	if err != nil {
		ctc_lib.Log.Errorf("error retrieving image config: %s", err.Error())
	}
	fullCommand := utils.SubstituteEnvVars(append([]string{ct.Command}, ct.Args...), config.Env)
	stdout, stderr, exitcode, err := driver.ProcessCommand(ct.EnvVars, fullCommand)
	result := &types.TestResult{
		Name:   ct.LogName(),
		Pass:   true,
		Errors: make([]string, 0),
		Stderr: stderr,
		Stdout: stdout,
	}
	if err != nil {
		result.Fail()
		result.Error(err.Error())
		return result
	}

	ct.CheckOutput(result, stdout, stderr, exitcode)
	return result
}

func (ct *CommandTest) CheckOutput(result *types.TestResult, stdout string, stderr string, exitCode int) {
	for _, errStr := range ct.ExpectedError {
		if !utils.CompileAndRunRegex(errStr, stderr, true) {
			result.Errorf("Expected string '%s' not found in error '%s'", errStr, stderr)
			result.Fail()
		}
	}
	for _, errStr := range ct.ExcludedError {
		if !utils.CompileAndRunRegex(errStr, stderr, false) {
			result.Errorf("Excluded string '%s' found in error '%s'", errStr, stderr)
			result.Fail()
		}
	}
	for _, outStr := range ct.ExpectedOutput {
		if !utils.CompileAndRunRegex(outStr, stdout, true) {
			result.Errorf("Expected string '%s' not found in output '%s'", outStr, stdout)
			result.Fail()
		}
	}
	for _, outStr := range ct.ExcludedOutput {
		if !utils.CompileAndRunRegex(outStr, stdout, false) {
			result.Errorf("Excluded string '%s' found in output '%s'", outStr, stdout)
			result.Fail()
		}
	}
	if ct.ExitCode != exitCode {
		result.Errorf("Test '%s' exited with incorrect error code. Expected: %d, Actual: %d", ct.Name, ct.ExitCode, exitCode)
		result.Fail()
	}
}
