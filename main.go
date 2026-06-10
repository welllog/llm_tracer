package main

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/getlantern/systray"
)

//go:embed static/*
var staticFS embed.FS

type AppConfig struct {
	ListenAddr             string `json:"listenAddr"`
	OpenaiBaseURL          string `json:"openaiBaseURL"`
	OpenaiAPIKey           string `json:"openaiAPIKey"`
	OpenaiResponsesBaseURL string `json:"openaiResponsesBaseURL"`
	OpenaiResponsesAPIKey  string `json:"openaiResponsesAPIKey"`
	AnthropicBaseURL       string `json:"anthropicBaseURL"`
	AnthropicAPIKey        string `json:"anthropicAPIKey"`
	DBPath                 string `json:"dbPath"`
}

type ProxyServer struct {
	config               AppConfig
	configLock           sync.RWMutex
	db                   *DBManager
	client               *http.Client
	proxyHTTPSrv         *http.Server
	proxyMu              sync.Mutex
	lastActiveSessionMap map[string]string // Client fingerprint -> SessionID
	activeMu             sync.RWMutex
	suspendedTools       []SuspendedTool
	suspMu               sync.Mutex
}

type SuspendedTool struct {
	LogID      int64
	SessionID  string
	ToolCallID string
	Arguments  string
	Provider   string
	ToolName   string
	ClientFingerprint string
	CreatedAt  time.Time
}

func semanticMatchScore(promptText, arguments string) float64 {
	promptText = strings.ToLower(promptText)
	arguments = strings.ToLower(arguments)

	if len(promptText) > 5 && (strings.Contains(arguments, promptText) || strings.Contains(promptText, arguments)) {
		shorterLen := len(promptText)
		if len(arguments) < shorterLen {
			shorterLen = len(arguments)
		}
		longerLen := len(promptText)
		if len(arguments) > longerLen {
			longerLen = len(arguments)
		}
		return 1 + float64(shorterLen)/float64(longerLen)
	}

	isStopWord := func(w string) bool {
		switch w {
		case "this", "that", "with", "from", "have", "some", "here", "there", "what", "your":
			return true
		}
		return false
	}

	tokenize := func(text string) map[string]bool {
		words := make(map[string]bool)
		f := func(c rune) bool {
			return !(unicode.IsLetter(c) || unicode.IsDigit(c))
		}
		parts := strings.FieldsFunc(text, f)
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if len(p) > 3 && !isStopWord(p) {
				words[p] = true
			}
		}
		return words
	}

	promptWords := tokenize(promptText)
	argWords := tokenize(arguments)

	if len(promptWords) == 0 || len(argWords) == 0 {
		return 0
	}

	matchCount := 0
	for w := range promptWords {
		if argWords[w] {
			matchCount++
		}
	}

	// 使用两方单词数量的较小值作为分母，防止因为子代理的超长指令提示词稀释了重合度
	minLen := len(promptWords)
	if len(argWords) < minLen {
		minLen = len(argWords)
	}

	ratio := float64(matchCount) / float64(minLen)
	if ratio >= 0.3 {
		return ratio
	}

	return 0
}

