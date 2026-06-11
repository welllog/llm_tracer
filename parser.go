package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ConversationHandle struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type RequestParseResult struct {
	Model    string             `json:"model"`
	Messages []ChatMessage      `json:"messages"`
	Tools    []ToolDef          `json:"tools"`
	Provider string             `json:"provider"`
	Handles  []ConversationHandle `json:"handles,omitempty"`
	Hints    map[string]string  `json:"hints,omitempty"`
}

type ResponseParseResult struct {
	Model   string               `json:"model"`
	Message ChatMessage          `json:"message"`
	Usage   TokenUsage           `json:"usage"`
	Handles []ConversationHandle `json:"handles,omitempty"`
}

// 通用消息定义
type ChatMessage struct {
	Role       string     `json:"role"`                 // system, user, assistant, tool
	Content    string     `json:"content,omitempty"`    // 文本内容
	Name       string     `json:"name,omitempty"`       // 可选名字
	Thinking   string     `json:"thinking,omitempty"`   // 思维链/思考过程
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"` // Assistant 调用的工具列表
	ToolCallID string     `json:"tool_call_id,omitempty"` // Tool 响应对应的 ID
}

type ToolCall struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // function
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolDef struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  string `json:"parameters,omitempty"` // JSON string representation
}

type TokenUsage struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	TotalTokens         int `json:"total_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens,omitempty"`
	CacheReadTokens     int `json:"cache_read_tokens,omitempty"`
	CachedTokens        int `json:"cached_tokens,omitempty"`
}

// ==========================================
// OpenAI 解析部分
// ==========================================

type openAIRequest struct {
	Model     string                  `json:"model"`
	Messages  []openAIMessage         `json:"messages"`
	Tools     []openAITool            `json:"tools"`
	Functions []struct {
		Name        string          `json:"name"`
		Description string          `json:"description,omitempty"`
		Parameters  json.RawMessage `json:"parameters,omitempty"`
	} `json:"functions"` // 适配旧版兼容性或部分特定 Agent 的习惯
	System    string                  `json:"system"`
}

type openAIResponsesRequest struct {
	Model              string          `json:"model"`
	Instructions       json.RawMessage `json:"instructions"`
	Input              json.RawMessage `json:"input"`
	Tools              []json.RawMessage `json:"tools"`
	PreviousResponseID string          `json:"previous_response_id"`
	ConversationID     string          `json:"conversation_id"`
	SessionID          string          `json:"session_id"`
	Metadata           struct {
		SessionID      string          `json:"session_id"`
		ConversationID string          `json:"conversation_id"`
		ThreadID       string          `json:"thread_id"`
		UserID         json.RawMessage `json:"user_id"`
		ClientID       string          `json:"client_id"`
		DeviceID       string          `json:"device_id"`
		Custom         map[string]any  `json:"custom"`
	} `json:"metadata"`
}

type openAIResponsesResponse struct {
	ID     string `json:"id"`
	Model  string `json:"model"`
	Output []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Role     string `json:"role"`
		Name     string `json:"name,omitempty"`
		CallID   string `json:"call_id,omitempty"`
		Arguments string `json:"arguments,omitempty"`
		Content  []struct {
			Type string `json:"type"`
			Text string `json:"text,omitempty"`
		} `json:"content,omitempty"`
	} `json:"output"`
	Usage struct {
		InputTokens              int `json:"input_tokens"`
		OutputTokens             int `json:"output_tokens"`
		TotalTokens              int `json:"total_tokens"`
		PromptTokens             int `json:"prompt_tokens"`
		CompletionTokens         int `json:"completion_tokens"`
		InputTokensDetails       *struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"input_tokens_details,omitempty"`
		PromptTokensDetails *struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details,omitempty"`
	} `json:"usage"`
}

type openAIMessage struct {
	Role             string           `json:"role"`
	Content          json.RawMessage  `json:"content"` // 可能为 string 或 array
	Name             string           `json:"name,omitempty"`
	ReasoningContent string           `json:"reasoning_content,omitempty"`
	ToolCalls        []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string           `json:"tool_call_id,omitempty"`
	FunctionCall     *struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function_call,omitempty"`
}

type openAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type openAITool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string          `json:"name"`
		Description string          `json:"description,omitempty"`
		Parameters  json.RawMessage `json:"parameters,omitempty"`
	} `json:"function"`
}

type openAIResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role             string           `json:"role"`
			Content          string           `json:"content"`
			ReasoningContent string           `json:"reasoning_content,omitempty"`
			ToolCalls        []openAIToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails *struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details,omitempty"`
	} `json:"usage"`
}

