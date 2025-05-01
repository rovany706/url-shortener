#!/bin/bash

jq -ncM 'while(true; .+1) | {method: "POST", url: "http://localhost:8080/api/shorten", body: {"url": . | "http://example\(.).com"} | @base64 }' | ./vegeta attack -rate 100/s -duration 10s -lazy -format=json
