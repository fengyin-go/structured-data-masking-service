package store

// PolicyRegistry 维护字段到脱敏策略的映射。
// 零值 PolicyRegistry 可直接使用：首次 Put 会惰性初始化内部 map。
type PolicyRegistry struct {
	policies map[string]string
}

// NewPolicyRegistry 创建一个空的策略注册表。
func NewPolicyRegistry() *PolicyRegistry {
	return &PolicyRegistry{policies: make(map[string]string)}
}

// Put 写入或更新指定字段的脱敏策略。
// 零值注册表首次调用时会自动初始化，因此无需构造即可使用。
func (r *PolicyRegistry) Put(field, policy string) {
	if r.policies == nil {
		r.policies = make(map[string]string)
	}
	r.policies[field] = policy
}

// Get 读取指定字段的脱敏策略，字段不存在时返回空字符串。
func (r *PolicyRegistry) Get(field string) string {
	return r.policies[field]
}
