# 📡 xmtpd Studio & Decentralized Node Operator Suite

An interactive **XMTP Node Daemon Simulator**, **P2P Envelope Stream Replicator**, and **gRPC Service Monitor** for **xmtpd (`xmtp/xmtpd`)**.

---

## 🌟 Key Features

- 📡 **Decentralized XMTP Node Daemon**: Simulate node-to-node message indexing and MLS envelope stream replication.
- 🚀 **gRPC Services Architecture**: Monitor Replication API (Port 9001), Payer API (Port 9002), and Store API (Port 9003).
- 🌐 **Interactive Web Studio**: Real-time envelope stream dashboard and node operator console on `http://localhost:3429`.
- ⌨️ **Universal CLI (`xmtpd-cli`)**: Terminal utility for processing envelopes and checking network health.

---

## 🚀 Quickstart

```bash
# Launch xmtpd Studio
npm start
# Open http://localhost:3429

# Or run via CLI
node bin/xmtpd-cli.js process
node bin/xmtpd-cli.js health
```
