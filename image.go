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
 * @Date:   2020-07-30T17:55:40+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-28T13:19:16+10:00
 */

package main

import (
	"bytes"
	"errors"
	"image"
	"math"

	"github.com/corona10/goimagehash"
)

type imageTrieRootMetadata struct {
	ImageSize image.Point `json:"imageSize"`
}

func min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}

	return b
}

func getBoxNumX(width int) int {
	return int(math.Ceil(float64(width) / float64(mismatchBoxSize.X)))
}

func getBoxNumY(height int) int {
	return int(math.Ceil(float64(height) / float64(mismatchBoxSize.Y)))
}

func (s *appState) imageSize() image.Point {
	return s.img.Bounds().Size()
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

func getBoxRectAt(p image.Point) image.Rectangle {
	return image.
		Rect(0, 0, mismatchBoxSize.X, mismatchBoxSize.Y).
		Add(image.Pt(p.X*mismatchBoxSize.X, p.Y*mismatchBoxSize.Y))
}

func (s *appState) getBoxImageAt(p image.Point) (si image.Image, er error) {
	boxRect := getBoxRectAt(p)

	if s.img.Bounds().Intersect(boxRect).Empty() {
		return
	}

	if img, ok := s.img.(subImager); ok {
		si = img.SubImage(boxRect)
		return
	}

	er = errors.New("failed to get subimage: `SubImage` not available")
	return
}

func imageHashToBytes(hash *goimagehash.ExtImageHash) (da []byte, er error) {
	buf := new(bytes.Buffer)

	err := hash.Dump(buf)
	if err != nil {
		er = err
		return
	}

	da = buf.Bytes()
	return
}

func bytesToImageHash(data []byte) (ha *goimagehash.ExtImageHash, er error) {
	buf := bytes.NewBuffer(data)

	return goimagehash.LoadExtImageHash(buf)
}

func (s *appState) getBoxImageHashAt(p image.Point) (ha *goimagehash.ExtImageHash, er error) {
	si, err := s.getBoxImageAt(p)
	if err != nil {
		er = err
		return
	}

	w, h := 16, 16

	if si == nil {
		imgSize := w * h
		lenOfUnit := 64

		var phash []uint64

		if imgSize%lenOfUnit == 0 {
			phash = make([]uint64, imgSize/lenOfUnit)
		} else {
			phash = make([]uint64, imgSize/lenOfUnit+1)
		}

		ha = goimagehash.NewExtImageHash(phash, goimagehash.PHash, imgSize)
		return
	}

	return goimagehash.ExtPerceptionHash(si, w, h)
}
