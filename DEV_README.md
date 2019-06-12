# pack3d


### Installation

First, install Go, set your GOPATH, and make sure $GOPATH/bin is on your PATH.

```
brew install go
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

Clone the repo Authentise/pack3d (IN YOUR GO/SRC DIRECTORY!)

Then go to the project folder and do


```
go get github.com/fogleman/fauxgl
```

```
go install cmd/pack3d/main.go
```