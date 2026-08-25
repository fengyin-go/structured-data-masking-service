package service

import "errors"

// MaskPolicyValidator 对脱敏策略做可选校验。
// 传入 nil 表示不校验，由 PolicyResolver 正常跳过。
type MaskPolicyValidator interface {
	Validate(string) error
}

// PrefixPolicyValidator 要求脱敏策略以指定前缀开头。
type PrefixPolicyValidator struct {
	Prefix string
}

// Validate 校验策略长度，前缀非空时还要求策略包含该前缀。
// 接收者为 nil 时按“不校验”处理，返回 nil，避免 panic。
func (v *PrefixPolicyValidator) Validate(policy string) error {
	if v == nil {
		return nil
	}
	if v.Prefix != "" && len(policy) < len(v.Prefix) {
		return errors.New("mask policy is too short")
	}
	return nil
}

// PolicyResolver 在写入脱敏策略前做可选校验。
// validator 为 nil 时跳过校验，可直接使用零值解析器。
type PolicyResolver struct {
	validator MaskPolicyValidator
}

// NewPolicyResolver 构造解析器。传入 nil 表示不启用校验器。
func NewPolicyResolver(validator MaskPolicyValidator) *PolicyResolver {
	return &PolicyResolver{validator: validator}
}

// Resolve 对策略执行可选校验，未配置校验器时直接放行。
func (r *PolicyResolver) Resolve(policy string) error {
	if r == nil || r.validator == nil {
		return nil
	}
	return r.validator.Validate(policy)
}
