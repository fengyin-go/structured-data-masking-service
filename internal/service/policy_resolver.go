package service

import "errors"

type MaskPolicyValidator interface {
	Validate(string) error
}

type PrefixPolicyValidator struct {
	Prefix string
}

func (v *PrefixPolicyValidator) Validate(policy string) error {
	if v == nil {
		panic("nil policy validator")
	}
	if v.Prefix != "" && len(policy) < len(v.Prefix) {
		return errors.New("mask policy is too short")
	}
	return nil
}

type PolicyResolver struct {
	validator MaskPolicyValidator
}

func NewPolicyResolver(validator *PrefixPolicyValidator) *PolicyResolver {
	return &PolicyResolver{validator: validator}
}

func (r *PolicyResolver) Resolve(policy string) error {
	if r.validator != nil {
		return r.validator.Validate(policy)
	}
	return nil
}
