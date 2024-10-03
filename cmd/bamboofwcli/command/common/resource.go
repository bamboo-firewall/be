package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/resouremanager"
)

type FileExtension string

const (
	FileExtensionYAML FileExtension = "yaml"
	FileExtensionYML  FileExtension = "yml"
	FileExtensionJSON FileExtension = "json"
)

func GetResourceMgrByType(resourceType string) (resouremanager.Resource, error) {
	switch strings.ToLower(resourceType) {
	case "hostendpoint", "hep":
		return resouremanager.NewHEP(), nil
	case "globalnetworkset", "gns":
		return resouremanager.NewGNS(), nil
	case "globalnetworkpolicy", "gnp":
		return resouremanager.NewGNP(), nil
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

type ResourceFile struct {
	Name    string
	Content interface{}
}

func GetResourceFilesByFileNames[T any](fileNames []string) ([]*ResourceFile, error) {
	var resources []*ResourceFile

	for _, fileName := range fileNames {
		resource, err := GetResourceFileByFileName[T](fileName)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func GetResourceFileByFileName[T any](fileName string) (*ResourceFile, error) {
	var input T
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open file %q: %w", fileName, err)
	}
	defer f.Close()

	contentFile, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading file %q: %w", fileName, err)
	}
	fileExtension := filepath.Ext(f.Name())
	switch FileExtension(strings.TrimLeft(fileExtension, ".")) {
	case FileExtensionYAML, FileExtensionYML:
		if err = yaml.Unmarshal(contentFile, &input); err != nil {
			return nil, fmt.Errorf("error parsing file %q: %w", fileName, err)
		}
	case FileExtensionJSON:
		if err = json.Unmarshal(contentFile, &input); err != nil {
			return nil, fmt.Errorf("error parsing file %q: %w", fileName, err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %q", fileExtension)
	}

	return &ResourceFile{
		Name:    fileName,
		Content: &input,
	}, nil
}
