# pack3d - Installation, Codebase and Development


## 1. Installation (macOS users)

First, install Go, set your GOPATH, and make sure $GOPATH/bin is on your PATH.

```
brew install go

mkdir ~/go
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
export GOBIN=$GOPATH/bin
```

Then go to the src folder and do:

```
go get github.com/fogleman/fauxgl
```

Clone the repo Authentise/pack3d in your ```go/src``` folder.

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

Bin file is run using, where frame_x, frame_y, frame_z are the build_crate size.
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
[
   {
      "Filename": "Box.stl",
      "Transformation": [
          [ -1,           0, -1.224e-16, -17.991 ],
          [  0,           1,  0,         -19.997 ],
          [  1.224e-16,   0, -1,          22.451 ],
          [  0,           0,  0,          1      ]
     ],
     "VolumeWithSpacing": 16015.625
   }
]
```

NB: Nautilus can only use a binary file built on a linux machine.

## 2. Codebase

What follows below is a quick peek into a few important files:

### cmd/pack3d/main.go
```go/src/github.com/pack3d/cmd/pack3d/main.go```

The ```func main()``` performs command line's arguments parsing and related actions:

- pack3d version which can be invoked in the CLI with ```pack3d --version```.
- new pack3d model creation: ```model := pack3d.NewModel()```
- build-volume variables set up. Notice ```annealingIterations = 2000000```
- STL geometries loading and addition of each loaded geometry to the pack3d model: ```model.Add(mesh, bvhDetail, count, spacing)```
- packing of the pack3d model.
	- ```model.Pack(annealingIterations, nil, singleStlSize, frameSize, packItemNum)```
	- there is a time limit to break out if the packing algotithm is struggling.
	- binary search. Uncommenting some code allows serving up an error-related json file.
- json data generation. More specifically a list of transformations applied to each original model inside of the build volume.
- creation of STL pack3d model file - requires uncommenting some code. ```model.Mesh().SaveSTL( ... )```
```model.TreeMesh().SaveSTL( ... )```

The output.json file contains a list of dictionaries - just one item was added (to the 3dpack-model) in this case:

```
[
  {
    "Filename": "logo.stl",
    "Transformation": [
      [  3.749E-33,  6.123E-17,  -1,          30.733 ],
      [ -1,          6.123E-17,   0,          26.604 ],
      [  6.123E-17,  1,           6.123E-17, -29.382 ],
      [  0,          0,           0,           1     ]
    ],
    "VolumeWithSpacing": 780.5089485754926
  }
]
```

The numbers, expressed in scientific notation, have been shortened for visualisation purposes. The Transformation's value represents an affine transformation (for homogeneous coordinates). The VolumeWithSpacing's value is the minimum distance between two items in the built-volumn.

### pack3d/model.go

```go/src/github.com/pack3d/pack3d/main.go```

The model class owns methods some of which have already been highlighted above, in the file main.go. Some of those and other methods are recapped below here for convenience.

|Method|Purpose|
|-|-|
|Add( ... )|Add geometry to Model|
|Pack( ... )|Model packing|
|DoMove( ... ), UndoMove( ... )|Item rotation / translation|
|BoundingBox()| get BoundingBox |
|TreeMesh()|get TreeMesh|
|Volume()|get Volume |
|Copy()| Copy |
|...|...|


## 3. Development

__WARNING__: After creating a PR, you can select which master (fogleman/master or Authentise/master) you want to merge into.