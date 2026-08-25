package model

import (
	"strings"
)

const (
	TaskStatusPending = "pending"
	TaskStatusRunning = "running"
	TaskStatusDone    = "done"
	TaskStatusFailed  = "failed"
)

var validTaskStatuses = map[string]bool{
	TaskStatusPending: true,
	TaskStatusRunning: true,
	TaskStatusDone:    true,
	TaskStatusFailed:  true,
}

var taskTransitions = map[string]map[string]bool{
	TaskStatusPending: {TaskStatusRunning: true, TaskStatusFailed: true},
	TaskStatusRunning: {TaskStatusDone: true, TaskStatusFailed: true},
	TaskStatusDone:    {},
	TaskStatusFailed:  {},
}

// CanTransition 校验任务状态是否允许从 from 流转到 to。
func CanTransition(from, to string) bool {
	if m, ok := taskTransitions[from]; ok {
		return m[to]
	}
	return false
}

// MaskTask 表示一次脱敏任务。
type MaskTask struct {
	ID           string   `json:"id"`
	DataSourceID string   `json:"data_source_id"`
	RuleIDs      []string `json:"rule_ids"`
	Status       string   `json:"status"`
	Total        int      `json:"total"`
	MaskedCount  int      `json:"masked_count"`
	ErrorMsg     string   `json:"error_msg"`
	CreatedAt    int64    `json:"created_at"`
	UpdatedAt    int64    `json:"updated_at"`
}

func (t *MaskTask) Validate() error {
	t.DataSourceID = strings.TrimSpace(t.DataSourceID)
	if t.DataSourceID == "" {
		return NewValidationError("data_source_id", "数据源 ID 不能为空")
	}
	if len(t.RuleIDs) == 0 {
		return NewValidationError("rule_ids", "至少选择一个脱敏规则")
	}
	if t.Status == "" {
		t.Status = TaskStatusPending
	}
	if !validTaskStatuses[t.Status] {
		return NewValidationError("status", "任务状态不合法")
	}
	if t.Total < 0 {
		return NewValidationError("total", "总记录数不能为负数")
	}
	if t.MaskedCount < 0 {
		return NewValidationError("masked_count", "脱敏记录数不能为负数")
	}
	return nil
}

// MaskTaskFilter 用于过滤脱敏任务。
type MaskTaskFilter struct {
	Status  string
	Keyword string
}

func (f MaskTaskFilter) Match(t *MaskTask) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.DataSourceID), k) &&
			!strings.Contains(strings.ToLower(t.ErrorMsg), k) {
			return false
		}
	}
	return true
}
