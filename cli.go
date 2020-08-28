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
 * @Date:   2020-08-28T13:06:22+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-28T13:09:52+10:00
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/SouthbankSoftware/proofable/pkg/colorcli"
	anchorPB "github.com/SouthbankSoftware/proofable/pkg/protos/anchor"
)

func (s *appState) parseCliOptions() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Proofable Image builds trust into your image\n\nUsage:\n\tproofable-image [flags] path/to/image.{png,jpeg,gif}\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output())
	}

	flag.StringVar(&s.imgcertPath,
		"imgcert-path", "", "specify the image certificate path. If not specified, \"path/to/image.{png,jpeg,gif}.imgcert\" will be used")

	var anchorTypeStr string

	flag.StringVar(&anchorTypeStr, "anchor-type", "HEDERA_MAINNET", "specify the anchor type. Please refer to https://github.com/SouthbankSoftware/proofable/blob/master/docs/anchor.md#anchortype for all available anchor types")

	flag.BoolVar(&s.outputDotGraph, "output-dot-graph", false, "indicate whether to output the image certificate's Graphviz Dot Graph (.dot)")

	flag.Uint64Var(&s.distanceTolerance, "distance-tolerance", 1, "specify the tolerance of the difference distance when comparing two pixel boxes")

	flag.Parse()

	if args := flag.Args(); len(args) == 1 {
		s.imagePath = args[0]
	} else {
		exitWithHelpAndErr(errors.New("the image path must be provided"))
	}

	if s.imgcertPath == "" {
		s.imgcertPath = s.imagePath + ".imgcert"
	}

	if v, ok := anchorPB.Anchor_Type_value[anchorTypeStr]; ok {
		s.anchorType = anchorPB.Anchor_Type(v)
	} else {
		exitWithHelpAndErr(fmt.Errorf("invalid anchor type: %s", anchorTypeStr))
	}
}

func exitWithHelpAndErr(err error) {
	flag.Usage()
	fmt.Fprintln(flag.CommandLine.Output(), err)
	os.Exit(2)
}

func exitWithErr(err error) {
	colorcli.Faillnf("%s", err)
	os.Exit(1)
}
