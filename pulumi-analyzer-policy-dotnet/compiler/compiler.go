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

package compiler

import (
	"fmt"
	policyAnalyzer "github.com/pulumi/pulumi/sdk/v3/go/analyzer-policy-common"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/executable"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// This function takes a file target to specify where to compile to.
// If `outfile` is "", the binary is compiled to a new temporary file.
// This function returns the path of the file that was produced.
func CompileProgram(config *policyAnalyzer.CompileConfig) (*policyAnalyzer.CompileResult, error) {
	csprojFileSearchPattern := filepath.Join(config.ProgramDirectory, "*.csproj")
	matches, err := filepath.Glob(csprojFileSearchPattern)
	if err != nil || len(matches) == 0 {
		return nil, fmt.Errorf("Failed to find csproj files for 'dotnet build' matching %s", csprojFileSearchPattern)
	}

	// Extract the project name from the first matched .csproj file
	projectName := strings.TrimSuffix(filepath.Base(matches[0]), filepath.Ext(matches[0]))

	if config.OutFile == "" {
		// If no outfile is supplied, write the .NET assembly to a temporary directory.
		config.OutFile, err = os.MkdirTemp("", "pulumi-dotnet.*")
		if err != nil {
			return nil, fmt.Errorf("unable to create dotnet program temp directory: %w", err)
		}
	}

	dotnetbin, err := executable.FindExecutable("dotnet")
	if err != nil {
		return nil, fmt.Errorf("unable to find 'dotnet' executable: %w", err)
	}

	logging.V(5).Infof("Attempting to build .NET program in %s with: %s build -o %s", config.ProgramDirectory, dotnetbin, config.OutFile)
	buildCmd := exec.Command(dotnetbin, "build", matches[0], "-o", config.OutFile, "--configuration", "Release")
	buildCmd.Dir = config.ProgramDirectory
	buildCmd.Stdout, buildCmd.Stderr = os.Stdout, os.Stderr

	if err := buildCmd.Run(); err != nil {
		return nil, fmt.Errorf("unable to run `dotnet build`: %w", err)
	}

	// Ensure outfile points to the actual executable
	// On Linux/Mac, the executable doesn't have an extension; on Windows, it has a .exe extension
	var exePath string
	if os.PathSeparator == '/' {
		exePath = filepath.Join(config.OutFile, projectName) // Unix-based systems
	} else {
		exePath = filepath.Join(config.OutFile, projectName+".exe") // Windows
	}

	return &policyAnalyzer.CompileResult{Program: exePath}, nil
}
