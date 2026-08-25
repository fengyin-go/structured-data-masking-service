// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"datamasking/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// MaskRule
	CreateMaskRule(r *model.MaskRule) error
	GetMaskRule(id string) (*model.MaskRule, error)
	ListMaskRules() []*model.MaskRule
	UpdateMaskRule(r *model.MaskRule) error
	DeleteMaskRule(id string) error

	// DataSource
	CreateDataSource(d *model.DataSource) error
	GetDataSource(id string) (*model.DataSource, error)
	ListDataSources() []*model.DataSource
	UpdateDataSource(d *model.DataSource) error
	DeleteDataSource(id string) error

	// MaskTask
	CreateMaskTask(t *model.MaskTask) error
	GetMaskTask(id string) (*model.MaskTask, error)
	ListMaskTasks() []*model.MaskTask
	UpdateMaskTask(t *model.MaskTask) error
	DeleteMaskTask(id string) error

	// MaskRecord
	CreateMaskRecord(r *model.MaskRecord) error
	GetMaskRecord(id string) (*model.MaskRecord, error)
	ListMaskRecords() []*model.MaskRecord
	ListMaskRecordsByTask(taskID string) []*model.MaskRecord
	DeleteMaskRecord(id string) error
}
