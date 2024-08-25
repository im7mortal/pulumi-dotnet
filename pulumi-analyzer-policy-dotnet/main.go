// Copyright 2016-2023, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	dotnetcompiler "github.com/pulumi/pulumi/sdk/go/pulumi-analyzer-policy-dotnet/compiler"
	policyAnalyzer "github.com/pulumi/pulumi/sdk/v3/go/analyzer-policy-common"
	"os"
)

func init() {
	// Specify the file path where both stdout and stderr will be redirected
	filePath := "/home/user/dev/pulumi/pulumi-policy/DOTNET.txt"

	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Redirect stdout and stderr to the file
	os.Stdout = file
	os.Stderr = file
}

// Launches the language host, which in turn fires up an RPC server implementing the LanguageRuntimeServer endpoint.
func main() {
	policyAnalyzer.Main(&policyAnalyzer.MainConfig{CompileTargetFunc: dotnetcompiler.CompileProgram})
}
