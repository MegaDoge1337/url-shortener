import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 1,
  duration: '30s',
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

  let id = JSON.parse(res1.body).id;
  let res2 = http.get(`http://localhost:8082/urls/id/${id}`);

  let alias = JSON.parse(res1.body).alias;
  let res3 = http.get(`http://localhost:8082/urls/alias/${alias}`);

  check(res1, { "status is 201": (res1) => res1.status === 201 });
  check(res2, { "status is 200": (res2) => res2.status === 200 });
  check(res3, { "status is 200": (res3) => res3.status === 200 });
  sleep(1);
}
