/*
Copyright © 2021 Benjamín García Roqués <benjamingarciaroques@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lazy",
	Short: "Lazy is an application that lets you manipulate programming files easily.",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lazy.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".lazy" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".lazy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// getDir gets the directory name of the given name and returns it.
func getDirPath(file string) string {
	dot := getExtensionIndex(file)

	// get the home dir ("~/") to append to it
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(dir, file[dot+1:]+"_projects") + "\\"
	} else {
		// get the documents path and append the new name to it
		path = filepath.Join(dir, "Documents", file[dot+1:]+"_projects") + "/"
	}

	return path
}

// getExtensionIndex gets position of the dot in the given filepath's extension.
func getExtensionIndex(filepath string) int {
	dot := strings.Index(filepath, ".")
	index := dot

	hasMoreThanOneDot := strings.Count(filepath, string(filepath[dot])) > 1
	if hasMoreThanOneDot {
		// get the index of the other dot and add to the total index count
		dot = getExtensionIndex(filepath[dot+1:])
		index += dot
	}

	return index
}

// getBasePath gets root directory of the file where it is called and returns the path to it.
func getBasePath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "..")

	return basepath
}

// hasExtension returns if the given file has an extension or not.
func hasExtension(f string) bool {
	if err := getExtensionIndex(f); err == -1 {
		return false
	}

	return true
}
