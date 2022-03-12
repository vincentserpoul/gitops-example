import http from "k6/http";
import { check } from "k6";
import { Rate } from "k6/metrics";

export let errorRate = new Rate("errors");

export default function () {
  const environment = `${__ENV.APP_ENVIRONMENT}`;


  var url = "http://localhost:3001/word";
  if(environment == "dev"){
    url = "https://orchestrator.orchestration.dev:8443/word";
  }

  var params = {};

  check(http.post(url, "bench", params), {
    "status is 201": (r) => r.status == 201,
  }) || errorRate.add(1);
}
