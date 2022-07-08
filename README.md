# vscode-cmake-configure

A tool for configuring a CMake project on the command line, using Visual Studio Code settings file `.vscode/settings.json`.

NOTE: this tool does not require VSCode to be available.

[![Go](https://github.com/MawKKe/vscode-cmake-configure/actions/workflows/go.yml/badge.svg)](https://github.com/MawKKe/vscode-cmake-configure/actions/workflows/go.yml)

## Motivation

Let's say are working on a CMake project with VSCode. You have saved various settings into `.vscode/settings.json`, such as:

```
{
    "editor.formatOnSave": true,
    "cmake.configureOnOpen": true,
    "cmake.configureArgs": [
        "-GNinja"
    ],
    // A comment here
    "cmake.configureSettings": {
		"CMAKE_CXX_COMPILER": "clang++",
		"CMAKE_CXX_FLAGS_INIT": "-fdiagnostics-color=always -O0",
		"CMAKE_CXX_STANDARD_REQUIRED": "ON",
		"CMAKE_CXX_STANDARD": "17" // just an example, not the right way to do this
    },
    "cmake.ctestArgs": []
     // ... and rest of your settings.json
}
```
You usually configure and build the project using VSCode CMake extension that automatically collects these options and runs the appropriate `cmake` commands for you.

*However*, now you would like to configure the project without VScode, perhaps because you are working via terminal, or VSCode is unavailable for some reason. You realize that to configure the project, **you will need to copy all the options manually** (one-by-one) from `.vscode/settings.json` into your command line:

    $ cmake \
        -DCMAKE_CXX_COMPILER=clang++ \
        -DCMAKE_CXX_FLAGS_INIT='-fdiagnostics-color=always -O0' \
        -DCMAKE_CXX_STANDARD_REQUIRED=ON \
        -DCMAKE_CXX_STANDARD=17 \
        -B mybuild \
        .

    $ cmake --build mybuild

while this way of working gets the job done, it is rather tiresome, wouldn't you agree?

**Alternatively** you could use this simple helper program:

    $ vscode-cmake-configure -B mybuild .
    $ cmake --build mybuild

Did you notice how easy that was? No copy-pasting, no messing around. Just one command, and presto! You are ready to build! You can even augment the `settings.json` options by specifying them as arguments to `vscode-cmake-configure` (see Usage below).

## Usage

Typical usage:

    $ vscode-cmake-configure <normal-cmake-arguments>

for example:

    $ vscode-cmake-configure -B path/to/mybuild -DCMAKE_FOO_VAR=BAZ path/to/src

behind the scenes this command will collect its command line arguments, and options/settings from `.vscode/settings.json` and call `cmake` using them

Note that the program will always print the full constructed  `cmake` command before its execution. This is to ensure that you (the programmer) is always aware what is going on.

If you wish to *only* see the command without actually executing it, specify environment variable `VCC_DRY_RUN` :

    $ env VCC_DRY_RUN=1 vscode-cmake-configure -B mybuild .

By default the program expects the VSCode settings to exist in `$PWD/.vscode/settings.json`, but you may specify alternative path via environment variable `VCC_VSCODE_SETTINGS`:

    $ env VCC_VSCODE_SETTINGS=path/to/mysettings.json vscode-cmake-configure ...

## Install

    $ go install github.com/MawKKe/vscode-cmake-configure@latest

## Dependencies

This program is written with Go 1.18, but it may compile with earlier versions.
The program expects the `cmake` executable to be found in your `$PATH`.

## Note about JSON and comments

VSCode allows `.vscode/settings.json` to contain C-style comments, which is not valid in standard JSON. This program also supports commented JSON, allowing you to use the settings file as is.

The support for commented JSON is provided by third party library (https://github.com/tidwall/jsonc).

## License

Copyright 2022 Markus Holmström (MawKKe)

The works under this repository are licenced under Apache License 2.0.
See file `LICENSE` for more information.

## Contributing

This project is hosted at https://github.com/MawKKe/vscode-cmake-configure

You are welcome to leave bug reports, fixes and feature requests. Thanks!
