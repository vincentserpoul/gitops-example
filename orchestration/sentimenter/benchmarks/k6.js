import http from "k6/http";
import { check } from "k6";
import { Rate } from "k6/metrics";

export let errorRate = new Rate("errors");

export default function () {
  const environment = `${__ENV.APP_ENVIRONMENT}`;


  var url = "http://localhost:3001/";
  if(environment == "dev"){
    url = "https://sentimenter.orchestration.dev:8443/";
  }

  var params = {};

  check(http.post(url, {text: "hello la polka"}, params), {
    "status is 201": (r) => r.status == 201,
  }) || errorRate.add(1);
}