func semanticMatch(promptText, arguments string) bool {
	return semanticMatchScore(promptText, arguments) > 0
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "********"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func (s *ProxyServer) loadConfig() error {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	appDir := getAppDir()

	// 1. 设置默认值
	s.config = AppConfig{
		ListenAddr:             ":1238",
		OpenaiBaseURL:          "https://api.openai.com",
		OpenaiResponsesBaseURL: "https://api.openai.com",
		AnthropicBaseURL:       "https://api.anthropic.com",
		DBPath:                 filepath.Join(appDir, "llm_tracer.db"),
	}

	// 2. 从文件读取
	configPath := filepath.Join(appDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err == nil {
		var fileCfg AppConfig
		if err := json.Unmarshal(data, &fileCfg); err == nil {
			if fileCfg.ListenAddr != "" {
				s.config.ListenAddr = fileCfg.ListenAddr
			}
			if fileCfg.OpenaiBaseURL != "" {
				s.config.OpenaiBaseURL = fileCfg.OpenaiBaseURL
			}
			s.config.OpenaiAPIKey = fileCfg.OpenaiAPIKey
			if fileCfg.OpenaiResponsesBaseURL != "" {
				s.config.OpenaiResponsesBaseURL = fileCfg.OpenaiResponsesBaseURL
			}
			s.config.OpenaiResponsesAPIKey = fileCfg.OpenaiResponsesAPIKey
			if fileCfg.AnthropicBaseURL != "" {
				s.config.AnthropicBaseURL = fileCfg.AnthropicBaseURL
			}
			s.config.AnthropicAPIKey = fileCfg.AnthropicAPIKey
			if fileCfg.DBPath != "" {
				s.config.DBPath = fileCfg.DBPath
			}
		}
	}

	// 3. 环境变量覆盖
	if env := os.Getenv("LISTEN_ADDR"); env != "" {
		s.config.ListenAddr = env
	}
	if env := os.Getenv("OPENAI_BASE_URL"); env != "" {
		s.config.OpenaiBaseURL = env
	}
	if env := os.Getenv("OPENAI_API_KEY"); env != "" {
		s.config.OpenaiAPIKey = env
	}
	if env := os.Getenv("ANTHROPIC_BASE_URL"); env != "" {
		s.config.AnthropicBaseURL = env
	}
	if env := os.Getenv("ANTHROPIC_API_KEY"); env != "" {
		s.config.AnthropicAPIKey = env
	}
	if env := os.Getenv("DB_PATH"); env != "" {
		s.config.DBPath = env
	}

	return nil
}

func (s *ProxyServer) saveConfig(newCfg AppConfig) error {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	// 比较并保留被掩码的 API Key
	if strings.Contains(newCfg.OpenaiAPIKey, "...") || newCfg.OpenaiAPIKey == "********" {
		newCfg.OpenaiAPIKey = s.config.OpenaiAPIKey
	}
	if strings.Contains(newCfg.OpenaiResponsesAPIKey, "...") || newCfg.OpenaiResponsesAPIKey == "********" {
		newCfg.OpenaiResponsesAPIKey = s.config.OpenaiResponsesAPIKey
	}
	if strings.Contains(newCfg.AnthropicAPIKey, "...") || newCfg.AnthropicAPIKey == "********" {
		newCfg.AnthropicAPIKey = s.config.AnthropicAPIKey
	}

	s.config = newCfg

	data, err := json.MarshalIndent(newCfg, "", "  ")
	if err != nil {
		return err
	}

	appDir := getAppDir()
	configPath := filepath.Join(appDir, "config.json")
	return os.WriteFile(configPath, data, 0644)
}

func (s *ProxyServer) startProxyServiceLocked() error {
	s.configLock.RLock()
	addr := s.config.ListenAddr
	s.configLock.RUnlock()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/chat/completions", s.handleProxyOpenAI)
	mux.HandleFunc("POST /v1/responses", s.handleProxyOpenAIResponses)
	mux.HandleFunc("POST /v1/messages", s.handleProxyAnthropic)

	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Api-Key, anthropic-version, anthropic-beta")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})

	s.proxyHTTPSrv = &http.Server{
		Addr:    addr,
		Handler: proxyHandler,
	}

	log.Printf("LLM Proxy service starting on %s", addr)
	go func() {
		if err := s.proxyHTTPSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("LLM Proxy service exited with error: %v", err)
		}
	}()

	return nil
}

func (s *ProxyServer) stopProxyServiceLocked() error {
	if s.proxyHTTPSrv != nil {
		log.Printf("Stopping LLM Proxy service on %s", s.proxyHTTPSrv.Addr)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.proxyHTTPSrv.Shutdown(ctx)
		s.proxyHTTPSrv = nil
		return err
	}
	return nil
}

func (s *ProxyServer) startProxyService() error {
	s.proxyMu.Lock()
	defer s.proxyMu.Unlock()
	return s.startProxyServiceLocked()
}

func (s *ProxyServer) stopProxyService() error {
	s.proxyMu.Lock()
	defer s.proxyMu.Unlock()
	return s.stopProxyServiceLocked()
}

func (s *ProxyServer) restartProxyService() error {
	s.proxyMu.Lock()
	defer s.proxyMu.Unlock()

	if err := s.stopProxyServiceLocked(); err != nil {
		log.Printf("Error stopping proxy service during restart: %v", err)
	}
	return s.startProxyServiceLocked()
}

var (
	globalServer     *ProxyServer
	globalConsoleURL string
	consoleListener  net.Listener
	consoleSrv       *http.Server
)

// 获取配置和数据库隐藏目录 ~/.llm_tracer
func getAppDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	dir := filepath.Join(home, ".llm_tracer")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// 探测并分配空闲端口
func findFreePort(addr string) (string, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		if strings.HasPrefix(addr, ":") {
			portStr = addr[1:]
			host = ""
		} else {
			return addr, err
		}
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return addr, err
	}

	for i := 0; i < 100; i++ {
		testPort := port + i
		testAddr := net.JoinHostPort(host, strconv.Itoa(testPort))
		ln, err := net.Listen("tcp", testAddr)
		if err == nil {
			ln.Close()
			return testAddr, nil
		}
	}
	return addr, fmt.Errorf("no free port found starting from %d", port)
}

// 跨平台打开默认浏览器
func openBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", etc.
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()
}

