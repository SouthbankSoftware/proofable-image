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
# @Last modified time: 2020-08-28T18:15:19+10:00

NAME := proofable-image

all: build

.PHONY: run build clean build-all build-image archive

run:
	go run . -output-dot-graph images/pineapple.png
build:
	go build .
clean:
	rm -f $(NAME)*
	docker rmi -f $(NAME) 2> /dev/null
build-all: build-image
	go run src.techknowlogick.com/xgo -image=$(NAME) --targets=linux/amd64,windows/amd64,darwin/amd64 github.com/SouthbankSoftware/proofable-image
build-image:
ifeq ($(shell docker image inspect $(NAME) &> /dev/null; echo $$?),1)
	docker build -t $(NAME) - < Dockerfile
endif
archive:
	gtar -czvf $(NAME)_darwin_amd64.tar.gz \
		--transform "flags=r;s|$(NAME)-darwin-10.6-amd64|$(NAME)|" \
		--owner=root $(NAME)-darwin-10.6-amd64
	gtar -czvf $(NAME)_linux_amd64.tar.gz \
	--transform "flags=r;s|$(NAME)-linux-amd64|$(NAME)|" \
	--owner=root $(NAME)-linux-amd64
	cp $(NAME)-windows-4.0-amd64.exe $(NAME).exe
	zip -r $(NAME)_windows_amd64.zip $(NAME).exe
	rm $(NAME).exe
