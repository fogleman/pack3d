# pack3d - Installation, Codebase, Development and Deployment

The Go Development can be carried out straight into the Nautilus container, same as for the Python development. There are a few ways to compile/run that Go code:

   1. set up the Go dev environment in the Nautilus docker container (Ubuntu).
   2. Set up the Go dev environment in macOS with the `go` folder inside the HOME folder.

The first option is probably preferable since the deployment of compiled/executable Go code will eventually need to go through the first step anyway, but someone might find it helpful to develop in macOS.

## 1. Go Development Setup (Ubuntu 18.04.6 LTS - using `fish`)

(TODO: the addition of Go-related env vars would be nice to avoid long paths.)

Start off `fish` in the Nautilus docker container.

Install golang (go1.16.6.linux-amd64.tar.gz is 123.07MB)

```
wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz
```

Untarring might take a few minutes.
The untarred ./go folder is about 395MB.

```
tar -C /src -xzf go1.16.6.linux-amd64.tar.gz
```

The tar file can be removed, also check the version of the Go compiler.

```
rm go1.16.6.linux-amd64.tar.gz
/src/go/bin/go version
```

This should return the line below:

`go version go1.16.6 linux/amd64`

Get the necessary repos from github:

```
cd go/src
/src/go/bin/go get github.com/fogleman/fauxgl
```

clone the pack3d codebase (do this in macOS's iTerm so you can clone pack3d using the github tokens):

```
cd ./Authentise/nautilus/go/src
git clone git@github.com:Authentise/pack3d.git
```

back in fish:

```
/src/go/bin/go mod vendor
```

### Compile Go code:

```
cd /src/go/src/pack3d

cd cmd/pack3d
/src/go/bin/go get     # ignore this line if it errs!
/src/go/bin/go install

cd ../binpack/
/src/go/bin/go get
/src/go/bin/go install
```

Now, the (compiled) executable of `pack3d` is available in the `/src/go/bin/`.

### Manual testing of Go code:

Co-packing tests. These ones might quickly become obsolete as newer pack3d features are added in.

```
cd /src/go/src/pack3d
```

`/src/go/bin/pack3d --input_config_json_filename=tests/jenkins_tests/input_jenkins_test_1.json --output_packing_json_filename=tests/jenkins_tests/output_jenkins_test_1`

`/src/go/bin/pack3d --input_config_json_filename=tests/jenkins_tests/input_jenkins_test_2.json --output_packing_json_filename=tests/jenkins_tests/output_jenkins_test_2`

---

## 2. Go Development Setup (macOS - Homebrew) [OPTIONAL]

The `go` folder in the home folder is to contain the projects' source code and bin/exec files.

```
export GOPATH="${HOME}/go"
export GOROOT="$(brew --prefix golang)/libexec"
export PATH="$PATH:${GOPATH}/bin:${GOROOT}/bin"
test -d "${GOPATH}" || mkdir "${GOPATH}"
test -d "${GOPATH}/src/github.com" || mkdir -p "${GOPATH}/src/github.com"
```

```
brew install go
```

The installed version of `go` is assumed to be 1.16+ in macOS:

```
go env -w GO111MODULE=auto  # necessary for go 1.16+ to use the legacy `go get` incantation.
```

Retrieve the `fogleman/fauxgl` code necessary for 3d renders within pack3d.
Unused in Nautilus at the moment of writing this doc but a necessary dependency of pack3d.

```
cd go/src
go get github.com/fogleman/fauxgl
```

Clone the Authentise/pack3d repository into your `go/src` folder.

```
git clone git@github.com:Authentise/pack3d.git
```

Install pack3d.

```
cd pack3d
cd cmd/pack3d
go get; go install
```

```
cd ../..
cd cmd/binpack
go get; go install
```

pack3d can be invoked directly from the command line now, see an example below here:

`pack3d --input_config_json_filename=tests/ch32838_test/input.json --output_packing_json_filename=tests/ch32838_test/output/output`

See folder pack3d/tests/ch32838_test for an example of input files, or the folder pack3d/tests/jenkins_tests.

NB: Nautilus can only use a binary file built on a linux machine and not from macOS. See 4. Deployment.

---

## 3. Quick glance at the pack3d codebase

### cmd/pack3d/main.go
```go/src/github.com/pack3d/cmd/pack3d/main.go```

The ```func main()``` performs command line's arguments parsing and related actions:

- pack3d version which can be invoked in the CLI.
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

The numbers, expressed in scientific notation, have been shortened for visualisation purposes. Very small numbers like: 6.123E-17 and 3.749E-33 can be safely approximated with 0.0.

The Transformation's value represents an affine transformation (for homogeneous coordinates). The VolumeWithSpacing's value is the minimum distance between two items in the built-volumn.

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


## 4. Development

__WARNING__: After creating a PR, you can select which master (fogleman/master or Authentise/master) you want to merge into.

TODO

## 5. Deployment

Pack3d build location

`/src/go/bin/pack3d`

This build needs to be copied into the Nautilus bin when deplyoing.
