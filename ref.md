# Reference
## FS Layout

```bash
> tree -L 2
.
├── Dockerfile
├── app
│   ├── build
│   │   └── app
│   └── cmd
├── docker-compose.yml
├── goapp
│   ├── build
│   ├── cmd
│   │   └── goapp
│   │       └── main.go
│   ├── configs
│   │   └── config.yaml
│   ├── docs
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   └── models
│   ├── pkg
│   └── scripts
├── ref.md
└── rename_app.sh
```

## Rename the app
```bash
./rename_app.sh <old_app_name> <new_app_name> # Default old_app_name is goapp
```

## Attach a shell to the container
```bash
docker run -it --entrypoint /bin/sh go-template-app
```