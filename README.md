# chroma-viz
Real time broadcast graphics application built in [Golang][go] using [GTK][gotk] go bindings.

## Features

Communicates with [Chroma Engine][chroma-engine] over tcp to render graphics.

- Custom Templates which correspond to assets in [Chroma Engine][chroma-engine]
- Templates can be added to the show to become pages.
- Each page has its own set of properties, editable through the editor
- [Chroma Engine][chroma-engine] combines the template and the data set in the editor to display the graphic

![Chroma_Engine](data/chroma-viz.png)

## Disclaimer

This is a personal project, not an application intended for production.

[go]: https://github.com/golang/go
[gotk]: https://github.com/gotk3/gotk3
[chroma-engine]: https://github.com/jchilds0/chroma-engine
