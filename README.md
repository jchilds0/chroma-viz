# Chroma Graphics
Real time broadcast graphics application built in [Golang][go] using [GTK][gotk] go bindings.

## Installation

- Copy `default_conf.json` in `artist/` and `viz/` and rename to `conf.json`.
- Download and build [Chroma Engine][chroma-engine].
The preview window uses `cgo` to create a GtkGLRender window using [Chroma Engine][chroma-engine].
The expected file structure is shown below, either copy or symlink files from [Chroma Engine][chroma-engine].
For an alternative structure, modify the compiler flags in `library/preview.go`.

```
chroma-viz
├───library
│   └───preview.go
├───chroma-engine
│   ├───libchroma.a
│   └───chroma-engine.h
...
```

- Setup a sql database (e.g. `mariadb`)

### Chroma Hub 

- (First Install) Import the db schema in `library/hub/chroma_hub.sql` to the sql db or start Chroma Hub in the next step with the -c flag.
- Run Chroma Hub 
```
go run ./cmd/chroma-hub -u <username> -p <password>
```
where `username` and `password` correspond to the user login for the sql database. 

- Chroma Hub creates a REST api at the port specified on startup, which is makes assets available to Chroma Viz, Artist and Engine.

### Chroma Viz

- Run Chroma Viz
```
go run ./cmd/chroma-viz
```
### Chroma Artist 

- Run Chroma Artist
```
go run ./cmd/chroma-artist
```

## Features

Chroma Viz retrieves a list of templates from Chroma Hub on startup.
Chroma Viz communicates with [Chroma Engine][chroma-engine] over tcp to render graphics.

- Pages can be easily created from templates by double clicking on the template in the template list.
- Each page has its own set of properties, editable through the editor
- [Chroma Engine][chroma-engine] combines the template and the data set in the editor to render the graphic.

https://github.com/user-attachments/assets/2203a13e-ccde-4edd-8170-44f922fc1997

Chroma Artist can be used to design templates, which can be imported to Chroma Hub.

- Tree View for creating the heirachy of geometry elements
- Keyframes for creating animations, by setting a geometry attribute to a value, the value given by the user, or a value of another attribute.
- Import/Export templates and assets to Chroma Hub.
 
https://github.com/user-attachments/assets/6b88397d-30d7-447f-b158-35ad6523b273

## Disclaimer

This is a personal project, not an application intended for production.

[go]: https://github.com/golang/go
[gotk]: https://github.com/gotk3/gotk3
[chroma-engine]: https://github.com/jchilds0/chroma-engine
