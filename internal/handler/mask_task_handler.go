package handler

import (
	"net/http"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerMaskTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tasks", s.createMaskTask)
	mux.HandleFunc("GET /api/tasks", s.listMaskTasks)
	mux.HandleFunc("GET /api/tasks/{id}", s.getMaskTask)
	mux.HandleFunc("PUT /api/tasks/{id}/status", s.updateMaskTaskStatus)
	mux.HandleFunc("DELETE /api/tasks/{id}", s.deleteMaskTask)
	mux.HandleFunc("POST /api/tasks/{id}/run", s.runMaskTask)
}

type createMaskTaskRequest struct {
	DataSourceID string   `json:"data_source_id"`
	RuleIDs      []string `json:"rule_ids"`
	Total        int      `json:"total"`
}

func (s *Server) createMaskTask(w http.ResponseWriter, r *http.Request) {
	var req createMaskTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateMaskTask(model.MaskTask{
		DataSourceID: req.DataSourceID,
		RuleIDs:      req.RuleIDs,
		Total:        req.Total,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listMaskTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MaskTaskFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMaskTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetMaskTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

type updateMaskTaskStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) updateMaskTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMaskTaskStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateMaskTaskStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMaskTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) runMaskTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.UpdateMaskTaskStatus(id, model.TaskStatusRunning)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}