// 解析 OpenAI 请求
func ParseOpenAIRequest(body []byte) (string, []ChatMessage, []ToolDef, error) {
	var req openAIRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return "", nil, nil, err
	}

	var messages []ChatMessage
	for _, m := range req.Messages {
		contentStr := parseOpenAIContent(m.Content)
		msg := ChatMessage{
			Role:       m.Role,
			Content:    contentStr,
			Name:       m.Name,
			Thinking:   m.ReasoningContent,
			ToolCallID: m.ToolCallID,
		}

		if len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				msg.ToolCalls = append(msg.ToolCalls, ToolCall{
					ID:        tc.ID,
					Type:      tc.Type,
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				})
			}
		} else if m.FunctionCall != nil {
			msg.ToolCalls = []ToolCall{{
				ID:        "legacy_fn",
				Type:      "function",
				Name:      m.FunctionCall.Name,
				Arguments: m.FunctionCall.Arguments,
			}}
		}

		messages = append(messages, msg)
	}

	var tools []ToolDef
	for _, t := range req.Tools {
		if t.Type == "function" {
			tools = append(tools, ToolDef{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  string(t.Function.Parameters),
			})
		}
	}
	// 补充 legacy functions 解析
	for _, f := range req.Functions {
		tools = append(tools, ToolDef{
			Name:        f.Name,
			Description: f.Description,
			Parameters:  string(f.Parameters),
		})
	}

	return req.Model, messages, tools, nil
}

func ParseOpenAIRequestWithHandles(body []byte) (RequestParseResult, error) {
	model, messages, tools, err := ParseOpenAIRequest(body)
	if err != nil {
		return RequestParseResult{}, err
	}
	result := RequestParseResult{
		Model:    model,
		Messages: messages,
		Tools:    tools,
		Provider: "openai",
		Handles:  ExtractConversationHandles(body),
		Hints:    ExtractRequestHints(body),
	}
	return result, nil
}

func ParseOpenAIResponsesRequest(body []byte) (RequestParseResult, error) {
	var req openAIResponsesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return RequestParseResult{}, err
	}

	result := RequestParseResult{
		Model:    req.Model,
		Provider: "openai-responses",
		Handles:  ExtractConversationHandles(body),
		Hints:    ExtractRequestHints(body),
	}

	if instructions := parseResponsesText(req.Instructions); instructions != "" {
		result.Messages = append(result.Messages, ChatMessage{Role: "system", Content: instructions})
	}

	result.Messages = append(result.Messages, parseResponsesInput(req.Input)...)
	result.Tools = parseResponsesTools(req.Tools)

	return result, nil
}

// 解析 OpenAI 响应
func ParseOpenAIResponse(body []byte) (string, ChatMessage, TokenUsage, error) {
	var resp openAIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", ChatMessage{}, TokenUsage{}, err
	}

	var msg ChatMessage
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		msg.Role = choice.Message.Role
		msg.Content = choice.Message.Content
		msg.Thinking = choice.Message.ReasoningContent
		for _, tc := range choice.Message.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, ToolCall{
				ID:        tc.ID,
				Type:      tc.Type,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			})
		}
	}

	usage := TokenUsage{
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}
	if resp.Usage.PromptTokensDetails != nil {
		usage.CachedTokens = resp.Usage.PromptTokensDetails.CachedTokens
	}

	return resp.Model, msg, usage, nil
}

func ParseOpenAIResponsesResponse(body []byte) (ResponseParseResult, error) {
	var resp openAIResponsesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ResponseParseResult{}, err
	}

	result := ResponseParseResult{
		Model: resp.Model,
		Message: ChatMessage{Role: "assistant"},
		Handles: []ConversationHandle{{Kind: "response_id", Value: strings.TrimSpace(resp.ID)}},
	}

	var textParts []string
	for _, item := range resp.Output {
		switch item.Type {
		case "message":
			if item.Role != "" {
				result.Message.Role = item.Role
			}
			for _, content := range item.Content {
				if content.Type == "output_text" || content.Type == "text" {
					if strings.TrimSpace(content.Text) != "" {
						textParts = append(textParts, content.Text)
					}
				}
			}
		case "function_call":
			toolCallID := strings.TrimSpace(item.CallID)
			if toolCallID == "" {
				toolCallID = strings.TrimSpace(item.ID)
			}
			result.Message.ToolCalls = append(result.Message.ToolCalls, ToolCall{
				ID:        toolCallID,
				Type:      "function",
				Name:      item.Name,
				Arguments: item.Arguments,
			})
		}
	}
	result.Message.Content = strings.Join(textParts, "\n")

	result.Usage = normalizeResponsesUsage(resp.Usage)
	result.Handles = compactHandles(result.Handles)

	return result, nil
}

