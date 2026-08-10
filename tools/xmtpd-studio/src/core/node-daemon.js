/**
 * xmtpd Node Daemon & Envelope Processor
 */

import crypto from 'crypto';

export class XmtpdNodeDaemon {
  constructor() {
    this.sequenceNumber = 1258400;
    this.envelopes = [];
  }

  /**
   * Process & Index a new XMTP Message Envelope
   */
  processEnvelope() {
    this.sequenceNumber += 1;
    const envelopeHash = '0x' + crypto.randomBytes(32).toString('hex');
    const senderPeer = 'node_' + crypto.randomBytes(4).toString('hex');
    const envelopeSize = Math.floor(Math.random() * 800) + 200;

    const envelope = {
      sequenceNumber: this.sequenceNumber,
      hash: envelopeHash,
      originPeer: senderPeer,
      sizeBytes: `${envelopeSize} B`,
      replicationStatus: 'REPLICATED_TO_14_NODES',
      timestamp: new Date().toISOString(),
    };

    this.envelopes.unshift(envelope);
    return envelope;
  }

  getRecentEnvelopes() {
    return this.envelopes.slice(0, 10);
  }
}

export const defaultXmtpdDaemon = new XmtpdNodeDaemon();
