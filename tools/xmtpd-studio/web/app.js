/**
 * xmtpd Studio Client Logic
 */

let isAutoStreaming = false;
let autoInterval = null;

document.addEventListener('DOMContentLoaded', () => {
  initTabs();
  loadConfig();
  initListeners();
});

function initTabs() {
  const tabs = document.querySelectorAll('.nav-tab');
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      document.querySelectorAll('.nav-tab').forEach(t => t.classList.toggle('active', t === tab));
      document.querySelectorAll('.tab-pane').forEach(p => p.classList.toggle('active', p.id === `tab-${tab.dataset.tab}`));
    });
  });
}

async function loadConfig() {
  try {
    const res = await fetch('/api/config');
    const data = await res.json();

    document.getElementById('header-seq').textContent = `Seq: #${data.topology.indexedEnvelopes.toLocaleString()}`;

    const grid = document.getElementById('services-container');
    grid.innerHTML = '';

    data.services.forEach(s => {
      const card = document.createElement('div');
      card.className = 'service-card';
      card.innerHTML = `
        <div class="service-title">${s.name}</div>
        <div class="service-port">Port: ${s.port}</div>
        <div class="text-muted" style="font-size: 0.85rem;">${s.description}</div>
      `;
      grid.appendChild(card);
    });
  } catch (e) {
    console.error(e);
  }
}

function initListeners() {
  document.getElementById('btn-process-envelope').addEventListener('click', processEnvelope);

  const autoBtn = document.getElementById('btn-toggle-auto');
  autoBtn.addEventListener('click', () => {
    if (isAutoStreaming) {
      clearInterval(autoInterval);
      isAutoStreaming = false;
      autoBtn.textContent = '▶️ Start Live Envelope Stream (5 envs/sec)';
      autoBtn.className = 'btn btn-gradient btn-lg';
    } else {
      isAutoStreaming = true;
      autoBtn.textContent = '⏸️ Pause Envelope Stream';
      autoBtn.className = 'btn btn-secondary btn-lg';
      processEnvelope();
      autoInterval = setInterval(processEnvelope, 200); // 200ms fast stream
    }
  });
}

async function processEnvelope() {
  try {
    const res = await fetch('/api/node/envelope', { method: 'POST' });
    const data = await res.json();
    if (data.success) {
      appendEnvelopeRow(data.envelope);
    }
  } catch (e) {
    console.warn(e);
  }
}

function appendEnvelopeRow(env) {
  const container = document.getElementById('envelopes-container');
  const empty = container.querySelector('.empty-state');
  if (empty) container.innerHTML = '';

  const row = document.createElement('div');
  row.className = 'ledger-row';
  row.innerHTML = `
    <div>
      <div style="font-weight: 700; color: #fff;">Sequence #${env.sequenceNumber.toLocaleString()}</div>
      <div class="mono text-muted" style="font-size: 0.72rem;">Hash: ${env.hash.slice(0, 16)}...</div>
    </div>
    <div style="text-align: right;">
      <div style="color: #ef4444; font-weight: 700; font-family: var(--font-mono);">${env.sizeBytes} Payload</div>
      <div class="mono text-muted" style="font-size: 0.72rem;">${env.replicationStatus}</div>
    </div>
  `;
  container.insertBefore(row, container.firstChild);
}