// 生成系统托盘的默认图标 (16x16 PNG)
func generateDefaultIcon() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	// 填充背景：深灰色
	bg := color.RGBA{40, 44, 52, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	// 在中间画一个绿色的小方块（代表服务正常）
	green := color.RGBA{16, 185, 129, 255} // Emerald 500
	for x := 5; x <= 10; x++ {
		for y := 5; y <= 10; y++ {
			img.Set(x, y, green)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

func onReady() {
	icon := generateDefaultIcon()
	if icon != nil {
		systray.SetIcon(icon)
	}
	systray.SetTitle("LLM Tracer")
	systray.SetTooltip("LLM Tracer Proxy & Recorder")

	mOpen := systray.AddMenuItem("打开控制台", "打开 LLM Tracer 控制台页面")
	mRestart := systray.AddMenuItem("重启代理服务", "重启底层 LLM 代理拦截服务")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "关闭服务并退出")

	// 自动打开默认浏览器
	go func() {
		time.Sleep(500 * time.Millisecond) // 等待监听就绪
		_ = openBrowser(globalConsoleURL)
	}()

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				_ = openBrowser(globalConsoleURL)
			case <-mRestart.ClickedCh:
				log.Println("Restarting proxy service from tray...")
				if err := globalServer.restartProxyService(); err != nil {
					log.Printf("Failed to restart proxy service: %v", err)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	log.Println("Stopping services and exiting...")
	if consoleSrv != nil {
		_ = consoleSrv.Close()
	}
	if globalServer != nil {
		_ = globalServer.stopProxyService()
		if globalServer.db != nil {
			globalServer.db.Close()
		}
	}
}

func main() {
	// 允许用户在命令行设置监听地址
	addrFlag := flag.String("listen", "", "override proxy listen address")
	consoleFlag := flag.String("console", ":56129", "override console listen address")
	flag.Parse()

	// 确保静态目录存在，以便 go build 即使没有构建前端时也不会因为 embed 找不到目录而报错。
	_ = os.MkdirAll("static", 0755)
	// 如果 static 下面是空的，写入一个占位文件，以便 go:embed static/* 不会报错
	placeholderPath := filepath.Join("static", "placeholder.txt")
	if _, err := os.Stat(placeholderPath); os.IsNotExist(err) {
		_ = os.WriteFile(placeholderPath, []byte("Vite frontend assets will go here"), 0644)
	}

	server := &ProxyServer{
		client: &http.Client{
			Timeout: 180 * time.Second,
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				DialContext:         (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
				ForceAttemptHTTP2:   true,
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		lastActiveSessionMap: make(map[string]string),
	}
	globalServer = server

	if err := server.loadConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if *addrFlag != "" {
		server.config.ListenAddr = *addrFlag
	}

	// 自动探测并分配空闲端口
	freeListenAddr, err := findFreePort(server.config.ListenAddr)
	if err == nil {
		if freeListenAddr != server.config.ListenAddr {
			log.Printf("Proxy port %s occupied, fallback to %s", server.config.ListenAddr, freeListenAddr)
		}
		server.config.ListenAddr = freeListenAddr
	}

	freeConsoleAddr, err := findFreePort(*consoleFlag)
	if err != nil {
		log.Fatalf("failed to find free console port: %v", err)
	}

	// 初始化 SQLite
	dbMgr, err := InitDB(server.config.DBPath)
	if err != nil {
		log.Fatalf("failed to initialize SQLite: %v", err)
	}
	server.db = dbMgr

	// 1. 启动代理服务
	if err := server.startProxyService(); err != nil {
		log.Fatalf("failed to start proxy service: %v", err)
	}

	// 2. 配置控制台服务
	consoleMux := http.NewServeMux()
	consoleMux.HandleFunc("GET /api/config", server.handleGetConfig)
	consoleMux.HandleFunc("POST /api/config", server.handleSetConfig)
	consoleMux.HandleFunc("GET /api/logs", server.handleGetLogs)
	consoleMux.HandleFunc("GET /api/logs/{id}", server.handleGetLogDetail)
	consoleMux.HandleFunc("GET /api/sessions/{id}/logs", server.handleGetSessionLogs)
	consoleMux.HandleFunc("GET /api/stats", server.handleGetStats)
	consoleMux.HandleFunc("DELETE /api/logs/{id}", server.handleDeleteLog)
	consoleMux.HandleFunc("DELETE /api/sessions/{id}/logs", server.handleDeleteSessionLogs)

	// 静态前端路由
	var staticHandler http.Handler
	subFS, err := fs.Sub(staticFS, "static")
	if err == nil {
		fileServer := http.FileServer(http.FS(subFS))
		staticHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 支持 SPA 路由回落，如果访问的不是静态文件，返回 index.html
			if r.URL.Path != "/" && !strings.Contains(r.URL.Path, ".") {
				r.URL.Path = "/"
			}
			fileServer.ServeHTTP(w, r)
		})
	} else {
		staticHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Waiting for Vite build in 'static' directory..."))
		})
	}

	// 注入 CORS 和通用头部中间件
	consoleHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Api-Key, anthropic-version, anthropic-beta")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/") {
			consoleMux.ServeHTTP(w, r)
		} else {
			staticHandler.ServeHTTP(w, r)
		}
	})

	consoleListener, err = net.Listen("tcp", freeConsoleAddr)
	if err != nil {
		log.Fatalf("failed to listen on console address %s: %v", freeConsoleAddr, err)
	}

	_, consolePort, _ := net.SplitHostPort(consoleListener.Addr().String())
	globalConsoleURL = fmt.Sprintf("http://localhost:%s", consolePort)

	consoleSrv = &http.Server{
		Handler: consoleHandler,
	}

	go func() {
		log.Printf("LLM Tracer Console started on %s", globalConsoleURL)
		log.Printf("SQLite database located at %s", server.config.DBPath)
		if err := consoleSrv.Serve(consoleListener); err != nil && err != http.ErrServerClosed {
			log.Printf("console server exited with error: %v", err)
		}
	}()

	// 启动系统托盘
	systray.Run(onReady, onExit)
}

// (using standard io/fs instead of custom struct)

// ==========================================
// 控制台端点处理器
// ==========================================

func (s *ProxyServer) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	s.configLock.RLock()
	cfg := s.config
	s.configLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cfg)
}

