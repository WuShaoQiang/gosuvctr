# Gosuv Controller

## Installation

### Build from source
```sh
go get github.com/WuShaoQiang/gosuvctr
mkdir -p ~/.gosuvctr/conf
mv $GOPATH/src/github.com/WuShaoQiang/gosuvctr/config.json ~/.gosuvctr/conf
vim ~/.gosuvctr/conf/config.json ##load your config file

cd $GOPATH/src/github.com/WuShaoQiang/gosuvctr
go build
```

### Example

After `go build`, there is a binary file called `gosuvctr`,you can move this binary file to `/user/local/bin`.

```sh
gosuvctr status
```

you should see the result like this:

```sh
PROGRAM NAME            STATUS
test                    fatal
shutdown server         stopped
```

## Configuration

config file should be stored in directory `$HOME/.gosuvctr/conf/`

- `config.json`

```json
{
    "admin":{
        "username":"",
        "password":""
    },
    "remoteAddr":"",
    "remotePort":""
}
```

