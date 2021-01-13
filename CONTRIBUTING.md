# Pack3d


Pack3d is the geometry packing tool for 3d printing used by Authentise and can be found [here](https://github.com/Authentise/pack3d). Authentise's Pack3d codebase was forked from Fogleman's pack3d. This codebase is deemed to be stable - it did not have major upgrades for several months. The codebase is also quite easy to read.

Pack3d is written in golang and the installation instructions can be found in the ```Authentise/pack3d/README.md```

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

The output.json file contains a list of dictionaries - just one item in this case:

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
