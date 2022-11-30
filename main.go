package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

type PackageJson struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Main        string                 `json:"main"`
	Scripts     map[string]interface{} `json:"scripts"`
	Repository  struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"repository"`
	Homepage        string                 `json:"homepage"`
	Keywords        []string               `json:"keywords"`
	Author          string                 `json:"author"`
	Email           string                 `json:"email"`
	License         string                 `json:"license"`
	Files           []string               `json:"files"`
	Dependencies    map[string]interface{} `json:"dependencies"`
	DevDependencies map[string]interface{} `json:"devDependencies"`
}

func findPackageJsonPath(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if _, err := os.Stat(path + "/package.json"); errors.Is(err, os.ErrNotExist) {
		if path != "/" {
			return findPackageJsonPath(path + "/..")
		}
		return ""
	}

	return path
}

func readPackageJson(path string) (PackageJson, string) {
	path = findPackageJsonPath(path)
	if path == "" {
		return PackageJson{}, ""
	}
	jsonFile, err := os.Open(path + "/package.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var packageJson PackageJson
	json.Unmarshal(byteValue, &packageJson)
	return packageJson, path + "/package.json"
}

func processPackageJson() {
	path, _ := os.Getwd()
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		path = os.Args[1]
	}
	packageJson, path := readPackageJson(path)

	if !reflect.DeepEqual(packageJson, PackageJson{}) {
		found := false
		fmt.Println("Found package.json in", path)
		if packageJson.Name != "" {
			fmt.Println("Name:", packageJson.Name)
			found = true
		}
		if packageJson.Version != "" {
			fmt.Println("Version:", packageJson.Version)
			found = true
		}
		if packageJson.Description != "" {
			fmt.Println("Description:", packageJson.Description)
			found = true
		}
		if packageJson.Main != "" {
			fmt.Println("Main:", packageJson.Main)
			found = true
		}
		if packageJson.Homepage != "" {
			fmt.Println("Homepage:", packageJson.Homepage)
			found = true
		}
		if packageJson.Author != "" {
			fmt.Println("Author:", packageJson.Author)
			found = true
		}
		if packageJson.Email != "" {
			fmt.Println("Email:", packageJson.Email)
			found = true
		}
		if packageJson.License != "" {
			fmt.Println("License:", packageJson.License)
			found = true
		}
		if packageJson.Scripts != nil {
			fmt.Println("Scripts:")
			keys := make([]string, 0, len(packageJson.Scripts))
			for k := range packageJson.Scripts {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			length := 0
			for _, key := range keys {
				if len(key) > length {
					length = len(key)
				}
			}
			for _, key := range keys {
				fmt.Println("  ", key, strings.Repeat(" ", length-len(key)), ":", packageJson.Scripts[key])
				found = true
			}
		}
		if found == false {
			fmt.Println("No information found")
		}
	} else {
		fmt.Println("No package.json found")
	}
}

func main() {
	cwdProcess := flag.Bool("c", false, "Process current working directory or defined path")
	showHelp := flag.Bool("h", false, "Show help")
	flag.Parse()
	if *showHelp {
		fmt.Println("Usage: package-json [path]")
	} else if *cwdProcess {
		processPackageJson()
	} else {
		processPackageJson()
	}
}
