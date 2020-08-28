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
 * @Date:   2020-07-30T00:25:06+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-28T13:03:38+10:00
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"time"

	"github.com/SouthbankSoftware/proofable/pkg/api"
	"github.com/SouthbankSoftware/proofable/pkg/colorcli"
	anchorPB "github.com/SouthbankSoftware/proofable/pkg/protos/anchor"
	apiPB "github.com/SouthbankSoftware/proofable/pkg/protos/api"
	"github.com/SouthbankSoftware/proofable/pkg/strutil"
	"golang.org/x/sync/errgroup"
)

func (s *appState) createImageCertificate(ctx context.Context) error {
	colorcli.Printf("Creating image certificate for %s...\n", colorcli.Green(s.imagePath))

	return api.WithTrie(ctx, s.apiClient, func(id, root string) error {
		imageSize := s.imageSize()

		metadata, err := api.MarshalToKeyValues(api.MetadataPrefix, &imageTrieRootMetadata{
			ImageSize: imageSize,
		})
		if err != nil {
			return err
		}

		root, err = api.SetTrieKeyValues(ctx, s.apiClient, id, root, metadata)
		if err != nil {
			return err
		}

		numX, numY := getBoxNumX(imageSize.X), getBoxNumY(imageSize.Y)

		kvCH := make(chan *apiPB.KeyValue, 3)

		eg, egCtx := errgroup.WithContext(ctx)

		count := 0

		eg.Go(func() error {
			defer close(kvCH)

			for j := 0; j < numY; j++ {
				for i := 0; i < numX; i++ {
					p := image.Pt(i, j)

					key, err := json.Marshal(&p)
					if err != nil {
						return err
					}

					hash, err := s.getBoxImageHashAt(p)
					if err != nil {
						return err
					}

					val, err := imageHashToBytes(hash)
					if err != nil {
						return err
					}

					kvCH <- &apiPB.KeyValue{
						Key:   key,
						Value: val,
					}
				}
			}

			return nil
		})

		eg.Go(func() error {
			kvCH := api.InterceptKeyValueStream(egCtx, kvCH,
				func(kv *apiPB.KeyValue) *apiPB.KeyValue {
					if bytes.HasPrefix(kv.Key, strutil.Bytes(api.MetadataPrefix)) {
						colorcli.Printf("%s -> %s\n",
							colorcli.HeaderWhite(
								strutil.String(strutil.BytesWithoutNullChar(kv.Key))),
							strutil.HexOrString(kv.Value))
					} else {
						count++
					}

					return kv
				})

			rt, err := api.SetTrieKeyValues(egCtx, s.apiClient, id, root, kvCH)
			if err != nil {
				return err
			}

			root = rt
			return nil
		})

		err = eg.Wait()
		if err != nil {
			return err
		}

		triePf, err := api.CreateTrieProof(ctx, s.apiClient, id, root, s.anchorType)
		if err != nil {
			return err
		}

		tpCH, errCH := api.SubscribeTrieProof(ctx, s.apiClient, id, triePf.GetId())

		for tp := range tpCH {
			colorcli.Printf("Anchoring proof: %s\n", tp.GetStatus())
			triePf = tp

			if tp.GetStatus() == anchorPB.Batch_ERROR {
				return errors.New(tp.GetError())
			}
		}

		err = <-errCH
		if err != nil {
			return err
		}

		err = api.ExportTrie(ctx, s.apiClient, id, s.imgcertPath)
		if err != nil {
			return err
		}

		colorcli.Oklnf("the image certificate has successfully been created at %s with %v pixel boxes and a root hash of %s, which is anchored to %s in block %v with transaction %s at %s, which can be viewed at %s",
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

		return nil
	})
}
