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
 * @Date:   2020-07-30T00:25:21+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-31T16:32:46+10:00
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"image"
	"io"
	"os"
	"time"

	"github.com/SouthbankSoftware/proofable/pkg/api"
	"github.com/SouthbankSoftware/proofable/pkg/colorcli"
	anchorPB "github.com/SouthbankSoftware/proofable/pkg/protos/anchor"
	apiPB "github.com/SouthbankSoftware/proofable/pkg/protos/api"
	"google.golang.org/grpc/status"
)

func unpackGRPCErr(err error) error {
	if s, ok := status.FromError(err); ok {
		return errors.New(s.Message())
	}

	return err
}

func (s *appState) verifyImage(ctx context.Context, prfFile io.Reader) {
	colorcli.Printf("Verifying %s against image certificate %s...\n",
		colorcli.Green(s.imagePath), colorcli.Green(s.imgcertPath))
	var (
		verifiable = false
		triePf     *apiPB.TrieProof
		count      = 0
	)

	err := api.WithImportedTrie(ctx, s.apiClient, "", s.imgcertPath, func(id, root string) error {
		tp, err := api.GetTrieProof(ctx, s.apiClient, id, "", root)
		if err != nil {
			return err
		}

		triePf = tp

		var dotOutputPath string

		if s.outputDotGraph {
			dotOutputPath = s.imgcertPath + ".dot"
		}

		kvCH, rpCH, errCH := api.VerifyTrieProof(ctx, s.apiClient, id, tp.GetId(),
			true, dotOutputPath)

		kvCH = api.InterceptKeyValueStream(ctx, kvCH,
			api.StripCompoundKeyAnchorTriePart)

		metadata := &imageTrieRootMetadata{}

		// read metadata
		err = api.UnmarshalFromKeyValues(
			api.MetadataPrefix,
			func() (kv *apiPB.KeyValue, er error) {
				keyValue, ok := <-kvCH
				if !ok {
					er = errors.New("stream is closed")
					return
				}

				kv = keyValue
				return
			},
			metadata,
		)
		if err != nil {
			return err
		}

		s.imgcertSize = metadata.ImageSize

		for kv := range kvCH {
			count++

			p := image.Point{}

			err := json.Unmarshal(kv.Key, &p)
			if err != nil {
				return err
			}

			hash, err := s.getBoxImageHashAt(p)
			if err != nil {
				return err
			}

			// return nil
			expectedHash, err := bytesToImageHash(kv.Value)
			if err != nil {
				return err
			}

			dist, err := hash.Distance(expectedHash)
			if err != nil {
				return err
			}

			if uint64(dist) > s.distanceTolerance {
				colorcli.Faillnf("pixel box mismatched at %s with distance %s",
					colorcli.Red(p), colorcli.Red(dist))

				s.mismatches = append(s.mismatches, p)
			}
		}

		err = <-errCH
		if err != nil {
			return err
		}

		verifiable = true
		rp := <-rpCH
		if !rp.GetVerified() {
			return errors.New(rp.GetError())
		}

		return nil
	})
	if err != nil {
		if verifiable {
			colorcli.Faillnf("the image certificate at %s with a root hash of %s is falsified: %s",
				colorcli.Red(s.imgcertPath),
				colorcli.Red(triePf.GetProofRoot()),
				unpackGRPCErr(err))
			return
		}

		colorcli.Faillnf("the image certificate at %s is unverifiable: %s",
			colorcli.Red(s.imgcertPath),
			unpackGRPCErr(err))
		os.Exit(1)
	}

	colorcli.Oklnf("the image certificate at %s with %v pixel boxes and a root hash of %s is anchored to %s in block %v with transaction %s at %s, which can be viewed at %s",
		colorcli.Green(s.imgcertPath),
		colorcli.Green(count),
		colorcli.Green(triePf.GetProofRoot()),
		colorcli.Green(triePf.GetAnchorType()),
		colorcli.Green(anchorPB.GetBlockNumberString(
			triePf.GetAnchorType().String(),
			triePf.GetBlockTime(),
			triePf.GetBlockTimeNano(),
			triePf.GetBlockNumber())),
		colorcli.Green(triePf.GetTxnId()),
		colorcli.Green(time.Unix(
			int64(triePf.GetBlockTime()),
			int64(triePf.GetBlockTimeNano())).Format(time.UnixDate)),
		colorcli.Green(triePf.GetTxnUri()))

	// also add area not covered by certificate as mismatches
	imageSize := s.imageSize()
	numX, numY := getBoxNumX(imageSize.X), getBoxNumY(imageSize.Y)
	numPfX, numPfY :=
		getBoxNumX(s.imgcertSize.X), getBoxNumY(s.imgcertSize.Y)

	for i := numPfX; i < numX; i++ {
		for j := 0; j < numY; j++ {
			s.mismatches = append(s.mismatches, image.Pt(i, j))
		}
	}

	l := min(numPfX, numX)
	for j := numPfY; j < numY; j++ {
		for i := 0; i < l; i++ {
			s.mismatches = append(s.mismatches, image.Pt(i, j))
		}
	}

	if l := len(s.mismatches); l > 0 {
		colorcli.Faillnf("%s is fasified: found %s image mismatches",
			colorcli.Red(s.imagePath), colorcli.Red(l))
		return
	}

	colorcli.Passlnf("%s is verified", colorcli.Green(s.imagePath))
}
