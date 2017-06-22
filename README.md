# pack3d

Tightly pack 3D models.

### Installation

First, install Go, set your GOPATH, and make sure $GOPATH/bin is on your PATH.

```
brew install go
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

Next, fetch and build the two binaries.

```
go get github.com/fogleman/pack3d/cmd/pack3d
go get github.com/fogleman/pack3d/cmd/binpack
```

### Usage Examples

Note that pack3d runs until stopped, writing its output to disk whenever a new best is found.

```
pack3d 2 3DBenchy.stl  # tightly pack 2 boats
pack3d 4 3DBenchy.stl  # tightly pack 4 boats
pack3d 1 *.stl         # tightly pack various meshes, one of each

# pack as many boats as possible into the printer volume, given a few different arrangements
binpack 1 3DBenchy.stl 2 3DBenchy-x2.stl 4 3DBenchy-x4.stl
```

### Examples

113 3DBenchy tug boats packed tightly

![3DBenchy](http://i.imgur.com/adjchjy.png)

27 R2-D2 droids, 8 parts each

![R2-D2](http://i.imgur.com/qE90ijK.png)
