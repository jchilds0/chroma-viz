# chroma-viz
Real time broadcast graphics application built in [Golang][go] using [GTK][gotk] go bindings.

## Features

[Chroma Hub][chroma-hub] sends a list of templates to Chroma Viz on startup in order to synchronize templates between Chroma Viz and [Chroma Engine][chroma-engine] instances.
Communicates with [Chroma Engine][chroma-engine] over tcp to render graphics.

- Templates can be added to the show to become pages.
- Each page has its own set of properties, editable through the editor
- [Chroma Engine][chroma-engine] combines the template and the data set in the editor to display the graphic
- Shows and Pages can be exported and imported.

[![Chroma_Viz](data/demo.png)](data/demo.mp4)

![Chroma_Engine](data/chroma-viz.png)

## Installation

- Install and build [Chroma Engine][chroma-engine].
- Install and run [Chroma Hub][chroma-hub].
- Set the `engDir` constants in `viz/preview.go` and `artist/preview.go` to the location of the Chroma Engine binary.
- Run 

```
go run main.go -mode [viz | artist] -hub [127.0.0.1:9000]
```

## Disclaimer

This is a personal project, not an application intended for production.

[go]: https://github.com/golang/go
[gotk]: https://github.com/gotk3/gotk3
[chroma-engine]: https://github.com/jchilds0/chroma-engine
[chroma-hub]: https://github.com/jchilds0/chroma-hub
