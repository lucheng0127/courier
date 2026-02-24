package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server ServerConfig  `yaml:"server"`
	DB     DBConfig      `yaml:"db"`
	Models []ModelConfig `yaml:"models"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DBConfig 数据库配置
type DBConfig struct {
	DataSourceName string `yaml:"data_source_name"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	Name     string `yaml:"name"`
	Provider string `yaml:"provider"`
	BaseURL  string `yaml:"base_url"`
	APIKey   string `yaml:"api_key"`
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 替换环境变量
	err = cfg.replaceEnvVars()
	if err != nil {
		return nil, fmt.Errorf("替换环境变量失败: %w", err)
	}

	// 验证配置
	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// replaceEnvVars 替换配置中的环境变量
func (c *Config) replaceEnvVars() error {
	envVarPattern := regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

	for i := range c.Models {
		matches := envVarPattern.FindAllStringSubmatch(c.Models[i].APIKey, -1)
		for _, match := range matches {
			if len(match) > 1 {
				envValue := os.Getenv(match[1])
				if envValue == "" {
					return fmt.Errorf("环境变量 %s 未设置", match[1])
				}
				c.Models[i].APIKey = strings.ReplaceAll(c.Models[i].APIKey, match[0], envValue)
			}
		}
	}
	return nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证模型名称唯一性
	modelNames := make(map[string]bool)
	for i, model := range c.Models {
		if model.Name == "" {
			return fmt.Errorf("模型 %d: name 不能为空", i)
		}
		if model.BaseURL == "" {
			return fmt.Errorf("模型 %s: base_url 不能为空", model.Name)
		}
		if model.APIKey == "" {
			return fmt.Errorf("模型 %s: api_key 不能为空", model.Name)
		}
		if modelNames[model.Name] {
			return fmt.Errorf("模型名称重复: %s", model.Name)
		}
		modelNames[model.Name] = true
	}
	return nil
}
