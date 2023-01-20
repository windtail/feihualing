package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed star-fill.svg
var starFillData []byte

var starFillSvg = &fyne.StaticResource{
	StaticName:    "star-fill.svg",
	StaticContent: starFillData,
}

//go:embed star-outline.svg
var starOutlineData []byte

var starOutlineSvg = &fyne.StaticResource{
	StaticName:    "star-outline.svg",
	StaticContent: starOutlineData,
}
