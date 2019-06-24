# jenkinz
`jenkinz` is a tool to retrieve every build for every job ever created and run on a given Jenkins instance. 


## Usage:
`jenkinz -d https://[jenkins-instance] [options]`


### Flags:
| Flag | Description | Example |
|------|-------------|---------|
| `-d` | url of jenkins instance | `jenkinz -d https://jenkins.example.com` |
| `-c` | Limit the number of workers that are spawned | `jenkinz -d https://jenkins.example.com -c 200` |
| `-timeout` | Timeout for the tool in seconds (default 30) | `jenkinz -d https://jenkins.example.com -timeout 10`|

## Installation:

### Via `git clone`

```
git clone git@github.com:lc/jenkinz
cd jenkinz && go build -o jenkinz main.go
```