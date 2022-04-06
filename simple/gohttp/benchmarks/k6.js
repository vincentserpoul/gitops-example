import http from "k6/http";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";
import { check } from "k6";
import { Rate } from "k6/metrics";

export let errorRate = new Rate("errors");

export default function () {
  const environment = `${__ENV.ENV}`;

  var url = "http://localhost:3003/v1/user";
  if(environment == "dev"){
    url = "https://gohttp.simple.dev:8443/v1/user";
  } else if(environment == "prod"){
    url = "https://simple-gohttp.do-gitops.tk/v1/user";
  }


  var params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  var payload = JSON.stringify(
    {
      id: uuidv4(),
      created_at: "1937-01-01T12:00:27.87+08:00"
    }
  )

  check(
    http.post(url, payload, params),
    {
      "status is 201": (r) => r.status == 201,
    }
  ) || errorRate.add(1);
}