// 增量解析 OpenAI SSE 流式数据
func ParseOpenAIStreamChunk(chunk []byte, currentResp *ChatMessage, currentUsage *TokenUsage) error {
	// 去掉 "data: " 前缀
	line := string(chunk)
	if !strings.HasPrefix(line, "data: ") {
		return nil
	}
	dataStr := strings.TrimPrefix(line, "data: ")
	dataStr = strings.TrimSpace(dataStr)
	if dataStr == "[DONE]" || dataStr == "" {
		return nil
	}

	var chunkObj struct {
		Choices []struct {
			Delta struct {
				Role             string           `json:"role"`
				Content          string           `json:"content"`
				ReasoningContent string           `json:"reasoning_content,omitempty"`
				ToolCalls        []struct {
					Index    int    `json:"index"`
					ID       string `json:"id"`
					Type     string `json:"type"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"delta"`
		} `json:"choices"`
		Usage *struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			TotalTokens         int `json:"total_tokens"`
			PromptTokensDetails *struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details,omitempty"`
		} `json:"usage"`
	}

	if err := json.Unmarshal([]byte(dataStr), &chunkObj); err != nil {
		return err
	}

	if chunkObj.Usage != nil {
		currentUsage.PromptTokens = chunkObj.Usage.PromptTokens
		currentUsage.CompletionTokens = chunkObj.Usage.CompletionTokens
		currentUsage.TotalTokens = chunkObj.Usage.TotalTokens
		if chunkObj.Usage.PromptTokensDetails != nil {
			currentUsage.CachedTokens = chunkObj.Usage.PromptTokensDetails.CachedTokens
		}
	}

	if len(chunkObj.Choices) > 0 {
		delta := chunkObj.Choices[0].Delta
		if delta.Role != "" {
			currentResp.Role = delta.Role
		}
		if delta.Content != "" {
			currentResp.Content += delta.Content
		}
		if delta.ReasoningContent != "" {
			currentResp.Thinking += delta.ReasoningContent
		}
		for _, tc := range delta.ToolCalls {
			// 流式 tool_calls 可能是分段的，需要根据 Index 拼合
			for len(currentResp.ToolCalls) <= tc.Index {
				currentResp.ToolCalls = append(currentResp.ToolCalls, ToolCall{})
			}
			t := &currentResp.ToolCalls[tc.Index]
			if tc.ID != "" {
				t.ID = tc.ID
			}
			if tc.Type != "" {
				t.Type = tc.Type
			}
			if tc.Function.Name != "" {
				t.Name += tc.Function.Name
			}
			if tc.Function.Arguments != "" {
				t.Arguments += tc.Function.Arguments
			}
		}
	}

	return nil
}

func ParseOpenAIResponsesStreamChunk(chunk []byte, currentResp *ChatMessage, currentUsage *TokenUsage, responseHandles *[]ConversationHandle) error {
	line := string(chunk)
	if !strings.HasPrefix(line, "data: ") {
		return nil
	}
	dataStr := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
	if dataStr == "" || dataStr == "[DONE]" {
		return nil
	}

	var event struct {
		Type string `json:"type"`
		Response *struct {
			ID string `json:"id"`
			Usage struct {
				InputTokens      int `json:"input_tokens"`
				OutputTokens     int `json:"output_tokens"`
				TotalTokens      int `json:"total_tokens"`
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				InputTokensDetails *struct {
					CachedTokens int `json:"cached_tokens"`
				} `json:"input_tokens_details,omitempty"`
				PromptTokensDetails *struct {
					CachedTokens int `json:"cached_tokens"`
				} `json:"prompt_tokens_details,omitempty"`
			} `json:"usage"`
		} `json:"response,omitempty"`
		Item *struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Role     string `json:"role"`
			Name     string `json:"name,omitempty"`
			CallID   string `json:"call_id,omitempty"`
			Arguments string `json:"arguments,omitempty"`
			Content  []struct {
				Type string `json:"type"`
				Text string `json:"text,omitempty"`
			} `json:"content,omitempty"`
		} `json:"item,omitempty"`
		OutputIndex *int `json:"output_index,omitempty"`
		ItemID      string `json:"item_id,omitempty"`
		Delta       string `json:"delta,omitempty"`
		ContentIndex *int `json:"content_index,omitempty"`
	}

	if err := json.Unmarshal([]byte(dataStr), &event); err != nil {
		return err
	}

	if event.Response != nil {
		if responseID := strings.TrimSpace(event.Response.ID); responseID != "" {
			*responseHandles = append(*responseHandles, ConversationHandle{Kind: "response_id", Value: responseID})
		}
		*currentUsage = normalizeResponsesUsage(event.Response.Usage)
	}

	switch event.Type {
	case "response.output_item.added":
		if event.Item == nil {
			return nil
		}
		if event.Item.Type == "message" {
			if event.Item.Role != "" {
				currentResp.Role = event.Item.Role
			}
			for _, content := range event.Item.Content {
				if (content.Type == "output_text" || content.Type == "text") && content.Text != "" {
					currentResp.Content += content.Text
				}
			}
		}
		if event.Item.Type == "function_call" {
			toolCallID := strings.TrimSpace(event.Item.CallID)
			if toolCallID == "" {
				toolCallID = strings.TrimSpace(event.Item.ID)
			}
			currentResp.ToolCalls = append(currentResp.ToolCalls, ToolCall{
				ID:        toolCallID,
				Type:      "function",
				Name:      event.Item.Name,
				Arguments: event.Item.Arguments,
			})
		}
	case "response.output_text.delta":
		if event.Delta != "" {
			currentResp.Role = "assistant"
			currentResp.Content += event.Delta
		}
	case "response.function_call_arguments.delta":
		if len(currentResp.ToolCalls) == 0 {
			currentResp.ToolCalls = append(currentResp.ToolCalls, ToolCall{Type: "function", ID: strings.TrimSpace(event.ItemID)})
		}
		currentResp.ToolCalls[len(currentResp.ToolCalls)-1].Arguments += event.Delta
	}

	return nil
}

