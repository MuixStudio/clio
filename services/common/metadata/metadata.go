package metadata

import (
	"runtime"
)

// ServiceMetadata 服务元数据
type ServiceMetadata struct {
	Name        string
	Version     string
	Description string
	BuildTime   string
	GitCommit   string
	GoVersion   string
}

var (
	// 必需字段
	Name        = "unknown"        // 服务名称
	Version     = "0.0.0"          // 版本号
	Description = "No description" // 服务描述

	// 构建信息（可选）
	BuildTime = "unknown"         // 构建时间
	GitCommit = "unknown"         // Git commit hash
	GoVersion = runtime.Version() // Go 版本
)

// Get 获取服务元数据
func Get() *ServiceMetadata {
	return &ServiceMetadata{
		Name:        Name,
		Version:     Version,
		Description: Description,
		BuildTime:   BuildTime,
		GitCommit:   GitCommit,
		GoVersion:   GoVersion,
	}
}

// GetName 获取服务名称
func GetName() string {
	return Name
}

// GetVersion 获取版本号
func GetVersion() string {
	return Version
}

// GetDescription 获取服务描述
func GetDescription() string {
	return Description
}

// GetBuildTime 获取构建时间
func GetBuildTime() string {
	return BuildTime
}

// GetGitCommit 获取 Git commit
func GetGitCommit() string {
	return GitCommit
}

// GetGoVersion 获取 Go 版本
func GetGoVersion() string {
	return GoVersion
}

// String 返回格式化的元数据字符串
func (m *ServiceMetadata) String() string {
	return m.Name + " v" + m.Version
}

// FullString 返回完整的元数据字符串
func (m *ServiceMetadata) FullString() string {
	s := m.Name + " v" + m.Version
	if m.Description != "" && m.Description != "No description" {
		s += " - " + m.Description
	}
	if m.BuildTime != "unknown" {
		s += "\nBuilt at: " + m.BuildTime
	}
	if m.GitCommit != "unknown" {
		s += "\nCommit: " + m.GitCommit
	}
	if m.GoVersion != "unknown" {
		s += "\nGo: " + m.GoVersion
	}
	return s
}
