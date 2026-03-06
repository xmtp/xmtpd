// HTTPS reverse proxy for the gateway sidecar
// Needed because WASM bindings have a TLS mismatch bug when gateway is HTTP
import { createServer } from "node:https";
import { readFileSync } from "node:fs";
import { request as httpRequest } from "node:http";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const HTTPS_PORT = 5056;
const GATEWAY_PORT = 5055;

const options = {
  key: readFileSync(join(__dirname, "certs/localhost-key.pem")),
  cert: readFileSync(join(__dirname, "certs/localhost.pem")),
};

// Stats tracking
const stats = { total: 0, byEndpoint: {} };

function logRequest(method, path, status, reqBytes, resBytes, durationMs) {
  const endpoint = path.split("/").slice(-1)[0] || path;
  stats.total++;
  stats.byEndpoint[endpoint] = (stats.byEndpoint[endpoint] || 0) + 1;

  const ts = new Date().toISOString().substring(11, 23);
  const arrow = status < 400 ? "->" : "!>";
  console.log(
    `[${ts}] ${method} ${path} ${arrow} ${status} | req:${reqBytes}B res:${resBytes}B ${durationMs}ms | #${stats.total} (${endpoint}: ${stats.byEndpoint[endpoint]})`,
  );
}

function printStats() {
  console.log("\n=== Gateway Proxy Stats ===");
  console.log(`Total requests: ${stats.total}`);
  for (const [ep, count] of Object.entries(stats.byEndpoint).sort((a, b) => b[1] - a[1])) {
    console.log(`  ${ep}: ${count}`);
  }
  console.log("===========================\n");
}

// Print stats every 30 seconds if there's activity
let lastTotal = 0;
setInterval(() => {
  if (stats.total > lastTotal) {
    printStats();
    lastTotal = stats.total;
  }
}, 30000);

// Print stats on SIGUSR1
process.on("SIGUSR1", printStats);

const server = createServer(options, (req, res) => {
  const start = Date.now();

  // Forward CORS headers for browser
  if (req.method === "OPTIONS") {
    logRequest(req.method, req.url, 204, 0, 0, Date.now() - start);
    res.writeHead(204, {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET,HEAD,POST,PUT,PATCH,DELETE",
      "Access-Control-Allow-Headers":
        "Content-Type,Accept,Authorization,X-Client-Version,X-App-Version,Baggage,DNT,Sec-CH-UA,Sec-CH-UA-Mobile,Sec-CH-UA-Platform,x-grpc-web,grpc-timeout,Sentry-Trace,User-Agent,x-libxmtp-version,x-app-version",
      "Access-Control-Expose-Headers":
        "grpc-status,grpc-message,grpc-status-details-bin",
      "Access-Control-Max-Age": "86400",
    });
    res.end();
    return;
  }

  let reqBytes = 0;
  let resBytes = 0;

  req.on("data", (chunk) => {
    reqBytes += chunk.length;
  });

  const proxyReq = httpRequest(
    {
      hostname: "localhost",
      port: GATEWAY_PORT,
      path: req.url,
      method: req.method,
      headers: req.headers,
    },
    (proxyRes) => {
      // Ensure CORS on all responses
      const headers = { ...proxyRes.headers };
      headers["access-control-allow-origin"] = "*";
      headers["access-control-expose-headers"] =
        "grpc-status,grpc-message,grpc-status-details-bin";
      res.writeHead(proxyRes.statusCode, headers);

      proxyRes.on("data", (chunk) => {
        resBytes += chunk.length;
      });

      proxyRes.on("end", () => {
        logRequest(req.method, req.url, proxyRes.statusCode, reqBytes, resBytes, Date.now() - start);
      });

      proxyRes.pipe(res);
    },
  );

  proxyReq.on("error", (err) => {
    logRequest(req.method, req.url, 502, reqBytes, 0, Date.now() - start);
    console.error(`  Proxy error: ${err.message}`);
    res.writeHead(502);
    res.end("Bad Gateway");
  });

  req.pipe(proxyReq);
});

server.listen(HTTPS_PORT, () => {
  console.log(`HTTPS proxy: https://localhost:${HTTPS_PORT} -> http://localhost:${GATEWAY_PORT}`);
  console.log(`Logging all requests. Send SIGUSR1 (kill -USR1 ${process.pid}) for stats summary.`);
});
