package store

type PolicyRegistry struct {
	policies map[string]string
}

func (r *PolicyRegistry) Put(field, policy string) {
	r.policies[field] = policy
}

func (r *PolicyRegistry) Get(field string) string {
	return r.policies[field]
}