func parseOpenAIContent(raw json.RawMessage) string {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	var blocks []map[string]any
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return string(raw)
	}
	var lines []string
	for _, block := range blocks {
		blockType, _ := block["type"].(string)
		switch blockType {
		case "text":
			if text, ok := block["text"].(string); ok && strings.TrimSpace(text) != "" {
				lines = append(lines, text)
			}
		case "image_url":
			imageUrl, ok := block["image_url"].(map[string]any)
			if ok {
				url, _ := imageUrl["url"].(string)
				if strings.HasPrefix(url, "data:") {
					lines = append(lines, fmt.Sprintf("[image: base64, size %d bytes]", len(url)))
				} else {
					lines = append(lines, fmt.Sprintf("[image: %s]", url))
				}
			}
		}
	}
	return strings.Join(lines, "\n")
}

func parseResponsesInput(raw json.RawMessage) []ChatMessage {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return nil
	}

	var plainText string
	if err := json.Unmarshal(raw, &plainText); err == nil {
		if strings.TrimSpace(plainText) == "" {
			return nil
		}
		return []ChatMessage{{Role: "user", Content: plainText}}
	}

	var rawArray []json.RawMessage
	if err := json.Unmarshal(raw, &rawArray); err != nil {
		return nil
	}

	var messages []ChatMessage
	for _, item := range rawArray {
		var msg struct {
			Role    string          `json:"role"`
			Type    string          `json:"type"`
			Content json.RawMessage `json:"content"`
			Text    string          `json:"text"`
		}
		if err := json.Unmarshal(item, &msg); err != nil {
			continue
		}

		if msg.Role != "" {
			text := parseResponsesText(msg.Content)
			if text == "" {
				text = strings.TrimSpace(msg.Text)
			}
			messages = append(messages, ChatMessage{Role: msg.Role, Content: text})
			continue
		}

		if msg.Type == "input_text" || msg.Type == "text" {
			text := strings.TrimSpace(msg.Text)
			if text != "" {
				messages = append(messages, ChatMessage{Role: "user", Content: text})
			}
		}
	}

	return coalesceAdjacentMessages(messages)
}

func parseResponsesText(raw json.RawMessage) string {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return ""
	}

	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}

	var blocks []map[string]any
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return ""
	}

	var lines []string
	for _, block := range blocks {
		blockType, _ := block["type"].(string)
		if blockType == "input_text" || blockType == "output_text" || blockType == "text" {
			if text, ok := block["text"].(string); ok && strings.TrimSpace(text) != "" {
				lines = append(lines, text)
			}
		}
	}
	return strings.Join(lines, "\n")
}

func parseResponsesTools(rawTools []json.RawMessage) []ToolDef {
	var tools []ToolDef
	for _, raw := range rawTools {
		var base struct {
			Type        string          `json:"type"`
			Name        string          `json:"name"`
			Description string          `json:"description,omitempty"`
			Parameters  json.RawMessage `json:"parameters,omitempty"`
		}
		if err := json.Unmarshal(raw, &base); err != nil {
			continue
		}
		if base.Type != "function" || strings.TrimSpace(base.Name) == "" {
			continue
		}
		tools = append(tools, ToolDef{
			Name:        base.Name,
			Description: base.Description,
			Parameters:  string(base.Parameters),
		})
	}
	return tools
}

func normalizeResponsesUsage(usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	TotalTokens              int `json:"total_tokens"`
	PromptTokens             int `json:"prompt_tokens"`
	CompletionTokens         int `json:"completion_tokens"`
	InputTokensDetails       *struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"input_tokens_details,omitempty"`
	PromptTokensDetails *struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details,omitempty"`
}) TokenUsage {
	result := TokenUsage{
		PromptTokens:     usage.InputTokens,
		CompletionTokens: usage.OutputTokens,
		TotalTokens:      usage.TotalTokens,
	}
	if result.PromptTokens == 0 {
		result.PromptTokens = usage.PromptTokens
	}
	if result.CompletionTokens == 0 {
		result.CompletionTokens = usage.CompletionTokens
	}
	if result.TotalTokens == 0 {
		result.TotalTokens = result.PromptTokens + result.CompletionTokens
	}
	if usage.InputTokensDetails != nil {
		result.CachedTokens = usage.InputTokensDetails.CachedTokens
	}
	if result.CachedTokens == 0 && usage.PromptTokensDetails != nil {
		result.CachedTokens = usage.PromptTokensDetails.CachedTokens
	}
	return result
}

