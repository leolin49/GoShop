docker run --name consul1 -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600 hashicorp/consul agent -server -bootstrap-expect=1 -ui -bind=0.0.0.0 -client=0.0.0.0
