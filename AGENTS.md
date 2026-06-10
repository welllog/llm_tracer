# 项目结构说明

本文件只保留项目结构和最小测试约定，供后续接手的 Agent 快速定位代码。

## 项目结构

- `web_recorder/`
  - `main.go` - 代理转发、SSE 流式响应扫描、REST API、前端托管
  - `parser.go` - OpenAI / Responses / Anthropic 请求与响应解析
  - `db.go` - SQLite 日志存储、查询、统计、会话句柄映射
  - `main_test.go` - 后端离线集成测试与回归测试
  - `frontend/` - React + Vite 前端
  - `static/` - 前端构建产物，供 `go:embed` 使用
  - `Makefile` - 构建与运行命令

## 测试约定

- 会话归属问题一旦修复，必须把触发问题的真实请求/响应日志样本脱敏后直接写进 `main_test.go`，作为回归测试夹具。
- 单测不得依赖 `data/llm_tracer.db` 或任何持久化 db 文件；测试必须使用临时数据库或纯内存样本。
