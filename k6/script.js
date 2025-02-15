import http from 'k6/http';
import { sleep, check } from 'k6';
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const URL = 'http://localhost:8080/api';

export const options = {
  scenarios: {
    main_scenario: {
      executor: 'constant-arrival-rate',
      rate: 1000 / 4, // cause every iteration makes 4 requests. total rps is 1000
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 20,
      maxVUs: 100,
    }
  },
  thresholds: {
    http_req_failed: ['rate<0.001'],
    http_req_duration: ['avg<50'],
  },
};

let users = [];

const auth = () => {
  const user = {
    username: randomString(15),
    password: randomString(12),
  }

  const res = http.post(
      `${URL}/auth`,
      JSON.stringify(user),
      {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
  );

  const ok = check(res, {
    "status code should be 200": res => res.status === 200
  })

  if (!ok) {
    console.log(res)
  }

  user.token = res.json().token
  users.push(user);
}

const buy = () => {
  const i = randomIntBetween(0, users.length - 1)
  const user = users[i];

  const res = http.get(
      `${URL}/buy/cup`,
      {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'Authorization': 'Bearer ' + user.token
        }
      }
  );

  const ok = check(res, {
    "status code should be 200": res => res.status === 200
  })

  if (!ok) {
    console.log(res)
  }
}

const sendCoin = () => {
  const user = users[randomIntBetween(0, users.length - 1)];
  let sendToUser = users[randomIntBetween(0, users.length - 1)]
  while (sendToUser.username === user.username) {
    sendToUser = users[randomIntBetween(0, users.length - 1)]
  }

  const res = http.post(
      `${URL}/sendCoin`,
      JSON.stringify({
        toUser: sendToUser,
        amount: 1,
      }),
      {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'Authorization': 'Bearer ' + user.token
        }
      }
  );

  const ok = check(res, {
    "status code should be 200": res => res.status === 200
  })

  if (!ok) {
    console.log(res)
  }
}

const info = () => {
  const user = users[randomIntBetween(0, users.length - 1)];

  const res = http.get(
      `${URL}/info`,
      {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'Authorization': 'Bearer ' + user.token
        }
      }
  );

  const ok = check(res, {
    "status code should be 200": res => res.status === 200
  })

  if (!ok) {
    console.log(res)
  }
}

export default function() {
  auth()
  sleep(randomIntBetween(0, 5))
  buy()
  sleep(randomIntBetween(0, 5))
  sendCoin()
  sleep(randomIntBetween(0, 5))
  info()
  sleep(randomIntBetween(0, 5))
}
