// loadtest.js
import http from "k6/http";
import {check, sleep} from "k6";
export let options = {stages: [{duration: "30s", target: 100}]}; // ramp to 50 VUs
export default function () {
  let res = http.get("http://localhost:8080/ping");

  check(res, {"status==200": (r) => r.status === 200});
  sleep(1);
}
