# pack3d

Pack3d is the geometry packing tool for 3d printing used by Authentise and can be found [here](https://github.com/Authentise/pack3d). Authentise's Pack3d codebase was forked from Fogleman's pack3d.

As it is now (mid January 2021), the Authentise/pack3d codebase is deemed to be stable - it did not have major upgrades for several months.

Pack3d is written in golang and the installation instructions can be found in the CONTRIBUTING.md


### Example

```
pack3d {100,100,100} 5 output 1 mesh1.stl 1 mesh2.stl
```

Output example:

```
[
   {
      "Filename": "Box.stl",
      "Transformation": [
          [ -1,           0,  -1.224e-16, -17.991 ],
          [  0,           1,   0,         -19.997 ],
          [  1.224e-16,   0,  -1,          22.451 ],
          [  0,           0,   0,          1      ]
     ],
     "VolumeWithSpacing": 16015.625
   }
]
```

