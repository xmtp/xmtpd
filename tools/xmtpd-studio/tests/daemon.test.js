/**
 * xmtpd Daemon Unit Tests
 */

import { defaultXmtpdDaemon } from '../src/core/node-daemon.js';
import { defaultReplicationEngine } from '../src/core/replication.js';

async function runDaemonTests() {
  console.log('Testing xmtpd Decentralized Node Daemon & Stream Replication...');

  // 1. Process Envelope
  const env = defaultXmtpdDaemon.processEnvelope();
  if (!env.hash || !env.sequenceNumber) {
    throw new Error('xmtpd envelope processing failed');
  }

  // 2. Replication Health
  const h = defaultReplicationEngine.verifyHealth();
  if (!h.connectedNodes || h.gRPCStreamStatus !== 'HEALTHY_STREAMING') {
    throw new Error('xmtpd replication health check failed');
  }

  console.log(`✅ xmtpd Node Envelope Processed (Seq #${env.sequenceNumber}) & P2P Health Verified!`);
}

runDaemonTests().catch(e => {
  console.error('❌ Daemon Test Failed:', e);
  process.exit(1);
});
