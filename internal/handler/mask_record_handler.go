package handler

import (
	"net/http"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerMaskRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/records", s.createMaskRecord)
	mux.HandleFunc("GET /api/records", s.listMaskRecords)
	mux.HandleFunc("GET /api/records/{id}", s.getMaskRecord)
	mux.HandleFunc("GET /api/tasks/{id}/records", s.listMaskRecordsByTask)
	mux.HandleFunc("DELETE /api/records/{id}", s.deleteMaskRecord)
}

type createMaskRecordRequest struct {
	TaskID   string `json:"task_id"`
	RuleID   string `json:"rule_id"`
	Original string `json:"original"`
	Masked   string `json:"masked"`
}

func (s *Server) createMaskRecord(w http.ResponseWriter, r *http.Request) {
	var req createMaskRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateMaskRecord(model.MaskRecord{
		TaskID:   req.TaskID,
		RuleID:   req.RuleID,
		Original: req.Original,
		Masked:   req.Masked,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listMaskRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MaskRecordFilter{
		TaskID:  r.URL.Query().Get("task_id"),
		RuleID:  r.URL.Query().Get("rule_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMaskRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMaskRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetMaskRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) listMaskRecordsByTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListMaskRecordsByTask(id, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) deleteMaskRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMaskRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
