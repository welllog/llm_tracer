const fs = require('fs');
const path = require('path');

const filePath = '/Users/chentairen/Downloads/go-project/ludao/llm_recorder/web_recorder/frontend/src/App.jsx';
let content = fs.readFileSync(filePath, 'utf8');

// 我们定位到 "滚动卡片列表" 这一行作为开始
const startPattern = '          {/* 滚动卡片列表 */}';
// 我们定位到 "分页组件" 这一行作为结束
const endPattern = '          {/* 分页组件 */}';

const startIndex = content.indexOf(startPattern);
const endIndex = content.indexOf(endPattern);

if (startIndex !== -1 && endIndex !== -1) {
  const replacement = `          {/* 滚动卡片列表 */}
          <div className="flex-1 overflow-y-auto flex flex-col gap-3 pr-1 min-h-0">
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
              <div className="flex flex-col gap-3">
                {groupedLogs.map((pathObj) => {
                  const isPathCollapsed = collapsedPaths[pathObj.path] === true; // 默认展开
                  const sessionEntries = Object.entries(pathObj.sessions);
                  const dateEntries = Object.entries(pathObj.dates);
                  const totalCount = sessionEntries.reduce((acc, curr) => acc + curr[1].length, 0) + 
                                     dateEntries.reduce((acc, curr) => acc + curr[1].length, 0);

                  return (
                    <div key={pathObj.path} className="flex flex-col gap-2 border-b border-slate-900 pb-3 last:border-b-0">
                      {/* 接口第一层节点 */}
                      <div 
                        onClick={() => togglePath(pathObj.path)}
                        className="flex items-center justify-between p-2 bg-slate-900/60 border border-slate-800/80 rounded-lg cursor-pointer hover:bg-slate-850/80 transition-all"
                      >
                        <div className="flex items-center gap-1.5 min-w-0">
                          <span className={\`transform transition-transform text-slate-400 \${!isPathCollapsed ? 'rotate-90' : ''}\`}>
                            <IconChevronRight />
                          </span>
                          <span className="text-[10px] font-bold px-1.5 py-0.5 rounded bg-slate-800 border border-slate-700 text-slate-300 font-mono uppercase tracking-wide">POST</span>
                          <span className="text-xs font-semibold text-slate-200 truncate font-mono" title={pathObj.path}>
                            {pathObj.path}
                          </span>
                        </div>
                        <span className="text-[10px] text-slate-500 font-mono shrink-0">({totalCount}条)</span>
                      </div>

                      {/* 接口下的会话/日期列表 */}
                      {!isPathCollapsed && (
                        <div className="flex flex-col gap-2 pl-2 border-l border-slate-900/60 ml-2">
                          {/* 会话列表 */}
                          {sessionEntries.map(([sessionId, sessionLogs]) => {
                            const groupKey = \`\${pathObj.path}-\${sessionId}\`;
                            const isGroupCollapsed = collapsedGroups[groupKey] !== false; // 默认折叠
                            const count = sessionLogs.length;
                            // 提取时间跨度
                            const times = sessionLogs.map(l => l.createdAt.split(' ')[1] || '').filter(t => t);
                            times.sort();
                            const timeRange = times.length > 0 ? \`\${times[0]} - \${times[times.length - 1]}\` : '';

                            return (
                              <div key={sessionId} className="flex flex-col gap-1.5">
                                <div 
                                  onClick={() => toggleGroup(groupKey)}
                                  className="flex items-center justify-between p-2.5 bg-indigo-950/15 border border-indigo-900/30 rounded-xl cursor-pointer hover:bg-indigo-950/25 transition-all my-1"
                                >
                                  <div className="flex items-center gap-1.5 min-w-0">
                                    <span className={\`transform transition-transform text-indigo-400 \${!isGroupCollapsed ? 'rotate-90' : ''}\`}>
                                      <IconChevronRight />
                                    </span>
                                    <span className="text-xs">💬</span>
                                    <span className="text-[11px] font-mono text-slate-300 font-semibold truncate" title={sessionId}>
                                      会话: {sessionId.substring(0, 10)}...
                                    </span>
                                  </div>
                                  <div className="flex items-center gap-2 shrink-0">
                                    {timeRange && (
                                      <span className="text-[8px] px-1.5 py-0.5 bg-slate-900/80 text-slate-500 rounded font-mono">
                                        {timeRange}
                                      </span>
                                    )}
                                    <span className="text-[9px] text-indigo-400/90 font-mono">({count}次)</span>
                                  </div>
                                </div>

                                {!isGroupCollapsed && (
                                  <div className="flex flex-col gap-1.5 pl-3 border-l border-indigo-950/30 ml-2">
                                    {sessionLogs.map(log => renderLogCard(log))}
                                  </div>
                                )}
                              </div>
                            );
                          })}

                          {/* 日期列表 */}
                          {dateEntries.map(([date, dateLogs]) => {
                            const groupKey = \`\${pathObj.path}-\${date}\`;
                            const isGroupCollapsed = collapsedGroups[groupKey] !== false; // 默认折叠
                            const count = dateLogs.length;

                            return (
                              <div key={date} className="flex flex-col gap-1.5">
                                <div 
                                  onClick={() => toggleGroup(groupKey)}
                                  className="flex items-center justify-between p-2.5 bg-slate-900/50 border border-slate-800 rounded-xl cursor-pointer hover:bg-slate-850/60 transition-all my-1"
                                >
                                  <div className="flex items-center gap-1.5">
                                    <span className={\`transform transition-transform text-slate-400 \${!isGroupCollapsed ? 'rotate-90' : ''}\`}>
                                      <IconChevronRight />
                                    </span>
                                    <span className="text-xs">📅</span>
                                    <span className="text-[11px] font-semibold text-slate-300">{date}</span>
                                  </div>
                                  <span className="text-[9px] text-slate-500 font-mono">({count}次)</span>
                                </div>

                                {!isGroupCollapsed && (
                                  <div className="flex flex-col gap-1.5 pl-3 border-l border-slate-800/60 ml-2">
                                    {dateLogs.map(log => renderLogCard(log))}
                                  </div>
                                )}
                              </div>
                            );
                          })}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
          
`;

  content = content.substring(0, startIndex) + replacement + content.substring(endIndex);
  fs.writeFileSync(filePath, content, 'utf8');
  console.log('Successfully refactored App.jsx tree view!');
} else {
  console.log('Patterns not found!', { startIndex, endIndex });
}
