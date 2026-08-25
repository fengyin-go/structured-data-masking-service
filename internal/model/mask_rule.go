package model

import (
	"strings"
)

const (
	RuleTypePhone  = "phone"
	RuleTypeIDCard = "idcard"
	RuleTypeEmail  = "email"
	RuleTypeCustom = "custom"
	RuleTypeRegex  = "regex"
)

var validRuleTypes = map[string]bool{
	RuleTypePhone:  true,
	RuleTypeIDCard: true,
	RuleTypeEmail:  true,
	RuleTypeCustom: true,
	RuleTypeRegex:  true,
}

// MaskRule 表示一条脱敏规则。
type MaskRule struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	FieldName    string `json:"field_name"`
	RuleType     string `json:"rule_type"`
	KeepPrefix   int    `json:"keep_prefix"`
	KeepSuffix   int    `json:"keep_suffix"`
	MaskChar     string `json:"mask_char"`
	RegexPattern string `json:"regex_pattern"`
	Enabled      bool   `json:"enabled"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

func (r *MaskRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.FieldName = strings.TrimSpace(r.FieldName)
	r.RuleType = strings.TrimSpace(r.RuleType)
	r.MaskChar = strings.TrimSpace(r.MaskChar)
	r.RegexPattern = strings.TrimSpace(r.RegexPattern)

	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.FieldName == "" {
		return NewValidationError("field_name", "字段名称不能为空")
	}
	if r.RuleType == "" {
		return NewValidationError("rule_type", "规则类型不能为空")
	}
	if !validRuleTypes[r.RuleType] {
		return NewValidationError("rule_type", "规则类型不合法")
	}
	if r.MaskChar == "" {
		r.MaskChar = "*"
	}
	if r.RuleType == RuleTypeRegex && r.RegexPattern == "" {
		return NewValidationError("regex_pattern", "正则类型规则必须提供正则表达式")
	}
	if r.KeepPrefix < 0 {
		return NewValidationError("keep_prefix", "保留前缀长度不能为负数")
	}
	if r.KeepSuffix < 0 {
		return NewValidationError("keep_suffix", "保留后缀长度不能为负数")
	}
	return nil
}

// MaskRuleFilter 用于过滤脱敏规则。
type MaskRuleFilter struct {
	RuleType string
	Enabled  *bool
	Keyword  string
}

func (f MaskRuleFilter) Match(r *MaskRule) bool {
	if f.RuleType != "" && r.RuleType != f.RuleType {
		return false
	}
	if f.Enabled != nil && r.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) &&
			!strings.Contains(strings.ToLower(r.FieldName), k) {
			return false
		}
	}
	return true
}
