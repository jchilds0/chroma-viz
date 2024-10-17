package hub

const (
	TEMPLATE_INSERT = "template-insert"
	TEMPLATE_SELECT = "template-select"
	TEMPLATE_DELETE = "template-delete"

	GEOMETRY_INSERT  = "geometry-insert"
	RECTANGLE_INSERT = "rectangle-insert"
	TEXT_INSERT      = "text-insert"
	CIRCLE_INSERT    = "circle-insert"
	ASSET_INSERT     = "asset-insert"
	CLOCK_INSERT     = "clock-insert"
	POLYGON_INSERT   = "polygon-insert"
	POINT_INSERT     = "point-insert"
	LIST_INSERT      = "list-insert"
	ROW_INSERT       = "row-insert"

	GEOMETRY_SELECT  = "geometry-select"
	RECTANGLE_SELECT = "rectangle-select"
	CIRCLE_SELECT    = "circle-select"
	TEXT_SELECT      = "text-select"
	ASSET_SELECT     = "asset-select"
	CLOCK_SELECT     = "clock-select"
	POLYGON_SELECT   = "polygon-select"
	POINT_SELECT     = "point-select"
	LIST_SELECT      = "list-select"
	ROW_SELECT       = "row-select"
)

var stmts = map[string]string{
	TEMPLATE_INSERT: `
        INSERT INTO template VALUES (?, ?, ?);
    `,
	TEMPLATE_SELECT: `
        SELECT t.Name, t.Layer, COUNT(*)
        FROM template t
        INNER JOIN geometry g 
        ON g.templateID = t.templateID
        WHERE t.templateID = ?;
    `,
	TEMPLATE_DELETE: `
        DELETE FROM template WHERE templateID = ?;
    `,

	GEOMETRY_INSERT: `
        INSERT INTO geometry VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?);
    `,
	RECTANGLE_INSERT: `
        INSERT INTO rectangle VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `,
	TEXT_INSERT: `
        INSERT INTO text VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `,
	CIRCLE_INSERT: `
        INSERT INTO circle VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
    `,
	ASSET_INSERT: `
        INSERT INTO asset VALUES (?, ?, ?, ?, ?);
    `,
	CLOCK_INSERT: `
        INSERT INTO clock VALUES (?, ?, ?, ?, ?, ?);
    `,
	POLYGON_INSERT: `
       INSERT INTO polygon VALUES (?, ?, ?, ?, ?);
    `,
	POINT_INSERT: `
       INSERT INTO point VALUES (?, ?, ?, ?);
    `,
	LIST_INSERT: `
        INSERT INTO list VALUES (?, ?, ?, ?, ?, ?, ?);
    `,
	ROW_INSERT: `
        INSERT INTO row VALUES (?, ?, ?);
    `,

	GEOMETRY_SELECT: `
        SELECT g.geometryID, g.geoNum, g.name, g.geoType, g.rel_x, g.rel_y, g.parent, g.mask
        FROM geometry g 
        WHERE g.templateID = ?;
    `,
	RECTANGLE_SELECT: `
        SELECT r.geometryID, r.width, r.height, r.rounding, r.red, r.green, r.blue, r.alpha
        FROM rectangle r 
        INNER JOIN geometry g
        ON r.geometryID = g.geometryID
        WHERE g.templateID = ?;
    `,
	CIRCLE_SELECT: `
        SELECT c.geometryID, c.inner_radius, c.outer_radius, c.start_angle, c.end_angle, c.red, c.green, c.blue, c.alpha
        FROM circle c
        INNER JOIN geometry g
        ON c.geometryID = g.geometryID
        WHERE g.templateID = ?;
    `,
	TEXT_SELECT: `
        SELECT t.geometryID, t.text, t.fontSize, t.red, t.green, t.blue, t.alpha
        FROM text t
        INNER JOIN geometry g
        ON g.geometryID = t.geometryID
        WHERE g.templateID = ?;
    `,
	ASSET_SELECT: `
        SELECT a.geometryID, a.directory, a.name, a.assetID, a.scale
        FROM asset a 
        INNER JOIN geometry g 
        ON a.geometryID = g.geometryID 
        WHERE g.templateID = ?;
    `,
	CLOCK_SELECT: `
        SELECT c.geometryID, c.scale, c.red, c.green, c.blue, c.alpha
        FROM clock c
        INNER JOIN geometry g
        ON g.geometryID = c.geometryID
        WHERE g.templateID = ?;
    `,
	POLYGON_SELECT: `
        SELECT p.geometryID, p.red, p.green, p.blue, p.alpha
        FROM polygon p
        INNER JOIN geometry g 
        ON p.geometryID = g.geometryID 
        WHERE g.templateID = ?;
    `,
	POINT_SELECT: `
        SELECT point.pointID, point.pos_x, point.pos_y
        FROM point
        INNER JOIN polygon
        INNER JOIN geometry g
        ON point.geometryID = polygon.geometryID
        AND polygon.geometryID = g.geometryID
        WHERE g.geoNum = ?
        AND g.templateID = ?;
    `,
	LIST_SELECT: `
        SELECT l.geometryID, l.red, l.green, l.blue, l.alpha, l.single_row, l.scale
        FROM list l
        INNER JOIN geometry g
        ON g.geometryID = l.geometryID
        WHERE g.templateID = ?;
    `,
	ROW_SELECT: `
        SELECT r.rowID, r.row
        FROM row r
        INNER JOIN list l 
        INNER JOIN geometry g
        ON r.geometryID = l.geometryID
        AND l.geometryID = g.geometryID
        WHERE g.geoNum = ?
        AND g.templateID = ?;
    `,
}
