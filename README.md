# pack3d

Pack3d is the geometry packing tool for 3d printing used by Authentise and can be found [here](https://github.com/Authentise/pack3d). Authentise's Pack3d codebase was forked from Fogleman's pack3d.

As it is now (mid January 2021), the Authentise/pack3d codebase is deemed to be stable - it did not have major upgrades for several months.

Pack3d is written in golang and the installation instructions can be found in the CONTRIBUTING.md


### Example

```
pack3d --json_file=input.json --filename=output
```

Output example:

```
{
   "build_volume": [1, 0, 0],
   "spacing": 2,
   "items": [
      {
         "filename": "./tests/coprint/Christmas_1.stl",
         "count": 1,
         "copack": [
            {
               "filename": "./tests/coprint/Merry.stl",
               "transformation": [
                  [1, 0, 0, -60],
                  [0, 1, 0, -343],
                  [0, 0, 1, 0],
                  [0, 0, 0, 1]
               ]
            }
         ]
      },
      {
         "filename": "./tests/coprint/logo.stl",
         "count": 1
      }
   ]
}
```
