package canvus

// Contains returns true if rectangle a fully contains rectangle b.
//
// Usage Example:
//   a := canvus.Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
//   b := canvus.Rectangle{X: 2, Y: 2, Width: 5, Height: 5}
//   ok := canvus.Contains(a, b) // true
func Contains(a, b Rectangle) bool {
	return b.X >= a.X && b.Y >= a.Y &&
		b.X+b.Width <= a.X+a.Width &&
		b.Y+b.Height <= a.Y+a.Height
}

// Touches returns true if rectangles a and b overlap or touch at any edge/corner.
//
// Usage Example:
//   a := canvus.Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
//   b := canvus.Rectangle{X: 9, Y: 9, Width: 5, Height: 5}
//   ok := canvus.Touches(a, b) // true
func Touches(a, b Rectangle) bool {
	return a.X < b.X+b.Width && a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}

// WidgetBoundingBox returns the bounding box (Rectangle) for a Widget.
//
// Usage Example:
//   w := canvus.Widget{Location: &canvus.Point{X: 1, Y: 2}, Size: &canvus.Size{Width: 3, Height: 4}}
//   rect := canvus.WidgetBoundingBox(w)
func WidgetBoundingBox(w Widget) Rectangle {
	x, y := 0.0, 0.0
	wVal, hVal := 0.0, 0.0
	if w.Location != nil {
		x = w.Location.X
		y = w.Location.Y
	}
	if w.Size != nil {
		wVal = w.Size.Width
		hVal = w.Size.Height
	}
	return Rectangle{X: x, Y: y, Width: wVal, Height: hVal}
}

// WidgetContains returns true if widget a fully contains widget b.
//
// Usage Example:
//   ok := canvus.WidgetContains(widgetA, widgetB)
func WidgetContains(a, b Widget) bool {
	return Contains(WidgetBoundingBox(a), WidgetBoundingBox(b))
}

// WidgetsTouch returns true if widgets a and b touch or overlap.
//
// Usage Example:
//   ok := canvus.WidgetsTouch(widgetA, widgetB)
func WidgetsTouch(a, b Widget) bool {
	return Touches(WidgetBoundingBox(a), WidgetBoundingBox(b))
}
