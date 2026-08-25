package handler

import (
	"net/http"
	"strconv"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerMaskRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rules", s.createMaskRule)
	mux.HandleFunc("GET /api/rules", s.listMaskRules)
	mux.HandleFunc("GET /api/rules/{id}", s.getMaskRule)
	mux.HandleFunc("PUT /api/rules/{id}", s.updateMaskRule)
	mux.HandleFunc("DELETE /api/rules/{id}", s.deleteMaskRule)
	mux.HandleFunc("POST /api/rules/batch", s.batchCreateMaskRules)
	mux.HandleFunc("POST /api/rules/mask", s.maskSingleValue)
	mux.HandleFunc("POST /api/rules/mask-batch", s.maskBatchValues)
}

type createMaskRuleRequest struct {
	Name         string `json:"name"`
	FieldName    string `json:"field_name"`
	RuleType     string `json:"rule_type"`
	KeepPrefix   int    `json:"keep_prefix"`
	KeepSuffix   int    `json:"keep_suffix"`
	MaskChar     string `json:"mask_char"`
	RegexPattern string `json:"regex_pattern"`
	Enabled      bool   `json:"enabled"`
}

func (s *Server) createMaskRule(w http.ResponseWriter, r *http.Request) {
	var req createMaskRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateMaskRule(model.MaskRule{
		Name:         req.Name,
		FieldName:    req.FieldName,
		RuleType:     req.RuleType,
		KeepPrefix:   req.KeepPrefix,
		KeepSuffix:   req.KeepSuffix,
		MaskChar:     req.MaskChar,
		RegexPattern: req.RegexPattern,
		Enabled:      req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listMaskRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	var enabled *bool
	if v := r.URL.Query().Get("enabled"); v != "" {
		b, _ := strconv.ParseBool(v)
		enabled = &b
	}
	filter := model.MaskRuleFilter{
		RuleType: r.URL.Query().Get("rule_type"),
		Enabled:  enabled,
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMaskRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMaskRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetMaskRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) updateMaskRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createMaskRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateMaskRule(id, model.MaskRule{
		Name:         req.Name,
		FieldName:    req.FieldName,
		RuleType:     req.RuleType,
		KeepPrefix:   req.KeepPrefix,
		KeepSuffix:   req.KeepSuffix,
		MaskChar:     req.MaskChar,
		RegexPattern: req.RegexPattern,
		Enabled:      req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteMaskRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMaskRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchCreateMaskRulesRequest struct {
	Items []createMaskRuleRequest `json:"items"`
}

func (s *Server) batchCreateMaskRules(w http.ResponseWriter, r *http.Request) {
	var req batchCreateMaskRulesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.MaskRule, len(req.Items))
	for i, item := range req.Items {
		inputs[i] = model.MaskRule{
			Name:         item.Name,
			FieldName:    item.FieldName,
			RuleType:     item.RuleType,
			KeepPrefix:   item.KeepPrefix,
			KeepSuffix:   item.KeepSuffix,
			MaskChar:     item.MaskChar,
			RegexPattern: item.RegexPattern,
			Enabled:      item.Enabled,
		}
	}
	items, err := s.svc.BatchCreateMaskRules(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, items)
}

type maskSingleValueRequest struct {
	RuleID   string `json:"rule_id"`
	Original string `json:"original"`
}

type maskSingleValueResponse struct {
	Original string `json:"original"`
	Masked   string `json:"masked"`
}

func (s *Server) maskSingleValue(w http.ResponseWriter, r *http.Request) {
	var req maskSingleValueRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.RuleID == "" {
		httpx.BadRequest(w, "rule_id 不能为空")
		return
	}
	if req.Original == "" {
		httpx.BadRequest(w, "original 不能为空")
		return
	}
	rule, err := s.svc.GetMaskRule(req.RuleID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	masked, err := s.svc.MaskValue(rule, req.Original)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, maskSingleValueResponse{Original: req.Original, Masked: masked})
}

type maskBatchValuesRequest struct {
	RuleID    string   `json:"rule_id"`
	Originals []string `json:"originals"`
}

type maskBatchValuesResponse struct {
	Results []maskSingleValueResponse `json:"results"`
}

func (s *Server) maskBatchValues(w http.ResponseWriter, r *http.Request) {
	var req maskBatchValuesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.RuleID == "" {
		httpx.BadRequest(w, "rule_id 不能为空")
		return
	}
	rule, err := s.svc.GetMaskRule(req.RuleID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	results := make([]maskSingleValueResponse, 0, len(req.Originals))
	for _, orig := range req.Originals {
		masked, merr := s.svc.MaskValue(rule, orig)
		if merr != nil {
			writeServiceError(w, merr)
			return
		}
		results = append(results, maskSingleValueResponse{Original: orig, Masked: masked})
	}
	httpx.OK(w, maskBatchValuesResponse{Results: results})
}
