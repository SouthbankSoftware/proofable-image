# proofable-image
# Copyright (C) 2020  Southbank Software Ltd.
# 
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
# 
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
# 
# 
# @Author: guiguan
# @Date:   2020-07-30T11:58:11+10:00
# @Last modified by:   guiguan
# @Last modified time: 2020-08-28T13:04:18+10:00

all: build

.PHONY: run build

run:
	go run . -output-dot-graph images/pineapple.png
build:
	go build .
build-all:
	go run src.techknowlogick.com/xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 --targets=linux/amd64,windows/amd64,darwin/amd64 --pkg . github.com/SouthbankSoftware/proofable-image
