# bezier-shading

**bezier-shading** is a desktop bezier shading application with [multi-os support](https://github.com/fyne-io/fyne/wiki/Supported-Platforms). It allows you to create bezier surfaces, and simulate shadows.

## installation

In order to build and run project you need to install:

- [Go 1.21](https://go.dev/doc/install)

- [Fyne](https://developer.fyne.io/)

- [Fyne dependencies](https://developer.fyne.io/started/)

For windows, recomennded C compiler is [TDM-GCC](https://jmeubank.github.io/tdm-gcc/download/).

## building & running

Linux

```sh
$ make run
```

Windows

```sh
go build -o bin\bezier-shading cmd\bezier-shading\main.go

.\bin\bezier-shading
```

## drawing

Circles are drawn using [midpoint circle algoritm](https://en.wikipedia.org/wiki/Midpoint_circle_algorithm).

Lines are drawn using [Bresenham'slinealgorithm](https://en.wikipedia.org/wiki/Bresenham's_line_algorithm).

## screenshots

TODO

## author

Author of this project is Jakub Rudnik.
