package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const piAgentSystemPromptFixture = `You are an expert coding assistant operating inside pi, a coding agent harness. You help users by reading files, executing commands, editing code, and writing new files.

Available tools:
- read: Read file contents
- bash: Execute bash commands (ls, grep, find, etc.)
- edit: Make precise file edits with exact text replacement, including multiple disjoint edits in one call
- write: Create or overwrite files
- web_search: Use for web research questions.`

const piAgentInitialChatRequestFixture = `{
	"model": "mimo-v2.5-pro",
	"messages": [
		{
			"role": "system",
			"content": "You are an expert coding assistant operating inside pi, a coding agent harness. You help users by reading files, executing commands, editing code, and writing new files.\n\nAvailable tools:\n- read: Read file contents\n- bash: Execute bash commands (ls, grep, find, etc.)\n- edit: Make precise file edits with exact text replacement, including multiple disjoint edits in one call\n- write: Create or overwrite files\n- web_search: Use for web research questions."
		},
		{
			"role": "user",
			"content": [
				{"type": "text", "text": "今天成都天气"}
			]
		}
	]
}`

const piAgentToolContinuationRequestFixture = `{
	"model": "mimo-v2.5-pro",
	"messages": [
		{
			"role": "system",
			"content": "You are an expert coding assistant operating inside pi, a coding agent harness. You help users by reading files, executing commands, editing code, and writing new files.\n\nAvailable tools:\n- read: Read file contents\n- bash: Execute bash commands (ls, grep, find, etc.)\n- edit: Make precise file edits with exact text replacement, including multiple disjoint edits in one call\n- write: Create or overwrite files\n- web_search: Use for web research questions."
		},
		{
			"role": "user",
			"content": [
				{"type": "text", "text": "今天成都天气"}
			]
		},
		{
			"role": "assistant",
			"content": null,
			"tool_calls": [
				{
					"id": "call_weather_1",
					"type": "function",
					"function": {
						"name": "search_weather",
						"arguments": "{\"city\":\"成都\"}"
					}
				},
				{
					"id": "call_weather_2",
					"type": "function",
					"function": {
						"name": "search_weather",
						"arguments": "{\"city\":\"成都\",\"source\":2}"
					}
				}
			]
		},
		{
			"role": "tool",
			"tool_call_id": "call_weather_1",
			"content": "2026年6月10日，成都天气以阴转多云为主，最高约34℃，最低约25℃。"
		},
		{
			"role": "tool",
			"tool_call_id": "call_weather_2",
			"content": "2026年6月10日，成都天气为阴，气温25℃到34℃，东风2级，空气质量优。"
		}
	]
}`

const piAgentInitialChatStreamFixture = `data: {"id":"8b3195d6ab0340bab53cf231dcd6c583","choices":[{"delta":{"content":"","role":"assistant","tool_calls":null,"reasoning_content":null},"finish_reason":null,"index":0}],"created":1781053362,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"8b3195d6ab0340bab53cf231dcd6c583","choices":[{"delta":{"content":null,"role":null,"tool_calls":[{"index":0,"id":"call_weather_1","function":{"arguments":"","name":"search_weather"},"type":"function"}],"reasoning_content":null},"finish_reason":null,"index":0}],"created":1781053369,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"8b3195d6ab0340bab53cf231dcd6c583","choices":[{"delta":{"content":null,"role":null,"tool_calls":[{"index":0,"id":null,"function":{"arguments":"{\"city\":\"成都\"}","name":null},"type":"function"}],"reasoning_content":null},"finish_reason":null,"index":0}],"created":1781053369,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"8b3195d6ab0340bab53cf231dcd6c583","choices":[{"delta":{"content":null,"role":null,"tool_calls":[{"index":1,"id":"call_weather_2","function":{"arguments":"","name":"search_weather"},"type":"function"}],"reasoning_content":null},"finish_reason":null,"index":0}],"created":1781053369,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"8b3195d6ab0340bab53cf231dcd6c583","choices":[{"delta":{"content":null,"role":null,"tool_calls":[{"index":1,"id":null,"function":{"arguments":"{\"city\":\"成都\",\"source\":2}","name":null},"type":"function"}],"reasoning_content":null},"finish_reason":null,"index":0}],"created":1781053369,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"8b3195d6ab0340bab53cf231dcd6c583","choices":[{"delta":{"content":null,"role":null,"tool_calls":null,"reasoning_content":null},"finish_reason":"tool_calls","index":0}],"created":1781053369,"model":"mimo-v2.5-pro","object":"chat.completion.chunk","usage":null}

data: {"id":"8b3195d6ab0340bab53cf231dcd6c583","choices":[],"created":1781053369,"model":"mimo-v2.5-pro","object":"chat.completion.chunk","usage":{"completion_tokens":250,"prompt_tokens":3728,"total_tokens":3978,"completion_tokens_details":{"reasoning_tokens":174},"prompt_tokens_details":{}}}

data: [DONE]
`

const piAgentToolContinuationStreamFixture = `data: {"id":"69d95a32d5a74da5a06ea33ecf297b7b","choices":[{"delta":{"content":"","role":"assistant","tool_calls":null,"reasoning_content":null},"finish_reason":null,"index":0}],"created":1781053373,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"69d95a32d5a74da5a06ea33ecf297b7b","choices":[{"delta":{"content":"2026年6月10日，成都天气为阴到多云，气温 25 到 34 度。","role":null,"tool_calls":null,"reasoning_content":null},"finish_reason":null,"index":0}],"created":1781053378,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"69d95a32d5a74da5a06ea33ecf297b7b","choices":[{"delta":{"content":null,"role":null,"tool_calls":null,"reasoning_content":null},"finish_reason":"stop","index":0}],"created":1781053378,"model":"mimo-v2.5-pro","object":"chat.completion.chunk"}

data: {"id":"69d95a32d5a74da5a06ea33ecf297b7b","choices":[],"created":1781053378,"model":"mimo-v2.5-pro","object":"chat.completion.chunk","usage":{"completion_tokens":196,"prompt_tokens":4940,"total_tokens":5136,"completion_tokens_details":{"reasoning_tokens":75},"prompt_tokens_details":{"cached_tokens":3712}}}

data: [DONE]
`

