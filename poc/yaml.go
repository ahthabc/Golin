package poc

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

//go:embed yaml-poc/*.yaml
var yamlFS embed.FS

type Config struct {
	Name          string            `yaml:"name"`          //漏洞名称
	Description   string            `yaml:"description"`   //漏洞描述
	Method        string            `yaml:"method"`        //请求类型
	Path          []string          `yaml:"path"`          //请求路径
	Body          string            `yaml:"body"`          //发送值
	Headers       map[string]string `yaml:"headers"`       //设置Headers
	Expression    Expression        `yaml:"expression"`    //返回值
	AlwaysExecute bool              `yaml:"alwaysExecute"` //是否直接执行不考虑app等组件
	Timeout       time.Duration     `yaml:"timeout"`       //等待时常
}

type Expression struct {
	Status      int      `yaml:"status"`       //返回的状态码
	ContentType string   `yaml:"content_type"` //返回头
	BodyALL     []string `yaml:"body_all"`     //必须包含所有特征
	BodyAny     []string `yaml:"body_any"`     //包含任意特征
	Time        float64  `yaml:"sleep"`        //总共耗时
}

// parseConfigs 解析yaml文件
func parseConfigs(dir string) ([]Config, error) {
	var configs []Config

	dirEntries, err := fs.ReadDir(yamlFS, dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".yaml") {
			data, err := fs.ReadFile(yamlFS, dir+"/"+entry.Name())
			if err != nil {
				return nil, err
			}

			var config Config
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				return nil, fmt.Errorf("解析 %s 失败: %w", entry.Name(), err)
			}

			configs = append(configs, config)
		}
	}

	return configs, nil
}
