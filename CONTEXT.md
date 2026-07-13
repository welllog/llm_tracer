# Domain Context

LLM Tracer：一个本地代理 + 控制台。代理转发 LLM 请求（OpenAI / Responses / Anthropic），落库记录每一轮交换；控制台按「会话」与「日志」两个视角浏览、检索、统计。

## 领域概念

### Session（会话）· 一等实体
一组属于同一对话的 LLM 轮次。**拥有自己的表 `sessions`**（见 db.go），不是从 `logs` 实时派生。

- **身份**：`session_id`（TEXT PRIMARY KEY）。代理总会赋一个真实 id；空值在写入时合成为 `standalone-<log_id>` 并回写 `logs.session_id`，使 log↔session 的 join 键即 `session_id` 本身。前端按 `standalone-` 前缀识别单条独立请求。
- **生命周期**：首条 log 到来时创建；**整组 log 被删尽**时销毁（删 sessions 行）。

### Log（请求日志）· `logs` 表 / `UnifiedLog`
单次 LLM 交换的完整记录：标量（id/provider/model/path/status/tokens/…）+ `prompt_json`（累积历史）+ `response_json` + `tools_json` + `raw_request`/`raw_response`。`session_seq` 标记其在会话中的**生命周期位次**。

### SessionMetadata / LogSummary
控制台列表视图的投影结构（`SessionMetadata` 用于会话列表，`LogSummary` 用于日志列表）。接口即测试面--这两个结构是稳定的契约。

## 不变量 · 字段分档

`sessions` 行的字段分两档，删一条 log 时行为不同（见 ADR-0001）：

- **Lifetime（生命周期累计，单删不变）**：`msg_count`、`sum_tokens`/`sum_input_uncached`/`sum_input_cached`/`sum_output`、`summary`、`summary_finalized`，以及 `logs.session_seq`。
  - `insert` 时增量更新；删单条 log **不回退**（消耗是历史事实）。
  - 仅在整组销毁（删尽会话）时随 sessions 行消失。
  - `session_seq` 是每条 log 的 lifetime 位次（insert 时 = `msg_count`，1 起编）；删中间留缺口，故 `session_seq ∈ [1, msg_count]`，恒有 `max(session_seq) ≤ msg_count`，等号当末轮仍在时成立。
- **Current（当前态，单删重算）**：`first_log_id`/`last_log_id`/`start_time`/`end_time`/`last_model`/`last_provider`。
  - 必须指向真实现存 log；删一条后按现存 logs 重算（标量子查询）。
- **Identity**：`session_id`（不变）。

## 摘要语义

`summary` 是会话**标题**：生命周期内**第一条 meaningful 的 user 消息**，一次设定后稳定（命中 `isMeaningfulSummary` 即冻结于 `summary_finalized`）。未命中 meaningful 前，首条非空 user 消息作暂存；命中后升级并锁定。反范式存储--删源 log 不清空标题。用内存中的 `log.Prompt` 经 `extractPromptSummary` 现算，零 DB 读。

## ADR
- `docs/adr/0001-session-as-first-class-entity.md` · Session 物化与 lifetime/current 字段分档。

## 文件地图
- `db.go` · 存储层（`sessions`/`logs`/`log_handles`）、`InsertLog`/`DeleteLog`/`DeleteSessionLogs`、`GetSessions`/`GetLogs`/`GetLogDetail`/`GetSessionLogs`/`GetStats`。
- `main.go` · 代理转发、SSE 扫描、REST API、前端托管。
- `parser.go` · OpenAI/Responses/Anthropic 请求与响应解析。
- `frontend/src/App.jsx` · 控制台前端。
