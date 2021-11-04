# pack3d

Pack3d is the geometry packing tool for 3d printing used by Authentise and can be found [here](https://github.com/Authentise/pack3d). Authentise's Pack3d codebase was forked from Fogleman's pack3d.

As it is now (mid January 2021), the Authentise/pack3d codebase is deemed to be stable - it did not have major upgrades for several months.

Pack3d is written in golang and the installation instructions can be found in the CONTRIBUTING.md


### Example

```
pack3d --input_config_json_filename=input.json --output_packing_json_filename=output
```

Notice the absence of the extension of the `output` file. This is because an `stl` file could optionally also be written as output.

Input example:

```
{
   "build_volume": [1, 0, 0],
   "spacing": 2,
   "items": [
      {
         "filename": "./tests/coprint/mesh_1.stl",
         "count": 2,
         "scale": 1.0,
         "copack": [
            {"filename": "./tests/coprint/mesh_2.stl"},
            {"filename": "./tests/coprint/mesh_3.stl"}
         ]
      },
      {
         "filename": "./tests/coprint/mesh_4.stl",
         "count": 2,
         "scale": 1.0,
      }
   ]
}
```

Output example (unrelated to the input example):

```
[
    {
        "Filename": "object1.stl",
        "Transformation":
        [
            [0, 0, -1, 0],
            [0, 1,  0, 0],
            [1, 0,  0, 0],
            [0, 0,  0, 1]
        ],
        "VolumeWithSpacing": 587.9278614633075
    }
    {
        "Filename": "object2.stl",
        "Transformation":
        [
            [0, 2, 0, -9.44458711204719],
            [0, 0, 2, -8.621629171058217],
            [2, 0, 0,  1.1454025845357874],
            [0, 0, 0,  1]
        ],
        "VolumeWithSpacing": 3908.6121061115296
    }
]
```
