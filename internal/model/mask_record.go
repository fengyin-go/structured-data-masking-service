package model

import (
	"strings"
)

// MaskRecord 表示一条脱敏记录。
type MaskRecord struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	RuleID    string `json:"rule_id"`
	Original  string `json:"original"`
	Masked    string `json:"masked"`
	CreatedAt int64  `json:"created_at"`
}

func (r *MaskRecord) Validate() error {
	r.TaskID = strings.TrimSpace(r.TaskID)
	r.RuleID = strings.TrimSpace(r.RuleID)
	r.Original = strings.TrimSpace(r.Original)

	if r.TaskID == "" {
		return NewValidationError("task_id", "任务 ID 不能为空")
	}
	if r.RuleID == "" {
		return NewValidationError("rule_id", "规则 ID 不能为空")
	}
	if r.Original == "" {
		return NewValidationError("original", "原始值不能为空")
	}
	return nil
}

// MaskRecordFilter 用于过滤脱敏记录。
type MaskRecordFilter struct {
	TaskID  string
	RuleID  string
	Keyword string
}

func (f MaskRecordFilter) Match(r *MaskRecord) bool {
	if f.TaskID != "" && r.TaskID != f.TaskID {
		return false
	}
	if f.RuleID != "" && r.RuleID != f.RuleID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Original), k) &&
			!strings.Contains(strings.ToLower(r.Masked), k) {
			return false
		}
	}
	return true
}
