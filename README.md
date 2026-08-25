# 数据脱敏服务 (datamasking)

纯 Go 标准库实现的内存级数据脱敏后端服务，支持脱敏规则管理、数据源管理、脱敏任务与记录以及核心脱敏算法。

## 运行

在 `origin/` 目录下执行：

```bash
/Users/fengyin/.workbuddy/binaries/go-dist/go/bin/go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/rules | 创建脱敏规则 |
| GET | /api/rules | 规则列表（支持 rule_type/enabled/keyword 过滤与分页） |
| GET | /api/rules/{id} | 规则详情 |
| PUT | /api/rules/{id} | 更新规则 |
| DELETE | /api/rules/{id} | 删除规则 |
| POST | /api/rules/batch | 批量创建规则 |
| POST | /api/rules/mask | 单值脱敏 |
| POST | /api/rules/mask-batch | 批量值脱敏 |
| POST | /api/sources | 创建数据源 |
| GET | /api/sources | 数据源列表（支持 type/keyword 过滤与分页） |
| GET | /api/sources/{id} | 数据源详情 |
| PUT | /api/sources/{id} | 更新数据源 |
| DELETE | /api/sources/{id} | 删除数据源 |
| POST | /api/tasks | 创建脱敏任务 |
| GET | /api/tasks | 任务列表（支持 status/keyword 过滤与分页） |
| GET | /api/tasks/{id} | 任务详情 |
| PUT | /api/tasks/{id}/status | 更新任务状态 |
| DELETE | /api/tasks/{id} | 删除任务 |
| POST | /api/tasks/{id}/run | 启动任务（pending→running） |
| POST | /api/records | 创建脱敏记录 |
| GET | /api/records | 脱敏记录列表（支持 task_id/rule_id/keyword 过滤与分页） |
| GET | /api/records/{id} | 记录详情 |
| GET | /api/tasks/{id}/records | 按任务查询记录 |
| DELETE | /api/records/{id} | 删除记录 |
| GET | /health | 健康检查 |

## 响应格式

统一 JSON 响应：

```json
{"code":0,"message":"ok","data":...}
```

错误码映射：400（校验错误）、404（记录不存在）、409（冲突）、500（服务器内部错误）。
