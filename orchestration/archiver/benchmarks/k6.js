import http from "k6/http";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";
import { check } from "k6";
import { Rate } from "k6/metrics";

export let errorRate = new Rate("errors");

export default function () {
  const environment = `${__ENV.APP_ENVIRONMENT}`;

  var url = "http://localhost:3001/v1/happycatfact";
  if(environment == "dev"){
    url = "https://archiver.orchestration.dev:8443/v1/happycatfact";
  }

  var params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  var payload = JSON.stringify(
    {
      id: uuidv4(),
      fact: "bench"
    }
  )

  check(
    http.post(url, payload, params),
    {
      "status is 201": (r) => r.status == 201,
    }
  ) || errorRate.add(1);
}
