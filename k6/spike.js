import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 1 },
    { duration: '10s', target: 50 },
    { duration: '10s', target: 1 },
    { duration: '30s', target: 1 },
  ],
};

export default function () {
  let data = JSON.stringify({
    url: 'https://google.com',
  });
  let res1 = http.post('http://localhost:8082/urls/', data, {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  let alias = JSON.parse(res1.body).alias;
  let res2 = http.get(`http://localhost:8082/urls/alias/${alias}`);

  check(res1, { "status is 201": (res1) => res1.status === 201 });
  check(res2, { "status is 200": (res2) => res2.status === 200 });
  sleep(1);
}
