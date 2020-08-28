/*
 * proofable-image
 * Copyright (C) 2020  Southbank Software Ltd.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * @Author: guiguan
 * @Date:   2020-07-30T00:23:27+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-28T13:03:43+10:00
 */

package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func getImageDrawingSize(imgSize image.Point, containerSize image.Point) f32.Point {
	imgRatio := float32(imgSize.X) / float32(imgSize.Y)

	if float32(containerSize.X)/float32(containerSize.Y) >= imgRatio {
		// bound by container height
		return f32.Pt(float32(containerSize.Y)*imgRatio, float32(containerSize.Y))
	}

	// bound by container width
	return f32.Pt(float32(containerSize.X), float32(containerSize.X)/imgRatio)
}

func (s *appState) drawImage(ops *op.Ops, containerSize image.Point) float32 {
	imageSize := s.imageSize()
	canvasSize := image.Pt(
		max(imageSize.X, s.proofImageSize.X), max(imageSize.Y, s.proofImageSize.Y))
	drawingSize := getImageDrawingSize(canvasSize, containerSize)
	clip.RRect{Rect: f32.Rect(0, 0, drawingSize.X, drawingSize.Y)}.Add(ops)

	scale := drawingSize.X / float32(canvasSize.X)

	paint.ColorOp{Color: color.RGBA{R: 232, G: 232, B: 232, A: 255}}.Add(ops)
	paint.PaintOp{Rect: f32.Rect(0, 0,
		float32(s.proofImageSize.X)*scale, float32(s.proofImageSize.Y)*scale)}.Add(ops)

	imageOp := paint.NewImageOp(s.img)
	imageOp.Add(ops)
	paint.PaintOp{Rect: f32.Rect(0, 0,
		float32(imageSize.X)*scale, float32(imageSize.Y)*scale)}.Add(ops)

	return scale
}

func (s *appState) drawMismatches(ops *op.Ops, scale float32) {
	macro := op.Record(ops)
	stack := op.Push(ops)
	paint.ColorOp{Color: color.RGBA{R: 255, A: 255}}.Add(ops)
	bounds := f32.Rect(0, 0, float32(mismatchBoxSize.X)*scale, float32(mismatchBoxSize.Y)*scale)
	clip.Border{Rect: bounds, Width: 3}.Add(ops)
	paint.PaintOp{Rect: bounds}.Add(ops)
	stack.Pop()
	box := macro.Stop()

	drawBox := func(p image.Point) {
		stack := op.Push(ops)
		op.Offset(f32.Pt(
			float32(p.X*mismatchBoxSize.X)*scale,
			float32(p.Y*mismatchBoxSize.Y)*scale)).Add(ops)
		box.Add(ops)
		stack.Pop()
	}

	// draw error boxes in mismatches
	for _, v := range s.mismatches {
		drawBox(v)
	}
}
