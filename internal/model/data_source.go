package model

import (
	"strings"
)

const (
	DataSourceTypeDatabase = "database"
	DataSourceTypeAPI      = "api"
	DataSourceTypeFile     = "file"
	DataSourceTypeMessage  = "message"
)

var validDataSourceTypes = map[string]bool{
	DataSourceTypeDatabase: true,
	DataSourceTypeAPI:      true,
	DataSourceTypeFile:     true,
	DataSourceTypeMessage:  true,
}

// DataSource 表示一个数据源。
type DataSource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func (d *DataSource) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	d.Type = strings.TrimSpace(d.Type)
	d.Description = strings.TrimSpace(d.Description)

	if d.Name == "" {
		return NewValidationError("name", "数据源名称不能为空")
	}
	if d.Type == "" {
		return NewValidationError("type", "数据源类型不能为空")
	}
	if !validDataSourceTypes[d.Type] {
		return NewValidationError("type", "数据源类型不合法")
	}
	return nil
}

// DataSourceFilter 用于过滤数据源。
type DataSourceFilter struct {
	Type    string
	Keyword string
}

func (f DataSourceFilter) Match(d *DataSource) bool {
	if f.Type != "" && d.Type != f.Type {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Name), k) &&
			!strings.Contains(strings.ToLower(d.Description), k) {
			return false
		}
	}
	return true
}
