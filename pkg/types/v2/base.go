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
	types "github.com/GoogleContainerTools/container-structure-test/pkg/types/unversioned"
)

type BaseTest struct {
	Name     string         `yaml:"name"`
	Setup    [][]string     `yaml:"setup"`
	Teardown [][]string     `yaml:"teardown"`
	EnvVars  []types.EnvVar `yaml:"envVars"`
}

func (bt *BaseTest) Validate(channel chan interface{}, res *types.TestResult) bool {

	if bt.Name == "" {
		res.Error("Please provide a valid name for every test")
	}
	res.Name = bt.Name
	if bt.Setup != nil {
		for _, c := range bt.Setup {
			if len(c) == 0 {
				res.Error("Error in setup command configuration encountered; please check formatting and remove all empty setup commands")
			}
		}
	}
	if bt.Teardown != nil {
		for _, c := range bt.Teardown {
			if len(c) == 0 {
				res.Error("Error in teardown command configuration encountered; please check formatting and remove all empty teardown commands")
			}
		}
	}
	if bt.EnvVars != nil {
		for _, envVar := range bt.EnvVars {
			if envVar.Key == "" || envVar.Value == "" {
				res.Error("Please provide non-empty keys and values for all specified env vars")
			}
		}
	}

	if len(res.Errors) > 0 {
		channel <- res
		return false
	}

	return true
}
