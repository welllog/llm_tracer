package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	SessionID      string `json:"sessionId"`
	SessionSummary string `json:"sessionSummary"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	TotalTokens    int    `json:"totalTokens"`
	MessageCount   int    `json:"messageCount"`
	Model          string `json:"model"`
	Provider       string `json:"provider"`
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
		parent_tool_call_id TEXT
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
	return err
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
		parent_tool_call_id, client_fingerprint
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")
	tx, err := mgr.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		query,
		log.Provider, log.Model, log.Path, log.StatusCode, log.DurationMs,
		log.InputTokens, log.OutputTokens, log.TotalTokens,
		string(promptJSON), string(respJSON), string(toolsJSON),
		log.RawRequest, log.RawResponse, log.ErrorMessage, createdAt, log.SessionID,
		log.CachedTokens, log.CacheReadTokens, log.CacheCreationTokens, log.ParentID,
		log.ParentToolCallID, log.ClientFingerprint,
	)
	if err != nil {
		return 0, err
	}

	logID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := insertConversationHandles(tx, logID, "request", log.RequestHandles); err != nil {
		return 0, err
	}
	if err := insertConversationHandles(tx, logID, "response", log.ResponseHandles); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return logID, nil
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

func (mgr *DBManager) GetLogs(page, pageSize int, provider, model, keyword string, statusFilter *int, excludeBranches bool) ([]LogSummary, int, error) {
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
	if statusFilter != nil {
		whereClause += " AND status_code = ?"
		args = append(args, *statusFilter)
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

	// 2. 获取分页数据
	dataQuery := fmt.Sprintf(`
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       CASE WHEN COALESCE(l1.session_id, '') != '' THEN
		            (SELECT COUNT(*) FROM logs l2 WHERE l2.session_id = l1.session_id AND l2.id <= l1.id)
		       ELSE 0 END as session_seq
		FROM logs l1
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
			&s.SessionSeq,
		)
		if err != nil {
			return nil, 0, err
		}
		s.normalizeAnthropicTokens()
		summaries = append(summaries, s)
	}

	// 提取当前页的所有 SessionID 并在 Go 层面进行去重
	sessionIDs := []string{}
	seen := make(map[string]bool)
	for _, s := range summaries {
		if s.SessionID != "" && !seen[s.SessionID] {
			seen[s.SessionID] = true
			sessionIDs = append(sessionIDs, s.SessionID)
		}
	}

	if len(sessionIDs) > 0 {
		placeholders := make([]string, len(sessionIDs))
		queryArgs := make([]any, len(sessionIDs))
		for i, id := range sessionIDs {
			placeholders[i] = "?"
			queryArgs[i] = id
		}

		// 使用窗口函数，查出每个会话最早的 3 条日志的 prompt_json
		query := fmt.Sprintf(`
			SELECT session_id, prompt_json
			FROM (
				SELECT COALESCE(session_id, '') as session_id, COALESCE(prompt_json, '') as prompt_json,
				       ROW_NUMBER() OVER (PARTITION BY session_id ORDER BY id ASC) as rn
				FROM logs
				WHERE session_id IN (%s)
			)
			WHERE rn <= 3
		`, strings.Join(placeholders, ","))

		subRows, err := mgr.db.Query(query, queryArgs...)
		if err == nil {
			summariesMap := make(map[string]string)
			meaningfulMap := make(map[string]bool)
			for subRows.Next() {
				var sID, promptStr string
				if err := subRows.Scan(&sID, &promptStr); err == nil && promptStr != "" {
					if meaningfulMap[sID] {
						continue
					}
					var msgs []ChatMessage
					if err := json.Unmarshal([]byte(promptStr), &msgs); err == nil {
						summary := extractPromptSummary(msgs)
						if summary != "" {
							if isMeaningfulSummary(summary) {
								summariesMap[sID] = summary
								meaningfulMap[sID] = true
							} else if summariesMap[sID] == "" {
								summariesMap[sID] = summary
							}
						}
					}
				}
			}
			subRows.Close()

			for i := range summaries {
				if sum, ok := summariesMap[summaries[i].SessionID]; ok {
					summaries[i].SessionSummary = sum
				}
			}
		}
	}

	return summaries, total, nil
}

