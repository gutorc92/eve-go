package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gutorc92/api-farm/metrics"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	uri      = "uri"
	database = "database"
	logLevel = "log-level"
	files    = "files"
)

// WebConfig defines the parametric information of a data-controller server instance
type WebConfig struct {
	Metrics   *metrics.Metrics
	Database  string
	Uri       string
	JsonFiles []string
}

func AddFlags(flags *pflag.FlagSet) {
	flags.StringP(uri, "u", "", "Mongo db uri")
	flags.StringP(database, "d", "", "Mongo database")
	flags.StringP(logLevel, "l", "info", "[optional] The loggin level for this service")
	flags.StringP(files, "f", "", "The file or directory to generate api service")
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func IsFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return false, nil
	}
	ext := filepath.Ext(path)
	fmt.Println("ext", ext)
	if ext == ".json" {
		return true, nil
	}
	return false, nil
}

func ListFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if isFile, err := IsFile(path); isFile && err == nil {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// Init initializes the web config with properties retrieved from Viper.
func (c *WebConfig) Init(v *viper.Viper) *WebConfig {
	c.Metrics = metrics.New()
	if v.GetString(database) != "" {
		c.Database = v.GetString(database)
	}
	if v.GetString(uri) != "" {
		c.Uri = v.GetString(uri)
	}
	if v.GetString(files) != "" {
		argumentFiles := v.GetString(files)
		isDirec, err := IsDirectory(argumentFiles)
		if err != nil {
			panic(err)
		}
		if isDirec {
			filesListed, err := ListFiles(argumentFiles)
			if err != nil {
				panic(filesListed)
			}
			c.JsonFiles = filesListed
		} else {
			isFileArgument, err := IsFile(argumentFiles)
			if isFileArgument && err == nil {
				c.JsonFiles = append(c.JsonFiles, argumentFiles)
			}
		}
	}
	return c
}
