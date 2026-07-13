package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DBManager struct {
	db *sql.DB
}

type UnifiedLog struct {
	ID                  int64         `json:"id"`
	Provider            string        `json:"provider"`
	Model               string        `json:"model"`
	Path                string        `json:"path"`
	StatusCode          int           `json:"statusCode"`
	DurationMs          int64         `json:"durationMs"`
	InputTokens         int           `json:"inputTokens"`
	OutputTokens        int           `json:"outputTokens"`
	TotalTokens         int           `json:"totalTokens"`
	CachedTokens        int           `json:"cachedTokens"`
	CacheReadTokens     int           `json:"cacheReadTokens"`
	CacheCreationTokens int           `json:"cacheCreationTokens"`
	Prompt              []ChatMessage `json:"prompt"`
	Response            ChatMessage   `json:"response"`
	Tools               []ToolDef     `json:"tools,omitempty"`
	RawRequest          string        `json:"rawRequest"`
	RawResponse         string        `json:"rawResponse"`
	ErrorMessage        string        `json:"errorMessage,omitempty"`
	CreatedAt           string        `json:"createdAt"`
	SessionID           string        `json:"sessionId,omitempty"`
	ParentID            *int64        `json:"parentId,omitempty"`
	ParentToolCallID    string        `json:"parentToolCallId,omitempty"`
	SessionSummary      string        `json:"sessionSummary,omitempty"`
	SessionSeq          int           `json:"sessionSeq,omitempty"`
	SubLogs             []UnifiedLog  `json:"subLogs,omitempty"`
	ClientFingerprint   string        `json:"-"`
	RequestHandles      []ConversationHandle `json:"-"`
	ResponseHandles     []ConversationHandle `json:"-"`
}

type LogSummary struct {
	ID                  int64  `json:"id"`
	Provider            string `json:"provider"`
	Model               string `json:"model"`
	Path                string `json:"path"`
	StatusCode          int    `json:"statusCode"`
	DurationMs          int64  `json:"durationMs"`
	InputTokens         int    `json:"inputTokens"`
	OutputTokens        int    `json:"outputTokens"`
	TotalTokens         int    `json:"totalTokens"`
	CachedTokens        int    `json:"cachedTokens"`
	CacheReadTokens     int    `json:"cacheReadTokens"`
	CacheCreationTokens int    `json:"cacheCreationTokens"`
	CreatedAt           string `json:"createdAt"`
	SessionID           string `json:"sessionId"`
	ParentID            *int64 `json:"parentId,omitempty"`
	ParentToolCallID    string `json:"parentToolCallId,omitempty"`
	SessionSummary      string `json:"sessionSummary,omitempty"`
	SessionSeq          int    `json:"sessionSeq,omitempty"`
}

type UsageStats struct {
	TotalCalls        int             `json:"totalCalls"`
	TotalInputTokens  int             `json:"totalInputTokens"`
	TotalOutputTokens int             `json:"totalOutputTokens"`
	TotalTokens       int             `json:"totalTokens"`
	TotalCachedTokens int             `json:"totalCachedTokens"`
	AvgDurationMs     int64           `json:"avgDurationMs"`
	SuccessCalls      int             `json:"successCalls"`
	SuccessRate       float64         `json:"successRate"`
	CallsByProvider   map[string]int  `json:"callsByProvider"`
	TokensByModel     map[string]int  `json:"tokensByModel"`
}

type SessionMetadata struct {
	SessionID                string `json:"sessionId"`
	SessionSummary           string `json:"sessionSummary"`
	StartTime                string `json:"startTime"`
	EndTime                  string `json:"endTime"`
	TotalTokens              int    `json:"totalTokens"`
	TotalInputUncachedTokens int    `json:"totalInputUncachedTokens"`
	TotalInputCachedTokens   int    `json:"totalInputCachedTokens"`
	TotalOutputTokens        int    `json:"totalOutputTokens"`
	MessageCount             int    `json:"messageCount"`
	Model                    string `json:"model"`
	Provider                 string `json:"provider"`
	FirstLogID               int64  `json:"firstLogId,omitempty"`
}

func (log *LogSummary) normalizeAnthropicTokens() {
	if log.Provider != "anthropic" {
		return
	}
	log.CachedTokens = log.CacheReadTokens
	log.InputTokens += log.CacheReadTokens
	log.TotalTokens = log.InputTokens + log.OutputTokens
}

func (log *UnifiedLog) normalizeAnthropicTokens() {
	if log.Provider != "anthropic" {
		return
	}
	log.CachedTokens = log.CacheReadTokens
	log.InputTokens += log.CacheReadTokens
	log.TotalTokens = log.InputTokens + log.OutputTokens
}