func (mgr *DBManager) GetSessionLogs(sessionID string) ([]UnifiedLog, error) {
	if sessionID == "" {
		return nil, nil
	}

	query := `
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       COALESCE(l1.prompt_json, ''), COALESCE(l1.response_json, ''), COALESCE(l1.tools_json, ''),
		       COALESCE(l1.raw_request, ''), COALESCE(l1.raw_response, ''), COALESCE(l1.error_message, ''),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       (SELECT COUNT(*) FROM logs l2 WHERE l2.session_id = l1.session_id AND l2.id <= l1.id) as session_seq
		FROM logs l1
		WHERE l1.session_id = ?
		ORDER BY l1.id ASC`

	rows, err := mgr.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []UnifiedLog
	for rows.Next() {
		var s UnifiedLog
		var promptStr, respStr, toolsStr string
		err := rows.Scan(
			&s.ID, &s.Provider, &s.Model, &s.Path, &s.StatusCode, &s.DurationMs,
			&s.InputTokens, &s.OutputTokens, &s.TotalTokens,
			&s.CachedTokens, &s.CacheReadTokens, &s.CacheCreationTokens,
			&promptStr, &respStr, &toolsStr, &s.RawRequest, &s.RawResponse, &s.ErrorMessage,
			&s.CreatedAt, &s.SessionID, &s.ParentID, &s.ParentToolCallID,
			&s.SessionSeq,
		)
		if err != nil {
			return nil, err
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
		if toolsStr != "" {
			if err := json.Unmarshal([]byte(toolsStr), &s.Tools); err != nil {
				log.Printf("GetSessionLogs: failed to unmarshal tools_json for log %d: %v", s.ID, err)
			}
		}

		s.normalizeAnthropicTokens()
		logs = append(logs, s)
	}

	return logs, nil
}


func (mgr *DBManager) GetLogDetail(id int64) (*UnifiedLog, error) {
	query := `
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       COALESCE(l1.prompt_json, ''), COALESCE(l1.response_json, ''), COALESCE(l1.tools_json, ''),
		       COALESCE(l1.raw_request, ''), COALESCE(l1.raw_response, ''), COALESCE(l1.error_message, ''),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       CASE WHEN COALESCE(l1.session_id, '') != '' THEN
		            (SELECT COUNT(*) FROM logs l2 WHERE l2.session_id = l1.session_id AND l2.id <= l1.id)
		       ELSE 0 END as session_seq
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

	// 填充详情的 sessionSummary
	if logUnified.SessionID != "" {
		subRows, err := mgr.db.Query(`
			SELECT COALESCE(prompt_json, '')
			FROM logs
			WHERE session_id = ?
			ORDER BY id ASC
			LIMIT 3
		`, logUnified.SessionID)
		if err == nil {
			for subRows.Next() {
				var promptStr string
				if subRows.Scan(&promptStr) == nil && promptStr != "" {
					var msgs []ChatMessage
					if err := json.Unmarshal([]byte(promptStr), &msgs); err == nil {
						summary := extractPromptSummary(msgs)
						if summary != "" {
							logUnified.SessionSummary = summary
							if isMeaningfulSummary(summary) {
								break
							}
						}
					}
				}
			}
			subRows.Close()
		}
	}

	// 级联查询子会话
	subQuery := `
		SELECT l1.id, l1.provider, l1.model, l1.path, l1.status_code, l1.duration_ms, l1.input_tokens, l1.output_tokens, l1.total_tokens,
		       COALESCE(l1.cached_tokens, 0), COALESCE(l1.cache_read_tokens, 0), COALESCE(l1.cache_creation_tokens, 0),
		       COALESCE(l1.prompt_json, ''), COALESCE(l1.response_json, ''), COALESCE(l1.tools_json, ''),
		       COALESCE(l1.raw_request, ''), COALESCE(l1.raw_response, ''), COALESCE(l1.error_message, ''),
		       datetime(l1.created_at, 'localtime'), COALESCE(l1.session_id, ''), l1.parent_id, COALESCE(l1.parent_tool_call_id, ''),
		       CASE WHEN COALESCE(l1.session_id, '') != '' THEN
		            (SELECT COUNT(*) FROM logs l2 WHERE l2.session_id = l1.session_id AND l2.id <= l1.id)
		       ELSE 0 END as session_seq
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

	// 1. 获取总数
	countQuery := `
		SELECT COUNT(DISTINCT CASE WHEN session_id IS NOT NULL AND session_id != '' THEN session_id ELSE 'standalone-' || id END)
		FROM logs
		WHERE 1=1`
	args := []any{}
	if keyword != "" {
		countQuery += " AND (raw_request LIKE ? OR raw_response LIKE ? OR error_message LIKE ?)"
		likeArg := "%" + keyword + "%"
		args = append(args, likeArg, likeArg, likeArg)
	}

	var total int
	err := mgr.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 2. 分页获取会话信息
	dataQuery := fmt.Sprintf(`
		WITH SessionGroups AS (
			SELECT
				CASE WHEN session_id IS NOT NULL AND session_id != '' THEN session_id ELSE 'standalone-' || id END as effective_session_id,
				MIN(id) as first_id,
				MAX(id) as last_id,
				datetime(MIN(created_at), 'localtime') as start_time,
				datetime(MAX(created_at), 'localtime') as end_time,
				SUM(total_tokens) as sum_tokens,
				COUNT(*) as msg_count,
				MAX(model) as last_model,
				MAX(provider) as last_provider
			FROM logs
			WHERE 1=1 %s
			GROUP BY effective_session_id
		)
		SELECT effective_session_id, first_id, last_id, start_time, end_time, sum_tokens, msg_count, last_model, last_provider
		FROM SessionGroups
		ORDER BY last_id DESC
		LIMIT ? OFFSET ?`, func() string {
		if keyword != "" {
			return " AND (raw_request LIKE ? OR raw_response LIKE ? OR error_message LIKE ?)"
		}
		return ""
	}())

	dataArgs := append(args, pageSize, offset)
	rows, err := mgr.db.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	type rawSession struct {
		SessionMetadata
		firstID int64
	}
	var raws []rawSession
	for rows.Next() {
		var s SessionMetadata
		var firstID, discardLastID int64
		err := rows.Scan(
			&s.SessionID, &firstID, &discardLastID, &s.StartTime, &s.EndTime,
			&s.TotalTokens, &s.MessageCount, &s.Model, &s.Provider,
		)
		if err != nil {
			return nil, 0, err
		}
		raws = append(raws, rawSession{SessionMetadata: s, firstID: firstID})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	rows.Close()

	// 必须在 rows 关闭后再查摘要，否则 SetMaxOpenConns(1) 会死锁
	for i := range raws {
		s := &raws[i]
		if strings.HasPrefix(s.SessionID, "standalone-") {
			var promptStr string
			err = mgr.db.QueryRow("SELECT prompt_json FROM logs WHERE id = ?", s.firstID).Scan(&promptStr)
			if err == nil && promptStr != "" {
				var msgs []ChatMessage
				if json.Unmarshal([]byte(promptStr), &msgs) == nil {
					s.SessionSummary = extractPromptSummary(msgs)
				}
			}
		} else {
			subRows, err := mgr.db.Query("SELECT prompt_json FROM logs WHERE session_id = ? ORDER BY id ASC LIMIT 3", s.SessionID)
			if err == nil {
				for subRows.Next() {
					var promptStr string
					if subRows.Scan(&promptStr) == nil && promptStr != "" {
						var msgs []ChatMessage
						if err := json.Unmarshal([]byte(promptStr), &msgs); err == nil {
							summary := extractPromptSummary(msgs)
							if summary != "" {
								s.SessionSummary = summary
								if isMeaningfulSummary(summary) {
									break
								}
							}
						}
					}
				}
				subRows.Close()
			}
		}
	}

	sessions := make([]SessionMetadata, len(raws))
	for i, r := range raws {
		sessions[i] = r.SessionMetadata
	}
	return sessions, total, nil
}

func (mgr *DBManager) DeleteLog(id int64) error {
	tx, err := mgr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除 log_handles 中关联的记录
	_, err = tx.Exec("DELETE FROM log_handles WHERE log_id = ?", id)
	if err != nil {
		return fmt.Errorf("delete log handles: %w", err)
	}

	// 2. 删除 logs 表中的记录
	_, err = tx.Exec("DELETE FROM logs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete log: %w", err)
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

	return tx.Commit()
}
