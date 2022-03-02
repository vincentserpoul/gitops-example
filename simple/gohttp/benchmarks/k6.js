import http from "k6/http";
import { check } from "k6";
import { Rate } from "k6/metrics";

export let errorRate = new Rate("errors");

export default function () {
  const environment = `${__ENV.APP_ENVIRONMENT}`;


  var url = "http://localhost:3000/user/123";
  if(environment == "dev"){
    url = "https://dev.gohttp.com:8443/user/123";
  }

  var params = {};

  check(http.get(url, params), {
    "status is 200": (r) => r.status == 200,
  }) || errorRate.add(1);
}
