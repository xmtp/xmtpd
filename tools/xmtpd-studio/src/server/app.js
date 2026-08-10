/**
 * xmtpd Studio Web Server
 */

import express from 'express';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';
import { XMTPD_CONFIG } from '../config.js';
import { defaultXmtpdDaemon } from '../core/node-daemon.js';
import { defaultReplicationEngine } from '../core/replication.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const WEB_ROOT = path.join(__dirname, '../../web');

const app = express();
const PORT = process.env.PORT || 3429;

app.use(cors());
app.use(express.json());
app.use(express.static(WEB_ROOT));

// 1. Config & Topology
app.get('/api/config', (req, res) => {
  res.json({
    node: XMTPD_CONFIG.node,
    topology: XMTPD_CONFIG.peerTopology,
    services: XMTPD_CONFIG.supportedServices,
  });
});

// 2. Process Envelope
app.post('/api/node/envelope', (req, res) => {
  const envelope = defaultXmtpdDaemon.processEnvelope();
  res.json({ success: true, envelope });
});

// 3. Envelopes History
app.get('/api/node/envelopes', (req, res) => {
  res.json(defaultXmtpdDaemon.getRecentEnvelopes());
});

// 4. Replication Health
app.get('/api/health', (req, res) => {
  res.json(defaultReplicationEngine.verifyHealth());
});

if (process.env.NODE_ENV !== 'test') {
  app.listen(PORT, () => {
    console.log(`\n======================================================`);
    console.log(`⚡ xmtpd Decentralized Node Operator Studio Running!`);
    console.log(`🌐 Web Dashboard: http://localhost:${PORT}`);
    console.log(`📡 Protocol: XMTP P2P Message Replication Daemon`);
    console.log(`======================================================\n`);
  });
}

export default app;