func InitDB(dbPath string) (*DBManager, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	// 设置 SQLite 的一些性能优化选项
	_, _ = db.Exec("PRAGMA journal_mode=WAL;")
	_, _ = db.Exec("PRAGMA synchronous=NORMAL;")
	// 提高忙等待上限，缓解并发写事务下的 "database is locked" 错误
	_, _ = db.Exec("PRAGMA busy_timeout = 15000;")
	// 写连接串行化，避免多写事务并发触发锁竞争；读连接不受限（WAL 下读不阻塞写）
	db.SetMaxOpenConns(1)

	mgr := &DBManager{db: db}
	if err := mgr.createTables(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return mgr, nil
}

func (mgr *DBManager) Close() error {
	if mgr.db != nil {
		return mgr.db.Close()
	}
	return nil
}

func (mgr *DBManager) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		provider TEXT NOT NULL,
		model TEXT NOT NULL,
		path TEXT NOT NULL,
		status_code INTEGER NOT NULL,
		duration_ms INTEGER NOT NULL,
		input_tokens INTEGER NOT NULL DEFAULT 0,
		output_tokens INTEGER NOT NULL DEFAULT 0,
		total_tokens INTEGER NOT NULL DEFAULT 0,
		cached_tokens INTEGER NOT NULL DEFAULT 0,
		cache_read_tokens INTEGER NOT NULL DEFAULT 0,
		cache_creation_tokens INTEGER NOT NULL DEFAULT 0,
		prompt_json TEXT,
		response_json TEXT,
		tools_json TEXT,
		raw_request TEXT,
		raw_response TEXT,
		error_message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		session_id TEXT,
		client_fingerprint TEXT,
		parent_id INTEGER,
		parent_tool_call_id TEXT,
		session_seq INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_logs_provider ON logs(provider);
	CREATE INDEX IF NOT EXISTS idx_logs_model ON logs(model);
	CREATE TABLE IF NOT EXISTS log_handles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		log_id INTEGER NOT NULL,
		source TEXT NOT NULL,
		handle_kind TEXT NOT NULL,
		handle_value TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS sessions (
		session_id         TEXT PRIMARY KEY,
		first_log_id       INTEGER NOT NULL,
		last_log_id        INTEGER NOT NULL,
		start_time         TEXT NOT NULL,
		end_time           TEXT NOT NULL,
		msg_count          INTEGER NOT NULL DEFAULT 0,
		sum_tokens         INTEGER NOT NULL DEFAULT 0,
		sum_input_uncached INTEGER NOT NULL DEFAULT 0,
		sum_input_cached   INTEGER NOT NULL DEFAULT 0,
		sum_output         INTEGER NOT NULL DEFAULT 0,
		last_model         TEXT NOT NULL DEFAULT '',
		last_provider      TEXT NOT NULL DEFAULT '',
		summary            TEXT NOT NULL DEFAULT '',
		summary_finalized  INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_sessions_last_log_id ON sessions(last_log_id DESC);
	`
	_, err := mgr.db.Exec(schema)
	if err != nil {
		return err
	}

	// 确保 session_id 列存在后，再创建其索引
	_, err = mgr.db.Exec("CREATE INDEX IF NOT EXISTS idx_logs_session_id ON logs(session_id);")
	_, _ = mgr.db.Exec("CREATE INDEX IF NOT EXISTS idx_logs_parent_id ON logs(parent_id);")
	_, _ = mgr.db.Exec("CREATE INDEX IF NOT EXISTS idx_logs_parent_tool_call_id ON logs(parent_tool_call_id);")
	_, _ = mgr.db.Exec("CREATE INDEX IF NOT EXISTS idx_logs_client_fingerprint ON logs(client_fingerprint);")
	_, _ = mgr.db.Exec("CREATE INDEX IF NOT EXISTS idx_log_handles_kind_value ON log_handles(handle_kind, handle_value);")
	_, _ = mgr.db.Exec("CREATE INDEX IF NOT EXISTS idx_log_handles_log_id ON log_handles(log_id);")
	_, _ = mgr.db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_log_handles_unique_entry ON log_handles(log_id, source, handle_kind, handle_value);")

	// === MIGRATION CODE（pre-C1 → C1，ADR-0001）=========================
	// 把旧库（无 sessions 表 / logs 无 session_seq 列）迁到物化会话架构。
	// 幂等：sessions 已填充时为 no-op；被清空则从 logs（真相源）自愈重建。
	// 待所有部署库都升级到 C1 后，可整体删除 backfillSessions + computeSessionSummaryInTx
	// + 本调用（runtime 路径 InsertLog/DeleteLog 不依赖它们）。
	if berr := mgr.backfillSessions(); berr != nil {
		return fmt.Errorf("backfill sessions: %w", berr)
	}
	return err
}

// backfillSessions 是 C1 迁移的入口（MIGRATION CODE，见 createTables 调用点注释）。
// 在 sessions 表为空时从 logs 重建会话实体（ADR-0001）。已填充则跳过（no-op）；
// sessions 被清空后下次启动即自愈重建。全程一个事务，保证原子。
func (mgr *DBManager) backfillSessions() error {
	var n int
	if err := mgr.db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}

	// 0. 老库 logs 表无 session_seq 列（CREATE TABLE IF NOT EXISTS 不会补）；补上。
	//    列已存在时 ALTER 返回错误，忽略即可。仅 pre-C1 迁移路径需要。
	_, _ = mgr.db.Exec("ALTER TABLE logs ADD COLUMN session_seq INTEGER NOT NULL DEFAULT 0")

	tx, err := mgr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 空 session_id -> 'standalone-<id>'，回写 logs.session_id，使 join 键统一。
	if _, err = tx.Exec("UPDATE logs SET session_id = 'standalone-' || id WHERE session_id IS NULL OR session_id = ''"); err != nil {
		return fmt.Errorf("backfill standalone keys: %w", err)
	}

	// 2. logs.session_seq = lifetime 位次（按 session_id 分组、id 升序编号；现存数据无删除，连续 1..N）。
	if _, err = tx.Exec(`
		UPDATE logs SET session_seq = (
			SELECT rn FROM (
				SELECT id, ROW_NUMBER() OVER (PARTITION BY session_id ORDER BY id ASC) AS rn FROM logs
			) WHERE id = logs.id
		)
	`); err != nil {
		return fmt.Errorf("backfill session_seq: %w", err)
	}

	// 3. sessions 聚合（last_model/provider 取最新一条 log 的值，而非 MAX(model)）。
	if _, err = tx.Exec(`
		WITH agg AS (
			SELECT session_id, MIN(id) first_id, MAX(id) last_id,
			       MIN(created_at) st, MAX(created_at) et,
			       COUNT(*) mc, SUM(total_tokens) st_tok,
			       SUM(CASE WHEN provider='anthropic' THEN input_tokens ELSE input_tokens - cached_tokens END) siu,
			       SUM(CASE WHEN provider='anthropic' THEN cache_read_tokens ELSE cached_tokens END) sic,
			       SUM(output_tokens) so_tok
			FROM logs GROUP BY session_id
		)
		INSERT INTO sessions (session_id, first_log_id, last_log_id, start_time, end_time,
			msg_count, sum_tokens, sum_input_uncached, sum_input_cached, sum_output,
			last_model, last_provider, summary, summary_finalized)
		SELECT a.session_id, a.first_id, a.last_id, a.st, a.et, a.mc, a.st_tok, a.siu, a.sic, a.so_tok,
		       l.model, l.provider, '', 0
		FROM agg a JOIN logs l ON l.id = a.last_id
	`); err != nil {
		return fmt.Errorf("backfill sessions aggregates: %w", err)
	}

	// 4. summary 回填：逐会话扫 prompt_json，取 lifetime 首条 meaningful（命中即短路）；无则首条非空作暂存。
	sRows, err := tx.Query("SELECT session_id FROM sessions")
	if err != nil {
		return err
	}
	var sessionIDs []string
	for sRows.Next() {
		var sid string
		if err = sRows.Scan(&sid); err == nil {
			sessionIDs = append(sessionIDs, sid)
		}
	}
	sRows.Close()

	for _, sid := range sessionIDs {
		summary, finalized, err := computeSessionSummaryInTx(tx, sid)
		if err != nil {
			return fmt.Errorf("compute summary for %s: %w", sid, err)
		}
		if _, err = tx.Exec("UPDATE sessions SET summary = ?, summary_finalized = ? WHERE session_id = ?",
			summary, finalized, sid); err != nil {
			return fmt.Errorf("update summary for %s: %w", sid, err)
		}
	}

	return tx.Commit()
}

// computeSessionSummaryInTx 是 backfillSessions 的专用 helper（MIGRATION CODE，
// 随 backfillSessions 一并移除）。按会话扫 logs（id 升序）取 lifetime 摘要：
// 首条 meaningful user 消息命中即返回(finalized=1)；否则记首条非空作暂存，扫完返回(finalized=0)。
func computeSessionSummaryInTx(tx *sql.Tx, sessionID string) (string, int, error) {
	rows, err := tx.Query("SELECT COALESCE(prompt_json,'') FROM logs WHERE session_id = ? ORDER BY id ASC", sessionID)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	var tentative string
	haveTentative := false
	for rows.Next() {
		var promptStr string
		if err := rows.Scan(&promptStr); err != nil {
			return "", 0, err
		}
		if promptStr == "" {
			continue
		}
		var msgs []ChatMessage
		if err := json.Unmarshal([]byte(promptStr), &msgs); err != nil {
			continue
		}
		cand := extractPromptSummary(msgs)
		if cand == "" {
			continue
		}
		if !haveTentative {
			tentative = cand
			haveTentative = true
		}
		if isMeaningfulSummary(cand) {
			return cand, 1, nil
		}
	}
	return tentative, 0, nil
}

func (mgr *DBManager) InsertLog(log *UnifiedLog) (int64, error) {
	promptJSON, err := json.Marshal(log.Prompt)
	if err != nil {
		return 0, fmt.Errorf("marshal prompt: %w", err)
	}

	respJSON, err := json.Marshal(log.Response)
	if err != nil {
		return 0, fmt.Errorf("marshal response: %w", err)
	}

	var toolsJSON []byte
	if len(log.Tools) > 0 {
		toolsJSON, err = json.Marshal(log.Tools)
		if err != nil {
			return 0, fmt.Errorf("marshal tools: %w", err)
		}
	}

	query := `
	INSERT INTO logs (
		provider, model, path, status_code, duration_ms,
		input_tokens, output_tokens, total_tokens,
		prompt_json, response_json, tools_json,
		raw_request, raw_response, error_message, created_at, session_id,
		cached_tokens, cache_read_tokens, cache_creation_tokens, parent_id,
		parent_tool_call_id, client_fingerprint, session_seq
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")
	tx, err := mgr.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 预读会话状态以算 session_seq（lifetime 位次 = msg_count+1）与 summary 规则。
	// standalone（空 session_id）必为新会话、位次 1，键待 insert 后用 logID 合成。
	sessionKey := log.SessionID
	isStandalone := strings.TrimSpace(sessionKey) == ""
	isNew := isStandalone
	var preMsgCount int64
	var preSummary string
	var preFinalized int
	sessionSeq := int64(1)
	if !isStandalone {
		e := tx.QueryRow("SELECT msg_count, summary, summary_finalized FROM sessions WHERE session_id = ?", sessionKey).
			Scan(&preMsgCount, &preSummary, &preFinalized)
		if e == sql.ErrNoRows {
			isNew = true
		} else if e != nil {
			return 0, e
		} else {
			sessionSeq = preMsgCount + 1
		}
	}

	res, err := tx.Exec(
		query,
		log.Provider, log.Model, log.Path, log.StatusCode, log.DurationMs,
		log.InputTokens, log.OutputTokens, log.TotalTokens,
		string(promptJSON), string(respJSON), string(toolsJSON),
		log.RawRequest, log.RawResponse, log.ErrorMessage, createdAt, log.SessionID,
		log.CachedTokens, log.CacheReadTokens, log.CacheCreationTokens, log.ParentID,
		log.ParentToolCallID, log.ClientFingerprint, sessionSeq,
	)
	if err != nil {
		return 0, err
	}

	logID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// standalone：合成 'standalone-<logID>' 回写 logs.session_id，使 join 键统一。
	if isStandalone {
		sessionKey = "standalone-" + strconv.FormatInt(logID, 10)
		if _, err := tx.Exec("UPDATE logs SET session_id = ? WHERE id = ?", sessionKey, logID); err != nil {
			return 0, fmt.Errorf("update standalone session_id: %w", err)
		}
	}

	if err := insertConversationHandles(tx, logID, "request", log.RequestHandles); err != nil {
		return 0, err
	}
	if err := insertConversationHandles(tx, logID, "response", log.ResponseHandles); err != nil {
		return 0, err
	}

	// 维护 sessions 实体（ADR-0001）：增量 upsert 聚合 + lifetime 摘要。
	if err := mgr.upsertSession(tx, log, logID, sessionKey, isNew, preSummary, preFinalized, createdAt); err != nil {
		return 0, fmt.Errorf("upsert session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return logID, nil
}

// upsertSession 在 InsertLog 事务内维护 sessions 行。新会话 INSERT；已存在 UPDATE：
// lifetime 字段（msg_count/sum_*/summary）按规则更新；current 字段（last_log_id/end_time/
// last_model/last_provider）每次刷新。token 增量按行 provider 分流（raw 值，与读时 normalize 解耦）。
func (mgr *DBManager) upsertSession(tx *sql.Tx, log *UnifiedLog, logID int64, sessionKey string, isNew bool, preSummary string, preFinalized int, createdAt string) error {
	// 用内存 log.Prompt 现算本轮候选（零 DB 读）。
	newSummary := extractPromptSummary(log.Prompt)
	newMeaningful := isMeaningfulSummary(newSummary)

	var sumInputUncached, sumInputCached int
	if log.Provider == "anthropic" {
		sumInputUncached = log.InputTokens
		sumInputCached = log.CacheReadTokens
	} else {
		sumInputUncached = log.InputTokens - log.CachedTokens
		sumInputCached = log.CachedTokens
	}

	if isNew {
		finalized := 0
		if newMeaningful {
			finalized = 1
		}
		_, err := tx.Exec(`INSERT INTO sessions
			(session_id, first_log_id, last_log_id, start_time, end_time, msg_count,
			 sum_tokens, sum_input_uncached, sum_input_cached, sum_output,
			 last_model, last_provider, summary, summary_finalized)
			VALUES (?, ?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?, ?)`,
			sessionKey, logID, logID, createdAt, createdAt,
			log.TotalTokens, sumInputUncached, sumInputCached, log.OutputTokens,
			log.Model, log.Provider, newSummary, finalized)
		return err
	}

	// lifetime 摘要规则：已冻结则不动；否则首条非空作暂存，命中 meaningful 升级并冻结。
	keepSummary := preSummary
	keepFinalized := preFinalized
	if preFinalized == 0 && newSummary != "" {
		if preSummary == "" {
			keepSummary = newSummary
			if newMeaningful {
				keepFinalized = 1
			}
		} else if newMeaningful {
			keepSummary = newSummary
			keepFinalized = 1
		}
	}

	_, err := tx.Exec(`UPDATE sessions SET
		last_log_id = ?, end_time = ?, msg_count = msg_count + 1,
		sum_tokens = sum_tokens + ?, sum_input_uncached = sum_input_uncached + ?,
		sum_input_cached = sum_input_cached + ?, sum_output = sum_output + ?,
		last_model = ?, last_provider = ?, summary = ?, summary_finalized = ?
		WHERE session_id = ?`,
		logID, createdAt,
		log.TotalTokens, sumInputUncached, sumInputCached, log.OutputTokens,
		log.Model, log.Provider, keepSummary, keepFinalized, sessionKey)
	return err
}

func insertConversationHandles(tx *sql.Tx, logID int64, source string, handles []ConversationHandle) error {
	handles = compactHandles(handles)
	if len(handles) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO log_handles (log_id, source, handle_kind, handle_value)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, handle := range handles {
		if _, err := stmt.Exec(logID, source, handle.Kind, handle.Value); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *DBManager) GetLogs(page, pageSize int, provider, model, keyword string, statusFilter string, excludeBranches bool) ([]LogSummary, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 构造过滤条件
	whereClause := "1=1"
	args := []any{}

	if excludeBranches {
		whereClause += " AND COALESCE(parent_tool_call_id, '') = ''"
	}
	if provider != "" {
		whereClause += " AND provider = ?"
		args = append(args, provider)
	}
	if model != "" {
		whereClause += " AND model LIKE ?"
		args = append(args, "%"+model+"%")
	}
	switch statusFilter {
	case "success":
		whereClause += " AND status_code >= 200 AND status_code < 300"
	case "error":
		whereClause += " AND (status_code < 200 OR status_code >= 300)"
	case "":
		// no filter
	default:
		// 兼容直接传数字状态码（如 "200"）
		if _, err := strconv.Atoi(statusFilter); err == nil {
			whereClause += " AND status_code = ?"
			args = append(args, statusFilter)
		}
	}
	if keyword != "" {
		whereClause += " AND (raw_request LIKE ? OR raw_response LIKE ? OR error_message LIKE ?)"
		likeArg := "%" + keyword + "%"
		args = append(args, likeArg, likeArg, likeArg)
	}

	// 1. 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM logs WHERE %s", whereClause)
	var total int
	err := mgr.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 2. 获取分页数据：session_seq 直读列；会话摘要 JOIN sessions 取（替代窗口查询 N+1）。
	dataQuery := fmt.Sprintf(`
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       l1.session_seq, COALESCE(s.summary, '')
		FROM logs l1
		LEFT JOIN sessions s ON s.session_id = l1.session_id
		WHERE %s
		ORDER BY l1.created_at DESC, l1.id DESC
		LIMIT ? OFFSET ?`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := mgr.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var summaries []LogSummary
	for rows.Next() {
		var s LogSummary
		err := rows.Scan(
			&s.ID, &s.Provider, &s.Model, &s.Path, &s.StatusCode,
			&s.DurationMs, &s.InputTokens, &s.OutputTokens, &s.TotalTokens,
			&s.CachedTokens, &s.CacheReadTokens, &s.CacheCreationTokens,
			&s.CreatedAt, &s.SessionID, &s.ParentID, &s.ParentToolCallID,
			&s.SessionSeq, &s.SessionSummary,
		)
		if err != nil {
			return nil, 0, err
		}
		s.normalizeAnthropicTokens()
		summaries = append(summaries, s)
	}

	return summaries, total, nil
}

// sessionLogsSlimQuery 故意不查 raw_request / raw_response / tools_json：会话列表
// 视图不需要它们，而每个轮次都携带累积历史，全量返回会让 payload 呈 O(N^2) 增长，
// 长会话下达到数百 MB 直接撑爆浏览器渲染进程（Chrome error code 5 / OOM）。
const sessionLogsSlimQuery = `
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       COALESCE(l1.prompt_json, ''), COALESCE(l1.response_json, ''), COALESCE(l1.error_message, ''),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       l1.session_seq
		FROM logs l1
		WHERE l1.session_id = ?
		ORDER BY l1.id ASC
		LIMIT ? OFFSET ?`

// sessionLogsFullQuery 保留原始全量字段，仅在 ?full=1 调试时使用。
const sessionLogsFullQuery = `
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       COALESCE(l1.prompt_json, ''), COALESCE(l1.response_json, ''), COALESCE(l1.tools_json, ''),
		       COALESCE(l1.raw_request, ''), COALESCE(l1.raw_response, ''), COALESCE(l1.error_message, ''),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       l1.session_seq
		FROM logs l1
		WHERE l1.session_id = ?
		ORDER BY l1.id ASC
		LIMIT ? OFFSET ?`

func (mgr *DBManager) GetSessionLogs(sessionID string, full bool, page, pageSize int) ([]UnifiedLog, int, error) {
	if sessionID == "" {
		return nil, 0, nil
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	// 总数 = 该会话现存 log 数（非 sessions.msg_count 的 lifetime 值，分页以现存为准）。
	var total int
	if err := mgr.db.QueryRow("SELECT COUNT(*) FROM logs WHERE session_id = ?", sessionID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := sessionLogsSlimQuery
	if full {
		query = sessionLogsFullQuery
	}

	rows, err := mgr.db.Query(query, sessionID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []UnifiedLog
	for rows.Next() {
		var s UnifiedLog
		var promptStr, respStr string
		var toolsStr, rawReqStr, rawRespStr string

		if full {
			err := rows.Scan(
				&s.ID, &s.Provider, &s.Model, &s.Path, &s.StatusCode, &s.DurationMs,
				&s.InputTokens, &s.OutputTokens, &s.TotalTokens,
				&s.CachedTokens, &s.CacheReadTokens, &s.CacheCreationTokens,
				&promptStr, &respStr, &toolsStr, &rawReqStr, &rawRespStr, &s.ErrorMessage,
				&s.CreatedAt, &s.SessionID, &s.ParentID, &s.ParentToolCallID,
				&s.SessionSeq,
			)
			if err != nil {
				return nil, 0, err
			}
			s.RawRequest = rawReqStr
			s.RawResponse = rawRespStr
		} else {
			err := rows.Scan(
				&s.ID, &s.Provider, &s.Model, &s.Path, &s.StatusCode, &s.DurationMs,
				&s.InputTokens, &s.OutputTokens, &s.TotalTokens,
				&s.CachedTokens, &s.CacheReadTokens, &s.CacheCreationTokens,
				&promptStr, &respStr, &s.ErrorMessage,
				&s.CreatedAt, &s.SessionID, &s.ParentID, &s.ParentToolCallID,
				&s.SessionSeq,
			)
			if err != nil {
				return nil, 0, err
			}
		}

		if promptStr != "" {
			if err := json.Unmarshal([]byte(promptStr), &s.Prompt); err != nil {
				log.Printf("GetSessionLogs: failed to unmarshal prompt_json for log %d: %v", s.ID, err)
			}
		}
		if respStr != "" {
			if err := json.Unmarshal([]byte(respStr), &s.Response); err != nil {
				log.Printf("GetSessionLogs: failed to unmarshal response_json for log %d: %v", s.ID, err)
			}
		}
		if full && toolsStr != "" {
			if err := json.Unmarshal([]byte(toolsStr), &s.Tools); err != nil {
				log.Printf("GetSessionLogs: failed to unmarshal tools_json for log %d: %v", s.ID, err)
			}
		}

		s.normalizeAnthropicTokens()
		if !full {
			trimSessionLogForList(&s)
		}
		logs = append(logs, s)
	}

	return logs, total, nil
}

// GetSessionByID 取单个会话的元数据（用于会话详情视图的总额/轮数等，来自 sessions 表）。
func (mgr *DBManager) GetSessionByID(sessionID string) (*SessionMetadata, error) {
	if sessionID == "" {
		return nil, nil
	}
	var s SessionMetadata
	err := mgr.db.QueryRow(`
		SELECT session_id, datetime(start_time,'localtime'), datetime(end_time,'localtime'),
		       sum_tokens, sum_input_uncached, sum_input_cached, sum_output,
		       msg_count, last_model, last_provider, summary, first_log_id
		FROM sessions WHERE session_id = ?`, sessionID).Scan(
		&s.SessionID, &s.StartTime, &s.EndTime,
		&s.TotalTokens, &s.TotalInputUncachedTokens, &s.TotalInputCachedTokens, &s.TotalOutputTokens,
		&s.MessageCount, &s.Model, &s.Provider, &s.SessionSummary, &s.FirstLogID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// sessionListContentCap 限制 slim 会话日志里单条消息文本的上限（rune 计）。
// 客户端展示时本就截断到 95 字符，此处只保留足够余量即可。
const sessionListContentCap = 4000

// trimSessionLogForList 裁剪会话列表用不到的字段并限制大文本：去掉 raw 请求/响应
// 与工具定义、把 prompt 收缩到最后一条 user/tool 输入消息、把响应文本按上限截断、
// 工具调用只保留名称。会话视图只渲染每轮标量 + 末条输入 + 响应摘要，因此不丢失可见信息，
// 同时避免长会话因每轮累积历史导致的 O(N^2) 体积爆炸。
func trimSessionLogForList(s *UnifiedLog) {
	s.RawRequest = ""
	s.RawResponse = ""
	s.Tools = nil

	// Prompt：只保留最后一条 user/tool 消息（即该轮输入）。
	lastIdx := -1
	for i := len(s.Prompt) - 1; i >= 0; i-- {
		if r := s.Prompt[i].Role; r == "user" || r == "tool" {
			lastIdx = i
			break
		}
	}
	if lastIdx >= 0 {
		m := s.Prompt[lastIdx]
		m.Content = capStringRunes(m.Content, sessionListContentCap)
		m.Thinking = ""
		m.ToolCalls = nil
		s.Prompt = []ChatMessage{m}
	} else {
		s.Prompt = nil
	}

	// Response：截断文本，工具调用只保留名称。
	s.Response.Content = capStringRunes(s.Response.Content, sessionListContentCap)
	s.Response.Thinking = capStringRunes(s.Response.Thinking, sessionListContentCap)
	if len(s.Response.ToolCalls) > 0 {
		tcs := make([]ToolCall, 0, len(s.Response.ToolCalls))
		for _, tc := range s.Response.ToolCalls {
			tcs = append(tcs, ToolCall{Name: tc.Name})
		}
		s.Response.ToolCalls = tcs
	}
}

// capStringRunes 按 rune 安全截断，保证结果不超过 max 个 rune（含结尾省略号），
// 避免切断多字节 UTF-8 序列。
func capStringRunes(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 1 {
		return "…"
	}
	return string(r[:max-1]) + "…"
}


func (mgr *DBManager) GetLogDetail(id int64) (*UnifiedLog, error) {
	query := `
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       COALESCE(l1.prompt_json, ''), COALESCE(l1.response_json, ''), COALESCE(l1.tools_json, ''),
		       COALESCE(l1.raw_request, ''), COALESCE(l1.raw_response, ''), COALESCE(l1.error_message, ''),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       l1.session_seq
		FROM logs l1
		WHERE l1.id = ?
	`

	var logUnified UnifiedLog
	var promptStr, respStr, toolsStr string

	err := mgr.db.QueryRow(query, id).Scan(
		&logUnified.ID, &logUnified.Provider, &logUnified.Model, &logUnified.Path, &logUnified.StatusCode, &logUnified.DurationMs,
		&logUnified.InputTokens, &logUnified.OutputTokens, &logUnified.TotalTokens,
		&logUnified.CachedTokens, &logUnified.CacheReadTokens, &logUnified.CacheCreationTokens,
		&promptStr, &respStr, &toolsStr, &logUnified.RawRequest, &logUnified.RawResponse, &logUnified.ErrorMessage, &logUnified.CreatedAt, &logUnified.SessionID, &logUnified.ParentID, &logUnified.ParentToolCallID,
		&logUnified.SessionSeq,
	)
	if err != nil {
		return nil, err
	}

	if promptStr != "" {
		if err := json.Unmarshal([]byte(promptStr), &logUnified.Prompt); err != nil {
			log.Printf("GetLogDetail: failed to unmarshal prompt_json for log %d: %v", logUnified.ID, err)
		}
	}

	if respStr != "" {
		if err := json.Unmarshal([]byte(respStr), &logUnified.Response); err != nil {
			log.Printf("GetLogDetail: failed to unmarshal response_json for log %d: %v", logUnified.ID, err)
		}
	}

	if toolsStr != "" {
		if err := json.Unmarshal([]byte(toolsStr), &logUnified.Tools); err != nil {
			log.Printf("GetLogDetail: failed to unmarshal tools_json for log %d: %v", logUnified.ID, err)
		}
	}

	logUnified.normalizeAnthropicTokens()

	// 填充详情的 sessionSummary：直读 sessions.summary（替代最早 3 轮 prompt 的 N+1 解析）。
	if logUnified.SessionID != "" {
		_ = mgr.db.QueryRow("SELECT COALESCE(summary,'') FROM sessions WHERE session_id = ?", logUnified.SessionID).Scan(&logUnified.SessionSummary)
	}

	// 级联查询子会话
	subQuery := `
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       COALESCE(l1.prompt_json, ''), COALESCE(l1.response_json, ''), COALESCE(l1.tools_json, ''),
		       COALESCE(l1.raw_request, ''), COALESCE(l1.raw_response, ''), COALESCE(l1.error_message, ''),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       l1.session_seq
		FROM logs l1
		WHERE l1.parent_id = ?
		ORDER BY l1.id ASC
	`
	rows, err := mgr.db.Query(subQuery, logUnified.ID)
	if err == nil {
		for rows.Next() {
			var subLog UnifiedLog
			var subPromptStr, subRespStr, subToolsStr string
			err := rows.Scan(
				&subLog.ID, &subLog.Provider, &subLog.Model, &subLog.Path, &subLog.StatusCode, &subLog.DurationMs,
				&subLog.InputTokens, &subLog.OutputTokens, &subLog.TotalTokens,
				&subLog.CachedTokens, &subLog.CacheReadTokens, &subLog.CacheCreationTokens,
				&subPromptStr, &subRespStr, &subToolsStr, &subLog.RawRequest, &subLog.RawResponse, &subLog.ErrorMessage, &subLog.CreatedAt, &subLog.SessionID, &subLog.ParentID, &subLog.ParentToolCallID,
				&subLog.SessionSeq,
			)
			if err == nil {
				if subPromptStr != "" {
					if err := json.Unmarshal([]byte(subPromptStr), &subLog.Prompt); err != nil {
						log.Printf("GetLogDetail: failed to unmarshal sub-log prompt_json for log %d: %v", subLog.ID, err)
					}
				}
				if subRespStr != "" {
					if err := json.Unmarshal([]byte(subRespStr), &subLog.Response); err != nil {
						log.Printf("GetLogDetail: failed to unmarshal sub-log response_json for log %d: %v", subLog.ID, err)
					}
				}
				if subToolsStr != "" {
					if err := json.Unmarshal([]byte(subToolsStr), &subLog.Tools); err != nil {
						log.Printf("GetLogDetail: failed to unmarshal sub-log tools_json for log %d: %v", subLog.ID, err)
					}
				}
				subLog.normalizeAnthropicTokens()
				logUnified.SubLogs = append(logUnified.SubLogs, subLog)
			}
		}
		rows.Close()
	}

	return &logUnified, nil
}

func isMeaningfulSummary(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	// 去掉标点符号
	s = strings.Trim(s, ".!?！？。，, ")
	if len(s) < 4 { // 太短的通常不是有意义的任务描述
		return false
	}
	greetings := []string{
		"你好", "hello", "hi", "hey", "hallo", "hola",
		"开始", "start", "run", "test", "testing",
		"你好呀", "在吗", "早上好", "中午好", "下午好", "晚上好",
	}
	for _, g := range greetings {
		if s == g {
			return false
		}
	}
	return true
}

func simplifyUserPrompt(content string) string {
	// 针对 Copilot / Claude Code 等框架常见的 XML 标签进行清洗
	// 1. 尝试提取特定标签内的内容
	tags := []string{"user_query", "query", "user_input", "task", "instruction", "input", "userRequest", "user_request"}
	for _, tag := range tags {
		startTag := "<" + tag + ">"
		endTag := "</" + tag + ">"
		if strings.Contains(content, startTag) {
			start := strings.Index(content, startTag) + len(startTag)
			end := strings.Index(content, endTag)
			if end > start {
				return simplifyUserPrompt(content[start:end]) // 递归清洗内部内容
			}
		}
	}

	// 2. 如果没找到配对标签，但内容就是个闭合标签（比如报错截图里的 </userRequest>），尝试去掉它
	// 使用正则或简单的 ReplaceAll
	cleanContent := content
	for _, tag := range tags {
		cleanContent = strings.ReplaceAll(cleanContent, "<"+tag+">", "")
		cleanContent = strings.ReplaceAll(cleanContent, "</"+tag+">", "")
	}
	cleanContent = strings.TrimSpace(cleanContent)
	if cleanContent == "" && content != "" {
		// 如果去完标签变空了，说明原内容可能只有标签或者标签嵌套有问题
		// 继续后续的启发式处理
	} else {
		content = cleanContent
	}

	// 常见的文本标记 (从后往前找最后一个)
	markers := []string{
		"User Query:", "User input:", "Query:", "Question:", "Task:", "Request:",
		"用户提问:", "用户输入:", "问题:", "任务:", "请求:",
		"Input:", "Prompt:", "User:", "用户:",
		"Goal:", "Objective:", "Instruction:",
	}
	contentLower := strings.ToLower(content)
	for _, marker := range markers {
		markerLower := strings.ToLower(marker)
		// 寻找最后一个匹配项，通常前面的可能是历史背景中的标记
		if idx := strings.LastIndex(contentLower, markerLower); idx != -1 {
			res := strings.TrimSpace(content[idx+len(marker):])
			// 如果提取出的内容包含明显的段落分隔，说明后面可能是其他注入内容
			if stopIdx := strings.Index(res, "\n\n"); stopIdx != -1 {
				res = res[:stopIdx]
			}
			// 过滤掉太短或太长（可能是大段代码注入）的提取结果
			if res != "" && len(res) > 2 && len(res) < 500 {
				return res
			}
		}
	}

	// 针对 Claude Code 等工具，如果被包裹在某些特定 XML 结构或分割符中
	if strings.Contains(content, "---") {
		parts := strings.Split(content, "---")
		if len(parts) >= 3 {
			// 如果有三个及以上部分，通常中间的部分（parts[1]）是真正的用户输入
			mid := strings.TrimSpace(parts[1])
			if len(mid) > 2 && len(mid) < 500 {
				return mid
			}
		}
		// 备选：找最后一个看起来像真实输入的部分
		for i := len(parts) - 1; i >= 0; i-- {
			p := strings.TrimSpace(parts[i])
			if len(p) > 5 && len(p) < 500 && !strings.Contains(p, "Context:") && !strings.Contains(p, "Instructions:") {
				return p
			}
		}
	}

	// 启发式：如果包含带有问号的行，这通常是用户真实意图
	lines := strings.Split(strings.TrimSpace(content), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if (strings.HasSuffix(line, "?") || strings.HasSuffix(line, "？")) && len(line) < 200 {
			return line
		}
	}

	if len(lines) > 3 {
		// 从后往前查找第一个非空且不带常见框架前缀的行
		for i := len(lines) - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			// 更加彻底地排除 XML 标签相关行
			if (strings.HasPrefix(line, "</") || strings.HasPrefix(line, "<")) && strings.HasSuffix(line, ">") {
				continue
			}
			// 排除纯标点符号
			if len(strings.Trim(line, ".!?！？。，, ")) < 2 {
				continue
			}

			if len(line) > 2 && len(line) < 300 &&
				!strings.HasPrefix(line, "Context:") &&
				!strings.HasPrefix(line, "History:") &&
				!strings.HasPrefix(line, "System:") &&
				!strings.HasPrefix(line, "Current:") &&
				!strings.HasSuffix(line, ":") { // 避免取到标题行
				return line
			}
		}
	}

	return content
}

func extractPromptSummary(msgs []ChatMessage) string {
	var candidate string
	// 从后往前查找最后一个 user 角色且内容非空、非工具结果的消息
	for i := len(msgs) - 1; i >= 0; i-- {
		m := msgs[i]
		role := strings.ToLower(m.Role)
		if role == "user" && strings.TrimSpace(m.Content) != "" && m.ToolCallID == "" {
			candidate = m.Content
			break
		}
	}
	if candidate == "" {
		// 回退：如果从后往前没找到 user，我们找最后一个非空且非 assistant/system 消息
		for i := len(msgs) - 1; i >= 0; i-- {
			m := msgs[i]
			if strings.TrimSpace(m.Content) != "" && strings.ToLower(m.Role) != "assistant" {
				candidate = m.Content
				break
			}
		}
	}
	if candidate == "" && len(msgs) > 0 {
		// 最终回退：取最后一条消息
		candidate = msgs[len(msgs)-1].Content
	}
	if candidate == "" {
		return ""
	}

	// 针对框架注入进行清理，提取真实用户输入
	candidate = simplifyUserPrompt(candidate)

	// 转换为单行，过滤换行符
	candidate = strings.ReplaceAll(candidate, "\r\n", " ")
	candidate = strings.ReplaceAll(candidate, "\n", " ")
	candidate = strings.ReplaceAll(candidate, "\r", " ")
	candidate = strings.TrimSpace(candidate)

	// 截取长度并加省略号
	runes := []rune(candidate)
	if len(runes) > 100 {
		return string(runes[:100]) + "..."
	}
	return candidate
}

func (mgr *DBManager) GetStats() (*UsageStats, error) {
	stats := &UsageStats{
		CallsByProvider: make(map[string]int),
		TokensByModel:   make(map[string]int),
	}

	// 1. 基础汇总统计
	summaryQuery := `
		SELECT COUNT(*),
		       COALESCE(SUM(input_tokens + CASE WHEN provider = 'anthropic' THEN cache_read_tokens ELSE 0 END), 0),
		       COALESCE(SUM(output_tokens), 0),
		       COALESCE(SUM(total_tokens + CASE WHEN provider = 'anthropic' THEN cache_read_tokens ELSE 0 END), 0),
		       COALESCE(AVG(duration_ms), 0),
		       COALESCE(SUM(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 ELSE 0 END), 0),
		       COALESCE(SUM(CASE WHEN provider = 'anthropic' THEN cache_read_tokens ELSE cached_tokens END), 0)
		FROM logs
	`
	var avgDuration float64
	err := mgr.db.QueryRow(summaryQuery).Scan(
		&stats.TotalCalls, &stats.TotalInputTokens, &stats.TotalOutputTokens,
		&stats.TotalTokens, &avgDuration, &stats.SuccessCalls, &stats.TotalCachedTokens,
	)
	if err != nil {
		return nil, err
	}
	stats.AvgDurationMs = int64(avgDuration)
	if stats.TotalCalls > 0 {
		stats.SuccessRate = float64(stats.SuccessCalls) / float64(stats.TotalCalls)
	} else {
		stats.SuccessRate = 0
	}

	// 2. 厂商统计
	providerRows, err := mgr.db.Query("SELECT provider, COUNT(*) FROM logs GROUP BY provider")
	if err != nil {
		return nil, err
	}
	defer providerRows.Close()
	for providerRows.Next() {
		var provider string
		var count int
		if err := providerRows.Scan(&provider, &count); err == nil {
			stats.CallsByProvider[provider] = count
		}
	}

	// 3. 模型 Token 统计
	modelRows, err := mgr.db.Query("SELECT model, SUM(total_tokens + CASE WHEN provider = 'anthropic' THEN COALESCE(cache_read_tokens, 0) ELSE 0 END) FROM logs GROUP BY model")
	if err != nil {
		return nil, err
	}
	defer modelRows.Close()
	for modelRows.Next() {
		var model string
		var totalTokens int
		if err := modelRows.Scan(&model, &totalTokens); err == nil {
			stats.TokensByModel[model] = totalTokens
		}
	}

	return stats, nil
}

type RecentLogForMatch struct {
	ID           int64
	SessionID    string
	PromptJSON   string
	ResponseJSON string
}


type LogHandleRef struct {
	ID        int64
	SessionID string
}

func (mgr *DBManager) GetRecentLogsForMatch(limit int, provider, clientFingerprint string) ([]RecentLogForMatch, error) {
	whereParts := []string{"status_code >= 200", "status_code < 300"}
	args := []any{}
	if strings.TrimSpace(provider) != "" {
		whereParts = append(whereParts, "provider = ?")
		args = append(args, provider)
	}
	if strings.TrimSpace(clientFingerprint) != "" {
		whereParts = append(whereParts, "COALESCE(client_fingerprint, '') = ?")
		args = append(args, clientFingerprint)
	}

	query := `
		SELECT id, COALESCE(session_id, ''), COALESCE(prompt_json, ''), COALESCE(response_json, '')
		FROM logs
		WHERE ` + strings.Join(whereParts, " AND ") + `
		ORDER BY id DESC
		LIMIT ?
	`
	args = append(args, limit)
	rows, err := mgr.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []RecentLogForMatch
	for rows.Next() {
		var r RecentLogForMatch
		if err := rows.Scan(&r.ID, &r.SessionID, &r.PromptJSON, &r.ResponseJSON); err == nil {
			list = append(list, r)
		}
	}
	return list, nil
}

func (mgr *DBManager) FindLatestLogByHandle(kind, value string) (*LogHandleRef, error) {
	if strings.TrimSpace(kind) == "" || strings.TrimSpace(value) == "" {
		return nil, nil
	}

	var ref LogHandleRef
	err := mgr.db.QueryRow(`
		SELECT l.id, COALESCE(l.session_id, '')
		FROM log_handles h
		JOIN logs l ON l.id = h.log_id
		WHERE h.handle_kind = ? AND h.handle_value = ?
		ORDER BY l.created_at DESC, l.id DESC
		LIMIT 1
	`, kind, value).Scan(&ref.ID, &ref.SessionID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (mgr *DBManager) GetLatestSessionByClientFingerprint(clientFingerprint string) (string, *int64, error) {
	if strings.TrimSpace(clientFingerprint) == "" {
		return "", nil, nil
	}

	var sessionID string
	var logID int64
	err := mgr.db.QueryRow(`
		SELECT COALESCE(session_id, ''), id
		FROM logs
		WHERE COALESCE(client_fingerprint, '') = ? AND COALESCE(session_id, '') != ''
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, clientFingerprint).Scan(&sessionID, &logID)
	if err == sql.ErrNoRows {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, err
	}
	return sessionID, &logID, nil
}

func (mgr *DBManager) GetLatestLogIDBySession(sessionID string) (*int64, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, nil
	}

	var id int64
	err := mgr.db.QueryRow(`
		SELECT id
		FROM logs
		WHERE session_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, sessionID).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (mgr *DBManager) GetSessions(page, pageSize int, keyword string) ([]SessionMetadata, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 关键字匹配 sessions.summary（不再 LIKE 扫 raw_request/raw_response 大字段）。
	whereClause := ""
	args := []any{}
	if keyword != "" {
		whereClause = " WHERE summary LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	// 1. 总数：sessions 行数（索引扫描，替代原 COUNT(DISTINCT ...) 全表扫）。
	var total int
	if err := mgr.db.QueryRow("SELECT COUNT(*) FROM sessions"+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2. 分页：idx_sessions_last_log_id 索引扫描，替代原 GROUP BY + 排序。
	dataQuery := fmt.Sprintf(`
		SELECT session_id, datetime(start_time,'localtime'), datetime(end_time,'localtime'),
		       sum_tokens, sum_input_uncached, sum_input_cached, sum_output,
		       msg_count, last_model, last_provider, summary, first_log_id
		FROM sessions%s
		ORDER BY last_log_id DESC
		LIMIT ? OFFSET ?`, whereClause)
	dataArgs := append(args, pageSize, offset)

	rows, err := mgr.db.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []SessionMetadata
	for rows.Next() {
		var s SessionMetadata
		if err := rows.Scan(
			&s.SessionID, &s.StartTime, &s.EndTime,
			&s.TotalTokens, &s.TotalInputUncachedTokens, &s.TotalInputCachedTokens, &s.TotalOutputTokens,
			&s.MessageCount, &s.Model, &s.Provider, &s.SessionSummary, &s.FirstLogID,
		); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return sessions, total, nil
}

func (mgr *DBManager) DeleteLog(id int64) error {
	tx, err := mgr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 0. 先读该 log 的 session_id（删除后无法再查）；不存在则幂等返回。
	var sessionID string
	switch e := tx.QueryRow("SELECT COALESCE(session_id,'') FROM logs WHERE id = ?", id).Scan(&sessionID); e {
	case nil:
	case sql.ErrNoRows:
		return nil
	default:
		return e
	}

	// 1. 删除 log_handles 中关联的记录
	if _, err = tx.Exec("DELETE FROM log_handles WHERE log_id = ?", id); err != nil {
		return fmt.Errorf("delete log handles: %w", err)
	}

	// 2. 删除 logs 表中的记录
	if _, err = tx.Exec("DELETE FROM logs WHERE id = ?", id); err != nil {
		return fmt.Errorf("delete log: %w", err)
	}

	// 3. 维护 sessions 实体（ADR-0001）：lifetime 字段（msg_count/sum_*/summary）不动，
	//    只重算结构字段（first/last/times/last_model/provider）；该会话 0 残留则销毁实体。
	if sessionID != "" {
		var remaining int
		if err = tx.QueryRow("SELECT COUNT(*) FROM logs WHERE session_id = ?", sessionID).Scan(&remaining); err != nil {
			return err
		}
		if remaining == 0 {
			if _, err = tx.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID); err != nil {
				return fmt.Errorf("delete empty session: %w", err)
			}
		} else {
			if _, err = tx.Exec(`
				UPDATE sessions SET
					first_log_id  = (SELECT MIN(id) FROM logs WHERE session_id = sessions.session_id),
					last_log_id   = (SELECT MAX(id) FROM logs WHERE session_id = sessions.session_id),
					start_time    = (SELECT MIN(created_at) FROM logs WHERE session_id = sessions.session_id),
					end_time      = (SELECT MAX(created_at) FROM logs WHERE session_id = sessions.session_id),
					last_model    = (SELECT model FROM logs WHERE id = (SELECT MAX(id) FROM logs WHERE session_id = sessions.session_id)),
					last_provider = (SELECT provider FROM logs WHERE id = (SELECT MAX(id) FROM logs WHERE session_id = sessions.session_id))
				WHERE session_id = ?`, sessionID); err != nil {
				return fmt.Errorf("recompute session structural fields: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (mgr *DBManager) DeleteSessionLogs(sessionID string) error {
	if sessionID == "" {
		return nil
	}

	tx, err := mgr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除 log_handles 中所有属于该会话 logs 的记录
	_, err = tx.Exec(`
		DELETE FROM log_handles
		WHERE log_id IN (SELECT id FROM logs WHERE session_id = ?)
	`, sessionID)
	if err != nil {
		return fmt.Errorf("delete session log handles: %w", err)
	}

	// 2. 删除 logs 表中该会话的记录
	_, err = tx.Exec("DELETE FROM logs WHERE session_id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("delete session logs: %w", err)
	}

	// 3. 销毁 sessions 实体（ADR-0001：整组删尽，lifetime 账随实体消失）
	if _, err = tx.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return tx.Commit()
}
