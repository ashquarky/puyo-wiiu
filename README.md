# PUYOPUYOTETRIS (Wii U) replacement server
Includes both the authentication and secure servers

## Compiling

### Setup
Install [Go](https://go.dev/doc/install) and [git](https://git-scm.com/downloads), then clone and enter the repository

```bash
$ git clone https://github.com/PretendoNetwork/puyopuyotetris
$ cd puyopuyotetris
```

### Compiling using `go`
To compile using Go, `go get` the required modules and then `go build` to your desired location. You may also want to tidy the go modules, though this is optional

```bash
$ go get -u
$ go mod tidy
$ go build -o build/puyopuyotetris
```

The server is now built to `build/puyopuyotetris`

When compiling with only Go, the authentication servers build string is not automatically set. This should not cause any issues with gameplay, but it means that the server build will not be visible in any packet dumps or logs a title may produce

To compile the servers with the authentication server build string, add `-ldflags "-X 'main.serverBuildString=BUILD_STRING_HERE'"` to the build command, or use `make` to compile the server

### Compiling using `make`
Compiling using `make` will read the local `.git` directory to create a dynamic authentication server build string, based on your repositories remote origin and current commit. It will also use the current folders name as the executables name

Install `make` onto your system (this varies by OS), and run `make` while inside the repository

```bash
$ make
```

The server is now built to `build/puyopuyotetris` with the authentication server build string already set

## Configuration
All configuration options are handled via environment variables

`.env` files are supported

| Name                                           | Description                                                                                     | Required                                      |
|------------------------------------------------|-------------------------------------------------------------------------------------------------|-----------------------------------------------|
| `PN_PUYOPUYOTETRIS_KERBEROS_PASSWORD`          | Password used as part of the internal server data in Kerberos tickets                           | No (Default password `password` will be used) |
| `PN_PUYOPUYOTETRIS_AUTHENTICATION_SERVER_PORT` | Port for the authentication server                                                              | Yes                                           |
| `PN_PUYOPUYOTETRIS_SECURE_SERVER_HOST`         | Host name for the secure server (should point to the same address as the authentication server) | Yes                                           |
| `PN_PUYOPUYOTETRIS_SECURE_SERVER_PORT`         | Port for the secure server                                                                      | Yes                                           |
| `PN_PUYOPUYOTETRIS_ACCOUNT_GRPC_HOST`          | Host name for your account server gRPC service                                                  | Yes                                           |
| `PN_PUYOPUYOTETRIS_ACCOUNT_GRPC_PORT`          | Port for your account server gRPC service                                                       | Yes                                           |
| `PN_PUYOPUYOTETRIS_ACCOUNT_GRPC_API_KEY`       | API key for your account server gRPC service                                                    | No (Assumed to be an open gRPC API)           |
| `PN_PUYOPUYOTETRIS_AES_KEY`                    | Account server token validation key                                                             | Yes                                           |
| `PN_PUYOPUYOTETRIS_LOCAL_AUTH`                 | Use local file for NEX credentials rather than talking to account server (INSECURE)             | No                                            |
| `PN_PUYOPUYOTETRIS_FRIENDS_GRPC_HOST`          | Host name for Friends gRPC service                                                              | Yes                                           |
| `PN_PUYOPUYOTETRIS_FRIENDS_GRPC_PORT`          | Port for Friends gRPC service                                                                   | Yes                                           |
| `PN_PUYOPUYOTETRIS_FRIENDS_GRPC_API_KEY`       | API key for Friends gRPC service                                                                | No (Assumed to be an open gRPC API)           |
| `PN_PUYOPUYOTETRIS_CONFIG_S3_ENDPOINT`         | URL for DataStore S3 storage                                                                    | Yes                                           |
| `PN_PUYOPUYOTETRIS_CONFIG_S3_ACCESS_KEY`       | Access key for S3                                                                               | Yes                                           |
| `PN_PUYOPUYOTETRIS_CONFIG_S3_SECRET`           | Secret key for S3                                                                               | Yes                                           |
| `PN_PUYOPUYOTETRIS_CONFIG_S3_BUCKET`           | Bucket for DataStore storage                                                                    | Yes                                           |
| `PN_PUYOPUYOTETRIS_CONFIG_S3_KEY_BASE`         | Base for Datastore keys                                                                         | YES                                           |
| `PN_PUYOPUYOTETRIS_CONFIG_S3_INSECURE`         | Use HTTP for S3 endpoints                                                                       | No (Default off)                              |