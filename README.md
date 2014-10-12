Proxy_Auth
==========

For Go Learning Challenge - Simple Web-API Server <http://community.topcoder.com/tc?module=ProjectDetail&pj=30046011>

proxy_auth is a simple authentication API.

Testing
-------

A simple server can be run with the following command : -

> go run main.go -p {server port} -a {server address} -u {json data}

For example :-

> go run main.go -p 8080 -a "localhost" -u users.json

A simple API tester is also included and can be run with the following command : -

> go run test.go -url http://{server address}:{server port}/api/2/domains/{domain}/proxyauth -testfile {json test file}

For example :-

> go run test.go -url http://localhost:8080/api/2/domains/{domain}/proxyauth -testfile test.json
