import { check, sleep } from 'k6';
import http from 'k6/http';
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { randomBytes } from 'k6/crypto';
import { b64encode } from 'k6/encoding';

export const options = {
  vus: 1,
  iterations: 10,
  // duration: '1m',

  thresholds: {
    http_req_duration: ['p(99)<1500'], // 99% of requests must complete below 1.5s
  },
};

export default function () {
  const env = {
    envelopes: [
      {
        contentTopic: "test-"+randomString(randomIntBetween(13, 20)),
        timestampNs: randomIntBetween(0, 100),
        message: b64encode(randomBytes(randomIntBetween(0, 100)))
      },
    ]
  };
  console.log(env);
  const res = http.post('http://localhost:5001/message/v1/publish', JSON.stringify(env));
  console.log(res);
  check(res, {
    'is status 200': (r) => r.status === 200,
    // 'protocol is HTTP/2': (r) => r.proto === 'HTTP/2.0',
  });
  sleep(0.2);
};
