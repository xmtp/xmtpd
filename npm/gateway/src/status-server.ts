import { createServer, type IncomingMessage, type Server, type ServerResponse } from "node:http";
import { EventEmitter } from "node:events";
import type { GatewayStats } from "./types.js";

export interface StatusServerOptions {
  port: number;
  gatewayPort: number;
  getStats: () => GatewayStats;
  logEmitter: EventEmitter;
  logBuffer: string[];
}

export interface StatusServer {
  server: Server;
  close: () => Promise<void>;
}

export function startStatusServer(options: StatusServerOptions): StatusServer {
  const { port, gatewayPort, getStats, logEmitter, logBuffer } = options;

  const server = createServer((req: IncomingMessage, res: ServerResponse) => {
    const url = new URL(req.url ?? "/", `http://localhost:${port}`);

    if (url.pathname === "/status") {
      res.writeHead(200, { "Content-Type": "application/json", "Access-Control-Allow-Origin": "*" });
      res.end(JSON.stringify(getStats()));
      return;
    }

    if (url.pathname === "/logs") {
      res.writeHead(200, {
        "Content-Type": "text/event-stream",
        "Cache-Control": "no-cache",
        "Connection": "keep-alive",
        "Access-Control-Allow-Origin": "*",
      });

      // Send buffered logs
      for (const line of logBuffer) {
        res.write(`data: ${JSON.stringify(line)}\n\n`);
      }

      const onLog = (line: string) => {
        res.write(`data: ${JSON.stringify(line)}\n\n`);
      };
      logEmitter.on("log", onLog);

      req.on("close", () => {
        logEmitter.off("log", onLog);
      });
      return;
    }

    // Dashboard HTML
    res.writeHead(200, { "Content-Type": "text/html" });
    res.end(dashboardHtml(port, gatewayPort));
  });

  server.listen(port);

  return {
    server,
    close: () => new Promise<void>((resolve) => {
      server.close(() => resolve());
    }),
  };
}

function dashboardHtml(statusPort: number, gatewayPort: number): string {
  return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>XMTP Gateway Status</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, monospace; background: #0d1117; color: #c9d1d9; padding: 24px; }
  .header { display: flex; align-items: center; gap: 12px; margin-bottom: 24px; }
  .header h1 { font-size: 20px; font-weight: 600; }
  .badge { display: inline-block; padding: 2px 10px; border-radius: 12px; font-size: 12px; font-weight: 600; }
  .badge.online { background: #238636; color: #fff; }
  .badge.offline { background: #da3633; color: #fff; }
  .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(140px, 1fr)); gap: 12px; margin-bottom: 24px; }
  .stat { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 16px; }
  .stat .label { font-size: 11px; color: #8b949e; text-transform: uppercase; letter-spacing: 0.5px; }
  .stat .value { font-size: 24px; font-weight: 700; margin-top: 4px; }
  .stat .value.green { color: #3fb950; }
  .stat .value.red { color: #f85149; }
  .stat .value.blue { color: #58a6ff; }
  .logs-section { background: #161b22; border: 1px solid #30363d; border-radius: 8px; }
  .logs-header { padding: 12px 16px; border-bottom: 1px solid #30363d; font-size: 14px; font-weight: 600; display: flex; justify-content: space-between; }
  .logs-header span { color: #8b949e; font-weight: 400; font-size: 12px; }
  #logs { height: 400px; overflow-y: auto; padding: 12px 16px; font-family: 'SF Mono', 'Fira Code', monospace; font-size: 12px; line-height: 1.6; }
  .log-line { white-space: pre-wrap; word-break: break-all; }
  .log-line.error { color: #f85149; }
  .log-line.warn { color: #d29922; }
  .log-line.info { color: #c9d1d9; }
  .log-line.debug { color: #8b949e; }
</style>
</head>
<body>

<div class="header">
  <h1>XMTP Gateway</h1>
  <span id="status-badge" class="badge offline">CONNECTING</span>
</div>

<div class="stats">
  <div class="stat">
    <div class="label">Uptime</div>
    <div class="value blue" id="uptime">--</div>
  </div>
  <div class="stat">
    <div class="label">Gateway Port</div>
    <div class="value" id="port">${gatewayPort}</div>
  </div>
  <div class="stat">
    <div class="label">Publishes</div>
    <div class="value green" id="publishes">0</div>
  </div>
  <div class="stat">
    <div class="label">Errors</div>
    <div class="value red" id="errors">0</div>
  </div>
  <div class="stat">
    <div class="label">Requests</div>
    <div class="value blue" id="requests">0</div>
  </div>
</div>

<div class="logs-section">
  <div class="logs-header">
    Live Logs
    <span id="log-count">0 lines</span>
  </div>
  <div id="logs"></div>
</div>

<script>
  const logsEl = document.getElementById('logs');
  const badge = document.getElementById('status-badge');
  let logCount = 0;
  let autoScroll = true;

  logsEl.addEventListener('scroll', () => {
    autoScroll = logsEl.scrollTop + logsEl.clientHeight >= logsEl.scrollHeight - 30;
  });

  function formatUptime(s) {
    const h = Math.floor(s / 3600);
    const m = Math.floor((s % 3600) / 60);
    const sec = Math.floor(s % 60);
    if (h > 0) return h + 'h ' + m + 'm ' + sec + 's';
    if (m > 0) return m + 'm ' + sec + 's';
    return sec + 's';
  }

  function addLog(text) {
    const div = document.createElement('div');
    div.className = 'log-line';
    if (text.includes('[gateway:err]') || text.toLowerCase().includes('error')) div.className += ' error';
    else if (text.toLowerCase().includes('warn')) div.className += ' warn';
    else if (text.toLowerCase().includes('debug')) div.className += ' debug';
    else div.className += ' info';
    div.textContent = text;
    logsEl.appendChild(div);
    logCount++;
    document.getElementById('log-count').textContent = logCount + ' lines';
    // Keep max 500 lines in DOM
    while (logsEl.children.length > 500) logsEl.removeChild(logsEl.firstChild);
    if (autoScroll) logsEl.scrollTop = logsEl.scrollHeight;
  }

  // SSE log stream
  const es = new EventSource('/logs');
  es.onmessage = (e) => addLog(JSON.parse(e.data));

  // Poll stats
  async function refreshStats() {
    try {
      const r = await fetch('/status');
      const s = await r.json();
      badge.textContent = s.online ? 'ONLINE' : 'OFFLINE';
      badge.className = 'badge ' + (s.online ? 'online' : 'offline');
      document.getElementById('uptime').textContent = formatUptime(s.uptimeSeconds);
      document.getElementById('publishes').textContent = s.publishes;
      document.getElementById('errors').textContent = s.errors;
      document.getElementById('requests').textContent = s.requests;
    } catch {
      badge.textContent = 'OFFLINE';
      badge.className = 'badge offline';
    }
  }
  refreshStats();
  setInterval(refreshStats, 2000);
</script>
</body>
</html>`;
}
