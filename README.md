# Infinimesh Platform
Infinimesh Platform is an opinionated Platform to connect IoT devices securely. It exposes simple to consume RESTful & gRPC APIs with both high-level (e.g. device shadow) and low-level (sending messages) concepts. Infinimesh Platform is open source and fully cloud native. No vendor lock-in - run it yourself on Kubernetes or use our SaaS offering (TBA).
## Build status
[![CircleCI](https://img.shields.io/circleci/project/github/infinimesh/infinimesh.svg)](https://circleci.com/gh/infinimesh/infinimesh/tree/master) 
[![GoReportCard](https://goreportcard.com/badge/github.com/infinimesh/infinimesh)](https://goreportcard.com/report/github.com/infinimesh/infinimesh) 

| Docker Image  | Build status  |
| ------------- |---------------|
| API Server | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/apiserver-rest/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/apiserver-rest) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/apiserver/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/apiserver) |
| Frontend | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/frontend/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/frontend) |
| Device Registry | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/device-registry/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/device-registry) |
| Telemetry Router | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/telemetry-router/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/telemetry-router) |
| MQTT-Bridge | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/mqtt-bridge/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/mqtt-bridge) |
| Shadow | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-delta-merger/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-delta-merger) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-api/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-api) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-persister/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-persister) |

## API Documentation
You can find swagger docs for the API server [here](https://infinimesh.github.io/infinimesh/swagger-ui/)
