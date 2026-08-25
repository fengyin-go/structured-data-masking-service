package service

import (
	"fmt"
	"regexp"
	"strings"

	"datamasking/internal/model"
)

// MaskValue 根据规则对原始值进行脱敏。
func (s *Service) MaskValue(rule *model.MaskRule, original string) (string, error) {
	if original == "" {
		return "", model.NewValidationError("original", "原始值不能为空")
	}
	switch rule.RuleType {
	case model.RuleTypePhone:
		return maskPhone(original, rule.KeepPrefix, rule.KeepSuffix, rule.MaskChar)
	case model.RuleTypeIDCard:
		return maskIDCard(original, rule.KeepPrefix, rule.KeepSuffix, rule.MaskChar)
	case model.RuleTypeEmail:
		return maskEmail(original, rule.KeepPrefix, rule.KeepSuffix, rule.MaskChar)
	case model.RuleTypeCustom:
		return maskCustom(original, rule.KeepPrefix, rule.KeepSuffix, rule.MaskChar)
	case model.RuleTypeRegex:
		return maskRegex(original, rule.RegexPattern, rule.MaskChar)
	default:
		return "", model.NewValidationError("rule_type", "不支持的规则类型")
	}
}

// MaskJSONBatch 对 JSON 对象列表中的指定字段按规则批量脱敏。
func (s *Service) MaskJSONBatch(items []map[string]interface{}, fieldRuleMap map[string]*model.MaskRule) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, len(items))
	for i, item := range items {
		newItem := make(map[string]interface{})
		for k, v := range item {
			newItem[k] = v
		}
		for field, rule := range fieldRuleMap {
			if val, ok := newItem[field]; ok {
				if str, ok := val.(string); ok {
					masked, err := s.MaskValue(rule, str)
					if err != nil {
						return nil, err
					}
					newItem[field] = masked
				}
			}
		}
		results[i] = newItem
	}
	return results, nil
}

func maskPhone(original string, keepPrefix, keepSuffix int, maskChar string) (string, error) {
	if len(original) < 7 {
		return maskCustom(original, keepPrefix, keepSuffix, maskChar)
	}
	prefix := 3
	suffix := 4
	if keepPrefix > 0 {
		prefix = keepPrefix
	}
	if keepSuffix > 0 {
		suffix = keepSuffix
	}
	if prefix+suffix >= len(original) {
		return strings.Repeat(maskChar, len(original)), nil
	}
	return original[:prefix] + strings.Repeat(maskChar, len(original)-prefix-suffix) + original[len(original)-suffix:], nil
}

func maskIDCard(original string, keepPrefix, keepSuffix int, maskChar string) (string, error) {
	if len(original) != 15 && len(original) != 18 {
		return maskCustom(original, keepPrefix, keepSuffix, maskChar)
	}
	prefix := 4
	suffix := 4
	if keepPrefix > 0 {
		prefix = keepPrefix
	}
	if keepSuffix > 0 {
		suffix = keepSuffix
	}
	if prefix+suffix >= len(original) {
		return strings.Repeat(maskChar, len(original)), nil
	}
	return original[:prefix] + strings.Repeat(maskChar, len(original)-prefix-suffix) + original[len(original)-suffix:], nil
}

func maskEmail(original string, keepPrefix, keepSuffix int, maskChar string) (string, error) {
	parts := strings.Split(original, "@")
	if len(parts) != 2 {
		return maskCustom(original, keepPrefix, keepSuffix, maskChar)
	}
	local := parts[0]
	domain := parts[1]
	prefix := 2
	suffix := 1
	if keepPrefix > 0 {
		prefix = keepPrefix
	}
	if keepSuffix > 0 {
		suffix = keepSuffix
	}
	if prefix+suffix >= len(local) {
		local = strings.Repeat(maskChar, len(local))
	} else {
		local = local[:prefix] + strings.Repeat(maskChar, len(local)-prefix-suffix) + local[len(local)-suffix:]
	}
	return local + "@" + domain, nil
}

func maskCustom(original string, keepPrefix, keepSuffix int, maskChar string) (string, error) {
	if keepPrefix < 0 {
		keepPrefix = 0
	}
	if keepSuffix < 0 {
		keepSuffix = 0
	}
	if keepPrefix+keepSuffix >= len(original) {
		return strings.Repeat(maskChar, len(original)), nil
	}
	return original[:keepPrefix] + strings.Repeat(maskChar, len(original)-keepPrefix-keepSuffix) + original[len(original)-keepSuffix:], nil
}

func maskRegex(original string, pattern string, maskChar string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", model.NewValidationError("regex_pattern", fmt.Sprintf("正则表达式编译失败: %v", err))
	}
	return re.ReplaceAllString(original, strings.Repeat(maskChar, 3)), nil
}

// ApplyMaskRules 依次应用多条规则对原始值脱敏，返回最终 masked 值和命中规则列表。
func (s *Service) ApplyMaskRules(original string, rules []*model.MaskRule) (string, []string, error) {
	current := original
	hitRules := make([]string, 0)
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		masked, err := s.MaskValue(rule, current)
		if err != nil {
			return "", nil, err
		}
		if masked != current {
			hitRules = append(hitRules, rule.ID)
			current = masked
		}
	}
	return current, hitRules, nil
}
