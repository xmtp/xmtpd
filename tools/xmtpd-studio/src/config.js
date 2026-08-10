/**
 * xmtpd Decentralized Node & Replication Configuration
 */

export const XMTPD_CONFIG = {
  node: {
    name: 'xmtpd Decentralized Node Daemon',
    language: 'Go / gRPC',
    network: 'XMTP Decentralized Mainnet / Testnet',
    syncProtocol: 'P2P Stream Replication & Storage Indexer',
  },
  peerTopology: {
    activeNodes: 14,
    syncLatencyMs: 45,
    indexedEnvelopes: 1258400,
    storageEngine: 'PostgreSQL + CockroachDB',
  },
  supportedServices: [
    { name: 'Replication API', port: 9001, description: 'gRPC stream for node-to-node envelope synchronization.' },
    { name: 'Payer API', port: 9002, description: 'Message origin verification & fee settlement endpoint.' },
    { name: 'Store API', port: 9003, description: 'Client message retrieval and querying service.' },
  ],
};
