package handler

import (
	"net/http"

	"datamasking/internal/model"
	"datamasking/pkg/httpx"
)

func (s *Server) registerDataSourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/sources", s.createDataSource)
	mux.HandleFunc("GET /api/sources", s.listDataSources)
	mux.HandleFunc("GET /api/sources/{id}", s.getDataSource)
	mux.HandleFunc("PUT /api/sources/{id}", s.updateDataSource)
	mux.HandleFunc("DELETE /api/sources/{id}", s.deleteDataSource)
}

type createDataSourceRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (s *Server) createDataSource(w http.ResponseWriter, r *http.Request) {
	var req createDataSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateDataSource(model.DataSource{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listDataSources(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DataSourceFilter{
		Type:    r.URL.Query().Get("type"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDataSources(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDataSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetDataSource(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) updateDataSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createDataSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateDataSource(id, model.DataSource{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteDataSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDataSource(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
