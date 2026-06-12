import React, { useState, useEffect, useCallback, useRef } from "react";

// ==========================================
// 极简 SVG 图标组件 (零依赖)
// ==========================================
const IconSettings = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="18"
    height="18"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.1a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" />
    <circle cx="12" cy="12" r="3" />
  </svg>
);

const IconSearch = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <circle cx="11" cy="11" r="8" />
    <line x1="21" y1="21" x2="16.65" y2="16.65" />
  </svg>
);

const IconRefresh = ({ className }) => (
  <svg
    className={className}
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
  </svg>
);

const IconCpu = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="4" y="4" width="16" height="16" rx="2" ry="2" />
    <rect x="9" y="9" width="6" height="6" />
    <line x1="9" y1="1" x2="9" y2="4" />
    <line x1="15" y1="1" x2="15" y2="4" />
    <line x1="9" y1="20" x2="9" y2="23" />
    <line x1="15" y1="20" x2="15" y2="23" />
    <line x1="20" y1="9" x2="23" y2="9" />
    <line x1="20" y1="15" x2="23" y2="15" />
    <line x1="1" y1="9" x2="4" y2="9" />
    <line x1="1" y1="15" x2="4" y2="15" />
  </svg>
);

const IconClock = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="14"
    height="14"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <circle cx="12" cy="12" r="10" />
    <polyline points="12 6 12 12 16 14" />
  </svg>
);

const IconFlame = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="14"
    height="14"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M8.5 14.5A2.5 2.5 0 0 0 11 12c0-1.38-.5-2-1-3-1.072-2.143-.224-4.054 2-6 .5 2.5 2 4.9 4 6.5 2 1.6 3 3.5 3 5.5a7 7 0 1 1-14 0c0-1.153.433-2.294 1-3a2.5 2.5 0 0 0 2.5 2.5z" />
  </svg>
);

const IconChevronDown = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="6 9 12 15 18 9" />
  </svg>
);

const IconChevronRight = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="9 18 15 12 9 6" />
  </svg>
);

const IconCopy = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="14"
    height="14"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
  </svg>
);

const IconCheck = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="14"
    height="14"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="20 6 9 17 4 12" />
  </svg>
);

const IconTerminal = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="4 17 10 11 4 5" />
    <line x1="12" y1="19" x2="20" y2="19" />
  </svg>
);

const IconDatabase = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <ellipse cx="12" cy="5" rx="9" ry="3" />
    <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
    <path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3" />
  </svg>
);

const IconWrench = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z" />
  </svg>
);

const IconEye = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
    <circle cx="12" cy="12" r="3" />
  </svg>
);

const IconEyeOff = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
    <line x1="1" y1="1" x2="23" y2="23" />
  </svg>
);

const IconTrash = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="14"
    height="14"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="3 6 5 6 21 6" />
    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    <line x1="10" y1="11" x2="10" y2="17" />
    <line x1="14" y1="11" x2="14" y2="17" />
  </svg>
);

// 智能模拟后端的路由路径拼接（进行去重处理）
function joinUpstreamURL(baseURL, pathSuffix) {
  if (!baseURL) return "";
  let base = baseURL.trim();
  if (base.endsWith("/")) {
    base = base.slice(0, -1);
  }
  let suffix = pathSuffix.trim();
  if (suffix.startsWith("/")) {
    suffix = suffix.slice(1);
  }

  if (base.endsWith("/v1") && suffix.startsWith("v1/")) {
    suffix = suffix.slice(3); // 剔除首部的 "v1/"
  }

  return base + "/" + suffix;
}

// 简单格式化 JSON string 或者进行着色处理
function formatJSON(jsonStr) {
  try {
    if (!jsonStr) return "";
    const obj = typeof jsonStr === "string" ? JSON.parse(jsonStr) : jsonStr;
    return JSON.stringify(obj, null, 2);
  } catch (e) {
    return jsonStr;
  }
}

// 极其精简的代码语法高亮
function CodeBlock({ code, language }) {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative mt-2 border border-slate-800 rounded-lg overflow-hidden bg-black/40">
      <div className="flex justify-between items-center px-4 py-1.5 bg-slate-900/60 border-b border-slate-800/80 text-xs text-slate-400">
        <span className="code-font">{language || "JSON"}</span>
        <button
          onClick={copyToClipboard}
          className="hover:text-cyan-400 flex items-center gap-1"
        >
          {copied ? <IconCheck /> : <IconCopy />}
          {copied ? "已复制" : "复制"}
        </button>
      </div>
      <pre className="p-4 overflow-auto code-font text-xs text-emerald-400/90 leading-relaxed max-h-[350px]">
        <code>{code}</code>
      </pre>
    </div>
  );
}

// 清理用户提问，提取核心意图用于摘要展示
function simplifyUserPrompt(content) {
  if (!content) return "";
  let text = content.trim();

  // 1. 清理 XML 标签（如 <user_query> 等）
  text = text.replace(/<[^>]+>/g, "");

  // 2. 清理常见的前缀前贴
  const markers = [
    "User Query:", "User input:", "Query:", "Question:", "Task:", "Request:",
    "用户提问:", "用户输入:", "问题:", "任务:", "请求:",
    "Input:", "Prompt:", "User:", "用户:",
    "Goal:", "Objective:", "Instruction:"
  ];
  
  for (const marker of markers) {
    const idx = text.toLowerCase().lastIndexOf(marker.toLowerCase());
    if (idx !== -1) {
      const remaining = text.substring(idx + marker.length).trim();
      if (remaining.length > 2) {
        text = remaining;
      }
    }
  }

  // 3. 换行符转空格，压缩空白
  text = text.replace(/\r\n|\n|\r/g, " ");
  text = text.replace(/\s+/g, " ");
  
  return text.trim();
}

