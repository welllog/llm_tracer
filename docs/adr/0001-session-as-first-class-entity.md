# ADR-0001 · Session 物化为一等实体，字段分 lifetime/current 两档

- 状态：Accepted
- 日期：2026-07-13
- 关联：`/improve-codebase-architecture` Candidate 1；`CONTEXT.md`

## 背景

`GetSessions` 此前每次翻页都从 `logs` 表实时 `GROUP BY`（`CASE` 表达式分组，无法用索引）+ 排序派生会话，`EXPLAIN` 为 `SCAN logs` + 两棵临时 B-Tree，实测 ~290 ms，外加单连接上 ~20 次 N+1 摘要查询，且零缓存。Session 是用户导航的主单元却无自己的存储--浅模块。

## 决策

1. **物化 `sessions` 表**（`session_id` PK + 聚合 + 摘要 + `first/last_log_id`）。在 `InsertLog` 事务内增量维护；读路径退化为 `SELECT … FROM sessions ORDER BY last_log_id DESC`。
2. **字段分两档**，删一条 log 时行为不同：
   - **Lifetime**（`msg_count`、`sum_tokens`/`sum_input_*`/`sum_output`、`summary`/`summary_finalized`、`logs.session_seq`）：生命周期累计，insert 增量更新，**单删不回退**（消耗是历史事实）。仅整组销毁时随 sessions 行消失。`session_seq` 为 lifetime 位次（insert 时=`msg_count`，1 起编），删中间留缺口，`session_seq ∈ [1, msg_count]`，恒有 `max(session_seq) ≤ msg_count`（等号当末轮仍在时成立）。
   - **Current**（`first/last_log_id`、`start/end_time`、`last_model/provider`）：指向真实现存 log，删一条后按现存 logs 重算。
3. **摘要语义**：生命周期内首条 `isMeaningful` 的 user 消息，命中即冻结；未命中前首条非空暂存；反范式存活（删源不清空）。
4. **空 `session_id`** 写入时合成 `standalone-<log_id>` 并回写 `logs.session_id`，使 join 键统一。
5. **会话搜索**改匹配 `sessions.summary`（不再 `LIKE` 扫 raw 大字段）；日志列表搜索仍搜 raw。

## 理由

- 单一真相源：聚合逻辑从 3 个读路径浓缩进 1 个写路径（删除测试通过）。
- Lifetime 分档避免「删一条记录就抹掉历史消耗」的失真；Current 分档保证 first/last 指向真实 log。两档职责清晰。
- 删中间 log 在 `session_seq` 留缺口，是有意的信息（暴露删除），而非误导性的重排连续序号。

## 后果

- 一次性 schema 迁移 + 回填（`createTables` 内 guarded，sessions 空则重建，自愈）。
- `DeleteLog` 需先读 `session_id` 再重算结构字段。
- `session_seq ∈ [1, msg_count]`（删中间留缺口）、「单删不减 lifetime」与「末轮仍在时 max(session_seq)==msg_count」成为新的不变量，需测试覆盖。
- 前端无需改动（接口形状不变；`session_seq` 缺口照显）。

## 不在本决策范围

- per-turn slim 摘要列（Candidate 2）、存储去重（Candidate 3）、连接池放开（Candidate 4）--各自独立。