func (s *ProxyServer) handleSetConfig(w http.ResponseWriter, r *http.Request) {
	var cfg AppConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	s.configLock.RLock()
	oldListenAddr := s.config.ListenAddr
	s.configLock.RUnlock()

	if err := s.saveConfig(cfg); err != nil {
		http.Error(w, fmt.Sprintf("failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	if cfg.ListenAddr != oldListenAddr {
		go func() {
			if err := s.restartProxyService(); err != nil {
				log.Printf("Failed to restart proxy service: %v", err)
			}
		}()
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *ProxyServer) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	provider := q.Get("provider")
	model := q.Get("model")
	keyword := q.Get("keyword")

	var statusFilter *int
	if statusStr := q.Get("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			statusFilter = &status
		}
	}

	excludeBranches := true
	if ebStr := q.Get("exclude_branches"); ebStr != "" {
		excludeBranches = ebStr == "true"
	}

	summaries, total, err := s.db.GetLogs(page, pageSize, provider, model, keyword, statusFilter, excludeBranches)
	if err != nil {
		http.Error(w, fmt.Errorf("query logs: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"list":     summaries,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (s *ProxyServer) handleGetSessionLogs(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		http.Error(w, "missing session id", http.StatusBadRequest)
		return
	}

	logs, err := s.db.GetSessionLogs(sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("get session logs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(logs)
}

func (s *ProxyServer) handleGetLogDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid log id", http.StatusBadRequest)
		return
	}

	detail, err := s.db.GetLogDetail(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("get log detail: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(detail)
}

func (s *ProxyServer) handleGetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("get stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}

func (s *ProxyServer) handleDeleteLog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid log id", http.StatusBadRequest)
		return
	}

	if err := s.db.DeleteLog(id); err != nil {
		http.Error(w, fmt.Sprintf("delete log: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *ProxyServer) handleDeleteSessionLogs(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		http.Error(w, "missing session id", http.StatusBadRequest)
		return
	}

	if err := s.db.DeleteSessionLogs(sessionID); err != nil {
		http.Error(w, fmt.Sprintf("delete session logs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// ==========================================
// 代理逻辑处理器
// ==========================================

func (s *ProxyServer) handleProxyOpenAI(w http.ResponseWriter, r *http.Request) {
	s.configLock.RLock()
	upstreamURL := s.config.OpenaiBaseURL
	apiKey := s.config.OpenaiAPIKey
	s.configLock.RUnlock()

	s.proxyAPI(w, r, "openai", upstreamURL, "/v1/chat/completions", apiKey)
}

func (s *ProxyServer) handleProxyOpenAIResponses(w http.ResponseWriter, r *http.Request) {
	s.configLock.RLock()
	upstreamURL := s.config.OpenaiResponsesBaseURL
	apiKey := s.config.OpenaiResponsesAPIKey
	// 如果没有配置，回退到 OpenAI 配置
	if upstreamURL == "" {
		upstreamURL = s.config.OpenaiBaseURL
	}
	if apiKey == "" {
		apiKey = s.config.OpenaiAPIKey
	}
	s.configLock.RUnlock()

	s.proxyAPI(w, r, "openai-responses", upstreamURL, "/v1/responses", apiKey)
}

func (s *ProxyServer) handleProxyAnthropic(w http.ResponseWriter, r *http.Request) {
	s.configLock.RLock()
	upstreamURL := s.config.AnthropicBaseURL
	apiKey := s.config.AnthropicAPIKey
	s.configLock.RUnlock()

	s.proxyAPI(w, r, "anthropic", upstreamURL, "/v1/messages", apiKey)
}

func (s *ProxyServer) proxyAPI(w http.ResponseWriter, r *http.Request, provider, upstreamBaseURL, pathSuffix, apiKey string) {
	startedAt := time.Now()

	// 1. 读取 Request Body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	// 2. 解析请求
	requestInfo := RequestParseResult{Provider: provider}
	parseResult, parseErr := ParseUnifiedRequestEnvelope(pathSuffix, reqBody)
	if parseErr != nil {
		log.Printf("Failed to parse request payload: %v", parseErr)
	} else {
		requestInfo = parseResult
	}

	resolvedProvider := provider
	if requestInfo.Provider != "" {
		resolvedProvider = requestInfo.Provider
	}
	model := requestInfo.Model
	promptMsgs := requestInfo.Messages
	tools := requestInfo.Tools
	sessionID, parentID, parentToolCallID, clientFingerprint := s.resolveSessionAndParent(r, resolvedProvider, promptMsgs, requestInfo.Handles, requestInfo.Hints)

	// 3. 构建上游请求
	targetURL := joinUpstreamURL(upstreamBaseURL, pathSuffix)
	upstreamReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, targetURL, bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "failed to create upstream request", http.StatusInternalServerError)
		return
	}

	// 复制头信息并剔除特定代理字段
	for key, values := range r.Header {
		if strings.EqualFold(key, "Connection") || strings.EqualFold(key, "Authorization") || strings.EqualFold(key, "X-Api-Key") {
			continue
		}
		for _, value := range values {
			upstreamReq.Header.Add(key, value)
		}
	}

	// 鉴权设置
	if provider == "anthropic" {
		upstreamReq.Header.Set("X-Api-Key", apiKey)
	} else {
		upstreamReq.Header.Set("Authorization", "Bearer "+apiKey)
	}
	upstreamReq.Header.Set("Content-Type", "application/json")

	// 4. 执行请求
	upstreamResp, err := s.client.Do(upstreamReq)
	if err != nil {
		// 写入一条错误日志
		s.logErrorExchange(resolvedProvider, model, pathSuffix, reqBody, err, time.Since(startedAt), sessionID, parentID, parentToolCallID, clientFingerprint)
		http.Error(w, fmt.Sprintf("upstream error: %v", err), http.StatusBadGateway)
		return
	}
	defer upstreamResp.Body.Close()

	// 拷贝上游响应头
	for key, values := range upstreamResp.Header {
		if strings.EqualFold(key, "Connection") || strings.EqualFold(key, "Content-Encoding") || strings.EqualFold(key, "Content-Length") {
			continue
		}
		w.Header().Del(key)
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 自动处理解压
	var bodyReader io.ReadCloser = upstreamResp.Body
	encoding := strings.ToLower(upstreamResp.Header.Get("Content-Encoding"))
	switch encoding {
	case "gzip":
		gr, err := gzip.NewReader(upstreamResp.Body)
		if err == nil {
			bodyReader = gr
		}
	case "deflate":
		bodyReader = flate.NewReader(upstreamResp.Body)
	}
	defer bodyReader.Close()

	// 判断是否是流式
	isStream := strings.Contains(strings.ToLower(upstreamResp.Header.Get("Content-Type")), "text/event-stream")

	w.WriteHeader(upstreamResp.StatusCode)

	var finalResponseMsg ChatMessage
	var finalUsage TokenUsage
	var responseHandles []ConversationHandle
	var rawResponseBuffer bytes.Buffer

	if isStream {
		// 流式代理：采用按行扫描的方式来解析 SSE 事件，并实时刷回给客户端
		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("ResponseWriter is not a Flusher, stream fallback")
		}

		writer := w
		reader := bufio.NewReader(bodyReader)
		blockIndex := 0

		for {
			line, err := reader.ReadBytes('\n')
			if len(line) > 0 {
				// 实时刷给客户端
				_, _ = writer.Write(line)
				if flusher != nil {
					flusher.Flush()
				}
				// 收集原始响应
				rawResponseBuffer.Write(line)

				// 解析该行
				trimmed := bytes.TrimSpace(line)
				if len(trimmed) > 0 {
					if strings.HasSuffix(pathSuffix, "/messages") {
						_ = ParseAnthropicStreamChunk(trimmed, &finalResponseMsg, &finalUsage, &blockIndex)
					} else if strings.HasSuffix(pathSuffix, "/responses") {
						_ = ParseOpenAIResponsesStreamChunk(trimmed, &finalResponseMsg, &finalUsage, &responseHandles)
					} else {
						_ = ParseOpenAIStreamChunk(trimmed, &finalResponseMsg, &finalUsage)
					}
				}
			}
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Printf("Error reading stream chunk: %v", err)
				}
				break
			}
		}
	} else {
		// 普通响应：一次性读取
		respBody, err := io.ReadAll(bodyReader)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return
		}

		// 写给客户端
		_, _ = w.Write(respBody)
		rawResponseBuffer.Write(respBody)

		// 解析
		responseParseResult, responseParseErr := ParseUnifiedResponseEnvelope(pathSuffix, respBody)
		if responseParseErr != nil {
			log.Printf("Failed to parse upstream response payload: %v", responseParseErr)
		} else {
			finalResponseMsg = responseParseResult.Message
			finalUsage = responseParseResult.Usage
			responseHandles = responseParseResult.Handles
			if model == "" && responseParseResult.Model != "" {
				model = responseParseResult.Model
			}
		}
	}

	// 5. 组装并写入日志记录数据库
	duration := time.Since(startedAt)

	// 过滤工具调用空槽
	var cleanToolCalls []ToolCall
	for _, tc := range finalResponseMsg.ToolCalls {
		if tc.Name != "" || tc.ID != "" {
			cleanToolCalls = append(cleanToolCalls, tc)
		}
	}
	finalResponseMsg.ToolCalls = cleanToolCalls

	logRecord := &UnifiedLog{
		Provider:            resolvedProvider,
		Model:               model,
		Path:                pathSuffix,
		StatusCode:          upstreamResp.StatusCode,
		DurationMs:          duration.Milliseconds(),
		InputTokens:         finalUsage.PromptTokens,
		OutputTokens:        finalUsage.CompletionTokens,
		TotalTokens:         finalUsage.TotalTokens,
		CachedTokens:        finalUsage.CachedTokens,
		CacheReadTokens:     finalUsage.CacheReadTokens,
		CacheCreationTokens: finalUsage.CacheCreationTokens,
		Prompt:              promptMsgs,
		Response:            finalResponseMsg,
		Tools:               tools,
		RawRequest:          string(reqBody),
		RawResponse:         rawResponseBuffer.String(),
		SessionID:           sessionID,
		ParentID:            parentID,
		ParentToolCallID:    parentToolCallID,
		ClientFingerprint:   clientFingerprint,
		RequestHandles:      requestInfo.Handles,
		ResponseHandles:     responseHandles,
	}

	// 如果在上游没抓到 usage，如果是 openai-responses 或特殊情况，我们可以在解析出的模型响应中做修正
	if logRecord.TotalTokens == 0 {
		logRecord.InputTokens = len(promptMsgs) * 5
		logRecord.TotalTokens = logRecord.InputTokens + logRecord.OutputTokens
	}

	insertID, dbErr := s.db.InsertLog(logRecord)
	if dbErr != nil {
		log.Printf("Database insertion error: %v", dbErr)
	} else if insertID > 0 && len(logRecord.Response.ToolCalls) > 0 {
		s.suspMu.Lock()
		for _, tc := range logRecord.Response.ToolCalls {
			s.suspendedTools = append(s.suspendedTools, SuspendedTool{
				LogID:             insertID,
				SessionID:         logRecord.SessionID,
				ToolCallID:        tc.ID,
				Arguments:         tc.Arguments,
				Provider:          logRecord.Provider,
				ToolName:          tc.Name,
				ClientFingerprint: clientFingerprint,
				CreatedAt:         time.Now(),
			})
		}
		s.suspMu.Unlock()
	}
}

func (s *ProxyServer) logErrorExchange(provider, model, path string, reqBody []byte, err error, duration time.Duration, sessionID string, parentID *int64, parentToolCallID, clientFingerprint string) {
	requestInfo := RequestParseResult{Provider: provider}
	if parsed, parseErr := ParseUnifiedRequestEnvelope(path, reqBody); parseErr == nil {
		requestInfo = parsed
		if model == "" {
			model = parsed.Model
		}
		if parsed.Provider != "" {
			provider = parsed.Provider
		}
	}
	logRecord := &UnifiedLog{
		Provider:         provider,
		Model:            model,
		Path:             path,
		StatusCode:       http.StatusBadGateway,
		DurationMs:       duration.Milliseconds(),
		Prompt:           requestInfo.Messages,
		Tools:            requestInfo.Tools,
		RawRequest:       string(reqBody),
		ErrorMessage:     err.Error(),
		SessionID:        sessionID,
		ParentID:         parentID,
		ParentToolCallID: parentToolCallID,
		ClientFingerprint: clientFingerprint,
		RequestHandles:   requestInfo.Handles,
	}
	_, _ = s.db.InsertLog(logRecord)
}

func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func isSideTask(prompt []ChatMessage) bool {
	if len(prompt) == 0 {
		return false
	}
	// 检查最后一条消息的内容
	lastMsg := prompt[len(prompt)-1]
	content := strings.ToLower(lastMsg.Content)

	keywords := []string{
		"summarize", "总结",
		"predict", "预测",
		"next step", "下一步",
		"generate title", "生成标题",
		"generate a short title",
		"session title", "会话标题",
		"chat_name", "naming", "命名",
	}

	for _, kw := range keywords {
		if strings.Contains(content, kw) {
			return true
		}
	}
	return false
}

func cleanDynamicContent(content string) string {
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	for _, line := range lines {
		lowerLine := strings.ToLower(line)
		// 过滤掉包含 x-anthropic-billing-header 的那一行
		if strings.Contains(lowerLine, "x-anthropic-billing-header:") {
			continue
		}
		// 过滤常见的动态注入信息
		if strings.Contains(lowerLine, "current date:") ||
			strings.Contains(lowerLine, "current time:") ||
			strings.Contains(lowerLine, "today is") ||
			strings.Contains(lowerLine, "working directory:") ||
			strings.Contains(lowerLine, "cwd:") ||
			strings.Contains(lowerLine, "current file:") ||
			strings.Contains(lowerLine, "active file:") ||
			strings.Contains(lowerLine, "opened files:") ||
			strings.Contains(lowerLine, "filename:") {
			continue
		}
		// 某些框架会注入包含 UUID 或临时路径的行
		if strings.Contains(lowerLine, "/var/folders/") || strings.Contains(lowerLine, "/tmp/") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}
	return strings.TrimSpace(strings.Join(cleanedLines, "\n"))
}

// cleanMessages 对不同接口的系统消息或内容进行策略清洗
func cleanMessages(provider string, messages []ChatMessage) []ChatMessage {
	var cleaned []ChatMessage
	for _, msg := range messages {
		role := strings.ToLower(msg.Role)
		cleanedContent := msg.Content

		if role == "system" {
			cleanedContent = cleanDynamicContent(msg.Content)
		} else if role == "user" {
			if simplified := simplifyUserPrompt(msg.Content); strings.TrimSpace(simplified) != "" {
				cleanedContent = simplified
			}
		}

		cleaned = append(cleaned, ChatMessage{
			Role:       msg.Role,
			Content:    cleanedContent,
			Name:       msg.Name,
			Thinking:   msg.Thinking,
			ToolCalls:  msg.ToolCalls,
			ToolCallID: msg.ToolCallID,
		})
	}
	return cleaned
}

func messagesEqual(provider string, a, b []ChatMessage, strict bool) bool {
	aCleaned := cleanMessages(provider, a)
	bCleaned := cleanMessages(provider, b)

	if !strict {
		// 模糊比对：过滤掉所有 system 消息再比对
		// 这样可以解决 Agent 在不同轮次注入不同 Context/System 指示带来的匹配失败
		filterNonSystem := func(msgs []ChatMessage) []ChatMessage {
			var res []ChatMessage
			for _, m := range msgs {
				if strings.ToLower(m.Role) != "system" {
					res = append(res, m)
				}
			}
			return res
		}
		aCleaned = filterNonSystem(aCleaned)
		bCleaned = filterNonSystem(bCleaned)
	}

	if len(aCleaned) != len(bCleaned) {
		return false
	}
	for i := range aCleaned {
		if !strings.EqualFold(aCleaned[i].Role, bCleaned[i].Role) {
			return false
		}
		if strings.TrimSpace(aCleaned[i].ToolCallID) != strings.TrimSpace(bCleaned[i].ToolCallID) {
			return false
		}
		if strings.TrimSpace(aCleaned[i].Content) != strings.TrimSpace(bCleaned[i].Content) {
			return false
		}
		if len(aCleaned[i].ToolCalls) != len(bCleaned[i].ToolCalls) {
			return false
		}
		for j := range aCleaned[i].ToolCalls {
			if strings.TrimSpace(aCleaned[i].ToolCalls[j].Name) != strings.TrimSpace(bCleaned[i].ToolCalls[j].Name) {
				return false
			}
			if strings.TrimSpace(aCleaned[i].ToolCalls[j].ID) != strings.TrimSpace(bCleaned[i].ToolCalls[j].ID) {
				return false
			}
		}
	}
	return true
}

func buildClientFingerprint(r *http.Request, provider string, hints map[string]string) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = strings.TrimSpace(r.RemoteAddr)
	}
	parts := []string{provider}
	if userAgent := strings.TrimSpace(strings.ToLower(r.UserAgent())); userAgent != "" {
		parts = append(parts, "ua="+userAgent)
	}
	if clientID := strings.TrimSpace(hints["client_id"]); clientID != "" {
		parts = append(parts, "client="+clientID)
	}
	if deviceID := strings.TrimSpace(hints["device_id"]); deviceID != "" {
		parts = append(parts, "device="+deviceID)
	}
	if len(parts) == 1 {
		if remote := strings.TrimSpace(r.RemoteAddr); remote != "" {
			parts = append(parts, "remote="+remote)
		}
		return strings.Join(parts, "|")
	}
	if host != "" {
		parts = append(parts, "host="+host)
	}
	return strings.Join(parts, "|")
}

func findConversationHandle(handles []ConversationHandle, kind string) (ConversationHandle, bool) {
	for _, handle := range handles {
		if handle.Kind == kind && strings.TrimSpace(handle.Value) != "" {
			return handle, true
		}
	}
	return ConversationHandle{}, false
}

func selectStableSessionHandle(handles []ConversationHandle) (ConversationHandle, bool) {
	for _, kind := range []string{"session_id", "conversation_id", "thread_id"} {
		if handle, ok := findConversationHandle(handles, kind); ok {
			return handle, true
		}
	}
	return ConversationHandle{}, false
}

func canonicalSessionIDForHandle(provider string, handle ConversationHandle) string {
	if handle.Kind == "session_id" {
		return handle.Value
	}
	return provider + ":" + handle.Kind + ":" + handle.Value
}

func seedSessionIDFromPreviousResponse(provider string, handle ConversationHandle) string {
	return provider + ":response-chain:" + handle.Value
}

func extractPrimaryUserPrompt(prompt []ChatMessage) string {
	for _, msg := range prompt {
		if strings.ToLower(msg.Role) != "user" {
			continue
		}
		candidate := simplifyUserPrompt(strings.TrimSpace(msg.Content))
		if strings.TrimSpace(candidate) != "" {
			return candidate
		}
		if strings.TrimSpace(msg.Content) != "" {
			return strings.TrimSpace(msg.Content)
		}
	}
	return ""
}

func isToolResultRole(role string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	return role == "tool" || role == "function"
}

func isToolResultSuffix(messages []ChatMessage) bool {
	if len(messages) == 0 {
		return false
	}
	for _, msg := range messages {
		if !isToolResultRole(msg.Role) {
			return false
		}
	}
	return true
}

func toolResultSuffixMatches(toolCalls []ToolCall, suffix []ChatMessage) bool {
	if len(toolCalls) == 0 || len(suffix) == 0 || !isToolResultSuffix(suffix) {
		return false
	}

	toolCallIDs := make(map[string]bool)
	for _, toolCall := range toolCalls {
		if strings.TrimSpace(toolCall.ID) != "" {
			toolCallIDs[strings.TrimSpace(toolCall.ID)] = true
		}
	}
	if len(toolCallIDs) == 0 {
		return true
	}

	matched := false
	for _, msg := range suffix {
		toolCallID := strings.TrimSpace(msg.ToolCallID)
		if toolCallID == "" {
			continue
		}
		if !toolCallIDs[toolCallID] {
			return false
		}
		matched = true
	}
	return matched
}

func (s *ProxyServer) matchHistory(provider, clientFingerprint string, prompt []ChatMessage) (string, *int64, bool) {
	if len(prompt) <= 1 {
		return "", nil, false
	}

	historyPrefix := prompt[:len(prompt)-1]
	tryMatch := func(recentLogs []RecentLogForMatch, allowFuzzy bool) (string, *int64, bool) {
		matchRecord := func(logRecord RecentLogForMatch, strict bool, candidate []ChatMessage) (string, *int64, bool) {
			var pMsgs []ChatMessage
			var rMsg ChatMessage
			if logRecord.PromptJSON != "" {
				_ = json.Unmarshal([]byte(logRecord.PromptJSON), &pMsgs)
			}
			if logRecord.ResponseJSON != "" {
				_ = json.Unmarshal([]byte(logRecord.ResponseJSON), &rMsg)
			}
			fullHistory := append(pMsgs, rMsg)

			if len(candidate) == len(fullHistory) && messagesEqual(provider, fullHistory, candidate, strict) {
				sessionID := logRecord.SessionID
				if sessionID == "" {
					sessionID = generateUUID()
				}
				parentID := logRecord.ID
				return sessionID, &parentID, true
			}

			if len(prompt) > len(fullHistory) && strings.ToLower(rMsg.Role) == "assistant" && len(rMsg.ToolCalls) > 0 {
				prefix := prompt[:len(fullHistory)]
				suffix := prompt[len(fullHistory):]
				if messagesEqual(provider, fullHistory, prefix, strict) && toolResultSuffixMatches(rMsg.ToolCalls, suffix) {
					sessionID := logRecord.SessionID
					if sessionID == "" {
						sessionID = generateUUID()
					}
					parentID := logRecord.ID
					return sessionID, &parentID, true
				}
			}

			return "", nil, false
		}

		for _, logRecord := range recentLogs {
			if sessionID, parentID, ok := matchRecord(logRecord, true, historyPrefix); ok {
				return sessionID, parentID, true
			}
		}

		if !allowFuzzy {
			return "", nil, false
		}

		for _, logRecord := range recentLogs {
			var pMsgs []ChatMessage
			var rMsg ChatMessage
			if logRecord.PromptJSON != "" {
				_ = json.Unmarshal([]byte(logRecord.PromptJSON), &pMsgs)
			}
			if logRecord.ResponseJSON != "" {
				_ = json.Unmarshal([]byte(logRecord.ResponseJSON), &rMsg)
			}
			fullHistory := append(pMsgs, rMsg)

			if sessionID, parentID, ok := matchRecord(logRecord, false, historyPrefix); ok {
				hasNonSystem := false
				for _, m := range fullHistory {
					if strings.ToLower(m.Role) != "system" {
						hasNonSystem = true
						break
					}
				}
				if hasNonSystem {
					return sessionID, parentID, true
				}
			}
		}

		return "", nil, false
	}

	recentLogs, err := s.db.GetRecentLogsForMatch(50, provider, clientFingerprint)
	if err == nil {
		if sessionID, parentID, ok := tryMatch(recentLogs, true); ok {
			return sessionID, parentID, true
		}
	} else {
		log.Printf("Failed to load recent scoped logs for session matching: %v", err)
	}

	if strings.TrimSpace(clientFingerprint) == "" {
		return "", nil, false
	}

	recentLogs, err = s.db.GetRecentLogsForMatch(50, provider, "")
	if err == nil {
		return tryMatch(recentLogs, false)
	}
	log.Printf("Failed to load provider recent logs for session matching: %v", err)
	return "", nil, false
}

func (s *ProxyServer) resolveSessionAndParent(r *http.Request, provider string, prompt []ChatMessage, handles []ConversationHandle, hints map[string]string) (string, *int64, string, string) {
	s.activeMu.Lock()
	defer s.activeMu.Unlock()

	if s.lastActiveSessionMap == nil {
		s.lastActiveSessionMap = make(map[string]string)
	}

	clientFingerprint := buildClientFingerprint(r, provider, hints)
	rememberSession := func(sessionID string) {
		if strings.TrimSpace(clientFingerprint) != "" && strings.TrimSpace(sessionID) != "" {
			s.lastActiveSessionMap[clientFingerprint] = sessionID
		}
	}

	if stableHandle, ok := selectStableSessionHandle(handles); ok {
		sessionID := canonicalSessionIDForHandle(provider, stableHandle)
		ref, err := s.db.FindLatestLogByHandle(stableHandle.Kind, stableHandle.Value)
		if err != nil {
			log.Printf("Failed to resolve stable conversation handle %s=%s: %v", stableHandle.Kind, stableHandle.Value, err)
		} else if ref != nil {
			if ref.SessionID != "" {
				sessionID = ref.SessionID
			}
			rememberSession(sessionID)
			parentID := ref.ID
			return sessionID, &parentID, "", clientFingerprint
		}
		rememberSession(sessionID)
		return sessionID, nil, "", clientFingerprint
	}

	if previousResponseHandle, ok := findConversationHandle(handles, "previous_response_id"); ok {
		ref, err := s.db.FindLatestLogByHandle("response_id", previousResponseHandle.Value)
		if err != nil {
			log.Printf("Failed to resolve previous response id %s: %v", previousResponseHandle.Value, err)
		} else if ref != nil {
			sessionID := ref.SessionID
			if sessionID == "" {
				sessionID = seedSessionIDFromPreviousResponse(provider, previousResponseHandle)
			}
			rememberSession(sessionID)
			parentID := ref.ID
			return sessionID, &parentID, "", clientFingerprint
		}
		sessionID := seedSessionIDFromPreviousResponse(provider, previousResponseHandle)
		rememberSession(sessionID)
		return sessionID, nil, "", clientFingerprint
	}

	isCold := true
	userCount := 0
	for _, m := range prompt {
		r := strings.ToLower(m.Role)
		if r == "assistant" || r == "tool" || r == "function" {
			isCold = false
			break
		}
		if r == "user" {
			userCount++
		}
	}

	if isCold && userCount <= 1 {
		s.suspMu.Lock()
		now := time.Now()
		var activeSuspended []SuspendedTool
		for _, st := range s.suspendedTools {
			if now.Sub(st.CreatedAt) < 180*time.Second {
				activeSuspended = append(activeSuspended, st)
			}
		}
		s.suspendedTools = activeSuspended

		promptText := extractPrimaryUserPrompt(prompt)

		if promptText != "" {
			matchedIndex := -1
			bestScore := 0.0
			secondBestScore := 0.0
			for i, st := range s.suspendedTools {
				score := semanticMatchScore(promptText, st.Arguments)
				if score > bestScore || (score == bestScore && matchedIndex != -1 && st.CreatedAt.After(s.suspendedTools[matchedIndex].CreatedAt)) {
					secondBestScore = bestScore
					matchedIndex = i
					bestScore = score
				} else if score > secondBestScore {
					secondBestScore = score
				}
			}

			if matchedIndex != -1 && bestScore >= 0.3 && (secondBestScore == 0 || bestScore >= secondBestScore+0.15) {
				st := s.suspendedTools[matchedIndex]
				s.suspendedTools = append(s.suspendedTools[:matchedIndex], s.suspendedTools[matchedIndex+1:]...)
				s.suspMu.Unlock()

				sessionID := st.SessionID
				rememberSession(sessionID)
				parentID := st.LogID
				return sessionID, &parentID, st.ToolCallID, clientFingerprint
			}
		}
		s.suspMu.Unlock()
	}

	if sessionID, parentID, ok := s.matchHistory(provider, clientFingerprint, prompt); ok {
		rememberSession(sessionID)
		return sessionID, parentID, "", clientFingerprint
	}

	if isSideTask(prompt) {
		if sid, ok := s.lastActiveSessionMap[clientFingerprint]; ok && sid != "" {
			return sid, nil, "", clientFingerprint
		}
		sid, _, err := s.db.GetLatestSessionByClientFingerprint(clientFingerprint)
		if err != nil {
			log.Printf("Failed to resolve client fingerprint fallback session for %s: %v", clientFingerprint, err)
		} else if sid != "" {
			rememberSession(sid)
			return sid, nil, "", clientFingerprint
		}
	}

	newSessionID := generateUUID()
	rememberSession(newSessionID)
	return newSessionID, nil, "", clientFingerprint
}

func joinUpstreamURL(baseURL, pathSuffix string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	pathSuffix = strings.TrimPrefix(pathSuffix, "/")

	// 智能去重：如果 baseURL 以 "/v1" 结尾，且 pathSuffix 以 "v1/" 开头，则剔除 pathSuffix 的 "v1/"
	if strings.HasSuffix(baseURL, "/v1") && strings.HasPrefix(pathSuffix, "v1/") {
		pathSuffix = strings.TrimPrefix(pathSuffix, "v1/")
	}

	return baseURL + "/" + pathSuffix
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
