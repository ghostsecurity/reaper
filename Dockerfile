# Build layer
FROM golang:latest AS build

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

ENV GOOS=linux

RUN go build -ldflags="-s -w" -o reaper ./cmd/reaper

# Run layer
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Add architecture detection
ARG TARGETARCH
ARG TARGETVARIANT

# Update and install certificates and basic tools
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    wget \
    gnupg \
    git \
    python3 \
    python3-pip \
    unzip \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

# Install Chrome dependencies first
RUN apt-get update && apt-get install -y \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libatspi2.0-0 \
    libcairo2 \
    libcups2 \
    libcurl4 \
    libdbus-1-3 \
    libdrm2 \
    libexpat1 \
    libgbm1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libudev1 \
    libvulkan1 \
    libx11-6 \
    libxcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxkbcommon0 \
    libxrandr2 \
    && rm -rf /var/lib/apt/lists/*

# Install Chromium from PPA
RUN apt-get update && \
    add-apt-repository ppa:saiarcot895/chromium-beta && \
    apt-get update && \
    apt-get install -y chromium-browser && \
    rm -rf /var/lib/apt/lists/*

# Install X11, VNC, and window manager
RUN apt-get update && apt-get install -y \
    xvfb \
    x11vnc \
    xterm \
    fluxbox \
    dbus-x11 \
    libx11-dev \
    libxcomposite-dev \
    libxdamage-dev \
    libxext-dev \
    libxfixes-dev \
    libxrandr-dev \
    && rm -rf /var/lib/apt/lists/*

# Install noVNC
RUN git clone --depth 1 https://github.com/novnc/noVNC.git /opt/novnc \
    && git clone --depth 1 https://github.com/novnc/websockify /opt/novnc/utils/websockify

RUN useradd -m -d /app -s /bin/bash app
WORKDIR /app

# Create start script with verbose logging and window manager
RUN echo '#!/bin/bash\n\
set -x\n\
echo "Starting Xvfb..."\n\
Xvfb :99 -screen 0 1024x768x16 -ac &\n\
sleep 2\n\
export DISPLAY=:99\n\
echo "Starting Fluxbox..."\n\
fluxbox &\n\
sleep 2\n\
echo "Starting x11vnc..."\n\
x11vnc -display :99 -forever -shared -verbose &\n\
sleep 2\n\
echo "Starting noVNC..."\n\
/opt/novnc/utils/novnc_proxy --vnc localhost:5900 --listen 6080 &\n\
sleep 5\n\
echo "Starting Chrome..."\n\
chromium-browser --proxy-server="127.0.0.1:8080" \
        --ignore-certificate-errors \
        --user-data-dir="$HOME/.chrome-reaper" \
        --no-sandbox \
        --disable-gpu \
        --disable-dev-shm-usage \
        --disable-software-rasterizer \
        "https://ghostsecurity.com" &\n\
sleep 2\n\
echo "Starting Reaper..."\n\
./reaper\n\
' > /app/start.sh \
    && chmod +x /app/start.sh

# Copy the binary and required files from build stage
COPY --from=build /app/reaper .
COPY --from=build /app/cmd/reaper/frontend ./frontend
COPY --from=build /app/cmd/reaper/dist ./dist

RUN chown -R app /app
USER app

ENV HOME=/app \
    DISPLAY=:99 \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    LC_ALL=C.UTF-8

CMD ["./start.sh"]
