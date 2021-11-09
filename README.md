# pack3d

Pack3d is the geometry packing tool for 3d printing used by Authentise and can be found [here](https://github.com/Authentise/pack3d). Authentise's Pack3d codebase was forked from Fogleman's pack3d.

Pack3d is written in golang and the installation instructions can be found in the CONTRIBUTING.md


## Invoking pack3d from the command line - example

```
pack3d --input_config_json_filename=input.json --output_packing_json_filename=output
```

Notice the absence of the extension of the `output` file. This is because an `stl` file could optionally also be written as output by pack3d.

## Input example:

```
{
    "build_volume": [100, 100, 100],
    "spacing": 5,
    "items": [
        {
            "filename": "tests/jenkins_tests/logo.stl",
            "count": 3,
            "scale": 1.0,
            "copack": [
                {
                    "filename": "tests/jenkins_tests/corner.stl"
                }
            ]
        },
        {
            "filename": "tests/jenkins_tests/logo.stl",
            "count": 2,
            "scale": 2.0,
            "copack": [
                {
                    "filename": "tests/jenkins_tests/corner.stl"
                }
            ]
        },
        {
            "filename": "tests/jenkins_tests/cube.stl",
            "count": 2,
            "scale": 4.0
        },
        {
            "filename": "tests/jenkins_tests/cube.stl",
            "count": 5,
            "scale": 1.0
        }
    ]
}
```

## Output example (related to the input example):

1. The co-packed objects have VolumeWithSpacing = 0. This is because their volume is already contemplated in the value of the main co-packing object's VolumeWithSpacing.

2. Notice the scaling visible in the 3x3 rotation matrix.

3. pack3d can either fail to pack a set of objects entirely - an error status is displayed in the command line, or pack3d can manage to pack fewer objects in such case the objects that did not make it into the build volume will have a null Transformation = `[0, 0, 0, 0], [0, 0, 0, 0], [0, 0, 0, 0], [0, 0, 0, 1]`.


```
[
    {
        "Filename": "tests/jenkins_tests/logo.stl",
        "Transformation": [
            [ 0, 1,  0, -38.08436791719043],
            [ 0, 0, -1, -3.36183820906897],
            [-1, 0,  0,  31.10764712961166],
            [ 0, 0,  0,  1]
        ],
        "VolumeWithSpacing": 964.5577062418419
    },
    {
        "Filename": "tests/jenkins_tests/corner.stl",
        "Transformation": [
            [ 0, 1,  0, -38.08436791719043],
            [ 0, 0, -1, -3.36183820906897],
            [-1, 0, 0,   31.10764712961166],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 0
    },
    {
        "Filename": "tests/jenkins_tests/logo.stl",
        "Transformation": [
            [-1,  0,  0, -35.064518578131974],
            [ 0,  0, -1, -12.783828938245192],
            [ 0, -1,  0, -24.15019322817524],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 964.5577062418419
    },
    {
        "Filename": "tests/jenkins_tests/corner.stl",
        "Transformation": [
            [-1,  0,  0, -35.064518578131974],
            [ 0,  0, -1, -12.783828938245192],
            [ 0, -1,  0, -24.15019322817524],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 0
    },
    {
        "Filename": "tests/jenkins_tests/logo.stl",
        "Transformation": [
            [ 0, 1, 0, -35.521257288959845],
            [-1, 0, 0, -4.392353279184597],
            [ 0, 0, 1, -9.788089697461157],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 964.5577062418419
    },
    {
        "Filename": "tests/jenkins_tests/corner.stl",
        "Transformation": [
            [ 0, 1, 0, -35.521257288959845],
            [-1, 0, 0, -4.392353279184597],
            [ 0, 0, 1, -9.788089697461157],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 0
    },
    {
        "Filename": "tests/jenkins_tests/logo.stl",
        "Transformation": [
            [ 0, 0, -2, -36.092921290618406],
            [-2, 0,  0, -4.735731505145346],
            [ 0, 2,  0, 8.056191563929794],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 5138.241184594143
    },
    {
        "Filename": "tests/jenkins_tests/corner.stl",
        "Transformation": [
            [ 0, 0, -2, -36.092921290618406],
            [-2, 0,  0, -4.735731505145346],
            [ 0, 2,  0, 8.056191563929794],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 0
    },
    {
        "Filename": "tests/jenkins_tests/logo.stl",
        "Transformation": [
            [ 0, 0, 2, -36.09345621544282],
            [ 2, 0, 0, -21.623185522659124],
            [ 0, 2, 0, -8.785596546107582],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 5138.241184594143
    },
    {
        "Filename": "tests/jenkins_tests/corner.stl",
        "Transformation": [
            [ 0, 0, 2, -36.09345621544282],
            [ 2, 0, 0, -21.623185522659124],
            [ 0, 2, 0, -8.785596546107582],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 0
    },
    {
        "Filename": "tests/jenkins_tests/cube.stl",
        "Transformation": [
            [ 0, 4, 0, -4.439943270386402],
            [ 0, 0, 4, -11.647772698025165],
            [ 4, 0, 0, -28.684157525681382],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 76765.625
    },
    {
        "Filename": "tests/jenkins_tests/cube.stl",
        "Transformation": [
            [ 0, 4, 0, -4.439858883393925],
            [ 0, 0, 4, -11.646377522333808],
            [ 4, 0, 0,  16.31737399791252],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 76765.625
    },
    {
        "Filename": "tests/jenkins_tests/cube.stl",
        "Transformation": [
            [ 0, 1, 0, -37.593934442604954],
            [ 0, 0, 1, -25.386707045922037],
            [ 1, 0, 0, -36.10923395402634],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 1953.125
    },
    {
        "Filename": "tests/jenkins_tests/cube.stl",
        "Transformation": [
            [ 0, 1, 0, -35.6471060803817],
            [-1, 0, 0, -26.355432623379993],
            [ 0, 0, 1,  13.959681356544618],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 1953.125
    },
    {
        "Filename": "tests/jenkins_tests/cube.stl",
        "Transformation": [
            [ 0, -1, 0, -35.80311229074046],
            [ 1,  0, 0, -3.854924315984657],
            [ 0,  0, 1, -39.99880616458849],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 1953.125
    },
    {
        "Filename": "tests/jenkins_tests/cube.stl",
        "Transformation": [
            [ 0, 0, -1, -35.48465191236189],
            [ 0, 1,  0,  2.5540333920688214],
            [ 1, 0,  0, -24.967430047438363],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 1953.125
    },
    {
        "Filename": "tests/jenkins_tests/cube.stl",
        "Transformation": [
            [ 0,  0, 1, -37.2963545459397],
            [ 0, -1, 0, -20.597736843343952],
            [ 1,  0, 0,  30.034134932756594],
            [ 0, 0, 0, 1]
        ],
        "VolumeWithSpacing": 1953.125
    }
]
]
```
