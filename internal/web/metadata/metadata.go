package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ServiceMetadata 服务元数据
type ServiceMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// DefaultMetadata 默认元数据
var DefaultMetadata = ServiceMetadata{
	Name:        "unknown-service",
	Version:     "0.0.0",
	Description: "No description",
}

// LoadServiceMetadata 加载服务元数据
// 会依次查找：当前目录 -> 父目录 -> 再上级目录
func LoadServiceMetadata() (*ServiceMetadata, error) {
	// 尝试查找的文件名列表
	filenames := []string{"service.json", ".service.json"}

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	// 向上查找最多3层
	searchDirs := []string{
		cwd,
		filepath.Dir(cwd),
		filepath.Dir(filepath.Dir(cwd)),
	}

	// 在每个目录中尝试所有文件名
	for _, dir := range searchDirs {
		for _, filename := range filenames {
			path := filepath.Join(dir, filename)
			metadata, err := loadMetadataFromFile(path)
			if err == nil {
				return metadata, nil
			}
		}
	}

	return nil, fmt.Errorf("service.json not found in current directory or parent directories")
}

// LoadServiceMetadataOrDefault 加载服务元数据，失败则返回默认值
func LoadServiceMetadataOrDefault() *ServiceMetadata {
	metadata, err := LoadServiceMetadata()
	if err != nil {
		return &DefaultMetadata
	}
	return metadata
}

// loadMetadataFromFile 从指定文件加载元数据
func loadMetadataFromFile(path string) (*ServiceMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata ServiceMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	// 验证必需字段
	if metadata.Name == "" {
		return nil, fmt.Errorf("service name is required in %s", path)
	}
	if metadata.Version == "" {
		return nil, fmt.Errorf("service version is required in %s", path)
	}

	return &metadata, nil
}

// MustLoadServiceMetadata 加载服务元数据，失败则 panic
// 仅用于初始化阶段
func MustLoadServiceMetadata() *ServiceMetadata {
	metadata, err := LoadServiceMetadata()
	if err != nil {
		panic(fmt.Sprintf("Failed to load service metadata: %v", err))
	}
	return metadata
}
