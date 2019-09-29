# WeatherClock API

simple service to getting forecast for IoT projects with data transformation and designed to support different forecast sources

[OpenAPI 3.0 doc](third_party/redoc-static.html) ([source](api/openapi.yaml))

## Build

### AWS Lambda

`make aws-lambda-build` and upload `aws-lambda-weatherclockapi.zip` (executable file called `main`)

### RestAPI server

* `make restapi-server-build`

* deploy to your server
* run `weatherclock-api -bind :5000` to listen 500 port on all available interfaces

### Command line tool

`make cmdtool-build` and `./weatherclock-cli --help`