export default function App() {
  // 日志列表状态
  const [logs, setLogs] = useState([]);
  const [totalLogs, setTotalLogs] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(15);

  // 筛选与搜索
  const [searchInput, setSearchInput] = useState("");
  const [searchKeyword, setSearchKeyword] = useState("");
  const [sessionLogs, setSessionLogs] = useState([]);
  const totalSessionTokens = sessionLogs.reduce((sum, log) => sum + (log.totalTokens || 0), 0);
  const totalSessionInputTokens = sessionLogs.reduce((sum, log) => sum + (log.inputTokens || 0), 0);
  const totalSessionOutputTokens = sessionLogs.reduce((sum, log) => sum + (log.outputTokens || 0), 0);
  const totalSessionCachedTokens = sessionLogs.reduce((sum, log) => sum + (log.cachedTokens || 0), 0);
  const totalSessionUncachedInputTokens = totalSessionInputTokens - totalSessionCachedTokens;
  const [isSessionLoading, setIsSessionLoading] = useState(false);
  const [providerFilter, setProviderFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");

  // 活跃日志详情
  const [selectedLogId, setSelectedLogId] = useState(null);
  const [selectedLog, setSelectedLog] = useState(null);
  const [isLogLoading, setIsLogLoading] = useState(false);
  const [isListLoading, setIsListLoading] = useState(false);
  const [selectedToolIndex, setSelectedToolIndex] = useState(0);
  const [toolSearchQuery, setToolSearchQuery] = useState("");

  // 历史详情弹窗状态
  const [modalLogId, setModalLogId] = useState(null);
  const [modalLog, setModalLog] = useState(null);
  const [isModalLoading, setIsModalLoading] = useState(false);
  const [modalActiveTab, setModalActiveTab] = useState("trace");
  const [modalSelectedToolIndex, setModalSelectedToolIndex] = useState(0);
  const [modalToolSearchQuery, setModalToolSearchQuery] = useState("");

  // 统计信息
  const [stats, setStats] = useState({
    totalCalls: 0,
    totalInputTokens: 0,
    totalOutputTokens: 0,
    totalTokens: 0,
    totalCachedTokens: 0,
    avgDurationMs: 0,
    successCalls: 0,
    successRate: 0,
    callsByProvider: {},
    tokensByModel: {},
  });

  // 配置状态
  const [config, setConfig] = useState({
    listenAddr: "",
    openaiBaseURL: "",
    openaiAPIKey: "",
    openaiResponsesBaseURL: "",
    openaiResponsesAPIKey: "",
    anthropicBaseURL: "",
    anthropicAPIKey: "",
    dbPath: "",
  });
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [showOpenAIKey, setShowOpenAIKey] = useState(false);
  const [showAnthropicKey, setShowAnthropicKey] = useState(false);
  const [showResponsesKey, setShowResponsesKey] = useState(false);
  const [isConfigSaving, setIsConfigSaving] = useState(false);
  const [copiedText, setCopiedText] = useState("");
  const [activeTab, setActiveTab] = useState("trace"); // trace, prompt, raw_req, raw_resp

  // 系统日志相关的状态
  const [isSystemLogsOpen, setIsSystemLogsOpen] = useState(false);
  const [systemLogs, setSystemLogs] = useState("");
  const [isSystemLogsLoading, setIsSystemLogsLoading] = useState(false);
  const [autoScroll, setAutoScroll] = useState(true);
  const logsContainerRef = useRef(null);

  // 折叠状态控制
  const [showSystemPrompt, setShowSystemPrompt] = useState(false);
  const [showToolsDef, setShowToolsDef] = useState(false);



  // 渲染子代理会话的内联时间链
  const renderSubLogsForTool = (toolCallId, targetLog = selectedLog) => {
    if (!targetLog || !targetLog.subLogs) return null;
    const matchingSubLogs = targetLog.subLogs.filter(
      (sl) => sl.parentToolCallId === toolCallId,
    );
    if (matchingSubLogs.length === 0) return null;

    return (
      <div className="mt-3 border-t border-slate-800/80 pt-3 flex flex-col gap-3 font-sans">
        <div className="text-[10px] font-bold text-amber-400/90 uppercase tracking-wider flex items-center gap-1.5">
          <span>↳ 子代理会话执行追踪 (Subagent Trace)</span>
          <span className="px-1.5 py-0.2 rounded bg-amber-950/50 text-[8px] border border-amber-900/30">
            共 {matchingSubLogs.length} 条调用
          </span>
        </div>
        <div className="flex flex-col gap-3 pl-3.5 border-l border-amber-500/20">
          {matchingSubLogs.map((subLog) => (
            <div
              key={subLog.id}
              className="p-3 bg-slate-900/30 border border-slate-800/50 rounded-xl text-xs flex flex-col gap-2.5"
            >
              <div className="flex justify-between items-center">
                <span className="text-[9px] font-mono text-slate-400">
                  分支日志 #{subLog.id} • {subLog.model}
                </span>
                <span className="text-[9px] px-1.5 py-0.5 rounded bg-black/30 font-mono text-orange-400/90">
                  {(subLog.durationMs / 1000).toFixed(2)}s •{" "}
                  {subLog.totalTokens} tkn
                </span>
              </div>

              {/* 渲染子日志中的消息流（排除system） */}
              <div className="flex flex-col gap-2">
                {subLog.prompt &&
                  subLog.prompt
                    .filter((p) => p.role !== "system")
                    .map((m, mIdx) => (
                      <div
                        key={mIdx}
                        className="bg-slate-950/30 p-2.5 rounded-lg border border-slate-900/60 leading-relaxed text-[11px]"
                      >
                        <div className="font-bold text-[9px] text-cyan-400/80 mb-1">
                          {m.role === "user"
                            ? "👤 User"
                            : m.role === "tool"
                              ? "🛠️ Tool Response"
                              : "🤖 Assistant"}
                          :
                        </div>
                        <div className="text-slate-300 whitespace-pre-wrap">
                          {m.content}
                        </div>
                        {m.tool_calls &&
                          m.tool_calls.map((tc, tcIdx) => (
                            <div
                              key={tcIdx}
                              className="mt-1.5 p-2 bg-slate-950/80 rounded border border-slate-900 font-mono text-[9px]"
                            >
                              <span className="text-orange-400/80">
                                Call: {tc.name}
                              </span>
                              <span className="text-slate-500 block break-all">
                                Args: {tc.arguments}
                              </span>
                            </div>
                          ))}
                      </div>
                    ))}

                {/* 渲染子日志最新回复 */}
                {subLog.response &&
                  (subLog.response.content ||
                    (subLog.response.tool_calls &&
                      subLog.response.tool_calls.length > 0)) && (
                    <div className="bg-slate-950/45 p-2.5 rounded-lg border border-slate-900/85 border-l border-l-cyan-500/40 leading-relaxed text-[11px]">
                      <div className="font-bold text-[9px] text-cyan-400 mb-1">
                        🤖 Assistant (最新回复):
                      </div>
                      {subLog.response.thinking && (
                        <div className="text-purple-400/80 italic text-[10px] mb-1">
                          Thinking: {subLog.response.thinking}
                        </div>
                      )}
                      <div className="text-slate-205 whitespace-pre-wrap">
                        {subLog.response.content}
                      </div>
                      {subLog.response.tool_calls &&
                        subLog.response.tool_calls.map((tc, tcIdx) => (
                          <div
                            key={tcIdx}
                            className="mt-1.5 p-2 bg-slate-950/80 rounded border border-slate-900 font-mono text-[9px]"
                          >
                            <span className="text-orange-400/80">
                              Call: {tc.name}
                            </span>
                            <span className="text-slate-500 block break-all">
                              Args: {tc.arguments}
                            </span>
                          </div>
                        ))}
                    </div>
                  )}
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  };

  // 轮询计数器与定时器引用
  const autoRefreshTimer = useRef(null);

  // 动态计算本地代理根地址
  const proxyBase = (() => {
    const addr = config.listenAddr || ":8080";
    if (addr.startsWith(":")) {
      return `http://localhost${addr}`;
    }
    if (!addr.startsWith("http://") && !addr.startsWith("https://")) {
      return `http://${addr}`;
    }
    return addr;
  })();

  // ==========================================
  // API 请求函数
  // ==========================================

  const fetchStats = useCallback(async () => {
    try {
      const res = await fetch("/api/stats");
      if (res.ok) {
        const data = await res.json();
        setStats(data);
      }
    } catch (err) {
      console.error("Failed to fetch stats:", err);
    }
  }, []);

  const fetchConfig = useCallback(async () => {
    try {
      const res = await fetch("/api/config");
      if (res.ok) {
        const data = await res.json();
        setConfig(data);
      }
    } catch (err) {
      console.error("Failed to fetch config:", err);
    }
  }, []);

  const fetchLogs = useCallback(
    async (page = 1) => {
      setIsListLoading(true);
      try {
        let url = `/api/logs?page=${page}&pageSize=${pageSize}`;
        if (providerFilter !== "all") url += `&provider=${providerFilter}`;
        if (statusFilter !== "all") {
          const filterVal = statusFilter === "success" ? 200 : 502; // 后端根据特定状态码过滤
          url += `&status=${filterVal}`;
        }
        if (searchKeyword.trim() !== "")
          url += `&keyword=${encodeURIComponent(searchKeyword)}`;

        const res = await fetch(url);
        if (res.ok) {
          const data = await res.json();
          setLogs(data.list || []);
          setTotalLogs(data.total || 0);
          setCurrentPage(data.page || 1);
        }
      } catch (err) {
        console.error("Failed to fetch logs:", err);
      } finally {
        setIsListLoading(false);
      }
    },
    [providerFilter, statusFilter, searchKeyword, pageSize],
  );

  const fetchSessionLogs = useCallback(async (sessionId) => {
    if (!sessionId) {
      setSessionLogs([]);
      return;
    }
    setIsSessionLoading(true);
    try {
      const res = await fetch(`/api/sessions/${sessionId}/logs`);
      if (res.ok) {
        const data = await res.json();
        setSessionLogs(data || []);
      }
    } catch (err) {
      console.error("Failed to fetch session logs:", err);
    } finally {
      setIsSessionLoading(false);
    }
  }, []);

  const fetchLogDetail = useCallback(async (id) => {
    if (!id) return;
    setIsLogLoading(true);
    try {
      const res = await fetch(`/api/logs/${id}`);
      if (res.ok) {
        const data = await res.json();
        setSelectedLog(data);
        if (data.sessionId) {
          fetchSessionLogs(data.sessionId);
        } else {
          setSessionLogs([]);
        }
      }
    } catch (err) {
      console.error("Failed to fetch log detail:", err);
    } finally {
      setIsLogLoading(false);
    }
  }, [fetchSessionLogs]);

  // 删除单个日志
  const handleDeleteLog = async (e, id) => {
    e.stopPropagation(); // 阻止卡片点击展开详情
    if (!window.confirm("确定要删除这条请求日志吗？该操作不可逆。")) {
      return;
    }

    try {
      const res = await fetch(`/api/logs/${id}`, {
        method: "DELETE",
      });
      if (res.ok) {
        // 如果当前选中的是正在被删除的日志，重置选中状态
        if (selectedLogId === id) {
          setSelectedLogId(null);
        }
        // 重新获取日志列表与用量统计
        fetchLogs(currentPage);
        fetchStats();
      } else {
        alert("删除请求日志失败");
      }
    } catch (err) {
      console.error("Failed to delete log:", err);
      alert("删除发生异常");
    }
  };

  // 删除整个会话的所有请求日志
  const handleDeleteSessionLogs = async (sessionId) => {
    if (!sessionId) return;
    if (!window.confirm(`确定要删除会话 (ID: ${sessionId}) 中的所有请求日志吗？该操作不可逆。`)) {
      return;
    }

    try {
      const res = await fetch(`/api/sessions/${sessionId}/logs`, {
        method: "DELETE",
      });
      if (res.ok) {
        // 如果当前选中的日志属于此会话，重置选中状态
        if (selectedLog && selectedLog.sessionId === sessionId) {
          setSelectedLogId(null);
        }
        // 重新获取日志列表与用量统计
        fetchLogs(currentPage);
        fetchStats();
      } else {
        alert("删除会话日志失败");
      }
    } catch (err) {
      console.error("Failed to delete session logs:", err);
      alert("删除发生异常");
    }
  };

  // 保存配置
  const handleSaveConfig = async (e) => {
    e.preventDefault();
    setIsConfigSaving(true);
    try {
      const res = await fetch("/api/config", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(config),
      });
      if (res.ok) {
        setIsSettingsOpen(false);
        fetchConfig();
      } else {
        alert("配置保存失败");
      }
    } catch (err) {
      console.error("Failed to save config:", err);
      alert("保存发生异常");
    } finally {
      setIsConfigSaving(false);
    }
  };

  // 获取系统运行日志
  const fetchSystemLogs = useCallback(async () => {
    setIsSystemLogsLoading(true);
    try {
      const res = await fetch("/api/system-logs");
      if (res.ok) {
        const text = await res.text();
        setSystemLogs(text);
      }
    } catch (err) {
      console.error("Failed to fetch system logs:", err);
    } finally {
      setIsSystemLogsLoading(false);
    }
  }, []);

  // 清空系统运行日志
  const handleClearSystemLogs = async () => {
    if (!window.confirm("确定要物理清空系统运行日志文件吗？该操作将重置本地 system.log 文件。")) {
      return;
    }
    try {
      const res = await fetch("/api/system-logs/clear", {
        method: "POST",
      });
      if (res.ok) {
        setSystemLogs("--- system log cleared ---\n");
      } else {
        alert("清空系统日志失败");
      }
    } catch (err) {
      console.error("Failed to clear system logs:", err);
      alert("清空系统日志发生异常");
    }
  };

  // 仅在系统日志模态框打开时进行 2.5 秒的轮询，并在关闭或卸载时销毁
  useEffect(() => {
    let timer = null;
    if (isSystemLogsOpen) {
      fetchSystemLogs();
      timer = setInterval(fetchSystemLogs, 2500);
    }
    return () => {
      if (timer) {
        clearInterval(timer);
      }
    };
  }, [isSystemLogsOpen, fetchSystemLogs]);

  // 控制系统日志终端滚屏到最下方
  useEffect(() => {
    if (autoScroll && logsContainerRef.current) {
      logsContainerRef.current.scrollTop = logsContainerRef.current.scrollHeight;
    }
  }, [systemLogs, autoScroll]);

  // ==========================================
  // 生命周期副作用
  // ==========================================

  // 搜索词防抖
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearchKeyword(searchInput);
      setCurrentPage(1); // 搜索词改变时重置为第一页
    }, 400);

    return () => clearTimeout(timer);
  }, [searchInput]);

  // 初始化加载
  useEffect(() => {
    fetchStats();
    fetchConfig();
    fetchLogs(1);
  }, [fetchStats, fetchConfig, fetchLogs]);

  // 轮询更新列表和用量统计（每 5 秒一次）
  useEffect(() => {
    autoRefreshTimer.current = setInterval(() => {
      // 仅在第一页且没有搜索过滤时自动刷新日志列表，避免漂移打断历史阅读
      if (currentPage === 1 && searchKeyword.trim() === "") {
        fetchLogs(1);
      }
      fetchStats();
    }, 5000);

    return () => {
      if (autoRefreshTimer.current) {
        clearInterval(autoRefreshTimer.current);
      }
    };
  }, [fetchLogs, fetchStats, currentPage, searchKeyword]);

  // 当选择的日志 ID 改变时加载详情
  useEffect(() => {
    if (selectedLogId) {
      fetchLogDetail(selectedLogId);
      setSelectedToolIndex(0);
      setToolSearchQuery("");
    } else {
      setSelectedLog(null);
      setSessionLogs([]);
    }
  }, [selectedLogId, fetchLogDetail]);

  // 当弹窗日志 ID 改变时加载弹窗详情
  useEffect(() => {
    if (modalLogId) {
      (async () => {
        setIsModalLoading(true);
        try {
          const res = await fetch(`/api/logs/${modalLogId}`);
          if (res.ok) {
            const data = await res.json();
            setModalLog(data);
            setModalActiveTab("trace"); // 默认切到 trace tab
            setModalSelectedToolIndex(0);
            setModalToolSearchQuery("");
          }
        } catch (err) {
          console.error("Failed to fetch modal log detail:", err);
        } finally {
          setIsModalLoading(false);
        }
      })();
    } else {
      setModalLog(null);
    }
  }, [modalLogId]);

  // 监听 ESC 键关闭弹窗
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === "Escape") {
        setModalLogId(null);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  // 页码切换
  const handlePageChange = (newPage) => {
    if (newPage >= 1 && newPage <= Math.ceil(totalLogs / pageSize)) {
      fetchLogs(newPage);
    }
  };

  // 手动刷新数据
  const handleManualRefresh = () => {
    fetchLogs(currentPage);
    fetchStats();
    if (selectedLogId) {
      fetchLogDetail(selectedLogId);
    }
  };



  const renderDetailModal = () => {
    if (!modalLogId) return null;

    return (
      <div
        className="fixed inset-0 bg-slate-950/75 backdrop-blur-sm z-50 flex items-center justify-center p-4 modal-fade-in"
        onClick={() => setModalLogId(null)}
      >
        <div
          className="glass-panel w-full max-w-5xl h-[85vh] overflow-hidden flex flex-col border border-slate-800/85 shadow-[0_20px_50px_rgba(0,0,0,0.85)] modal-scale-up"
          onClick={(e) => e.stopPropagation()}
        >
          {/* 头部 */}
          <div className="flex justify-between items-center px-6 py-4 border-b border-slate-900 bg-slate-950/20 shrink-0">
            <div className="flex flex-col gap-1 min-w-0">
              <div className="flex flex-wrap items-center gap-2.5">
                <span className="text-base font-bold font-mono text-cyan-400 tracking-tight truncate">
                  {isModalLoading ? "加载中..." : (modalLog?.model || "请求详情")}
                </span>
                {!isModalLoading && modalLog && (
                  <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-900 border border-slate-800 text-slate-400 uppercase font-bold tracking-wider shrink-0">
                    {modalLog.provider}
                  </span>
                )}
              </div>
              {!isModalLoading && modalLog && (
                <span className="text-[10px] text-slate-500 font-mono truncate">
                  {modalLog.sessionId
                    ? `会话 #${modalLog.sessionSeq} (ID: ${modalLog.id})`
                    : `ID: ${modalLog.id}`}{" "}
                  • Path: {modalLog.path}
                </span>
              )}
            </div>
            <button
              onClick={() => setModalLogId(null)}
              className="text-slate-450 hover:text-slate-200 text-2xl font-light transition-colors active:scale-95 leading-none w-8 h-8 flex items-center justify-center rounded-full hover:bg-slate-900/50"
            >
              &times;
            </button>
          </div>

          {/* 内容区 */}
          <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-6 min-h-0">
            {isModalLoading ? (
              <div className="flex-1 flex flex-col items-center justify-center text-slate-500 gap-4">
                <IconRefresh className="animate-spin text-cyan-400 w-8 h-8" />
                <p className="text-xs font-medium tracking-wide">
                  加载请求结构化详情中...
                </p>
              </div>
            ) : !modalLog ? (
              <div className="flex-1 flex flex-col items-center justify-center text-red-400 gap-3">
                <span>加载日志失败，请稍后重试</span>
              </div>
            ) : (
              <div className="flex-1 flex flex-col overflow-hidden gap-6 min-h-0">
                {/* 详情头部元数据卡片 (Grid 布局，突出核心指标) */}
                <div className="p-5 bg-slate-950/35 border border-slate-900/60 rounded-2xl flex flex-col gap-4 lg:flex-row lg:justify-between lg:items-center shrink-0">
                  <div className="flex flex-col gap-1 min-w-0">
                    <span className="text-xs font-semibold text-slate-305 font-sans">
                      运行指标摘要 (Run Metrics)
                    </span>
                    <span className="text-[10px] text-slate-500 font-mono">
                      开始时间: {modalLog.createdAt}
                    </span>
                  </div>

                  <div
                    className={`grid grid-cols-2 ${modalLog.cachedTokens > 0 || modalLog.cacheReadTokens > 0 || modalLog.cacheCreationTokens > 0 ? "sm:grid-cols-5" : "sm:grid-cols-4"} gap-3 text-xs font-mono shrink-0`}
                  >
                    <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-orange-950/15 border border-orange-900/15">
                      <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5 font-sans">
                        请求耗时
                      </span>
                      <span className="font-bold text-sm text-orange-400">
                        {(modalLog.durationMs / 1000).toFixed(2)}s
                      </span>
                    </div>
                    <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-cyan-950/15 border border-cyan-900/15">
                      <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5 font-sans">
                        输入 Tokens
                      </span>
                      <span className="font-bold text-sm text-cyan-400">
                        {modalLog.inputTokens}
                      </span>
                      {modalLog.cachedTokens > 0 && (
                        <span className="text-[8px] text-slate-450 mt-0.5">
                          实际计算:{" "}
                          {modalLog.inputTokens - modalLog.cachedTokens}
                        </span>
                      )}
                    </div>
                    {(modalLog.cachedTokens > 0 ||
                      modalLog.cacheReadTokens > 0 ||
                      modalLog.cacheCreationTokens > 0) && (
                      <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-indigo-950/20 border border-indigo-900/30 text-indigo-400">
                        <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5 font-sans">
                          缓存 Tokens
                        </span>
                        <span className="font-bold text-sm text-indigo-400">
                          {modalLog.cachedTokens}
                        </span>
                        <span className="text-[8px] text-slate-450 mt-0.5 truncate max-w-full">
                          {modalLog.provider === "openai" ||
                          modalLog.provider === "openai-responses"
                            ? "命中缓存"
                            : `读缓存 (另写 ${modalLog.cacheCreationTokens})`}
                        </span>
                      </div>
                    )}
                    <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-purple-950/15 border border-purple-900/15">
                      <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5 font-sans">
                        输出 Tokens
                      </span>
                      <span className="font-bold text-sm text-purple-400">
                        {modalLog.outputTokens}
                      </span>
                    </div>
                    <div
                      className={`flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl border ${
                        modalLog.statusCode >= 200 &&
                        modalLog.statusCode < 300
                          ? "bg-emerald-950/15 border-emerald-900/15 text-emerald-400"
                          : "bg-red-950/15 border-red-900/15 text-red-400"
                      }`}
                    >
                      <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5 font-sans">
                        状态代码
                      </span>
                      <span className="font-bold text-sm">
                        {modalLog.statusCode}
                      </span>
                    </div>
                  </div>
                </div>

                {/* 时序控制 Tab */}
                <div className="flex border-b border-slate-900 pb-1 gap-6 shrink-0">
                  <button
                    onClick={() => setModalActiveTab("trace")}
                    className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                      modalActiveTab === "trace"
                        ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                        : "border-transparent text-slate-400 hover:text-slate-200"
                    }`}
                  >
                    时序追踪 (Trace)
                  </button>
                  <button
                    onClick={() => setModalActiveTab("prompt")}
                    className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                      modalActiveTab === "prompt"
                        ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                        : "border-transparent text-slate-400 hover:text-slate-200"
                    }`}
                  >
                    请求消息 (
                    {modalLog.prompt
                      ? modalLog.prompt.filter((p) => p.role !== "system").length
                      : 0}
                    )
                  </button>
                  {modalLog.tools && modalLog.tools.length > 0 && (
                    <button
                      onClick={() => setModalActiveTab("tools")}
                      className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                        modalActiveTab === "tools"
                          ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                          : "border-transparent text-slate-400 hover:text-slate-200"
                      }`}
                    >
                      定义的工具 ({modalLog.tools.length})
                    </button>
                  )}
                  <button
                    onClick={() => setModalActiveTab("raw_req")}
                    className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                      modalActiveTab === "raw_req"
                        ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                        : "border-transparent text-slate-400 hover:text-slate-200"
                    }`}
                  >
                    原始请求 JSON
                  </button>
                  <button
                    onClick={() => setModalActiveTab("raw_resp")}
                    className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                      modalActiveTab === "raw_resp"
                        ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                        : "border-transparent text-slate-400 hover:text-slate-200"
                    }`}
                  >
                    原始响应 Raw
                  </button>
                </div>

                {/* Tab 内容区 */}
                <div className="flex-1 overflow-y-auto pr-1.5 min-h-0 flex flex-col gap-4 scroll-isolated">
                  {/* 1. 时序追踪 (Trace) */}
                  {modalActiveTab === "trace" && (
                    <div className="flex flex-col gap-5">
                      {/* 可视化分段 Timeline */}
                      <div className="p-4 bg-slate-950/20 border border-slate-900 rounded-xl flex flex-col gap-3">
                        <div className="flex justify-between items-center text-xs text-slate-400 font-mono">
                          <span className="font-semibold text-slate-300 font-sans">
                            时序可视化时间线
                          </span>
                          <span className="text-cyan-400 font-bold">
                            {(modalLog.durationMs / 1000).toFixed(2)}s
                          </span>
                        </div>

                        {/* 分段时间条 */}
                        <div className="h-7 w-full bg-slate-950/80 rounded-lg overflow-hidden flex relative p-0.5 border border-slate-900">
                          <div
                            className="h-full bg-indigo-600/40 rounded-l-md border-r border-slate-900/30 flex items-center justify-center text-[8px] font-bold text-indigo-300 shrink-0 font-sans"
                            style={{
                              width: `${Math.min(15, Math.max(5, (150 / modalLog.durationMs) * 100))}%`,
                            }}
                          >
                            连接
                          </div>
                          <div
                            className="h-full timeline-shimmer rounded-r-md"
                            style={{ flex: 1 }}
                          />
                          <div className="absolute inset-0 flex justify-end items-center pr-4 text-[9px] font-bold text-white uppercase tracking-wider drop-shadow-md font-sans">
                            <span>
                              传输流 •{" "}
                              {(modalLog.durationMs / 1000).toFixed(2)}s
                            </span>
                          </div>
                        </div>

                        {/* 如果有工具调用，显示时序序列图 */}
                        {modalLog.response &&
                          modalLog.response.tool_calls &&
                          modalLog.response.tool_calls.length > 0 && (
                            <div className="mt-1.5 flex flex-col gap-2.5">
                              <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider font-sans">
                                执行的 Tool Call Spans
                              </span>
                              <div className="flex flex-col gap-2">
                                {modalLog.response.tool_calls.map(
                                  (tc, idx) => (
                                    <div
                                      key={idx}
                                      className="flex items-center justify-between px-3 py-2.5 bg-slate-950/40 border border-slate-900 rounded-xl text-xs"
                                    >
                                      <div className="flex items-center gap-2 min-w-0 font-sans">
                                        <span className="text-[9px] font-bold px-2 py-0.5 bg-orange-950/80 border border-orange-900/20 text-orange-400 rounded-md shrink-0">
                                          TOOL
                                        </span>
                                        <span className="font-mono text-slate-200 font-semibold truncate">
                                          {tc.name}
                                        </span>
                                      </div>
                                      <span className="text-[9px] text-slate-500 font-mono shrink-0">
                                        id: {tc.id}
                                      </span>
                                    </div>
                                  ),
                                )}
                              </div>
                            </div>
                          )}
                      </div>
                    </div>
                  )}

                  {/* 2. 对话结构化展示 */}
                  {modalActiveTab === "prompt" && (
                    <div className="flex flex-col gap-4">
                      {/* 折叠 System Prompt */}
                      <div className="flex flex-col gap-2.5">
                        <button
                          onClick={() => setShowSystemPrompt(!showSystemPrompt)}
                          className="flex justify-between items-center px-4 py-3 bg-slate-950/40 hover:bg-slate-900/60 border border-slate-900/85 rounded-xl text-xs text-slate-300 transition-all active:scale-[0.99] font-sans"
                        >
                          <span className="flex items-center gap-2">
                            <IconTerminal />
                            <span className="font-semibold">System 提示词</span>
                          </span>
                          <span className="text-slate-455">
                            {showSystemPrompt ? (
                              <IconChevronDown />
                            ) : (
                              <IconChevronRight />
                            )}
                          </span>
                        </button>
                        {showSystemPrompt && (
                          <div className="p-4 bg-slate-950/60 border border-slate-900 rounded-xl text-xs font-mono text-slate-400 whitespace-pre-wrap leading-relaxed">
                            {modalLog.prompt &&
                            modalLog.prompt.filter((p) => p.role === "system")
                              .length > 0
                              ? modalLog.prompt
                                  .filter((p) => p.role === "system")
                                  .map((p) => p.content)
                                  .join("\n")
                              : "无系统提示词"}
                          </div>
                        )}
                      </div>

                      {/* 消息时间线 */}
                      <div className="flex flex-col gap-4 font-sans">
                        <h3 className="text-xs uppercase font-bold tracking-wider text-slate-400">
                          消息流记录
                        </h3>
                        {modalLog.prompt &&
                          modalLog.prompt
                            .filter((p) => p.role !== "system")
                            .map((m, idx) => {
                              const isUser = m.role === "user";
                              const isTool = m.role === "tool";
                              const isAssistant = m.role === "assistant";

                              return (
                                <div
                                  key={idx}
                                  className={`p-5 rounded-2xl border flex flex-col gap-3.5 shadow-sm ${
                                    isUser
                                      ? "bg-cyan-950/10 border-cyan-950/40"
                                      : isTool
                                        ? "bg-orange-950/10 border-orange-950/40"
                                        : "bg-slate-950/25 border-slate-900"
                                  }`}
                                >
                                  <div className="flex justify-between items-center">
                                    <span
                                      className={`text-[10px] font-bold uppercase tracking-wider px-2.5 py-1 rounded-md border ${
                                        isUser
                                          ? "bg-cyan-950/60 border-cyan-900/30 text-cyan-400"
                                          : isTool
                                            ? "bg-orange-950/60 border-orange-900/30 text-orange-400"
                                            : "bg-slate-900/60 border-slate-800/80 text-cyan-400"
                                      }`}
                                    >
                                      {isUser
                                        ? "👤 User"
                                        : isTool
                                          ? "🛠️ Tool Response"
                                          : "🤖 Assistant"}
                                      {m.name && (
                                        <span className="text-slate-400 font-mono text-[9px] ml-1">
                                          ({m.name})
                                        </span>
                                      )}
                                    </span>
                                    {m.tool_call_id && (
                                      <span className="text-[9px] text-slate-500 font-mono">
                                        ID: {m.tool_call_id}
                                      </span>
                                    )}
                                  </div>

                                  {/* 思维链 */}
                                  {isAssistant && m.thinking && (
                                    <div className="p-3.5 bg-slate-900/40 border border-slate-850/60 rounded-xl text-xs text-slate-450/90 italic whitespace-pre-wrap border-l-2 border-l-purple-500/60">
                                      <div className="text-[9px] font-bold text-purple-400 uppercase tracking-wider not-italic mb-1 font-sans">
                                        思考过程 (Thinking)
                                      </div>
                                      {m.thinking}
                                    </div>
                                  )}

                                  {m.content && (
                                    <div className="text-sm leading-relaxed tracking-wide text-slate-200 whitespace-pre-wrap mt-0.5">
                                      {m.content}
                                    </div>
                                  )}

                                  {/* 工具调用展示 */}
                                  {isAssistant &&
                                    m.tool_calls &&
                                    m.tool_calls.length > 0 && (
                                      <div className="mt-1 flex flex-col gap-2 font-mono">
                                        <div className="text-[9px] font-bold text-orange-400 uppercase tracking-wider font-sans">
                                          触发工具调用 (Tool Calls)
                                        </div>
                                        {m.tool_calls.map((tc, tcIdx) => (
                                          <div
                                            key={tcIdx}
                                            className="p-3 bg-slate-950/40 border border-slate-900 rounded-xl text-xs text-slate-300 font-sans"
                                          >
                                            <div>
                                              <span className="text-slate-500 font-mono">
                                                Method:
                                              </span>{" "}
                                              <span className="text-orange-400 font-semibold font-mono">
                                                {tc.name}
                                              </span>
                                            </div>
                                            <div className="mt-1">
                                              <span className="text-slate-500 font-mono">
                                                Args:
                                              </span>{" "}
                                              <span className="text-emerald-400/90 font-mono">
                                                {tc.arguments}
                                              </span>
                                            </div>
                                            {renderSubLogsForTool(tc.id, modalLog)}
                                          </div>
                                        ))}
                                      </div>
                                    )}
                                </div>
                              );
                            })}

                        {/* 渲染当次最新大模型回复 */}
                        {modalLog.response &&
                          (modalLog.response.content ||
                            modalLog.response.thinking ||
                            (modalLog.response.tool_calls &&
                              modalLog.response.tool_calls.length > 0)) && (
                            <div className="p-5 rounded-2xl border flex flex-col gap-3.5 shadow-sm bg-slate-950/25 border-slate-900 border-l-[3px] border-l-cyan-500/50">
                              <div className="flex justify-between items-center">
                                <span className="text-[10px] font-bold uppercase tracking-wider px-2.5 py-1 rounded-md border bg-slate-900/60 border-slate-800/80 text-cyan-400">
                                  🤖 Assistant (最新回复)
                                </span>
                              </div>

                              {/* 思维链 */}
                              {modalLog.response.thinking && (
                                <div className="p-3.5 bg-slate-900/40 border border-slate-850/60 rounded-xl text-xs text-slate-450/90 italic whitespace-pre-wrap border-l-2 border-l-purple-500/60">
                                  <div className="text-[9px] font-bold text-purple-400 uppercase tracking-wider not-italic mb-1">
                                    思考过程 (Thinking)
                                  </div>
                                  {modalLog.response.thinking}
                                </div>
                              )}

                              {/* 消息正文 */}
                              {modalLog.response.content && (
                                <div className="text-sm leading-relaxed tracking-wide text-slate-200 whitespace-pre-wrap mt-0.5">
                                  {modalLog.response.content}
                                </div>
                              )}

                              {/* 工具调用展示 */}
                              {modalLog.response.tool_calls &&
                                modalLog.response.tool_calls.length > 0 && (
                                  <div className="mt-1 flex flex-col gap-2 font-mono">
                                    <div className="text-[9px] font-bold text-orange-400 uppercase tracking-wider font-sans">
                                      触发工具调用 (Tool Calls)
                                    </div>
                                    {modalLog.response.tool_calls.map(
                                      (tc, tcIdx) => (
                                        <div
                                          key={tcIdx}
                                          className="p-3 bg-slate-950/40 border border-slate-900 rounded-xl text-xs text-slate-300 font-sans"
                                        >
                                          <div>
                                            <span className="text-slate-500 font-mono">
                                              Method:
                                            </span>{" "}
                                            <span className="text-orange-400 font-semibold font-mono">
                                              {tc.name}
                                            </span>
                                          </div>
                                          <div className="mt-1">
                                            <span className="text-slate-500 font-mono">
                                              Args:
                                            </span>{" "}
                                            <span className="text-emerald-400/90 font-mono">
                                              {tc.arguments}
                                            </span>
                                          </div>
                                          {renderSubLogsForTool(tc.id, modalLog)}
                                        </div>
                                      ),
                                    )}
                                  </div>
                                )}
                            </div>
                          )}
                      </div>
                    </div>
                  )}

                  {/* 3. 原始请求 Raw */}
                  {modalActiveTab === "raw_req" && (
                    <CodeBlock
                      code={formatJSON(modalLog.rawRequest)}
                      language="Raw HTTP Request Body"
                    />
                  )}

                  {/* 4. 原始响应 Raw */}
                  {modalActiveTab === "raw_resp" && (
                    <CodeBlock
                      code={formatJSON(modalLog.rawResponse)}
                      language="Raw Upstream Response"
                    />
                  )}

                  {/* 5. 已定义工具 (Tools) */}
                  {modalActiveTab === "tools" &&
                    modalLog.tools &&
                    modalLog.tools.length > 0 && (
                      <div className="flex flex-col lg:flex-row gap-5 flex-1 min-h-0 h-full overflow-hidden font-sans">
                        {/* 左侧：工具搜索与列表 */}
                        <div className="w-full lg:w-[240px] flex flex-col gap-3 shrink-0 bg-slate-950/20 border border-slate-900/40 p-3.5 rounded-2xl min-h-0 overflow-y-auto">
                          <div className="relative flex items-center">
                            <span className="absolute left-3.5 text-slate-500">
                              <IconSearch />
                            </span>
                            <input
                              type="text"
                              placeholder="搜索工具名..."
                              value={modalToolSearchQuery}
                              onChange={(e) => {
                                setModalToolSearchQuery(e.target.value);
                                setModalSelectedToolIndex(0);
                              }}
                              className="w-full pl-9 pr-3.5 py-1.5 bg-slate-950/40 focus:bg-slate-950/85 text-[11px] border border-slate-850/85 focus:border-cyan-500/60 rounded-lg outline-none placeholder:text-slate-500 text-slate-200"
                            />
                          </div>

                          {/* 工具列表 */}
                          <div className="flex flex-col gap-1.5 overflow-y-auto flex-1 pr-1 font-mono text-[11px]">
                            {modalLog.tools
                              .map((tool, index) => ({ tool, index }))
                              .filter(({ tool }) =>
                                tool.name
                                  .toLowerCase()
                                  .includes(modalToolSearchQuery.toLowerCase()),
                              )
                              .map(({ tool, index }) => {
                                const isSelected = modalSelectedToolIndex === index;
                                return (
                                  <button
                                    key={index}
                                    onClick={() => setModalSelectedToolIndex(index)}
                                    className={`text-left px-3 py-2 rounded-lg transition-all truncate border shrink-0 ${
                                      isSelected
                                        ? "bg-cyan-950/25 border-cyan-500/50 text-cyan-400 font-bold"
                                        : "bg-transparent border-transparent text-slate-400 hover:text-slate-200 hover:bg-slate-900/30"
                                    }`}
                                    title={tool.name}
                                  >
                                    🔧 {tool.name}
                                  </button>
                                );
                              })}
                            {modalLog.tools.filter((t) =>
                              t.name
                                .toLowerCase()
                                .includes(modalToolSearchQuery.toLowerCase()),
                            ).length === 0 && (
                              <div className="text-center text-slate-500 py-6 text-[10px] font-sans">
                                未找到相关工具
                              </div>
                            )}
                          </div>
                        </div>

                        {/* 右侧：被选中工具的 Schema 和详情 */}
                        {(() => {
                          const selectedTool =
                            modalLog.tools[modalSelectedToolIndex] ||
                            modalLog.tools[0];
                          if (!selectedTool) return null;
                          return (
                            <div className="flex-1 flex flex-col gap-3 min-w-0 overflow-y-auto">
                              <div className="p-4 bg-slate-950/20 border border-slate-900 rounded-xl flex flex-col gap-2.5">
                                <div className="flex items-center gap-2">
                                  <span className="text-[10px] font-bold px-2 py-0.5 bg-cyan-950/80 border border-cyan-900/20 text-cyan-400 rounded-md">
                                    TOOL
                                  </span>
                                  <span className="text-sm font-bold text-slate-200 font-mono">
                                    {selectedTool.name}
                                  </span>
                                </div>
                                {selectedTool.description && (
                                  <p className="text-xs text-slate-400 leading-relaxed bg-slate-950/45 p-3 rounded-lg border border-slate-900/60">
                                    <span className="font-semibold text-slate-350 block mb-1">
                                      功能描述：
                                    </span>
                                    {selectedTool.description}
                                  </p>
                                )}
                              </div>
                              <div className="flex flex-col shrink-0">
                                <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider block mb-1">
                                  Parameters Schema
                                </span>
                                <CodeBlock
                                  code={formatJSON(selectedTool.parameters)}
                                  language="parameters.schema.json"
                                />
                              </div>
                            </div>
                          );
                        })()}
                      </div>
                    )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    );
  };

  const renderLogCard = (log) => {
    const isSelected = selectedLogId === log.id;
    const isSuccess = log.statusCode >= 200 && log.statusCode < 300;
    const provName =
      log.provider === "openai"
        ? "OpenAI"
        : log.provider === "anthropic"
          ? "Anthropic"
          : log.provider === "openai-responses"
            ? "Responses"
            : log.provider;
    const tokenStr =
      log.totalTokens >= 1000
        ? `${(log.totalTokens / 1000).toFixed(1)}K`
        : log.totalTokens;

    const uncachedTokens = log.totalTokens - log.cachedTokens;
    const uncachedTokenStr =
      uncachedTokens >= 1000
        ? `${(uncachedTokens / 1000).toFixed(1)}K`
        : uncachedTokens;
    const actualTokenTitle =
      log.cachedTokens > 0
        ? `未缓存实际计算 Token: ${uncachedTokens}`
        : `实际计算 Token: ${uncachedTokens}`;

    const isBranch = !!log.parentToolCallId;

    return (
      <div
        key={log.id}
        onClick={(e) => {
          e.stopPropagation();
          setSelectedLogId(log.id);
        }}
        className={`pl-3 pr-4 py-3.5 rounded-xl cursor-pointer border transition-all flex flex-col gap-2 my-2.5 group ${
          isBranch
            ? "ml-4 bg-slate-900/35 border-amber-900/30"
            : "bg-slate-950/25 border-slate-900/80 hover:bg-slate-900/40 hover:border-slate-800/80"
        } ${
          isSelected
            ? "bg-slate-900/80 border-cyan-500/60 shadow-[0_0_15px_rgba(6,182,212,0.15)] translate-x-0.5"
            : ""
        } ${isSuccess ? "border-l-[3px] border-l-emerald-500/80" : "border-l-[3px] border-l-red-500/80"}`}
      >
        <div className="flex justify-between items-center">
          <div className="flex items-center gap-2">
            {isBranch && (
              <span className="text-[9px] font-bold text-amber-400 bg-amber-950/40 border border-amber-900/30 px-1.5 py-0.5 rounded flex items-center gap-1 shrink-0">
                ↳ 子分支 #{log.parentId}
              </span>
            )}
            <span className="text-[9px] uppercase font-bold text-cyan-400/90 tracking-wide px-1.5 py-0.5 rounded bg-cyan-950/40 border border-cyan-900/30">
              {provName}
            </span>
            <span className="text-[9px] text-slate-500 font-mono">
              {log.sessionId ? `会话 #${log.sessionSeq}` : `#${log.id}`}
            </span>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <span className="text-[9px] text-slate-400/80 code-font group-hover:hidden">
              {log.createdAt ? log.createdAt.split(" ")[1] : ""}
            </span>
            <button
              title="删除请求"
              onClick={(e) => handleDeleteLog(e, log.id)}
              className="hidden group-hover:flex items-center justify-center p-1 text-slate-500 hover:text-red-400 rounded hover:bg-slate-900 transition-colors"
            >
              <IconTrash />
            </button>
          </div>
        </div>

        <div className="text-[11px] font-semibold text-slate-200 truncate font-mono">
          {log.model}
        </div>

        <div className="flex justify-between items-center text-[9px] text-slate-400 code-font whitespace-nowrap">
          <div className="flex items-center gap-3">
            {/* 耗时 (时钟图标) */}
            <div
              className="flex items-center gap-1 text-orange-400/80"
              title="请求耗时"
            >
              <IconClock />
              <span>{(log.durationMs / 1000).toFixed(2)}s</span>
            </div>

            {/* 总吞吐 (CPU图标) */}
            <div
              className="flex items-center gap-1 text-purple-400/80"
              title={`总吞吐 Token: ${log.totalTokens}`}
            >
              <IconCpu />
              <span>{tokenStr}</span>
            </div>

            {/* 实际计算/未缓存 (火焰图标) */}
            {log.totalTokens > 0 && (
              <div
                className="flex items-center gap-0.5 text-orange-400/80"
                title={actualTokenTitle}
              >
                <IconFlame />
                <span>{uncachedTokenStr}</span>
              </div>
            )}
          </div>
          <span
            className={`font-bold px-1 py-0.5 rounded bg-black/10 shrink-0 ${isSuccess ? "text-emerald-400" : "text-red-400"}`}
          >
            {log.statusCode}
          </span>
        </div>
      </div>
    );
  };

  return (
    <div className="flex flex-col h-screen w-full text-slate-100 antialiased bg-transparent p-5 md:p-6 gap-5 md:gap-6 overflow-hidden">
      {/* ==========================================
          页头 & 汇总统计
         ========================================== */}
      <header className="flex justify-between items-center glass-panel px-6 py-4 h-20 shrink-0 shimmer-glow border-slate-850 shadow-lg">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-cyan-400 via-indigo-500 to-purple-600 flex items-center justify-center font-bold text-xl text-slate-950 shadow-[0_0_20px_rgba(6,182,212,0.4)] transition-transform hover:scale-105">
            LR
          </div>
          <div>
            <h1 className="text-xl font-bold bg-gradient-to-r from-cyan-400 via-sky-400 to-indigo-400 bg-clip-text text-transparent tracking-tight">
              LLM Proxy & Tracer
            </h1>
            <p className="text-[10px] md:text-xs text-slate-400/90 font-medium tracking-wide">
              本地代理、结构化拦截与耗时追踪面板
            </p>
          </div>
        </div>

        {/* 顶部实时概览 */}
        <div className="hidden lg:flex items-center gap-8 bg-slate-950/35 border border-slate-900 px-6 py-2.5 rounded-2xl">
          <div className="flex flex-col items-center border-r border-slate-900/60 pr-6">
            <span className="text-[10px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
              请求总量
            </span>
            <span className="text-base font-bold code-font text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.2)]">
              {stats.totalCalls}
            </span>
          </div>
          <div className="flex flex-col items-center border-r border-slate-900/60 pr-6">
            <span className="text-[10px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
              成功率
            </span>
            <span className="text-base font-bold code-font text-emerald-400 drop-shadow-[0_0_8px_rgba(16,185,129,0.2)]">
              {(stats.successRate * 100).toFixed(1)}%
            </span>
          </div>
          <div className="flex flex-col items-center border-r border-slate-900/60 pr-6">
            <span className="text-[10px] text-slate-500 font-medium uppercase tracking-wider mb-1">
              Token 消耗 (总/未缓存)
            </span>
            <div className="flex items-center gap-2.5 text-xs font-bold font-mono">
              <div
                className="flex items-center gap-1 text-purple-400 drop-shadow-[0_0_8px_rgba(139,92,246,0.15)]"
                title={`总吞吐 Token: ${stats.totalTokens}`}
              >
                <IconCpu />
                <span className="text-sm">
                  {stats.totalTokens >= 1000000
                    ? `${(stats.totalTokens / 1000000).toFixed(2)}M`
                    : stats.totalTokens >= 1000
                      ? `${(stats.totalTokens / 1000).toFixed(1)}K`
                      : stats.totalTokens}
                </span>
              </div>
              <span className="text-slate-700/80">/</span>
              {(() => {
                const uncached = stats.totalTokens - stats.totalCachedTokens;
                return (
                  <div
                    className="flex items-center gap-0.5 text-orange-400 drop-shadow-[0_0_8px_rgba(249,115,22,0.15)]"
                    title={`未缓存计费 Token: ${uncached}`}
                  >
                    <IconFlame />
                    <span className="text-sm">
                      {uncached >= 1000000
                        ? `${(uncached / 1000000).toFixed(2)}M`
                        : uncached >= 1000
                          ? `${(uncached / 1000).toFixed(1)}K`
                          : uncached}
                    </span>
                  </div>
                );
              })()}
            </div>
          </div>
          <div className="flex flex-col items-center">
            <span className="text-[10px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
              平均时延
            </span>
            <span className="text-base font-bold code-font text-orange-400 drop-shadow-[0_0_8px_rgba(249,115,22,0.2)]">
              {(stats.avgDurationMs / 1000).toFixed(2)}s
            </span>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={handleManualRefresh}
            className="p-2.5 rounded-xl bg-slate-900/50 border border-slate-800/80 hover:border-cyan-500/50 hover:bg-slate-800/40 hover:text-cyan-400 flex items-center justify-center transition-all active:scale-95 btn-glow-cyan"
            title="手动刷新"
          >
            <IconRefresh />
          </button>

          <button
            onClick={() => setIsSystemLogsOpen(true)}
            className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-slate-900/50 hover:bg-slate-800/40 border border-slate-800/80 hover:border-cyan-500/50 hover:text-cyan-400 text-sm font-semibold transition-all active:scale-95 btn-glow-cyan"
          >
            <IconTerminal />
            <span>系统日志</span>
          </button>

          <button
            onClick={() => setIsSettingsOpen(true)}
            className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-slate-900/50 hover:bg-slate-800/40 border border-slate-800/80 hover:border-cyan-500/50 hover:text-cyan-400 text-sm font-semibold transition-all active:scale-95 btn-glow-cyan"
          >
            <IconSettings />
            <span>上游配置</span>
          </button>
        </div>
      </header>

      {/* ==========================================
          主工作区 (左列表，右详情)
         ========================================== */}
      <main className="flex flex-1 overflow-hidden gap-5 md:gap-6 min-h-0">
        {/* 左栏：日志卡片列表 */}
        <section className="w-full md:w-[400px] lg:w-[440px] flex flex-col glass-panel p-5 gap-5 min-h-0 flex-shrink-0 border-slate-850 shadow-lg">
          {/* 搜索与过滤 */}
          <div className="flex flex-col gap-3 bg-slate-950/20 border border-slate-900/50 p-3.5 rounded-2xl">
            <div className="relative flex items-center">
              <span className="absolute left-3.5 text-cyan-400/80">
                <IconSearch />
              </span>
              <input
                type="text"
                placeholder="搜索请求/响应内容..."
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    setSearchKeyword(searchInput);
                    setCurrentPage(1);
                  }
                }}
                className="w-full pl-10 pr-3.5 py-2.5 bg-slate-950/40 focus:bg-slate-950/85 text-xs border border-slate-850/85 focus:border-cyan-500/60 rounded-xl transition-all outline-none placeholder:text-slate-500 text-slate-200"
              />
            </div>

            <div className="flex gap-2.5">
              <select
                value={providerFilter}
                onChange={(e) => {
                  setProviderFilter(e.target.value);
                  setCurrentPage(1);
                }}
                className="flex-1 bg-slate-950/40 focus:bg-slate-950/85 text-xs border border-slate-850/85 focus:border-cyan-500/60 rounded-xl p-2.5 outline-none transition-all text-slate-300"
              >
                <option value="all">所有厂商</option>
                <option value="openai">OpenAI</option>
                <option value="anthropic">Anthropic</option>
                <option value="openai-responses">OpenAI Responses</option>
              </select>

              <select
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value);
                  setCurrentPage(1);
                }}
                className="flex-1 bg-slate-950/40 focus:bg-slate-950/85 text-xs border border-slate-850/85 focus:border-cyan-500/60 rounded-xl p-2.5 outline-none transition-all text-slate-300"
              >
                <option value="all">所有状态</option>
                <option value="success">成功 (2xx)</option>
                <option value="error">失败 (非2xx)</option>
              </select>
            </div>
          </div>

          {/* 滚动卡片列表 */}
          <div className="flex-1 overflow-y-auto flex flex-col gap-3 pr-1.5 min-h-0 scroll-isolated">
            {isListLoading && logs.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-20 text-slate-500 text-xs gap-2">
                <IconRefresh className="animate-spin text-cyan-400" />
                <span>加载请求日志中...</span>
              </div>
            ) : logs.length === 0 ? (
              <div className="text-center py-20 text-slate-500 text-xs">
                暂无匹配的请求记录
              </div>
            ) : (
              <div className="flex flex-col gap-1">
                {logs.map((log) => renderLogCard(log))}
              </div>
            )}
          </div>

          {/* 分页组件 */}
          {totalLogs > pageSize && (
            <div className="flex justify-between items-center pt-2 border-t border-slate-850 text-xs text-slate-400">
              <span>共 {totalLogs} 条</span>
              <div className="flex gap-1">
                <button
                  onClick={() => handlePageChange(currentPage - 1)}
                  disabled={currentPage === 1}
                  className="px-2.5 py-1 rounded bg-slate-800 disabled:opacity-30 disabled:cursor-not-allowed hover:bg-slate-700"
                >
                  上一页
                </button>
                <span className="px-2 py-1 code-font text-cyan-400">
                  {currentPage} / {Math.ceil(totalLogs / pageSize)}
                </span>
                <button
                  onClick={() => handlePageChange(currentPage + 1)}
                  disabled={currentPage === Math.ceil(totalLogs / pageSize)}
                  className="px-2.5 py-1 rounded bg-slate-800 disabled:opacity-30 disabled:cursor-not-allowed hover:bg-slate-700"
                >
                  下一页
                </button>
              </div>
            </div>
          )}
        </section>

        {/* 右栏：追踪和调试面板 */}
        <section className="flex-1 flex flex-col glass-panel p-6 min-h-0 border-slate-850 shadow-lg min-w-0">
          {!selectedLogId ? (
            <div className="flex-1 flex flex-col items-center justify-center text-slate-500 gap-4">
              <div className="w-12 h-12 rounded-full bg-slate-900/60 border border-slate-800 flex items-center justify-center text-slate-400">
                <IconTerminal />
              </div>
              <p className="text-xs font-medium tracking-wide">
                请在左侧选择一个请求来查看详细的时序、请求和响应 Trace 细节
              </p>
            </div>
          ) : isLogLoading ? (
            <div className="flex-1 flex flex-col items-center justify-center text-slate-500 gap-4">
              <IconRefresh className="animate-spin text-cyan-400 w-8 h-8" />
              <p className="text-xs font-medium tracking-wide">
                加载请求结构化详情中...
              </p>
            </div>
          ) : !selectedLog ? (
            <div className="flex-1 flex flex-col items-center justify-center text-red-400 gap-3">
              <span>加载日志失败，请稍后重试</span>
            </div>
          ) : (
            <div className="flex-1 flex flex-col overflow-hidden gap-6 min-h-0">
              {/* 详情头部元数据卡片 (Grid 布局，突出核心指标) */}
              <div className="p-5 bg-slate-950/35 border border-slate-900/60 rounded-2xl flex flex-col gap-4 lg:flex-row lg:justify-between lg:items-center">
                <div className="flex flex-col gap-2 min-w-0">
                  <div className="flex flex-wrap items-center gap-2.5">
                    <span className="text-base font-bold font-mono text-cyan-400 tracking-tight truncate">
                      {selectedLog.model}
                    </span>
                    <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-900 border border-slate-800 text-slate-400 uppercase font-bold tracking-wider shrink-0">
                      {selectedLog.provider}
                    </span>
                  </div>
                  <span className="text-[10px] text-slate-500 font-mono truncate">
                    {selectedLog.sessionId
                      ? `会话 #${selectedLog.sessionSeq} (ID: ${selectedLog.id})`
                      : `ID: ${selectedLog.id}`}{" "}
                    • Path: {selectedLog.path}
                  </span>
                </div>

                <div
                  className={`grid grid-cols-2 ${selectedLog.cachedTokens > 0 || selectedLog.cacheReadTokens > 0 || selectedLog.cacheCreationTokens > 0 ? "sm:grid-cols-5" : "sm:grid-cols-4"} gap-3 text-xs font-mono shrink-0`}
                >
                  <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-orange-950/15 border border-orange-900/15">
                    <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
                      请求耗时
                    </span>
                    <span className="font-bold text-sm text-orange-400">
                      {(selectedLog.durationMs / 1000).toFixed(2)}s
                    </span>
                  </div>
                  <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-cyan-950/15 border border-cyan-900/15">
                    <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
                      输入 Tokens
                    </span>
                    <span className="font-bold text-sm text-cyan-400">
                      {selectedLog.inputTokens}
                    </span>
                    {selectedLog.cachedTokens > 0 && (
                      <span className="text-[8px] text-slate-450 mt-0.5">
                        实际计算:{" "}
                        {selectedLog.inputTokens - selectedLog.cachedTokens}
                      </span>
                    )}
                  </div>
                  {(selectedLog.cachedTokens > 0 ||
                    selectedLog.cacheReadTokens > 0 ||
                    selectedLog.cacheCreationTokens > 0) && (
                    <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-indigo-950/20 border border-indigo-900/30 text-indigo-400">
                      <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
                        缓存 Tokens
                      </span>
                      <span className="font-bold text-sm text-indigo-400">
                        {selectedLog.cachedTokens}
                      </span>
                      <span className="text-[8px] text-slate-450 mt-0.5 truncate max-w-full">
                        {selectedLog.provider === "openai" ||
                        selectedLog.provider === "openai-responses"
                          ? "命中缓存"
                          : `读缓存 (另写 ${selectedLog.cacheCreationTokens})`}
                      </span>
                    </div>
                  )}
                  <div className="flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl bg-purple-950/15 border border-purple-900/15">
                    <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
                      输出 Tokens
                    </span>
                    <span className="font-bold text-sm text-purple-400">
                      {selectedLog.outputTokens}
                    </span>
                  </div>
                  <div
                    className={`flex flex-col items-start lg:items-end px-3.5 py-2 rounded-xl border ${
                      selectedLog.statusCode >= 200 &&
                      selectedLog.statusCode < 300
                        ? "bg-emerald-950/15 border-emerald-900/15 text-emerald-400"
                        : "bg-red-950/15 border-red-900/15 text-red-400"
                    }`}
                  >
                    <span className="text-[9px] text-slate-500 font-medium uppercase tracking-wider mb-0.5">
                      状态代码
                    </span>
                    <span className="font-bold text-sm">
                      {selectedLog.statusCode}
                    </span>
                  </div>
                </div>
              </div>

              {/* 时序控制 Tab */}
              <div className="flex border-b border-slate-900 pb-1 gap-6 shrink-0">
                <button
                  onClick={() => setActiveTab("trace")}
                  className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                    activeTab === "trace"
                      ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                      : "border-transparent text-slate-400 hover:text-slate-200"
                  }`}
                >
                  时序追踪 (Trace)
                </button>
                <button
                  onClick={() => setActiveTab("prompt")}
                  className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                    activeTab === "prompt"
                      ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                      : "border-transparent text-slate-400 hover:text-slate-200"
                  }`}
                >
                  请求消息 (
                  {selectedLog.prompt
                    ? selectedLog.prompt.filter((p) => p.role !== "system")
                        .length
                    : 0}
                  )
                </button>
                {selectedLog.sessionId && (
                  <button
                    onClick={() => setActiveTab("session")}
                    className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                      activeTab === "session"
                        ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                        : "border-transparent text-slate-400 hover:text-slate-200"
                    }`}
                  >
                    会话历史 (Session){totalSessionTokens > 0 ? ` • ${totalSessionTokens} T` : ""}
                  </button>
                )}
                {selectedLog.tools && selectedLog.tools.length > 0 && (
                  <button
                    onClick={() => setActiveTab("tools")}
                    className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                      activeTab === "tools"
                        ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                        : "border-transparent text-slate-400 hover:text-slate-200"
                    }`}
                  >
                    定义的工具 ({selectedLog.tools.length})
                  </button>
                )}
                <button
                  onClick={() => setActiveTab("raw_req")}
                  className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                    activeTab === "raw_req"
                      ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                      : "border-transparent text-slate-400 hover:text-slate-200"
                  }`}
                >
                  原始请求 JSON
                </button>
                <button
                  onClick={() => setActiveTab("raw_resp")}
                  className={`pb-2.5 text-xs font-bold border-b-2 transition-all ${
                    activeTab === "raw_resp"
                      ? "border-cyan-400 text-cyan-400 drop-shadow-[0_0_8px_rgba(6,182,212,0.3)]"
                      : "border-transparent text-slate-400 hover:text-slate-200"
                  }`}
                >
                  原始响应 Raw
                </button>
              </div>

              {/* Tab 内容区 */}
              <div className="flex-1 overflow-y-auto pr-1.5 min-h-0 flex flex-col gap-4 scroll-isolated">
                {/* 1. 时序追踪 (Trace) */}
                {activeTab === "trace" && (
                  <div className="flex flex-col gap-5">
                    {/* 可视化分段 Timeline */}
                    <div className="p-4 bg-slate-950/20 border border-slate-900 rounded-xl flex flex-col gap-3">
                      <div className="flex justify-between items-center text-xs text-slate-400 font-mono">
                        <span className="font-semibold text-slate-300">
                          时序可视化时间线
                        </span>
                        <span className="text-cyan-400 font-bold">
                          {(selectedLog.durationMs / 1000).toFixed(2)}s
                        </span>
                      </div>

                      {/* 分段时间条 (两段式网络连接与传输模型响应) */}
                      <div className="h-7 w-full bg-slate-950/80 rounded-lg overflow-hidden flex relative p-0.5 border border-slate-900">
                        {/* 150ms 假定的连接建立 */}
                        <div
                          className="h-full bg-indigo-600/40 rounded-l-md border-r border-slate-900/30 flex items-center justify-center text-[8px] font-bold text-indigo-300 shrink-0"
                          style={{
                            width: `${Math.min(15, Math.max(5, (150 / selectedLog.durationMs) * 100))}%`,
                          }}
                        >
                          连接
                        </div>
                        {/* 剩余的大部分为响应传输段 */}
                        <div
                          className="h-full timeline-shimmer rounded-r-md"
                          style={{ flex: 1 }}
                        />
                        <div className="absolute inset-0 flex justify-end items-center pr-4 text-[9px] font-bold text-white uppercase tracking-wider drop-shadow-md">
                          <span>
                            传输流 •{" "}
                            {(selectedLog.durationMs / 1000).toFixed(2)}s
                          </span>
                        </div>
                      </div>

                      {/* 如果有工具调用，显示时序序列图 */}
                      {selectedLog.response &&
                        selectedLog.response.tool_calls &&
                        selectedLog.response.tool_calls.length > 0 && (
                          <div className="mt-1.5 flex flex-col gap-2.5">
                            <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider">
                              执行的 Tool Call Spans
                            </span>
                            <div className="flex flex-col gap-2">
                              {selectedLog.response.tool_calls.map(
                                (tc, idx) => (
                                  <div
                                    key={idx}
                                    className="flex items-center justify-between px-3 py-2.5 bg-slate-950/40 border border-slate-900 rounded-xl text-xs"
                                  >
                                    <div className="flex items-center gap-2 min-w-0">
                                      <span className="text-[9px] font-bold px-2 py-0.5 bg-orange-950/80 border border-orange-900/20 text-orange-400 rounded-md shrink-0">
                                        TOOL
                                      </span>
                                      <span className="font-mono text-slate-200 font-semibold truncate">
                                        {tc.name}
                                      </span>
                                    </div>
                                    <span className="text-[9px] text-slate-500 font-mono shrink-0">
                                      id: {tc.id}
                                    </span>
                                  </div>
                                ),
                              )}
                            </div>
                          </div>
                        )}
                    </div>
                  </div>
                )}

                {/* 2. 对话结构化展示 */}
                {activeTab === "prompt" && (
                  <div className="flex flex-col gap-4">
                    {/* 折叠 System Prompt 或定义工具 */}
                    <div className="flex flex-col gap-2.5">
                      <button
                        onClick={() => setShowSystemPrompt(!showSystemPrompt)}
                        className="flex justify-between items-center px-4 py-3 bg-slate-950/40 hover:bg-slate-900/60 border border-slate-900/85 rounded-xl text-xs text-slate-300 transition-all active:scale-[0.99]"
                      >
                        <span className="flex items-center gap-2">
                          <IconTerminal />
                          <span className="font-semibold">System 提示词</span>
                        </span>
                        <span className="text-slate-450">
                          {showSystemPrompt ? (
                            <IconChevronDown />
                          ) : (
                            <IconChevronRight />
                          )}
                        </span>
                      </button>
                      {showSystemPrompt && (
                        <div className="p-4 bg-slate-950/60 border border-slate-900 rounded-xl text-xs font-mono text-slate-400 whitespace-pre-wrap leading-relaxed">
                          {selectedLog.prompt &&
                          selectedLog.prompt.filter((p) => p.role === "system")
                            .length > 0
                            ? selectedLog.prompt
                                .filter((p) => p.role === "system")
                                .map((p) => p.content)
                                .join("\n")
                            : "无系统提示词"}
                        </div>
                      )}
                    </div>

                    {/* 消息时间线 */}
                    <div className="flex flex-col gap-4">
                      <h3 className="text-xs uppercase font-bold tracking-wider text-slate-400">
                        消息流记录
                      </h3>
                      {selectedLog.prompt &&
                        selectedLog.prompt
                          .filter((p) => p.role !== "system")
                          .map((m, idx) => {
                            const isUser = m.role === "user";
                            const isTool = m.role === "tool";
                            const isAssistant = m.role === "assistant";

                            return (
                              <div
                                key={idx}
                                className={`p-5 rounded-2xl border flex flex-col gap-3.5 shadow-sm ${
                                  isUser
                                    ? "bg-cyan-950/10 border-cyan-950/40"
                                    : isTool
                                      ? "bg-orange-950/10 border-orange-950/40"
                                      : "bg-slate-950/25 border-slate-900"
                                }`}
                              >
                                <div className="flex justify-between items-center">
                                  <span
                                    className={`text-[10px] font-bold uppercase tracking-wider px-2.5 py-1 rounded-md border ${
                                      isUser
                                        ? "bg-cyan-950/60 border-cyan-900/30 text-cyan-400"
                                        : isTool
                                          ? "bg-orange-950/60 border-orange-900/30 text-orange-400"
                                          : "bg-slate-900/60 border-slate-800/80 text-cyan-400"
                                    }`}
                                  >
                                    {isUser
                                      ? "👤 User"
                                      : isTool
                                        ? "🛠️ Tool Response"
                                        : "🤖 Assistant"}
                                    {m.name && (
                                      <span className="text-slate-400 font-mono text-[9px] ml-1">
                                        ({m.name})
                                      </span>
                                    )}
                                  </span>
                                  {m.tool_call_id && (
                                    <span className="text-[9px] text-slate-500 font-mono">
                                      ID: {m.tool_call_id}
                                    </span>
                                  )}
                                </div>

                                {/* 如果有思维链，进行展示 */}
                                {isAssistant && m.thinking && (
                                  <div className="p-3.5 bg-slate-900/40 border border-slate-850/60 rounded-xl text-xs text-slate-450/90 italic whitespace-pre-wrap border-l-2 border-l-purple-500/60">
                                    <div className="text-[9px] font-bold text-purple-400 uppercase tracking-wider not-italic mb-1">
                                      思考过程 (Thinking)
                                    </div>
                                    {m.thinking}
                                  </div>
                                )}

                                {m.content && (
                                  <div className="text-sm leading-relaxed tracking-wide text-slate-200 whitespace-pre-wrap mt-0.5">
                                    {m.content}
                                  </div>
                                )}

                                {/* 工具调用展示 */}
                                {isAssistant &&
                                  m.tool_calls &&
                                  m.tool_calls.length > 0 && (
                                    <div className="mt-1 flex flex-col gap-2 font-mono">
                                      <div className="text-[9px] font-bold text-orange-400 uppercase tracking-wider">
                                        触发工具调用 (Tool Calls)
                                      </div>
                                      {m.tool_calls.map((tc, tcIdx) => (
                                        <div
                                          key={tcIdx}
                                          className="p-3 bg-slate-950/40 border border-slate-900 rounded-xl text-xs text-slate-300"
                                        >
                                          <div>
                                            <span className="text-slate-500">
                                              Method:
                                            </span>{" "}
                                            <span className="text-orange-400 font-semibold">
                                              {tc.name}
                                            </span>
                                          </div>
                                          <div className="mt-1">
                                            <span className="text-slate-500">
                                              Args:
                                            </span>{" "}
                                            <span className="text-emerald-400/90">
                                              {tc.arguments}
                                            </span>
                                          </div>
                                          {renderSubLogsForTool(tc.id)}
                                        </div>
                                      ))}
                                    </div>
                                  )}
                              </div>
                            );
                          })}

                      {/* 渲染当次最新大模型回复 */}
                      {selectedLog.response &&
                        (selectedLog.response.content ||
                          selectedLog.response.thinking ||
                          (selectedLog.response.tool_calls &&
                            selectedLog.response.tool_calls.length > 0)) && (
                          <div className="p-5 rounded-2xl border flex flex-col gap-3.5 shadow-sm bg-slate-950/25 border-slate-900 border-l-[3px] border-l-cyan-500/50">
                            <div className="flex justify-between items-center">
                              <span className="text-[10px] font-bold uppercase tracking-wider px-2.5 py-1 rounded-md border bg-slate-900/60 border-slate-800/80 text-cyan-400">
                                🤖 Assistant (最新回复)
                              </span>
                            </div>

                            {/* 如果有思维链，进行展示 */}
                            {selectedLog.response.thinking && (
                              <div className="p-3.5 bg-slate-900/40 border border-slate-850/60 rounded-xl text-xs text-slate-450/90 italic whitespace-pre-wrap border-l-2 border-l-purple-500/60">
                                <div className="text-[9px] font-bold text-purple-400 uppercase tracking-wider not-italic mb-1">
                                  思考过程 (Thinking)
                                </div>
                                {selectedLog.response.thinking}
                              </div>
                            )}

                            {/* 消息正文 */}
                            {selectedLog.response.content && (
                              <div className="text-sm leading-relaxed tracking-wide text-slate-200 whitespace-pre-wrap mt-0.5">
                                {selectedLog.response.content}
                              </div>
                            )}

                            {/* 工具调用展示 */}
                            {selectedLog.response.tool_calls &&
                              selectedLog.response.tool_calls.length > 0 && (
                                <div className="mt-1 flex flex-col gap-2 font-mono">
                                  <div className="text-[9px] font-bold text-orange-400 uppercase tracking-wider">
                                    触发工具调用 (Tool Calls)
                                  </div>
                                  {selectedLog.response.tool_calls.map(
                                    (tc, tcIdx) => (
                                      <div
                                        key={tcIdx}
                                        className="p-3 bg-slate-950/40 border border-slate-900 rounded-xl text-xs text-slate-300"
                                      >
                                        <div>
                                          <span className="text-slate-500">
                                            Method:
                                          </span>{" "}
                                          <span className="text-orange-400 font-semibold">
                                            {tc.name}
                                          </span>
                                        </div>
                                        <div className="mt-1">
                                          <span className="text-slate-500">
                                            Args:
                                          </span>{" "}
                                          <span className="text-emerald-400/90">
                                            {tc.arguments}
                                          </span>
                                        </div>
                                        {renderSubLogsForTool(tc.id)}
                                      </div>
                                    ),
                                  )}
                                </div>
                              )}
                          </div>
                        )}
                    </div>
                  </div>
                )}

                {/* 会话历史 (Session Trace) */}
                {activeTab === "session" && selectedLog.sessionId && (
                  <div className="flex flex-col gap-5 select-none font-sans">
                    <div className="text-[10px] text-slate-400 font-bold uppercase tracking-wider mb-1 flex justify-between items-center bg-slate-950/20 px-4 py-2 border border-slate-900/50 rounded-xl">
                      <div className="flex items-center gap-2">
                        <span>🔗 会话轮次目录 / 点击可快速定位详情</span>
                        <span className="font-mono text-cyan-400">
                          ({sessionLogs.length} 轮对话 • 未缓存输入 {totalSessionUncachedInputTokens} / 缓存 {totalSessionCachedTokens} / 输出 {totalSessionOutputTokens} / 共 {totalSessionTokens} Tokens)
                        </span>
                      </div>
                      <button
                        onClick={() => handleDeleteSessionLogs(selectedLog.sessionId)}
                        className="flex items-center gap-1 text-[10px] font-bold text-red-400 hover:text-red-300 px-2.5 py-1 rounded-lg bg-red-950/20 hover:bg-red-950/40 border border-red-900/30 transition-all active:scale-95 cursor-pointer"
                        title="删除该会话所有请求日志"
                      >
                        <IconTrash />
                        <span>删除会话所有请求</span>
                      </button>
                    </div>

                    {isSessionLoading ? (
                      <div className="flex flex-col items-center justify-center py-20 text-slate-500 text-xs gap-2">
                        <IconRefresh className="animate-spin text-cyan-400" />
                        <span>加载会话目录中...</span>
                      </div>
                    ) : sessionLogs.length === 0 ? (
                      <div className="text-center py-20 text-slate-500 text-xs">
                        该会话暂无历史记录
                      </div>
                    ) : (
                      <div className="flex flex-col gap-3 relative pl-4 border-l border-slate-800/60 ml-3">
                        {sessionLogs.map((sLog, idx) => {
                          const isCurrent = sLog.id === selectedLog.id;
                          const isSuccess = sLog.statusCode >= 200 && sLog.statusCode < 300;
                          
                           // 提取输入消息 (逆序查找最新 user 消息或 tool 消息)
                           let inputMsg = null;
                           if (sLog.prompt && sLog.prompt.length > 0) {
                             for (let i = sLog.prompt.length - 1; i >= 0; i--) {
                               const m = sLog.prompt[i];
                               if (m.role === "user" || m.role === "tool") {
                                 inputMsg = m;
                                 break;
                               }
                             }
                           }

                           let inputText = "";
                           let inputRoleLabel = "";
                           let inputRoleBadgeColor = "";
                           if (inputMsg) {
                             if (inputMsg.role === "user") {
                               inputRoleLabel = "👤 User";
                               inputRoleBadgeColor = "text-cyan-400 bg-cyan-950/30 border-cyan-900/30";
                               inputText = simplifyUserPrompt(inputMsg.content || "");
                             } else {
                               inputRoleLabel = `🛠️ Tool (${inputMsg.name || "unknown"})`;
                               inputRoleBadgeColor = "text-orange-400 bg-orange-950/30 border-orange-900/30";
                               inputText = inputMsg.content || "";
                             }
                           } else {
                             inputRoleLabel = "💬 Info";
                             inputRoleBadgeColor = "text-slate-400 bg-slate-900/40 border-slate-800/40";
                             inputText = sLog.sessionSummary || "";
                           }

                           const shortInput = inputText && inputText.trim() !== ""
                             ? (inputText.length > 95 ? inputText.substring(0, 95) + "..." : inputText)
                             : "无详细输入";

                           // 提取输出消息
                           let outputText = "";
                           let outputLabel = "";
                           let outputBadgeColor = "";
                           if (!isSuccess) {
                             outputLabel = "❌ Error";
                             outputBadgeColor = "text-red-400 bg-red-950/30 border-red-900/30";
                             outputText = sLog.errorMessage || `HTTP Status ${sLog.statusCode}`;
                           } else if (sLog.response) {
                             const resp = sLog.response;
                             const hasToolCalls = resp.tool_calls && resp.tool_calls.length > 0;
                             const hasContent = resp.content && resp.content.trim() !== "";
                             const hasThinking = resp.thinking && resp.thinking.trim() !== "";

                             if (hasToolCalls) {
                               outputLabel = "🔧 Tool Call";
                               outputBadgeColor = "text-yellow-400 bg-yellow-950/30 border-yellow-900/30";
                               const toolNames = resp.tool_calls.map(tc => tc.name).join(", ");
                               outputText = `调用工具: [${toolNames}]`;
                             } else if (hasContent) {
                               outputLabel = "🤖 Assistant";
                               outputBadgeColor = "text-emerald-400 bg-emerald-950/30 border-emerald-900/30";
                               outputText = resp.content.trim();
                             } else if (hasThinking) {
                               outputLabel = "💭 Thinking";
                               outputBadgeColor = "text-purple-400 bg-purple-950/30 border-purple-900/30";
                               outputText = resp.thinking.trim();
                             } else {
                               outputLabel = "🤖 Assistant";
                               outputBadgeColor = "text-slate-450 bg-slate-900/30 border-slate-800/30";
                               outputText = "无输出内容";
                             }
                           } else {
                             outputLabel = "🤖 Assistant";
                             outputBadgeColor = "text-slate-455 bg-slate-900/30 border-slate-800/30";
                             outputText = "无输出内容";
                           }

                           const shortOutput = outputText && outputText.trim() !== ""
                             ? (outputText.length > 95 ? outputText.substring(0, 95) + "..." : outputText)
                             : "无详细输出";

                           return (
                             <div
                               key={sLog.id}
                               onClick={() => setModalLogId(sLog.id)}
                               className={`p-3.5 rounded-xl border cursor-pointer relative transition-all flex flex-col gap-2 ${
                                 isCurrent
                                   ? "bg-slate-900/70 border-cyan-500/60 shadow-[0_0_12px_rgba(6,182,212,0.15)] translate-x-0.5"
                                   : "bg-slate-950/25 border-slate-900/80 hover:bg-slate-900/40 hover:border-slate-800/80"
                               }`}
                             >
                               {/* 挂件：步骤条时间线小点/圆环 */}
                               <div
                                 className={`absolute -left-[22px] top-6 w-2.5 h-2.5 rounded-full border-2 transition-all ${
                                   isCurrent
                                     ? "bg-cyan-400 border-slate-950 scale-110"
                                     : "bg-slate-800 border-slate-950 hover:bg-slate-700"
                                 }`}
                               />

                               {/* 行一：轮次标题与模型 */}
                               <div className="flex justify-between items-center">
                                 <div className="flex items-center gap-2">
                                   <span className={`text-[10px] font-mono font-bold px-2 py-0.5 rounded transition-all ${
                                     isCurrent ? "bg-cyan-950 text-cyan-400" : "bg-slate-900 text-slate-450"
                                   }`}>
                                     Turn #{sLog.sessionSeq || idx + 1}
                                   </span>
                                   <span className="text-[9px] font-mono text-slate-500 truncate max-w-[120px] md:max-w-[180px]">
                                     {sLog.model}
                                   </span>
                                 </div>
                                 <div className="flex items-center gap-2.5 text-[9px] font-mono text-slate-450">
                                   <span>{sLog.createdAt ? sLog.createdAt.split(" ")[1] : ""}</span>
                                   <span>{(sLog.durationMs / 1000).toFixed(2)}s</span>
                                   {sLog.totalTokens > 0 && (
                                     <>
                                       <span>•</span>
                                       <span>
                                         {sLog.totalTokens} tkn
                                         {sLog.cachedTokens > 0 && ` (未缓存 ${sLog.inputTokens - sLog.cachedTokens})`}
                                       </span>
                                     </>
                                   )}
                                   <span className={`font-bold px-1 py-0.2 rounded ${isSuccess ? "text-emerald-400 bg-emerald-950/20" : "text-red-400 bg-red-950/20"}`}>
                                     {sLog.statusCode}
                                   </span>
                                 </div>
                               </div>

                               {/* 行二：精准展示的输入输出流程图 (IN/OUT 极具 Premium 感) */}
                               <div className="flex flex-col gap-2 mt-1 select-text">
                                 {/* 输入段 */}
                                 <div className="flex items-start gap-2 text-[11px] leading-relaxed">
                                   <span className={`shrink-0 text-[8px] font-bold px-1.5 py-0.5 rounded border ${inputRoleBadgeColor} font-mono w-[60px] text-center`}>
                                     IN
                                   </span>
                                   <div className="flex-1 text-slate-300 min-w-0 truncate">
                                     <span className="text-slate-455 font-semibold mr-1">{inputRoleLabel}:</span>
                                     <span className={isCurrent ? "text-slate-100 font-medium" : "text-slate-350 hover:text-slate-200"}>
                                       {shortInput}
                                     </span>
                                   </div>
                                 </div>

                                 {/* 指示箭头 */}
                                 <div className="pl-[23px] text-slate-650 flex items-center -my-1">
                                   <span className="text-[9px] text-slate-600">↓</span>
                                 </div>

                                 {/* 输出段 */}
                                 <div className="flex items-start gap-2 text-[11px] leading-relaxed">
                                   <span className={`shrink-0 text-[8px] font-bold px-1.5 py-0.5 rounded border ${outputBadgeColor} font-mono w-[60px] text-center`}>
                                     OUT
                                   </span>
                                   <div className="flex-1 text-slate-300 min-w-0 truncate">
                                     <span className="text-slate-455 font-semibold mr-1">{outputLabel}:</span>
                                     <span className={isCurrent ? "text-slate-100 font-medium" : "text-slate-350 hover:text-slate-200"}>
                                       {shortOutput}
                                     </span>
                                   </div>
                                 </div>
                               </div>
                             </div>
                           );
                        })}
                      </div>
                    )}
                  </div>
                )}

                {/* 3. 原始请求 Raw */}
                {activeTab === "raw_req" && (
                  <CodeBlock
                    code={formatJSON(selectedLog.rawRequest)}
                    language="Raw HTTP Request Body"
                  />
                )}

                {/* 4. 原始响应 Raw */}
                {activeTab === "raw_resp" && (
                  <CodeBlock
                    code={formatJSON(selectedLog.rawResponse)}
                    language="Raw Upstream Response"
                  />
                )}

                {/* 5. 已定义工具 (Tools) */}
                {activeTab === "tools" &&
                  selectedLog.tools &&
                  selectedLog.tools.length > 0 && (
                    <div className="flex flex-col lg:flex-row gap-5 flex-1 min-h-0 h-full overflow-hidden">
                      {/* 左侧：工具搜索与列表 */}
                      <div className="w-full lg:w-[240px] flex flex-col gap-3 shrink-0 bg-slate-950/20 border border-slate-900/40 p-3.5 rounded-2xl min-h-0 overflow-y-auto">
                        <div className="relative flex items-center">
                          <span className="absolute left-3.5 text-slate-500">
                            <IconSearch />
                          </span>
                          <input
                            type="text"
                            placeholder="搜索工具名..."
                            value={toolSearchQuery}
                            onChange={(e) => {
                              setToolSearchQuery(e.target.value);
                              setSelectedToolIndex(0);
                            }}
                            className="w-full pl-9 pr-3.5 py-1.5 bg-slate-950/40 focus:bg-slate-950/85 text-[11px] border border-slate-850/85 focus:border-cyan-500/60 rounded-lg outline-none placeholder:text-slate-500 text-slate-200"
                          />
                        </div>

                        {/* 工具列表 */}
                        <div className="flex flex-col gap-1.5 overflow-y-auto flex-1 pr-1 font-mono text-[11px]">
                          {selectedLog.tools
                            .map((tool, index) => ({ tool, index }))
                            .filter(({ tool }) =>
                              tool.name
                                .toLowerCase()
                                .includes(toolSearchQuery.toLowerCase()),
                            )
                            .map(({ tool, index }) => {
                              const isSelected = selectedToolIndex === index;
                              return (
                                <button
                                  key={index}
                                  onClick={() => setSelectedToolIndex(index)}
                                  className={`text-left px-3 py-2 rounded-lg transition-all truncate border shrink-0 ${
                                    isSelected
                                      ? "bg-cyan-950/25 border-cyan-500/50 text-cyan-400 font-bold"
                                      : "bg-transparent border-transparent text-slate-400 hover:text-slate-200 hover:bg-slate-900/30"
                                  }`}
                                  title={tool.name}
                                >
                                  🔧 {tool.name}
                                </button>
                              );
                            })}
                          {selectedLog.tools.filter((t) =>
                            t.name
                              .toLowerCase()
                              .includes(toolSearchQuery.toLowerCase()),
                          ).length === 0 && (
                            <div className="text-center text-slate-500 py-6 text-[10px]">
                              未找到相关工具
                            </div>
                          )}
                        </div>
                      </div>

                      {/* 右侧：被选中工具的 Schema 和详情 */}
                      {(() => {
                        const selectedTool =
                          selectedLog.tools[selectedToolIndex] ||
                          selectedLog.tools[0];
                        if (!selectedTool) return null;
                        return (
                          <div className="flex-1 flex flex-col gap-3 min-w-0 overflow-y-auto">
                            <div className="p-4 bg-slate-950/20 border border-slate-900 rounded-xl flex flex-col gap-2.5">
                              <div className="flex items-center gap-2">
                                <span className="text-[10px] font-bold px-2 py-0.5 bg-cyan-950/80 border border-cyan-900/20 text-cyan-400 rounded-md">
                                  TOOL
                                </span>
                                <span className="text-sm font-bold text-slate-200 font-mono">
                                  {selectedTool.name}
                                </span>
                              </div>
                              {selectedTool.description && (
                                <p className="text-xs text-slate-400 leading-relaxed bg-slate-950/45 p-3 rounded-lg border border-slate-900/60">
                                  <span className="font-semibold text-slate-350 block mb-1">
                                    功能描述：
                                  </span>
                                  {selectedTool.description}
                                </p>
                              )}
                            </div>
                            <div className="flex flex-col shrink-0">
                              <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider block mb-1">
                                Parameters Schema
                              </span>
                              <CodeBlock
                                code={formatJSON(selectedTool.parameters)}
                                language="parameters.schema.json"
                                />
                            </div>
                          </div>
                        );
                      })()}
                    </div>
                  )}
              </div>
            </div>
          )}
        </section>
      </main>

      {/* ==========================================
          System Logs Modal 系统运行日志查看器 (终端质感)
         ========================================== */}
      {isSystemLogsOpen && (
        <div className="fixed inset-0 bg-slate-950/70 z-50 flex items-center justify-center p-4 modal-fade-in">
          <div className="glass-panel w-full max-w-4xl h-[80vh] overflow-hidden flex flex-col border border-slate-800/85 shadow-[0_20px_50px_rgba(0,0,0,0.85)] modal-scale-up">
            <div className="flex justify-between items-center px-5 py-4 border-b border-slate-900 shrink-0">
              <h2 className="text-base font-bold text-slate-100 flex items-center gap-2">
                <IconTerminal />
                <span>系统运行日志</span>
                {isSystemLogsLoading && (
                  <span className="flex items-center ml-1">
                    <IconRefresh className="animate-spin text-cyan-400 w-3.5 h-3.5" />
                  </span>
                )}
              </h2>
              <button
                onClick={() => setIsSystemLogsOpen(false)}
                className="text-slate-450 hover:text-slate-200 text-2xl font-light transition-colors active:scale-95 leading-none"
              >
                &times;
              </button>
            </div>

            {/* 控制面板 */}
            <div className="flex flex-wrap items-center justify-between gap-3 px-5 py-3 bg-slate-950/30 border-b border-slate-900/60 shrink-0 text-xs">
              <div className="flex items-center gap-4">
                <button
                  onClick={fetchSystemLogs}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-900/60 hover:bg-slate-800/40 border border-slate-800/80 hover:border-cyan-500/50 hover:text-cyan-400 transition-all active:scale-95"
                >
                  <IconRefresh className={isSystemLogsLoading ? "animate-spin" : ""} />
                  <span>刷新</span>
                </button>
                <button
                  onClick={handleClearSystemLogs}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-950/20 hover:bg-red-950/40 border border-red-900/30 text-red-400 hover:text-red-300 transition-all active:scale-95"
                >
                  <IconTrash />
                  <span>清空日志</span>
                </button>
              </div>

              <div className="flex items-center gap-2">
                <label className="flex items-center gap-2 cursor-pointer text-slate-400 hover:text-slate-200 select-none">
                  <input
                    type="checkbox"
                    checked={autoScroll}
                    onChange={(e) => setAutoScroll(e.target.checked)}
                    className="accent-cyan-500 cursor-pointer rounded"
                  />
                  <span>自动滚动到底部</span>
                </label>
              </div>
            </div>

            {/* 终端日志展示区 */}
            <div
              ref={logsContainerRef}
              className="flex-1 p-5 overflow-y-auto bg-black/85 code-font text-xs text-slate-350 leading-relaxed scroll-isolated whitespace-pre-wrap select-text"
            >
              {systemLogs ? (
                systemLogs
              ) : (
                <div className="text-slate-600 text-center py-20">暂无系统运行日志记录</div>
              )}
            </div>

            <div className="px-5 py-3.5 border-t border-slate-900 bg-slate-950/40 flex justify-between items-center text-[10px] text-slate-500 shrink-0 font-sans">
              <span>日志保存于 ~/.llm_tracer/system.log (仅返回最新 100KB)</span>
              <span className="flex items-center gap-1.5">
                <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                <span>实时接收中 (每2.5秒更新)</span>
              </span>
            </div>
          </div>
        </div>
      )}

      {/* ==========================================
          Settings Modal 配置模态框 (极光玻璃态悬浮窗)
         ========================================== */}
      {isSettingsOpen && (
        <div className="fixed inset-0 bg-slate-950/70 z-50 flex items-center justify-center p-4 modal-fade-in">
          <div className="glass-panel w-full max-w-4xl overflow-hidden flex flex-col border border-slate-800/85 shadow-[0_20px_50px_rgba(0,0,0,0.85)] modal-scale-up">
            <div className="flex justify-between items-center px-5 py-4 border-b border-slate-900 shrink-0">
              <h2 className="text-base font-bold text-slate-100 flex items-center gap-2">
                <IconSettings />
                <span>上游模型服务配置 & 接入指南</span>
              </h2>
              <button
                onClick={() => setIsSettingsOpen(false)}
                className="text-slate-450 hover:text-slate-200 text-2xl font-light transition-colors active:scale-95 leading-none"
              >
                &times;
              </button>
            </div>

            <div className="flex flex-col md:flex-row overflow-hidden min-h-0 flex-1">
              {/* 左侧：服务配置表单 */}
              <form
                onSubmit={handleSaveConfig}
                className="flex-1 p-5 flex flex-col gap-5 overflow-y-auto max-h-[70vh] scroll-isolated border-r border-slate-900/60"
              >
                <div className="flex flex-col gap-1.5">
                  <label className="text-xs font-semibold text-slate-300">
                    本地代理监听地址
                  </label>
                  <input
                    type="text"
                    value={config.listenAddr}
                    onChange={(e) =>
                      setConfig({ ...config, listenAddr: e.target.value })
                    }
                    className="bg-slate-950/50 border border-slate-900 rounded-xl px-3.5 py-2.5 text-xs outline-none focus:border-cyan-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(6,182,212,0.15)] text-slate-200 transition-all font-mono"
                    required
                  />
                  <span className="text-[10px] text-slate-500 font-medium">
                    例如 :1238，配置后需重启代理生效
                  </span>
                  <div className="mt-1 px-3 py-2 bg-slate-950/60 border border-slate-900/60 rounded-xl flex flex-col gap-1">
                    <span className="text-[9px] text-slate-500 uppercase font-bold tracking-wider">
                      本地代理 API 根路径 (Base URL)
                    </span>
                    <code className="text-[10px] text-cyan-400 font-mono select-all break-all">
                      {proxyBase}/v1
                    </code>
                  </div>
                </div>

                <div className="border-t border-slate-900/60 pt-4 flex flex-col gap-3">
                  <span className="text-xs font-bold text-cyan-400 tracking-wide uppercase">
                    OpenAI Upstream 路由
                  </span>
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[10px] text-slate-550 font-bold">
                      上游 Base URL
                    </label>
                    <input
                      type="text"
                      value={config.openaiBaseURL}
                      onChange={(e) =>
                        setConfig({ ...config, openaiBaseURL: e.target.value })
                      }
                      className="bg-slate-950/50 border border-slate-900 rounded-xl px-3.5 py-2.5 text-xs outline-none focus:border-cyan-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(6,182,212,0.15)] text-slate-200 transition-all font-mono"
                      required
                    />
                  </div>
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[10px] text-slate-550 font-bold">
                      API Key
                    </label>
                    <div className="relative flex items-center">
                      <input
                        type={showOpenAIKey ? "text" : "password"}
                        value={config.openaiAPIKey}
                        onChange={(e) =>
                          setConfig({ ...config, openaiAPIKey: e.target.value })
                        }
                        className="w-full bg-slate-950/50 border border-slate-900 rounded-xl pl-3.5 pr-10 py-2.5 text-xs outline-none focus:border-cyan-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(6,182,212,0.15)] text-slate-200 transition-all font-mono"
                        placeholder="sk-..."
                      />
                      <button
                        type="button"
                        onClick={() => setShowOpenAIKey(!showOpenAIKey)}
                        className="absolute right-3 text-slate-400 hover:text-slate-200 transition-colors focus:outline-none flex items-center justify-center"
                      >
                        {showOpenAIKey ? <IconEyeOff /> : <IconEye />}
                      </button>
                    </div>
                  </div>
                  {/* 实时本地代理 URL 展示 */}
                  <div className="mt-1 px-3 py-2.5 bg-slate-950/40 border border-slate-900/60 rounded-xl flex flex-col gap-1 text-[10px] font-mono">
                    <div className="text-[9px] text-slate-500 uppercase font-bold tracking-wider font-sans mb-0.5">
                      OpenAI 本地代理客户端配置
                    </div>
                    <div className="flex justify-between gap-2">
                      <span className="text-slate-500">Base URL:</span>
                      <span className="text-cyan-400/90 select-all break-all">
                        {proxyBase}/v1
                      </span>
                    </div>
                    <div className="flex justify-between gap-2 mt-0.5">
                      <span className="text-slate-500">Chat URL:</span>
                      <span className="text-cyan-400/90 select-all break-all">
                        {proxyBase}/v1/chat/completions
                      </span>
                    </div>
                    <div className="flex justify-between gap-2 mt-0.5">
                      <span className="text-slate-500">Responses:</span>
                      <span className="text-cyan-400/90 select-all break-all">
                        {proxyBase}/v1/responses
                      </span>
                    </div>

                    <div className="border-t border-slate-900/40 my-1.5 pt-1.5" />

                    <div className="text-[9px] text-slate-500 uppercase font-bold tracking-wider font-sans mb-0.5">
                      实际发往上游 API 校验
                    </div>
                    <div className="flex justify-between gap-2">
                      <span className="text-slate-500">Chat URL:</span>
                      <span className="text-emerald-400/90 break-all">
                        {joinUpstreamURL(
                          config.openaiBaseURL,
                          "/v1/chat/completions",
                        ) || "(未配置上游地址)"}
                      </span>
                    </div>
                    <div className="flex justify-between gap-2 mt-0.5">
                      <span className="text-slate-500">Responses:</span>
                      <span className="text-emerald-400/90 break-all">
                        {joinUpstreamURL(
                          config.openaiResponsesBaseURL || config.openaiBaseURL,
                          "/v1/responses",
                        ) || "(未配置上游地址)"}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="border-t border-slate-900/60 pt-4 flex flex-col gap-3">
                  <span className="text-xs font-bold text-cyan-400 tracking-wide uppercase flex items-center gap-1.5">
                    <span>OpenAI Responses Upstream 路由 (可选)</span>
                    <span className="px-1.5 py-0.2 rounded bg-slate-900 text-[8px] text-slate-400 border border-slate-800">可选</span>
                  </span>
                  <p className="text-[10px] text-slate-400 leading-relaxed -mt-1 font-sans">
                    用于单独配置 OpenAI GPT-4o Realtime / Responses 接口的上游地址。若留空，则自动回退至上方的 OpenAI Upstream 配置。
                  </p>
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[10px] text-slate-550 font-bold">
                      上游 Base URL
                    </label>
                    <input
                      type="text"
                      value={config.openaiResponsesBaseURL}
                      onChange={(e) =>
                        setConfig({ ...config, openaiResponsesBaseURL: e.target.value })
                      }
                      className="bg-slate-950/50 border border-slate-900 rounded-xl px-3.5 py-2.5 text-xs outline-none focus:border-cyan-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(6,182,212,0.15)] text-slate-200 transition-all font-mono"
                      placeholder="https://api.openai.com"
                    />
                  </div>
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[10px] text-slate-550 font-bold">
                      API Key
                    </label>
                    <div className="relative flex items-center">
                      <input
                        type={showResponsesKey ? "text" : "password"}
                        value={config.openaiResponsesAPIKey}
                        onChange={(e) =>
                          setConfig({ ...config, openaiResponsesAPIKey: e.target.value })
                        }
                        className="w-full bg-slate-950/50 border border-slate-900 rounded-xl pl-3.5 pr-10 py-2.5 text-xs outline-none focus:border-cyan-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(6,182,212,0.15)] text-slate-200 transition-all font-mono"
                        placeholder="sk-..."
                      />
                      <button
                        type="button"
                        onClick={() => setShowResponsesKey(!showResponsesKey)}
                        className="absolute right-3 text-slate-400 hover:text-slate-200 transition-colors focus:outline-none flex items-center justify-center"
                      >
                        {showResponsesKey ? <IconEyeOff /> : <IconEye />}
                      </button>
                    </div>
                  </div>
                </div>

                <div className="border-t border-slate-900/60 pt-4 flex flex-col gap-3">
                  <span className="text-xs font-bold text-purple-400 tracking-wide uppercase">
                    Anthropic Upstream 路由
                  </span>
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[10px] text-slate-550 font-bold">
                      上游 Base URL
                    </label>
                    <input
                      type="text"
                      value={config.anthropicBaseURL}
                      onChange={(e) =>
                        setConfig({ ...config, anthropicBaseURL: e.target.value })
                      }
                      className="bg-slate-950/50 border border-slate-900 rounded-xl px-3.5 py-2.5 text-xs outline-none focus:border-purple-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(139,92,246,0.15)] text-slate-200 transition-all font-mono"
                      required
                    />
                  </div>
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[10px] text-slate-550 font-bold">
                      API Key
                    </label>
                    <div className="relative flex items-center">
                      <input
                        type={showAnthropicKey ? "text" : "password"}
                        value={config.anthropicAPIKey}
                        onChange={(e) =>
                          setConfig({ ...config, anthropicAPIKey: e.target.value })
                        }
                        className="w-full bg-slate-950/50 border border-slate-900 rounded-xl pl-3.5 pr-10 py-2.5 text-xs outline-none focus:border-purple-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(139,92,246,0.15)] text-slate-200 transition-all font-mono"
                        placeholder="sk-ant-..."
                      />
                      <button
                        type="button"
                        onClick={() => setShowAnthropicKey(!showAnthropicKey)}
                        className="absolute right-3 text-slate-400 hover:text-slate-200 transition-colors focus:outline-none flex items-center justify-center"
                      >
                        {showAnthropicKey ? <IconEyeOff /> : <IconEye />}
                      </button>
                    </div>
                  </div>
                  {/* 实时本地代理 URL 展示 */}
                  <div className="mt-1 px-3 py-2.5 bg-slate-950/40 border border-slate-900/60 rounded-xl flex flex-col gap-1 text-[10px] font-mono">
                    <div className="text-[9px] text-slate-500 uppercase font-bold tracking-wider font-sans mb-0.5">
                      Anthropic 本地代理客户端配置
                    </div>
                    <div className="flex justify-between gap-2">
                      <span className="text-slate-500">Base URL:</span>
                      <span className="text-purple-400/90 select-all break-all">
                        {proxyBase}/v1
                      </span>
                    </div>
                    <div className="flex justify-between gap-2 mt-0.5">
                      <span className="text-slate-500">Messages:</span>
                      <span className="text-purple-400/90 select-all break-all">
                        {proxyBase}/v1/messages
                      </span>
                    </div>

                    <div className="border-t border-slate-900/40 my-1.5 pt-1.5" />

                    <div className="text-[9px] text-slate-500 uppercase font-bold tracking-wider font-sans mb-0.5">
                      实际发往上游 API 校验
                    </div>
                    <div className="flex justify-between gap-2">
                      <span className="text-slate-500">Messages:</span>
                      <span className="text-emerald-400/90 break-all">
                        {joinUpstreamURL(
                          config.anthropicBaseURL,
                          "/v1/messages",
                        ) || "(未配置上游地址)"}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="border-t border-slate-900/60 pt-4 flex flex-col gap-3">
                  <span className="text-xs font-bold text-emerald-400 tracking-wide uppercase">
                    SQLite 数据路径
                  </span>
                  <div className="flex flex-col gap-1.5">
                    <input
                      type="text"
                      value={config.dbPath}
                      onChange={(e) =>
                        setConfig({ ...config, dbPath: e.target.value })
                      }
                      className="bg-slate-950/50 border border-slate-900 rounded-xl px-3.5 py-2.5 text-xs outline-none focus:border-emerald-500/50 focus:bg-slate-950/85 focus:shadow-[0_0_10px_rgba(16,185,129,0.15)] text-slate-200 transition-all font-mono"
                      required
                    />
                  </div>
                </div>

                <div className="flex justify-end gap-3 mt-4 pt-4 border-t border-slate-900 shrink-0">
                  <button
                    type="button"
                    onClick={() => setIsSettingsOpen(false)}
                    className="btn-secondary text-xs px-4 py-2.5 hover:bg-slate-800 transition-all duration-200 font-semibold"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    className="btn-primary btn-glow-cyan text-xs px-5 py-2.5 transition-all duration-200"
                    disabled={isConfigSaving}
                  >
                    {isConfigSaving ? "正在保存..." : "保存更改"}
                  </button>
                </div>
              </form>

              {/* 右侧：客户端接入指南 */}
              <div className="w-full md:w-[380px] p-5 flex flex-col gap-4 overflow-y-auto max-h-[70vh] scroll-isolated bg-slate-950/20 text-xs shrink-0">
                <div className="flex items-center gap-2 pb-2 border-b border-slate-900">
                  <span className="text-base">🚀</span>
                  <span className="font-bold text-slate-100">客户端接入指南</span>
                </div>
                <p className="text-[11px] text-slate-400 leading-relaxed">
                  通过将您的 LLM 开发库（Python/JS 等）的 <code className="text-cyan-400">baseURL</code> 指向本地代理，即可实现日志的无感劫持与录制：
                </p>

                {/* 环境变量复制 */}
                <div className="flex flex-col gap-2 bg-slate-900/40 p-3 rounded-xl border border-slate-900/60">
                  <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider">
                    方式一：系统环境变量
                  </span>
                  <div className="flex flex-col gap-1.5 font-mono text-[10px]">
                    <div className="flex justify-between items-center bg-slate-950/60 p-2 rounded-lg border border-slate-900">
                      <span className="text-cyan-400 truncate mr-2 font-semibold">export OPENAI_BASE_URL={proxyBase}/v1</span>
                      <button
                        type="button"
                        onClick={() => {
                          navigator.clipboard.writeText(`export OPENAI_BASE_URL=${proxyBase}/v1`);
                          setCopiedText("env_openai");
                          setTimeout(() => setCopiedText(""), 2000);
                        }}
                        className="text-slate-300 hover:text-cyan-400 text-[9px] transition-colors px-1.5 py-0.5 rounded bg-slate-900/80 shrink-0 font-sans"
                      >
                        {copiedText === "env_openai" ? "已复制 ✔️" : "复制"}
                      </button>
                    </div>
                    <div className="flex justify-between items-center bg-slate-950/60 p-2 rounded-lg border border-slate-900">
                      <span className="text-purple-400 truncate mr-2 font-semibold">export ANTHROPIC_BASE_URL={proxyBase}/v1</span>
                      <button
                        type="button"
                        onClick={() => {
                          navigator.clipboard.writeText(`export ANTHROPIC_BASE_URL=${proxyBase}/v1`);
                          setCopiedText("env_anthropic");
                          setTimeout(() => setCopiedText(""), 2000);
                        }}
                        className="text-slate-300 hover:text-purple-400 text-[9px] transition-colors px-1.5 py-0.5 rounded bg-slate-900/80 shrink-0 font-sans"
                      >
                        {copiedText === "env_anthropic" ? "已复制 ✔️" : "复制"}
                      </button>
                    </div>
                  </div>
                </div>

                {/* Python 客户端复制 */}
                <div className="flex flex-col gap-2 bg-slate-900/40 p-3 rounded-xl border border-slate-900/60">
                  <div className="flex justify-between items-center">
                    <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider font-sans">
                      方式二：Python (openai-python)
                    </span>
                    <button
                      type="button"
                      onClick={() => {
                        const code = `import openai\n\n# 将 base_url 指向 LLM Tracer 本地代理\nopenai.base_url = "${proxyBase}/v1"\n# apiKey 可填任意非空字符串\nopenai.api_key = "anything"\n\nresponse = openai.chat.completions.create(\n    model="gpt-4o",\n    messages=[{"role": "user", "content": "Hello!"}]\n)`;
                        navigator.clipboard.writeText(code);
                        setCopiedText("py_openai");
                        setTimeout(() => setCopiedText(""), 2000);
                      }}
                      className="text-slate-400 hover:text-cyan-400 transition-colors px-1.5 py-0.5 rounded bg-slate-900/80 shrink-0 text-[9px] font-sans"
                    >
                      {copiedText === "py_openai" ? "已复制 ✔️" : "复制代码"}
                    </button>
                  </div>
                  <pre className="bg-slate-950/80 border border-slate-900 rounded-lg p-2.5 text-[9px] font-mono text-emerald-400/90 overflow-x-auto leading-normal">
{`import openai

openai.base_url = "${proxyBase}/v1"
openai.api_key = "anything"`}
                  </pre>
                </div>

                {/* Node.js 客户端复制 */}
                <div className="flex flex-col gap-2 bg-slate-900/40 p-3 rounded-xl border border-slate-900/60">
                  <div className="flex justify-between items-center">
                    <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider font-sans">
                      方式三：JS / TS (openai-node)
                    </span>
                    <button
                      type="button"
                      onClick={() => {
                        const code = `import OpenAI from 'openai';\n\nconst openai = new OpenAI({\n  baseURL: '${proxyBase}/v1',\n  apiKey: 'anything'\n});\n\nconst completion = await openai.chat.completions.create({\n  model: 'gpt-4o',\n  messages: [{ role: 'user', content: 'Hello!' }],\n});`;
                        navigator.clipboard.writeText(code);
                        setCopiedText("node_openai");
                        setTimeout(() => setCopiedText(""), 2000);
                      }}
                      className="text-slate-400 hover:text-cyan-400 transition-colors px-1.5 py-0.5 rounded bg-slate-900/80 shrink-0 text-[9px] font-sans"
                    >
                      {copiedText === "node_openai" ? "已复制 ✔️" : "复制代码"}
                    </button>
                  </div>
                  <pre className="bg-slate-950/80 border border-slate-900 rounded-lg p-2.5 text-[9px] font-mono text-emerald-400/90 overflow-x-auto leading-normal">
{`import OpenAI from 'openai';

const openai = new OpenAI({
  baseURL: '${proxyBase}/v1',
  apiKey: 'anything'
});`}
                  </pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
      {modalLogId && renderDetailModal()}
    </div>
  );
}
