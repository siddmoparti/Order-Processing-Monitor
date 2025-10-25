// Script to configure Grafana with Prometheus data source
const https = require('https');
const http = require('http');

const grafanaUrl = 'http://localhost:3000';
const prometheusUrl = 'http://localhost:9090';

// Create Prometheus data source
const dataSourceConfig = {
  name: 'Prometheus',
  type: 'prometheus',
  url: prometheusUrl,
  access: 'proxy',
  isDefault: true
};

const options = {
  hostname: 'localhost',
  port: 3000,
  path: '/api/datasources',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Basic YWRtaW46YWRtaW4=' // admin:admin in base64
  }
};

const req = http.request(options, (res) => {
  console.log(`Status: ${res.statusCode}`);
  res.on('data', (d) => {
    console.log('Response:', d.toString());
  });
});

req.on('error', (e) => {
  console.error(`Problem with request: ${e.message}`);
});

req.write(JSON.stringify(dataSourceConfig));
req.end();