func coalesceAdjacentMessages(messages []ChatMessage) []ChatMessage {
	if len(messages) < 2 {
		return messages
	}
	result := []ChatMessage{messages[0]}
	for _, msg := range messages[1:] {
		last := &result[len(result)-1]
		if last.Role == msg.Role && len(last.ToolCalls) == 0 && len(msg.ToolCalls) == 0 && last.ToolCallID == "" && msg.ToolCallID == "" {
			if strings.TrimSpace(msg.Content) != "" {
				if strings.TrimSpace(last.Content) != "" {
					last.Content += "\n"
				}
				last.Content += msg.Content
			}
			continue
		}
		result = append(result, msg)
	}
	return result
}

// ==========================================
// Anthropic 解析部分
// ==========================================

type anthropicRequest struct {
	Model    string             `json:"model"`
	System   json.RawMessage    `json:"system"` // 可以是 string 或 array
	Messages []anthropicMessage `json:"messages"`
	Tools    []anthropicTool    `json:"tools"`
}

type anthropicMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"` // 可以是 string 或 array
}

type anthropicTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

type anthropicResponse struct {
	Model   string                 `json:"model"`
	Content []anthropicContentBlock `json:"content"`
	Usage   struct {
		InputTokens              int `json:"input_tokens"`
		OutputTokens             int `json:"output_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
	} `json:"usage"`
}

type anthropicContentBlock struct {
	Type  string          `json:"type"` // text, tool_use, thinking
	Text  string          `json:"text,omitempty"`
	Input json.RawMessage `json:"input,omitempty"` // for tool_use
	Name  string          `json:"name,omitempty"`  // for tool_use
	ID    string          `json:"id,omitempty"`    // for tool_use
	Signature string      `json:"signature,omitempty"` // for thinking
	Thinking  string      `json:"thinking,omitempty"` // for thinking block (Claude 3.7)
}

// 解析 Anthropic 请求
func ParseAnthropicRequest(body []byte) (string, []ChatMessage, []ToolDef, error) {
	var req anthropicRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return "", nil, nil, err
	}

	var messages []ChatMessage

	// 系统 prompt 作为第一个 system 消息放入
	if len(req.System) > 0 {
		sysPrompt := parseAnthropicContent(req.System)
		if sysPrompt != "" {
			messages = append(messages, ChatMessage{
				Role:    "system",
				Content: sysPrompt,
			})
		}
	}

	for _, m := range req.Messages {
		msg := ChatMessage{
			Role: m.Role,
		}

		// 检查 content 是否是 block 数组，如果是，提取文本与 tool_use/tool_result
		var blocks []map[string]any
		if err := json.Unmarshal(m.Content, &blocks); err == nil {
			var textParts []string
			var thinkingParts []string
			for _, block := range blocks {
				bType, _ := block["type"].(string)
				if bType == "text" {
					if text, ok := block["text"].(string); ok {
						textParts = append(textParts, text)
					}
				} else if bType == "thinking" {
					if thinking, ok := block["thinking"].(string); ok {
						thinkingParts = append(thinkingParts, thinking)
					}
				} else if bType == "tool_use" {
					id, _ := block["id"].(string)
					name, _ := block["name"].(string)
					input, _ := json.Marshal(block["input"])
					msg.ToolCalls = append(msg.ToolCalls, ToolCall{
						ID:        id,
						Type:      "function",
						Name:      name,
						Arguments: string(input),
					})
				} else if bType == "tool_result" {
					toolUseID, _ := block["tool_use_id"].(string)
					msg.ToolCallID = toolUseID
					// content block 可以继续嵌套
					if contentRaw, ok := block["content"]; ok {
						contentBytes, _ := json.Marshal(contentRaw)
						textParts = append(textParts, parseAnthropicContent(contentBytes))
					}
				}
			}
			msg.Content = strings.Join(textParts, "\n")
			msg.Thinking = strings.Join(thinkingParts, "\n")
		} else {
			// 普通字符串
			var str string
			if err := json.Unmarshal(m.Content, &str); err == nil {
				msg.Content = str
			} else {
				msg.Content = string(m.Content)
			}
		}

		messages = append(messages, msg)
	}

	var tools []ToolDef
	for _, t := range req.Tools {
		tools = append(tools, ToolDef{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  string(t.InputSchema),
		})
	}

	return req.Model, messages, tools, nil
}

func ParseAnthropicRequestWithHandles(body []byte) (RequestParseResult, error) {
	model, messages, tools, err := ParseAnthropicRequest(body)
	if err != nil {
		return RequestParseResult{}, err
	}
	return RequestParseResult{
		Model:    model,
		Messages: messages,
		Tools:    tools,
		Provider: "anthropic",
		Handles:  ExtractConversationHandles(body),
		Hints:    ExtractRequestHints(body),
	}, nil
}

// 解析 Anthropic 响应
func ParseAnthropicResponse(body []byte) (string, ChatMessage, TokenUsage, error) {
	var resp anthropicResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", ChatMessage{}, TokenUsage{}, err
	}

	var msg ChatMessage
	msg.Role = "assistant"

	var textParts []string
	var thinkingParts []string
	for _, b := range resp.Content {
		if b.Type == "text" {
			textParts = append(textParts, b.Text)
		} else if b.Type == "thinking" {
			thinkingParts = append(thinkingParts, b.Thinking)
		} else if b.Type == "tool_use" {
			msg.ToolCalls = append(msg.ToolCalls, ToolCall{
				ID:        b.ID,
				Type:      "function",
				Name:      b.Name,
				Arguments: string(b.Input),
			})
		}
	}
	msg.Content = strings.Join(textParts, "\n")
	msg.Thinking = strings.Join(thinkingParts, "\n")

	usage := TokenUsage{
		PromptTokens:        resp.Usage.InputTokens,
		CompletionTokens:    resp.Usage.OutputTokens,
		TotalTokens:         resp.Usage.InputTokens + resp.Usage.OutputTokens,
		CacheCreationTokens: resp.Usage.CacheCreationInputTokens,
		CacheReadTokens:     resp.Usage.CacheReadInputTokens,
	}

	return resp.Model, msg, usage, nil
}

func ParseAnthropicResponseWithHandles(body []byte) (ResponseParseResult, error) {
	model, msg, usage, err := ParseAnthropicResponse(body)
	if err != nil {
		return ResponseParseResult{}, err
	}
	return ResponseParseResult{
		Model:   model,
		Message: msg,
		Usage:   usage,
		Handles: ExtractResponseHandles(body),
	}, nil
}

// 增量解析 Anthropic SSE 流式数据
func ParseAnthropicStreamChunk(chunk []byte, currentResp *ChatMessage, currentUsage *TokenUsage, currentBlockIndex *int) error {
	line := string(chunk)
	if !strings.HasPrefix(line, "event: ") && !strings.HasPrefix(line, "data: ") {
		return nil
	}

	// 提取出 "data: " 那一行
	var dataStr string
	lines := strings.Split(line, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "data: ") {
			dataStr = strings.TrimPrefix(l, "data: ")
			break
		}
	}
	if dataStr == "" {
		return nil
	}

	var chunkObj struct {
		Type         string `json:"type"`
		Index        *int   `json:"index,omitempty"`
		ContentBlock *struct {
			Type     string          `json:"type"`
			Text     string          `json:"text,omitempty"`
			Thinking string          `json:"thinking,omitempty"`
			ID       string          `json:"id,omitempty"`
			Name     string          `json:"name,omitempty"`
			Input    json.RawMessage `json:"input,omitempty"`
		} `json:"content_block,omitempty"`
		Delta *struct {
			Type         string `json:"type"`
			Text         string `json:"text,omitempty"`
			Thinking     string `json:"thinking,omitempty"`
			PartialJson  string `json:"partial_json,omitempty"`
			InputTokens  int    `json:"input_tokens,omitempty"`
			OutputTokens int    `json:"output_tokens,omitempty"`
		} `json:"delta,omitempty"`
		Message *struct {
			Usage struct {
				InputTokens              int `json:"input_tokens"`
				OutputTokens             int `json:"output_tokens"`
				CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
				CacheReadInputTokens     int `json:"cache_read_input_tokens"`
			} `json:"usage"`
		} `json:"message,omitempty"`
		Usage *struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		} `json:"usage,omitempty"`
	}

	if err := json.Unmarshal([]byte(dataStr), &chunkObj); err != nil {
		return err
	}

	// 1. 处理 Token 计数
	if chunkObj.Message != nil {
		if chunkObj.Message.Usage.InputTokens > 0 {
			currentUsage.PromptTokens = chunkObj.Message.Usage.InputTokens
		}
		if chunkObj.Message.Usage.OutputTokens > 0 {
			currentUsage.CompletionTokens = chunkObj.Message.Usage.OutputTokens
		}
		if chunkObj.Message.Usage.CacheCreationInputTokens > 0 {
			currentUsage.CacheCreationTokens = chunkObj.Message.Usage.CacheCreationInputTokens
		}
		if chunkObj.Message.Usage.CacheReadInputTokens > 0 {
			currentUsage.CacheReadTokens = chunkObj.Message.Usage.CacheReadInputTokens
		}
		currentUsage.TotalTokens = currentUsage.PromptTokens + currentUsage.CompletionTokens
	}
	if chunkObj.Usage != nil {
		if chunkObj.Usage.InputTokens > 0 {
			currentUsage.PromptTokens = chunkObj.Usage.InputTokens
		}
		if chunkObj.Usage.OutputTokens > 0 {
			currentUsage.CompletionTokens = chunkObj.Usage.OutputTokens
		}
		if chunkObj.Usage.CacheCreationInputTokens > 0 {
			currentUsage.CacheCreationTokens = chunkObj.Usage.CacheCreationInputTokens
		}
		if chunkObj.Usage.CacheReadInputTokens > 0 {
			currentUsage.CacheReadTokens = chunkObj.Usage.CacheReadInputTokens
		}
		currentUsage.TotalTokens = currentUsage.PromptTokens + currentUsage.CompletionTokens
	}
	if chunkObj.Delta != nil {
		if chunkObj.Delta.InputTokens > 0 {
			currentUsage.PromptTokens = chunkObj.Delta.InputTokens
		}
		if chunkObj.Delta.OutputTokens > 0 {
			currentUsage.CompletionTokens = chunkObj.Delta.OutputTokens
		}
		if currentUsage.PromptTokens > 0 || currentUsage.CompletionTokens > 0 {
			currentUsage.TotalTokens = currentUsage.PromptTokens + currentUsage.CompletionTokens
		}
	}

	// 2. 处理 Block Start / Delta
	if chunkObj.Index != nil {
		*currentBlockIndex = *chunkObj.Index
	}

	if chunkObj.ContentBlock != nil {
		b := chunkObj.ContentBlock
		if b.Type == "tool_use" {
			for len(currentResp.ToolCalls) <= *currentBlockIndex {
				currentResp.ToolCalls = append(currentResp.ToolCalls, ToolCall{})
			}
			currentResp.ToolCalls[*currentBlockIndex] = ToolCall{
				ID:   b.ID,
				Type: "function",
				Name: b.Name,
			}
		}
	}

	if chunkObj.Delta != nil && chunkObj.Delta.Type != "" {
		d := chunkObj.Delta
		switch d.Type {
		case "text_delta":
			currentResp.Content += d.Text
		case "thinking_delta":
			currentResp.Thinking += d.Thinking
		case "input_json_delta":
			for len(currentResp.ToolCalls) <= *currentBlockIndex {
				currentResp.ToolCalls = append(currentResp.ToolCalls, ToolCall{})
			}
			currentResp.ToolCalls[*currentBlockIndex].Arguments += d.PartialJson
		}
	}

	currentResp.Role = "assistant"
	return nil
}

func parseAnthropicContent(raw json.RawMessage) string {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	var blocks []map[string]any
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return string(raw)
	}
	var lines []string
	for _, block := range blocks {
		blockType, _ := block["type"].(string)
		if blockType == "text" {
			if text, ok := block["text"].(string); ok && strings.TrimSpace(text) != "" {
				lines = append(lines, text)
			}
		}
	}
	return strings.Join(lines, "\n")
}

// 统一解析入口，方便代理过程调用
func ParseUnifiedRequest(path string, body []byte) (string, []ChatMessage, []ToolDef, string, error) {
	result, err := ParseUnifiedRequestEnvelope(path, body)
	if err != nil {
		return "", nil, nil, "", err
	}
	return result.Model, result.Messages, result.Tools, result.Provider, nil
}

func ParseUnifiedRequestEnvelope(path string, body []byte) (RequestParseResult, error) {
	if strings.HasSuffix(path, "/chat/completions") {
		return ParseOpenAIRequestWithHandles(body)
	} else if strings.HasSuffix(path, "/responses") {
		return ParseOpenAIResponsesRequest(body)
	} else if strings.HasSuffix(path, "/messages") {
		return ParseAnthropicRequestWithHandles(body)
	}
	return RequestParseResult{}, errors.New("unsupported endpoint path")
}

func ExtractRequestSessionID(body []byte) string {
	handles := ExtractConversationHandles(body)
	for _, handle := range handles {
		if handle.Kind == "session_id" {
			return handle.Value
		}
	}
	return ""
}

func ExtractConversationHandles(body []byte) []ConversationHandle {
	var envelope struct {
		PreviousResponseID string `json:"previous_response_id"`
		ConversationID     string `json:"conversation_id"`
		ThreadID           string `json:"thread_id"`
		SessionID          string `json:"session_id"`
		Metadata struct {
			SessionID string          `json:"session_id"`
			ConversationID string     `json:"conversation_id"`
			ThreadID string           `json:"thread_id"`
			UserID    json.RawMessage `json:"user_id"`
			ClientID  string          `json:"client_id"`
			DeviceID  string          `json:"device_id"`
			Custom    map[string]any  `json:"custom"`
		} `json:"metadata"`
		Intent string `json:"intent"` // 有些 Agent 将 Session 信息包含在 intent 或特定字段
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil
	}

	var handles []ConversationHandle
	appendHandle := func(kind, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		handles = append(handles, ConversationHandle{Kind: kind, Value: value})
	}

	appendHandle("session_id", envelope.Metadata.SessionID)
	appendHandle("session_id", envelope.SessionID)
	appendHandle("conversation_id", envelope.ConversationID)
	appendHandle("conversation_id", envelope.Metadata.ConversationID)
	appendHandle("thread_id", envelope.ThreadID)
	appendHandle("thread_id", envelope.Metadata.ThreadID)
	appendHandle("previous_response_id", envelope.PreviousResponseID)

	if envelope.Metadata.Custom != nil {
		if sid, ok := envelope.Metadata.Custom["session_id"].(string); ok {
			appendHandle("session_id", sid)
		}
		if sid, ok := envelope.Metadata.Custom["conversation_id"].(string); ok {
			appendHandle("conversation_id", sid)
		}
		if sid, ok := envelope.Metadata.Custom["thread_id"].(string); ok {
			appendHandle("thread_id", sid)
		}
		if sid, ok := envelope.Metadata.Custom["previous_response_id"].(string); ok {
			appendHandle("previous_response_id", sid)
		}
	}

	if len(envelope.Metadata.UserID) == 0 {
		return compactHandles(handles)
	}

	extractFromMap := func(value map[string]any) {
		if rawSessionID, ok := value["session_id"].(string); ok {
			appendHandle("session_id", rawSessionID)
		}
		if rawConversationID, ok := value["conversation_id"].(string); ok {
			appendHandle("conversation_id", rawConversationID)
		}
		if rawThreadID, ok := value["thread_id"].(string); ok {
			appendHandle("thread_id", rawThreadID)
		}
		if rawPreviousResponseID, ok := value["previous_response_id"].(string); ok {
			appendHandle("previous_response_id", rawPreviousResponseID)
		}
	}

	var userObject map[string]any
	if err := json.Unmarshal(envelope.Metadata.UserID, &userObject); err == nil {
		extractFromMap(userObject)
	}

	var userString string
	if err := json.Unmarshal(envelope.Metadata.UserID, &userString); err != nil || strings.TrimSpace(userString) == "" {
		return compactHandles(handles)
	}
	if err := json.Unmarshal([]byte(userString), &userObject); err != nil {
		return compactHandles(handles)
	}
	extractFromMap(userObject)

	return compactHandles(handles)
}

func ExtractRequestHints(body []byte) map[string]string {
	var envelope struct {
		Metadata struct {
			ClientID string          `json:"client_id"`
			DeviceID string          `json:"device_id"`
			UserID   json.RawMessage `json:"user_id"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil
	}
	hints := map[string]string{}
	if strings.TrimSpace(envelope.Metadata.ClientID) != "" {
		hints["client_id"] = strings.TrimSpace(envelope.Metadata.ClientID)
	}
	if strings.TrimSpace(envelope.Metadata.DeviceID) != "" {
		hints["device_id"] = strings.TrimSpace(envelope.Metadata.DeviceID)
	}
	var userObject map[string]any
	if err := json.Unmarshal(envelope.Metadata.UserID, &userObject); err == nil {
		if deviceID, ok := userObject["device_id"].(string); ok && strings.TrimSpace(deviceID) != "" {
			hints["device_id"] = strings.TrimSpace(deviceID)
		}
		if clientID, ok := userObject["client_id"].(string); ok && strings.TrimSpace(clientID) != "" {
			hints["client_id"] = strings.TrimSpace(clientID)
		}
	}
	if len(hints) == 0 {
		return nil
	}
	return hints
}

