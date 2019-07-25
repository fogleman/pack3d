# pack3d


### Installation

First, install Go, set your GOPATH, and make sure $GOPATH/bin is on your PATH.

```
brew install go
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
export GOBIN=$GOPATH/bin
```

Clone the repo Authentise/pack3d (IN YOUR GO/SRC DIRECTORY!)

Then go to the project folder and do


```
go get github.com/fogleman/fauxgl
```

From source directory,
```
cd cmd/pack3d
go get
go install
```


From source directory,
```
cd cmd/binpack
go get
go install
```

Bin file is run using,
```
pack3d {frame_x,frame_y,frame_z} mini_spacing output_file_name model_num model_file
```

For example,
```
pack3d {100,100,100} 5 output 1 mesh1.stl 1 mesh2.stl
```

After running `pack3d`, it will generate a json file. The format of the json file is:

```
{"Filename":  , "Transformation":  , "VolumeWithSpacing":  }
```

For example:

```
[{"Filename":"Box.stl","Transformation":[[-1,0,-1.2246467991473515e-16,-17.991808688673732],[0,1,0,-19.997626237452177],[1.2246467991473515e-16,0,-1,22.451572094004227],[0,0,0,1]],"VolumeWithSpacing":16015.625}]
```
