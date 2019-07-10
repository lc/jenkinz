# jenkinz
`jenkinz` is a tool to retrieve every build for every job ever created and run on a given Jenkins instance. 


## Usage:
`jenkinz -d https://[jenkins-instance] [options]`


### Flags:
| Flag | Description | Example |
|------|-------------|---------|
| `-d` | URL of Jenkins | `jenkinz -d https://jenkins.example.com` |
| `-creds` | Credentials for Jenkins instance (format user:apikey) | `jenkinz -d url -creds "admin:c129e5db6b5e3abdff6eb9b0008ad7f2"` |
| `-c` | Limit the number of workers that are spawned | `jenkinz -d url -c 200` |
| `-timeout` | Timeout for the tool in seconds (default 30) | `jenkinz -d url -timeout 10`|

#### Getting a Jenkins API Token
If anonymous read is disabled but you have credentials, you can generate an API Key by navigating to:
- `http://[jenkins]/user/username/configure`
- Under `API Token` click `Add New Token`

## Installation:
### Via `go get`
```
go get -u github.com/lc/jenkinz
```

### Via `git clone`

```
git clone git@github.com:lc/jenkinz
cd jenkinz && go build -o $GOPATH/bin/jenkinz main.go
```
