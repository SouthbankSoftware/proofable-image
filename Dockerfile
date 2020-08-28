FROM techknowlogick/xgo:latest

LABEL maintainer="guiguan <guan@southbanksoftware.com>"

# https://gioui.org/doc/install#linux
RUN \
    apt-get update && \
    apt-get install -y libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev \
    --no-install-recommends