func TestProxyAndLogging(t *testing.T) {
	// 1. 创建临时数据库目录
	tempDir, err := os.MkdirTemp("", "llm_tracer_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")

	// 2. 初始化数据库
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	// 3. 启动 Mock AI 上游服务器
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查鉴权头
		auth := r.Header.Get("Authorization")
		apiKey := r.Header.Get("X-Api-Key")

		if !strings.HasPrefix(auth, "Bearer ") && apiKey == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		var bodyMap map[string]any
		_ = json.Unmarshal(bodyBytes, &bodyMap)

		isStream, _ := bodyMap["stream"].(bool)

		if strings.HasSuffix(r.URL.Path, "/chat/completions") {
			w.Header().Set("Content-Type", "application/json")
			if isStream {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")
				w.WriteHeader(http.StatusOK)

				// 模拟发送两个 SSE 块加一个 [DONE]
				_, _ = w.Write([]byte("data: {\"id\":\"chatcmpl-123\",\"object\":\"chat.completion.chunk\",\"created\":1677652288,\"model\":\"gpt-4\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"Hello\"},\"finish_reason\":null}]}\n\n"))
				time.Sleep(10 * time.Millisecond)
				_, _ = w.Write([]byte("data: {\"id\":\"chatcmpl-123\",\"object\":\"chat.completion.chunk\",\"created\":1677652288,\"model\":\"gpt-4\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\" world!\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":5,\"total_tokens\":15}}\n\n"))
				time.Sleep(10 * time.Millisecond)
				_, _ = w.Write([]byte("data: [DONE]\n\n"))
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"id": "chatcmpl-123",
					"object": "chat.completion",
					"created": 1677652288,
					"model": "gpt-4",
					"choices": [{
						"index": 0,
						"message": {
							"role": "assistant",
							"content": "Hello world!"
						},
						"finish_reason": "stop"
					}],
					"usage": {
						"prompt_tokens": 9,
						"completion_tokens": 12,
						"total_tokens": 21
					}
				}`))
			}
		} else if strings.HasSuffix(r.URL.Path, "/messages") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "msg_013Zva5t9ec2ndvT23yKzxxa",
				"type": "message",
				"role": "assistant",
				"model": "claude-3-opus-20240229",
				"content": [
					{
						"type": "text",
						"text": "Hello, Claude response!"
					}
				],
				"usage": {
					"input_tokens": 15,
					"output_tokens": 25
				}
			}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockUpstream.Close()

	// 4. 创建并初始化 ProxyServer
	proxySrv := &ProxyServer{
		client: http.DefaultClient,
		db:     dbMgr,
		config: AppConfig{
			ListenAddr:             "127.0.0.1:0", // 随机端口
			OpenaiBaseURL:          mockUpstream.URL,
			OpenaiAPIKey:           "test-openai-key",
			AnthropicBaseURL:       mockUpstream.URL,
			AnthropicAPIKey:        "test-anthropic-key",
			DBPath:                 dbPath,
		},
		lastActiveSessionMap: make(map[string]string),
	}

	// 5. 模拟一个 HTTP Request 进行测试 (非流式 OpenAI)
	reqJSON := `{
		"model": "gpt-4",
		"messages": [
			{"role": "system", "content": "You are helpful"},
			{"role": "user", "content": "Hi"}
		],
		"metadata": {"client_id": "test-openai-client"}
	}`

	wRec := httptest.NewRecorder()
	httpReq := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(reqJSON)))
	httpReq.Header.Set("Content-Type", "application/json")

	proxySrv.handleProxyOpenAI(wRec, httpReq)

	// 验证代理返回
	resp := wRec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBytes, _ := io.ReadAll(resp.Body)
	var respMap map[string]any
	if err := json.Unmarshal(respBytes, &respMap); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if respMap["id"] != "chatcmpl-123" {
		t.Errorf("Expected id chatcmpl-123, got %v", respMap["id"])
	}

	// 验证数据库日志记录
	logs, total, err := dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if err != nil {
		t.Fatalf("Failed to fetch logs from DB: %v", err)
	}
	if total != 1 {
		t.Fatalf("Expected 1 log in database, got %d", total)
	}

	logSummary := logs[0]
	if logSummary.Provider != "openai" || logSummary.Model != "gpt-4" {
		t.Errorf("Unexpected log summary: %+v", logSummary)
	}
	if logSummary.TotalTokens != 21 {
		t.Errorf("Expected 21 tokens, got %d", logSummary.TotalTokens)
	}

	// 检查详情
	detail, err := dbMgr.GetLogDetail(logSummary.ID)
	if err != nil {
		t.Fatalf("Failed to fetch log detail: %v", err)
	}
	if detail.Response.Content != "Hello world!" {
		t.Errorf("Expected response content 'Hello world!', got %q", detail.Response.Content)
	}
	if len(detail.Prompt) != 2 || detail.Prompt[0].Role != "system" {
		t.Errorf("Unexpected prompt parsed: %+v", detail.Prompt)
	}

	// 6. 测试流式 OpenAI 代理的拦截记录
	wRecStream := httptest.NewRecorder()
	reqJSONStream := `{
		"model": "gpt-4",
		"messages": [
			{"role": "user", "content": "Stream request"}
		],
		"metadata": {"client_id": "test-openai-client"},
		"stream": true
	}`
	httpReqStream := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(reqJSONStream)))
	httpReqStream.Header.Set("Content-Type", "application/json")

	proxySrv.handleProxyOpenAI(wRecStream, httpReqStream)

	respStream := wRecStream.Result()
	if respStream.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for stream, got %d", respStream.StatusCode)
	}

	// 验证是否有 2 条日志在库中
	_, total, _ = dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if total != 2 {
		t.Fatalf("Expected 2 logs in DB after stream, got %d", total)
	}

	// 验证最后一条流式日志
	allLogs, _, _ := dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	latestLogSummary := allLogs[0] // 倒序
	latestLog, err := dbMgr.GetLogDetail(latestLogSummary.ID)
	if err != nil {
		t.Fatalf("Failed to fetch stream detail: %v", err)
	}

	if latestLog.Response.Content != "Hello world!" {
		t.Errorf("Expected aggregated stream content 'Hello world!', got %q", latestLog.Response.Content)
	}
	if latestLog.TotalTokens != 15 {
		t.Errorf("Expected 15 tokens aggregated from stream, got %d", latestLog.TotalTokens)
	}
}

func TestStatsAggregations(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "llm_tracer_test_stats_*")
	defer os.RemoveAll(tempDir)
	dbPath := filepath.Join(tempDir, "test.db")

	dbMgr, _ := InitDB(dbPath)
	defer dbMgr.Close()

	// 插入几条模拟日志
	_, _ = dbMgr.InsertLog(&UnifiedLog{
		Provider:     "openai",
		Model:        "gpt-4",
		Path:         "/v1/chat/completions",
		StatusCode:   200,
		DurationMs:   1000,
		InputTokens:  10,
		OutputTokens: 20,
		TotalTokens:  30,
		Prompt:       []ChatMessage{{Role: "user", Content: "Hi"}},
		Response:     ChatMessage{Role: "assistant", Content: "Hello"},
	})

	_, _ = dbMgr.InsertLog(&UnifiedLog{
		Provider:     "anthropic",
		Model:        "claude-3-opus",
		Path:         "/v1/messages",
		StatusCode:   200,
		DurationMs:   2000,
		InputTokens:  50,
		OutputTokens: 50,
		TotalTokens:  100,
		Prompt:       []ChatMessage{{Role: "user", Content: "Hello"}},
		Response:     ChatMessage{Role: "assistant", Content: "Hi there"},
	})

	stats, err := dbMgr.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalCalls != 2 {
		t.Errorf("Expected 2 calls, got %d", stats.TotalCalls)
	}
	if stats.TotalTokens != 130 {
		t.Errorf("Expected 130 total tokens, got %d", stats.TotalTokens)
	}
	if stats.AvgDurationMs != 1500 {
		t.Errorf("Expected 1500ms avg duration, got %d", stats.AvgDurationMs)
	}
	if stats.CallsByProvider["openai"] != 1 || stats.CallsByProvider["anthropic"] != 1 {
		t.Errorf("Unexpected calls by provider: %+v", stats.CallsByProvider)
	}
}

// 模拟 DB sql.NullString 测试
func TestDBNullFields(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "llm_tracer_test_nulls_*")
	defer os.RemoveAll(tempDir)
	dbPath := filepath.Join(tempDir, "test_nulls.db")

	db, _ := sql.Open("sqlite", dbPath)
	defer db.Close()

	mgr := &DBManager{db: db}
	if err := mgr.createTables(); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// 插入一条有 NULL 的数据以测试 DB 能否正常 Scan 且不崩溃
	_, err := db.Exec(`
		INSERT INTO logs (provider, model, path, status_code, duration_ms, raw_request, raw_response, error_message)
		VALUES ('openai', 'gpt-4', '/v1/chat/completions', 502, 500, 'req', NULL, 'Bad Gateway')
	`)
	if err != nil {
		t.Fatalf("Failed to insert log: %v", err)
	}

	detail, err := mgr.GetLogDetail(1)
	if err != nil {
		t.Fatalf("GetLogDetail failed: %v", err)
	}

	if detail.ErrorMessage != "Bad Gateway" {
		t.Errorf("Expected 'Bad Gateway', got %q", detail.ErrorMessage)
	}
	if detail.RawResponse != "" {
		t.Errorf("Expected empty string for NULL field, got %q", detail.RawResponse)
	}
}

func TestJoinUpstreamURL(t *testing.T) {
	tests := []struct {
		baseURL    string
		pathSuffix string
		expected   string
	}{
		{
			baseURL:    "https://api.openai.com",
			pathSuffix: "/v1/chat/completions",
			expected:   "https://api.openai.com/v1/chat/completions",
		},
		{
			baseURL:    "https://api.openai.com/v1",
			pathSuffix: "/v1/chat/completions",
			expected:   "https://api.openai.com/v1/chat/completions",
		},
		{
			baseURL:    "https://api.openai.com/v1/",
			pathSuffix: "/v1/chat/completions",
			expected:   "https://api.openai.com/v1/chat/completions",
		},
		{
			baseURL:    "https://my-proxy.com/custom/v1",
			pathSuffix: "/v1/chat/completions",
			expected:   "https://my-proxy.com/custom/v1/chat/completions",
		},
		{
			baseURL:    "https://api.anthropic.com",
			pathSuffix: "/v1/messages",
			expected:   "https://api.anthropic.com/v1/messages",
		},
		{
			baseURL:    "https://api.anthropic.com/v1",
			pathSuffix: "/v1/messages",
			expected:   "https://api.anthropic.com/v1/messages",
		},
	}

	for _, tc := range tests {
		actual := joinUpstreamURL(tc.baseURL, tc.pathSuffix)
		if actual != tc.expected {
			t.Errorf("joinUpstreamURL(%q, %q) = %q; want %q", tc.baseURL, tc.pathSuffix, actual, tc.expected)
		}
	}
}

func TestParseAnthropicStreamChunk(t *testing.T) {
	var usage TokenUsage
	var msg ChatMessage
	blockIndex := 0

	// 1. 模拟 message_start，带有 input_tokens
	line1 := []byte(`data: {"type":"message_start","message":{"id":"msg_123","type":"message","role":"assistant","content":[],"usage":{"input_tokens":100,"output_tokens":0}}}`)
	err := ParseAnthropicStreamChunk(line1, &msg, &usage, &blockIndex)
	if err != nil {
		t.Fatalf("failed to parse chunk 1: %v", err)
	}
	if usage.PromptTokens != 100 {
		t.Errorf("expected 100 prompt tokens, got %d", usage.PromptTokens)
	}

	// 2. 模拟 message_delta，带有 output_tokens 和缓存 token
	line2 := []byte(`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"input_tokens":100,"output_tokens":41,"cache_read_input_tokens":2048}}`)
	err = ParseAnthropicStreamChunk(line2, &msg, &usage, &blockIndex)
	if err != nil {
		t.Fatalf("failed to parse chunk 2: %v", err)
	}
	if usage.CompletionTokens != 41 {
		t.Errorf("expected 41 completion tokens, got %d", usage.CompletionTokens)
	}
	if usage.PromptTokens != 100 {
		t.Errorf("expected prompt tokens to remain 100, got %d", usage.PromptTokens)
	}
	if usage.CacheReadTokens != 2048 {
		t.Errorf("expected 2048 cache read tokens, got %d", usage.CacheReadTokens)
	}
}

func TestConversationChaining(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "llm_tracer_chain_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "chain_test.db")
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	// 模拟 AI 上游服务器，它会记录请求的历史并回复
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		bodyBytes, _ := io.ReadAll(r.Body)
		var bodyMap map[string]any
		_ = json.Unmarshal(bodyBytes, &bodyMap)
		messages, _ := bodyMap["messages"].([]any)

		// 模拟回复
		content := "Mocked response content"
		var toolCallsStr string
		if len(messages) > 0 {
			lastMsg, ok := messages[len(messages)-1].(map[string]any)
			if ok {
				msgContent, _ := lastMsg["content"].(string)
				if strings.Contains(strings.ToLower(msgContent), "summarize") || strings.Contains(strings.ToLower(msgContent), "title") {
					content = "Summary Title"
				} else if strings.Contains(strings.ToLower(msgContent), "trigger-tool") {
					content = "Triggering subagent..."
					toolCallsStr = `,
					"tool_calls": [{
						"id": "call_subagent_999",
						"type": "function",
						"function": {
							"name": "run_research_task",
							"arguments": "{\"query\":\"deep dive into rust async concurrency\"}"
						}
					}]`
				} else {
					content = "Response to: " + msgContent
				}
			}
		}

		_, _ = w.Write([]byte(`{
			"id": "chatcmpl-chain",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "gpt-4",
			"choices": [{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "` + content + `"` + toolCallsStr + `
				},
				"finish_reason": "stop"
			}],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 10,
				"total_tokens": 20
			}
		}`))
	}))
	defer mockUpstream.Close()

	proxySrv := &ProxyServer{
		client: http.DefaultClient,
		db:     dbMgr,
		config: AppConfig{
			ListenAddr:    "127.0.0.1:0",
			OpenaiBaseURL: mockUpstream.URL,
			OpenaiAPIKey:  "test-key",
		},
		lastActiveSessionMap: make(map[string]string),
	}

	// ==========================================
	// 场景 1：两个不同的会话在本地并发交错进行（用不同的 RemoteAddr 模拟）
	// ==========================================

	// 会话 A 轮次 1
	wA1 := httptest.NewRecorder()
	reqA1 := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [{"role": "user", "content": "Hello from Session A"}],
		"metadata": {"client_id": "client-a"}
	}`)))
	reqA1.RemoteAddr = "127.0.0.1:50001" // 会话 A 使用端口 50001
	proxySrv.handleProxyOpenAI(wA1, reqA1)

	// 会话 B 轮次 1
	wB1 := httptest.NewRecorder()
	reqB1 := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [{"role": "user", "content": "Hello from Session B"}],
		"metadata": {"client_id": "client-b"}
	}`)))
	reqB1.RemoteAddr = "127.0.0.1:50002" // 会话 B 使用端口 50002
	proxySrv.handleProxyOpenAI(wB1, reqB1)

	// 从库中查询前两条日志，确保分配了不同的 Session ID 且都是根节点（ParentID 为 nil）
	logs, total, err := dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if err != nil || total != 2 {
		t.Fatalf("Expected 2 logs, got %d, err: %v", total, err)
	}

	logB1 := logs[0] // 倒序排，第一条是 B
	logA1 := logs[1] // 第二条是 A

	if logA1.SessionID == logB1.SessionID {
		t.Errorf("Session A and Session B should have different SessionIDs, got %q", logA1.SessionID)
	}
	if logA1.ParentID != nil || logB1.ParentID != nil {
		t.Errorf("Initial requests should have nil ParentID")
	}

	// 会话 A 轮次 2
	wA2 := httptest.NewRecorder()
	reqA2 := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "user", "content": "Hello from Session A"},
			{"role": "assistant", "content": "Response to: Hello from Session A"},
			{"role": "user", "content": "Next question A"}
		],
		"metadata": {"client_id": "client-a"}
	}`)))
	reqA2.RemoteAddr = "127.0.0.1:50001"
	proxySrv.handleProxyOpenAI(wA2, reqA2)

	// 会话 B 轮次 2
	wB2 := httptest.NewRecorder()
	reqB2 := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "user", "content": "Hello from Session B"},
			{"role": "assistant", "content": "Response to: Hello from Session B"},
			{"role": "user", "content": "Next question B"}
		],
		"metadata": {"client_id": "client-b"}
	}`)))
	reqB2.RemoteAddr = "127.0.0.1:50002"
	proxySrv.handleProxyOpenAI(wB2, reqB2)

	// 重新查库，此时应该有 4 条记录
	logs, total, _ = dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if total != 4 {
		t.Fatalf("Expected 4 logs, got %d", total)
	}

	logB2 := logs[0]
	logA2 := logs[1]

	// 校验会话 A 轮次 2
	if logA2.SessionID != logA1.SessionID {
		t.Errorf("Session A Round 2 should inherit SessionID %q, got %q", logA1.SessionID, logA2.SessionID)
	}
	if logA2.ParentID == nil || *logA2.ParentID != logA1.ID {
		t.Errorf("Session A Round 2 parent should be %d, got %v", logA1.ID, logA2.ParentID)
	}

	// 校验会话 B 轮次 2
	if logB2.SessionID != logB1.SessionID {
		t.Errorf("Session B Round 2 should inherit SessionID %q, got %q", logB1.SessionID, logB2.SessionID)
	}
	if logB2.ParentID == nil || *logB2.ParentID != logB1.ID {
		t.Errorf("Session B Round 2 parent should be %d, got %v", logB1.ID, logB2.ParentID)
	}

	// ==========================================
	// 场景 2：会话 A 收到总结/命名辅助请求（不影响主干）
	// ==========================================
	wASummary := httptest.NewRecorder()
	reqASummary := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "user", "content": "Hello from Session A"},
			{"role": "assistant", "content": "Response to: Hello from Session A"},
			{"role": "user", "content": "Next question A"},
			{"role": "assistant", "content": "Response to: Next question A"},
			{"role": "user", "content": "Summarize session title"}
		],
		"metadata": {"client_id": "client-a"}
	}`)))
	reqASummary.RemoteAddr = "127.0.0.1:50001"
	proxySrv.handleProxyOpenAI(wASummary, reqASummary)

	// 检验是否归档在会话 A 下，且 parentID 指向 A 轮次 2
	logs, total, _ = dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if total != 5 {
		t.Fatalf("Expected 5 logs, got %d", total)
	}
	logASummary := logs[0]
	if logASummary.SessionID != logA1.SessionID {
		t.Errorf("Summary request should be grouped in Session A %q, got %q", logA1.SessionID, logASummary.SessionID)
	}
	if logASummary.ParentID == nil || *logASummary.ParentID != logA2.ID {
		t.Errorf("Summary request parent should be %d, got %v", logA2.ID, logASummary.ParentID)
	}

	// ==========================================
	// 场景 3：无上下文的辅助请求（不同端口，但同 client_id，应能回落到相同客户端指纹）
	// ==========================================
	wAFallback := httptest.NewRecorder()
	reqAFallback := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "user", "content": "predict next step"}
		],
		"metadata": {"client_id": "client-a"}
	}`)))
	reqAFallback.RemoteAddr = "127.0.0.1:50002" // 本次端口变了
	proxySrv.handleProxyOpenAI(wAFallback, reqAFallback)

	logs, total, _ = dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	logAFallback := logs[0]
	if logAFallback.SessionID != logA1.SessionID {
		t.Errorf("Fallback summary request should be grouped in Session A %q, got %q", logA1.SessionID, logAFallback.SessionID)
	}
	if logAFallback.ParentID != nil {
		t.Errorf("No-context summary request should have nil ParentID, got %v", logAFallback.ParentID)
	}

	// ==========================================
	// 场景 4：Subagent 级联多路并发语义匹配测试
	// ==========================================

	// 1. 会话 A 触发了一个 Tool Call
	wATrigger := httptest.NewRecorder()
	reqATrigger := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "user", "content": "Hello from Session A"},
			{"role": "assistant", "content": "Response to: Hello from Session A"},
			{"role": "user", "content": "trigger-tool now"}
		],
		"metadata": {"client_id": "client-a"}
	}`)))
	reqATrigger.RemoteAddr = "127.0.0.1:50001"
	proxySrv.handleProxyOpenAI(wATrigger, reqATrigger)

	// 获取触发 Tool Call 的那条日志
	logs, total, _ = dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	logATrigger := logs[0]

	// 确保在内存挂载池中有该工具调用
	proxySrv.suspMu.Lock()
	if len(proxySrv.suspendedTools) == 0 {
		proxySrv.suspMu.Unlock()
		t.Fatalf("Expected suspended tool registered, got 0")
	}
	proxySrv.suspMu.Unlock()

	// 2. 模拟 Subagent 启动，以全新无状态会话开始，并带上 system 角色规范（如 Claude Code 的典型做法）
	wSubagent := httptest.NewRecorder()
	reqSubagent := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{
		"model": "gpt-4",
		"messages": [
			{"role": "system", "content": "You are a helpful assistant specialized in research."},
			{"role": "user", "content": "Rust async concurrency deep dive"}
		]
	}`)))
	reqSubagent.RemoteAddr = "127.0.0.1:50099"
	proxySrv.handleProxyOpenAI(wSubagent, reqSubagent)

	// 获取子代理产生的日志
	logs, total, _ = dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	logSubagent := logs[0]

	// 3. 校验子代理日志是否正确绑定到了主会话，且 parent_id / parent_tool_call_id 精准写入
	if logSubagent.SessionID != logA1.SessionID {
		t.Errorf("Subagent session should inherit from parent session %q, got %q", logA1.SessionID, logSubagent.SessionID)
	}
	if logSubagent.ParentID == nil || *logSubagent.ParentID != logATrigger.ID {
		t.Errorf("Subagent parent should be %d, got %v", logATrigger.ID, logSubagent.ParentID)
	}
	if logSubagent.ParentToolCallID != "call_subagent_999" {
		t.Errorf("Subagent parent_tool_call_id should be 'call_subagent_999', got %q", logSubagent.ParentToolCallID)
	}

	// 4. 验证级联详情接口
	detail, err := dbMgr.GetLogDetail(logATrigger.ID)
	if err != nil {
		t.Fatalf("Failed to get log detail: %v", err)
	}
	if len(detail.SubLogs) != 1 {
		t.Errorf("Expected 1 sub-log cascade loaded, got %d", len(detail.SubLogs))
	} else {
		subLogDetail := detail.SubLogs[0]
		if subLogDetail.ID != logSubagent.ID {
			t.Errorf("Cascade sub-log id should be %d, got %d", logSubagent.ID, subLogDetail.ID)
		}
	}
}

func TestSemanticMatchUnit(t *testing.T) {
	prompt := "Rust async concurrency deep dive"
	args := `{"query":"deep dive into rust async concurrency"}`
	if !semanticMatch(prompt, args) {
		t.Errorf("Expected semanticMatch to return true, got false")
	}
}

func TestExtractRequestSessionID(t *testing.T) {
	body := []byte(`{
		"metadata": {
			"user_id": "{\"device_id\":\"abc\",\"session_id\":\"claude-session-123\"}"
		}
	}`)
	if got := ExtractRequestSessionID(body); got != "claude-session-123" {
		t.Fatalf("expected nested metadata session id, got %q", got)
	}

	body = []byte(`{
		"metadata": {
			"session_id": "top-level-session",
			"user_id": {"session_id": "nested-session"}
		}
	}`)
	if got := ExtractRequestSessionID(body); got != "top-level-session" {
		t.Fatalf("expected top-level metadata session id, got %q", got)
	}

	body = []byte(`{
		"previous_response_id": "resp_prev_123",
		"conversation_id": "conversation-456",
		"metadata": {
			"thread_id": "thread-789"
		}
	}`)
	handles := ExtractConversationHandles(body)
	if len(handles) != 3 {
		t.Fatalf("expected 3 conversation handles, got %d", len(handles))
	}
}

func TestParseOpenAIResponsesEnvelope(t *testing.T) {
	body := []byte(`{
		"model": "gpt-4.1",
		"instructions": "You are helpful",
		"input": [
			{
				"role": "user",
				"content": [{"type": "input_text", "text": "Explain Go interfaces"}]
			}
		],
		"tools": [
			{"type": "function", "name": "search_docs", "description": "Search docs", "parameters": {"type": "object"}}
		],
		"previous_response_id": "resp_prev_123",
		"conversation_id": "conv_456",
		"metadata": {"client_id": "responses-client"}
	}`)

	requestInfo, err := ParseUnifiedRequestEnvelope("/v1/responses", body)
	if err != nil {
		t.Fatalf("failed to parse responses request: %v", err)
	}
	if requestInfo.Provider != "openai-responses" {
		t.Fatalf("expected provider openai-responses, got %q", requestInfo.Provider)
	}
	if len(requestInfo.Messages) != 2 {
		t.Fatalf("expected 2 normalized messages, got %d", len(requestInfo.Messages))
	}
	if requestInfo.Messages[0].Role != "system" || requestInfo.Messages[0].Content != "You are helpful" {
		t.Fatalf("unexpected system message: %+v", requestInfo.Messages[0])
	}
	if requestInfo.Messages[1].Role != "user" || requestInfo.Messages[1].Content != "Explain Go interfaces" {
		t.Fatalf("unexpected user message: %+v", requestInfo.Messages[1])
	}
	if len(requestInfo.Tools) != 1 || requestInfo.Tools[0].Name != "search_docs" {
		t.Fatalf("unexpected tools parsed: %+v", requestInfo.Tools)
	}

	responseBody := []byte(`{
		"id": "resp_789",
		"model": "gpt-4.1",
		"output": [
			{
				"type": "message",
				"role": "assistant",
				"content": [{"type": "output_text", "text": "Interfaces define behavior."}]
			},
			{
				"type": "function_call",
				"id": "fc_1",
				"call_id": "call_1",
				"name": "search_docs",
				"arguments": "{\"query\":\"go interfaces\"}"
			}
		],
		"usage": {"input_tokens": 10, "output_tokens": 6, "total_tokens": 16}
	}`)

	responseInfo, err := ParseUnifiedResponseEnvelope("/v1/responses", responseBody)
	if err != nil {
		t.Fatalf("failed to parse responses response: %v", err)
	}
	if responseInfo.Message.Content != "Interfaces define behavior." {
		t.Fatalf("unexpected response content: %q", responseInfo.Message.Content)
	}
	if len(responseInfo.Message.ToolCalls) != 1 || responseInfo.Message.ToolCalls[0].ID != "call_1" {
		t.Fatalf("unexpected response tool calls: %+v", responseInfo.Message.ToolCalls)
	}
	if responseInfo.Usage.TotalTokens != 16 {
		t.Fatalf("unexpected usage: %+v", responseInfo.Usage)
	}
}

func TestResponsesContinuationChaining(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "llm_tracer_responses_chain_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "responses_chain.db")
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyStr := string(bodyBytes)
		if strings.Contains(bodyStr, "resp_root_1") {
			_, _ = w.Write([]byte(`{
				"id": "resp_followup_2",
				"model": "gpt-4.1",
				"output": [{
					"type": "message",
					"role": "assistant",
					"content": [{"type": "output_text", "text": "Follow-up answer"}]
				}],
				"usage": {"input_tokens": 12, "output_tokens": 8, "total_tokens": 20}
			}`))
			return
		}
		_, _ = w.Write([]byte(`{
			"id": "resp_root_1",
			"model": "gpt-4.1",
			"output": [{
				"type": "message",
				"role": "assistant",
				"content": [{"type": "output_text", "text": "Root answer"}]
			}],
			"usage": {"input_tokens": 10, "output_tokens": 6, "total_tokens": 16}
		}`))
	}))
	defer mockUpstream.Close()

	proxySrv := &ProxyServer{
		client: http.DefaultClient,
		db:     dbMgr,
		config: AppConfig{
			ListenAddr:             "127.0.0.1:0",
			OpenaiResponsesBaseURL: mockUpstream.URL,
			OpenaiResponsesAPIKey:  "test-key",
		},
		lastActiveSessionMap: make(map[string]string),
	}

	firstReq := httptest.NewRequest("POST", "/v1/responses", bytes.NewReader([]byte(`{
		"model": "gpt-4.1",
		"instructions": "Be concise",
		"input": "First question",
		"metadata": {"client_id": "responses-client"}
	}`)))
	firstReq.RemoteAddr = "127.0.0.1:51001"
	firstRec := httptest.NewRecorder()
	proxySrv.handleProxyOpenAIResponses(firstRec, firstReq)

	secondReq := httptest.NewRequest("POST", "/v1/responses", bytes.NewReader([]byte(`{
		"model": "gpt-4.1",
		"input": "Second question",
		"previous_response_id": "resp_root_1",
		"metadata": {"client_id": "responses-client"}
	}`)))
	secondReq.RemoteAddr = "127.0.0.1:51999"
	secondRec := httptest.NewRecorder()
	proxySrv.handleProxyOpenAIResponses(secondRec, secondReq)

	logs, total, err := dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 logs, got %d", total)
	}

	secondLog := logs[0]
	firstLog := logs[1]
	if secondLog.SessionID != firstLog.SessionID {
		t.Fatalf("responses continuation should inherit session %q, got %q", firstLog.SessionID, secondLog.SessionID)
	}
	if secondLog.ParentID == nil || *secondLog.ParentID != firstLog.ID {
		t.Fatalf("responses continuation should point to parent %d, got %v", firstLog.ID, secondLog.ParentID)
	}

	secondDetail, err := dbMgr.GetLogDetail(secondLog.ID)
	if err != nil {
		t.Fatalf("failed to read response continuation detail: %v", err)
	}
	if secondDetail.Response.Content != "Follow-up answer" {
		t.Fatalf("unexpected follow-up content: %q", secondDetail.Response.Content)
	}
}

func TestProxyDoesNotFailOnResponseParseError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "llm_tracer_parse_error_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "parse_error.db")
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	badPayload := `{"id":"resp_bad"`
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(badPayload))
	}))
	defer mockUpstream.Close()

	proxySrv := &ProxyServer{
		client: http.DefaultClient,
		db:     dbMgr,
		config: AppConfig{
			ListenAddr:             "127.0.0.1:0",
			OpenaiResponsesBaseURL: mockUpstream.URL,
			OpenaiResponsesAPIKey:  "test-key",
		},
		lastActiveSessionMap: make(map[string]string),
	}

	req := httptest.NewRequest("POST", "/v1/responses", bytes.NewReader([]byte(`{
		"model": "gpt-4.1",
		"input": "Hello",
		"metadata": {"client_id": "parse-error-client"}
	}`)))
	req.RemoteAddr = "127.0.0.1:52001"
	rec := httptest.NewRecorder()

	proxySrv.handleProxyOpenAIResponses(rec, req)
	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected proxy to preserve upstream 200, got %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if string(respBody) != badPayload {
		t.Fatalf("expected raw payload to pass through, got %q", string(respBody))
	}

	logs, total, err := dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 log, got %d", total)
	}
	detail, err := dbMgr.GetLogDetail(logs[0].ID)
	if err != nil {
		t.Fatalf("failed to fetch detail: %v", err)
	}
	if detail.RawResponse != badPayload {
		t.Fatalf("expected raw response to be stored, got %q", detail.RawResponse)
	}
}

func TestToolResultContinuationChaining(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "llm_tracer_tool_continue_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "tool_continue.db")
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	// 该回归样本来自真实 Pi 会话日志，必须以内嵌夹具形式保存在测试中，不能依赖本地 data/llm_tracer.db。
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyStr := string(bodyBytes)
		if strings.Contains(bodyStr, "tool_call_id") {
			_, _ = w.Write([]byte(piAgentToolContinuationStreamFixture))
			return
		}
		_, _ = w.Write([]byte(piAgentInitialChatStreamFixture))
	}))
	defer mockUpstream.Close()

	proxySrv := &ProxyServer{
		client: http.DefaultClient,
		db:     dbMgr,
		config: AppConfig{
			ListenAddr:    "127.0.0.1:0",
			OpenaiBaseURL: mockUpstream.URL,
			OpenaiAPIKey:  "test-key",
		},
		lastActiveSessionMap: make(map[string]string),
	}

	firstReq := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(piAgentInitialChatRequestFixture)))
	firstReq.RemoteAddr = "127.0.0.1:53001"
	firstReq.Header.Set("User-Agent", "OpenAI/JS 6.26.0")
	firstRec := httptest.NewRecorder()
	proxySrv.handleProxyOpenAI(firstRec, firstReq)

	secondReq := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(piAgentToolContinuationRequestFixture)))
	secondReq.RemoteAddr = "127.0.0.1:53999"
	secondReq.Header.Set("User-Agent", "OpenAI/JS 6.26.0")
	secondRec := httptest.NewRecorder()
	proxySrv.handleProxyOpenAI(secondRec, secondReq)

	logs, total, err := dbMgr.GetLogs(1, 10, "", "", "", nil, false)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 logs, got %d", total)
	}

	secondLog := logs[0]
	firstLog := logs[1]
	if secondLog.SessionID != firstLog.SessionID {
		t.Fatalf("tool-result continuation should inherit session %q, got %q", firstLog.SessionID, secondLog.SessionID)
	}
	if secondLog.ParentID == nil || *secondLog.ParentID != firstLog.ID {
		t.Fatalf("tool-result continuation should point to parent %d, got %v", firstLog.ID, secondLog.ParentID)
	}

	firstDetail, err := dbMgr.GetLogDetail(firstLog.ID)
	if err != nil {
		t.Fatalf("failed to load first detail: %v", err)
	}
	if firstDetail.TotalTokens != 3978 || firstDetail.InputTokens != 3728 || firstDetail.OutputTokens != 250 {
		t.Fatalf("unexpected root token usage from real Pi fixture: %+v", firstDetail)
	}

	secondDetail, err := dbMgr.GetLogDetail(secondLog.ID)
	if err != nil {
		t.Fatalf("failed to load second detail: %v", err)
	}
	if secondDetail.TotalTokens != 5136 || secondDetail.InputTokens != 4940 || secondDetail.OutputTokens != 196 {
		t.Fatalf("unexpected continuation token usage from real Pi fixture: %+v", secondDetail)
	}
}

func TestRepairSessionIDsFromRawRequest(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "llm_tracer_repair_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "repair.db")
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	rawRequest := `{"metadata":{"user_id":"{\"device_id\":\"abc\",\"session_id\":\"claude-session-123\"}"}}`
	insert := func(sessionID string) int64 {
		logID, err := dbMgr.InsertLog(&UnifiedLog{
			Provider:    "anthropic",
			Model:       "claude-test",
			Path:        "/v1/messages",
			StatusCode:  http.StatusOK,
			Prompt:      []ChatMessage{{Role: "user", Content: "hello"}},
			Response:    ChatMessage{Role: "assistant", Content: "ok"},
			RawRequest:  rawRequest,
			RawResponse: `{"ok":true}`,
			SessionID:   sessionID,
		})
		if err != nil {
			t.Fatalf("failed to insert log: %v", err)
		}
		return logID
	}

	firstID := insert("wrong-session-a")
	secondID := insert("wrong-session-b")

	repaired, err := dbMgr.RepairSessionIDsFromRawRequest()
	if err != nil {
		t.Fatalf("failed to repair session ids: %v", err)
	}
	if repaired != 2 {
		t.Fatalf("expected 2 repaired rows, got %d", repaired)
	}

	assertSessionID := func(logID int64) {
		var sessionID string
		err := dbMgr.db.QueryRow(`SELECT COALESCE(session_id, '') FROM logs WHERE id = ?`, logID).Scan(&sessionID)
		if err != nil {
			t.Fatalf("failed to read session id for log %d: %v", logID, err)
		}
		if sessionID != "claude-session-123" {
			t.Fatalf("expected repaired session id for log %d, got %q", logID, sessionID)
		}
	}

	assertSessionID(firstID)
	assertSessionID(secondID)
}

func TestResolveSessionWithAuthoritativeMetadata(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "llm_tracer_authoritative_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "authoritative.db")
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	seedID, err := dbMgr.InsertLog(&UnifiedLog{
		Provider:    "anthropic",
		Model:       "claude-test",
		Path:        "/v1/messages",
		StatusCode:  http.StatusOK,
		Prompt:      []ChatMessage{{Role: "user", Content: "你好"}},
		Response:    ChatMessage{Role: "assistant", Content: "你好"},
		RawRequest:  `{"metadata":{"user_id":"{\"session_id\":\"claude-session-123\"}"}}`,
		RawResponse: `{"ok":true}`,
		SessionID:   "claude-session-123",
		RequestHandles: []ConversationHandle{{Kind: "session_id", Value: "claude-session-123"}},
	})
	if err != nil {
		t.Fatalf("failed to seed log: %v", err)
	}

	proxySrv := &ProxyServer{
		db:                   dbMgr,
		lastActiveSessionMap: make(map[string]string),
	}

	prompt := []ChatMessage{
		{Role: "system", Content: "Generate a concise, sentence-case title for this coding session."},
		{Role: "user", Content: "<session>用子代理看看今天的科技新闻呢</session>"},
	}

	body := []byte(`{"metadata":{"user_id":"{\"session_id\":\"claude-session-123\",\"client_id\":\"client-a\"}"}}`)
	httpReq := httptest.NewRequest("POST", "/v1/messages", bytes.NewReader(body))
	httpReq.RemoteAddr = "127.0.0.1:50001"
	sessionID, parentID, parentToolCallID, clientFingerprint := proxySrv.resolveSessionAndParent(httpReq, "anthropic", prompt, ExtractConversationHandles(body), ExtractRequestHints(body))
	if sessionID != "claude-session-123" {
		t.Fatalf("expected authoritative session id, got %q", sessionID)
	}
	if parentID == nil || *parentID != seedID {
		t.Fatalf("expected parent id %d, got %v", seedID, parentID)
	}
	if parentToolCallID != "" {
		t.Fatalf("expected empty parent tool call id, got %q", parentToolCallID)
	}
	if clientFingerprint == "" {
		t.Fatalf("expected non-empty client fingerprint")
	}
}

func TestGetSessionLogs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "llm_tracer_session_logs_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "session_logs.db")
	dbMgr, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer dbMgr.Close()

	// 插入属于同一个会话的几条日志
	sessionID := "test-session-999"
	_, _ = dbMgr.InsertLog(&UnifiedLog{
		Provider:   "openai",
		Model:      "gpt-4",
		Path:       "/v1/chat/completions",
		StatusCode: 200,
		SessionID:  sessionID,
		Prompt:     []ChatMessage{{Role: "user", Content: "Hello 1"}},
		Response:   ChatMessage{Role: "assistant", Content: "Hi 1"},
	})
	_, _ = dbMgr.InsertLog(&UnifiedLog{
		Provider:   "openai",
		Model:      "gpt-4",
		Path:       "/v1/chat/completions",
		StatusCode: 200,
		SessionID:  sessionID,
		Prompt:     []ChatMessage{{Role: "user", Content: "Hello 2"}},
		Response:   ChatMessage{Role: "assistant", Content: "Hi 2"},
	})

	// 插入另一条不属于该会话的日志
	_, _ = dbMgr.InsertLog(&UnifiedLog{
		Provider:   "openai",
		Model:      "gpt-4",
		Path:       "/v1/chat/completions",
		StatusCode: 200,
		SessionID:  "other-session",
		Prompt:     []ChatMessage{{Role: "user", Content: "Hello 3"}},
		Response:   ChatMessage{Role: "assistant", Content: "Hi 3"},
	})

	logs, err := dbMgr.GetSessionLogs(sessionID)
	if err != nil {
		t.Fatalf("failed to get session logs: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 session logs, got %d", len(logs))
	}
	if logs[0].Response.Content != "Hi 1" || logs[1].Response.Content != "Hi 2" {
		t.Errorf("unexpected logs order or contents: %+v", logs)
	}
}