func ExtractResponseHandles(body []byte) []ConversationHandle {
	var envelope struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil
	}
	if strings.TrimSpace(envelope.ID) == "" {
		return nil
	}
	return []ConversationHandle{{Kind: "response_id", Value: strings.TrimSpace(envelope.ID)}}
}

func compactHandles(handles []ConversationHandle) []ConversationHandle {
	if len(handles) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var compacted []ConversationHandle
	for _, handle := range handles {
		if strings.TrimSpace(handle.Value) == "" || strings.TrimSpace(handle.Kind) == "" {
			continue
		}
		key := handle.Kind + "\x00" + handle.Value
		if seen[key] {
			continue
		}
		seen[key] = true
		compacted = append(compacted, ConversationHandle{Kind: strings.TrimSpace(handle.Kind), Value: strings.TrimSpace(handle.Value)})
	}
	if len(compacted) == 0 {
		return nil
	}
	return compacted
}

func ParseUnifiedResponse(path string, body []byte) (string, ChatMessage, TokenUsage, error) {
	result, err := ParseUnifiedResponseEnvelope(path, body)
	if err != nil {
		return "", ChatMessage{}, TokenUsage{}, err
	}
	return result.Model, result.Message, result.Usage, nil
}

func ParseUnifiedResponseEnvelope(path string, body []byte) (ResponseParseResult, error) {
	if strings.HasSuffix(path, "/chat/completions") || strings.HasSuffix(path, "/responses") {
		if strings.HasSuffix(path, "/responses") {
			return ParseOpenAIResponsesResponse(body)
		}
		model, msg, usage, err := ParseOpenAIResponse(body)
		if err != nil {
			return ResponseParseResult{}, err
		}
		return ResponseParseResult{Model: model, Message: msg, Usage: usage, Handles: ExtractResponseHandles(body)}, nil
	} else if strings.HasSuffix(path, "/messages") {
		return ParseAnthropicResponseWithHandles(body)
	}
	return ResponseParseResult{}, errors.New("unsupported endpoint path")
}
