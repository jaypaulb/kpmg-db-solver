package canvus

import "testing"

func TestContains(t *testing.T) {
	a := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	b := Rectangle{X: 2, Y: 2, Width: 5, Height: 5}
	c := Rectangle{X: -1, Y: -1, Width: 12, Height: 12}
	d := Rectangle{X: 9, Y: 9, Width: 2, Height: 2}

	if !Contains(a, b) {
		t.Error("a should contain b")
	}
	if Contains(a, c) {
		t.Error("a should not contain c")
	}
	if Contains(a, d) {
		t.Error("a should not contain d (partially outside)")
	}
}

func TestTouches(t *testing.T) {
	a := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	b := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}
	c := Rectangle{X: 10, Y: 0, Width: 5, Height: 5}
	d := Rectangle{X: 20, Y: 20, Width: 5, Height: 5}

	if !Touches(a, b) {
		t.Error("a should touch b (overlap)")
	}
	if !Touches(a, c) {
		t.Error("a should touch c (edge)")
	}
	if Touches(a, d) {
		t.Error("a should not touch d (no overlap)")
	}
}

func TestWidgetGeometry(t *testing.T) {
	a := Widget{Location: &Point{X: 0, Y: 0}, Size: &Size{Width: 10, Height: 10}}
	b := Widget{Location: &Point{X: 2, Y: 2}, Size: &Size{Width: 5, Height: 5}}
	c := Widget{Location: &Point{X: 9, Y: 9}, Size: &Size{Width: 2, Height: 2}}
	d := Widget{Location: &Point{X: 20, Y: 20}, Size: &Size{Width: 2, Height: 2}}

	if !WidgetContains(a, b) {
		t.Error("a should contain b")
	}
	if WidgetContains(a, c) {
		t.Error("a should not contain c (partially outside)")
	}
	if !WidgetsTouch(a, b) {
		t.Error("a should touch b (overlap)")
	}
	if !WidgetsTouch(a, c) {
		t.Error("a should touch c (corner)")
	}
	if WidgetsTouch(a, d) {
		t.Error("a should not touch d (no overlap)")
	}
}
