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
 * @Date:   2020-07-28T17:12:14+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-28T13:18:30+10:00
 */

package main

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"github.com/SouthbankSoftware/proofable/pkg/api"
	"github.com/SouthbankSoftware/proofable/pkg/authcli"
	"github.com/SouthbankSoftware/proofable/pkg/colorcli"
	anchorPB "github.com/SouthbankSoftware/proofable/pkg/protos/anchor"
	apiPB "github.com/SouthbankSoftware/proofable/pkg/protos/api"
)

const (
	authEndpoint = "https://apigateway.provendb.com"
	apiHostPort  = "api.proofable.io:443"
	apiSecure    = true
)

var (
	// mismatchBoxSize is the mismatch box size in image pixels
	mismatchBoxSize = image.Pt(80, 80)
)

type appState struct {
	imagePath,
	imgcertPath string
	anchorType        anchorPB.Anchor_Type
	outputDotGraph    bool
	distanceTolerance uint64

	apiClient      apiPB.APIServiceClient
	img            image.Image
	proofImageSize image.Point
	mismatches     []image.Point
}

func main() {
	state := appState{}

	state.parseCliOptions()

	ctx := context.Background()

	// create a Proofable API client. If it is first time to authenticate in, an access token will
	// be created on your local machine
	creds, err := authcli.AuthenticateForGRPC(
		ctx,
		authEndpoint,
		apiSecure,
		"",
	)
	if err != nil {
		exitWithErr(err)
	}

	conn, cli, err := api.NewAPIClient(apiHostPort, creds)
	if err != nil {
		exitWithErr(err)
	}
	defer conn.Close()
	state.apiClient = cli

	// create a certificate for the image if it doesn't exist
	imgFile, err := os.Open(state.imagePath)
	if err != nil {
		exitWithErr(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		exitWithErr(err)
	}
	state.img = img

getCert:
	prfFile, err := os.Open(state.imgcertPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := state.createImageCertificate(ctx)
			if err != nil {
				exitWithErr(err)
			}

			colorcli.Oklnf("created at %s", colorcli.Green(state.imgcertPath))

			goto getCert
		}

		exitWithErr(err)
	}
	defer prfFile.Close()

	// verify the image against the certificate
	state.verifyImageProof(ctx, prfFile)

	// visualize the image tampering
	go func() {
		err := loop(app.NewWindow(
			app.Title(fmt.Sprintf("Proofable Image - %s", state.imagePath)),
		), &state)
		if err != nil {
			exitWithErr(err)
		}

		os.Exit(0)
	}()

	app.Main()
}

func loop(window *app.Window, state *appState) error {
	ops := new(op.Ops)

	for e := range window.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ops.Reset()
			state.drawMismatches(ops, state.drawImage(ops, e.Size))
			e.Frame(ops)
		}
	}

	return nil
}
