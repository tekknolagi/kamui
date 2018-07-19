#!/usr/bin/env io

Point := Object clone do (
    x := 0
    y := 0

    set := method(x, y,
        self x := x
        self y := y
        self
    )

    + := method(p,
        newx := x + (p x)
        newy := y + (p y)
        Point clone set(newx, newy)
    )

    * := method(p,
        newx := x * (p x)
        newy := y * (p y)
        Point clone set(newx, newy)
    )

    asString := method(
        "Point(" .. x .. ", " .. y .. ")"
    )
)

Box := Object clone do (
    _origin := Point clone
    _size := Point clone

    /* this is origin, size -- not origin, topRight */
    set := method(origin, size,
        self _origin := origin
        self _size := size
        self
    )

    size := method(
        // Point clone set ((size x)-(origin x)) ((size y)-(origin y))
        _size
    )

    origin := method(
        _origin
    )

    asString := method(
        "Box @ " .. origin .. ", size " .. size
    )
)

Image := Object clone do (
    dims := Point clone

    width := method(
        dims x
    )

    height := method(
        dims y
    )
)

KUIWidget := Object clone do (
	/* parent & children here are used only for layout */
	parent := nil
	children := List clone

	bbox := Box clone
	parentcenter := Point clone
	center := Point clone
	offset := Point clone

	zlayer := 0

	img := Image clone
	img size := method(
		return Point clone set(width, height)
	)
	_calculateBbox := method(
		e := try(
			parentBbox := parent bbox
		)
		e catch (
			parentBbox := Box clone set(Point clone, Point clone)
		)
		bboxOrigin := ((parentBbox origin) + ((parentBbox size)*(parentcenter))) + offset
		bboxSize := Point clone set (img width, img height)
		self bbox := Box clone set(bboxOrigin, bboxSize)
		?widgetLint
		return bbox
	)

	size := method( return bbox size )

	_calculateSize := method ( return size )

	calculateBbox := method(
		_calculateBbox
		children foreach(c, c ?calculateBbox)
	)
	calculateSize := method(
		children foreach(c, c ?calculateSize)
		return _calculateSize
	)
	setParent := method(p,
		self parent := p
		p children append(self)
	)
)

KUIContainer := KUIWidget clone do (
	widgetLint := method(
		if(children size < 1,
			children foreach(c,
				if(c zlayer <= zlayer, c zlayer := (zlayer + 1))
			)
			children foreach(c, c ?widgetLint)
			children foreach(c,
				self bbox := bbox Union(c bbox)
			)
		)
		return bbox
	)
)
