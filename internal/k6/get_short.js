import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

export let options = {
    scenarios: {
        contacts: {
            executor: 'constant-vus',
            vus: 100,
            duration: '10s',
        },
    },
};

export default function(needUrl) {
    var response = http.get(needUrl);

    check(response, {
        'Response. is status 200': (r) => r.status === 200,
    });

    sleep(1);
}

export function setup() {

    var res = http.post('http://localhost:8080/', 'http://ya.ru?q=' + randomString(8));

    return res.body
}
