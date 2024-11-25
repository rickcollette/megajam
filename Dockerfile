FROM fyneio/fyne-cross-images:linux

# Install libasound2-dev and other dependencies
RUN apt-get update && apt-get install -y libasound2-dev

# Create and expose pkg-config path for ALSA
RUN mkdir -p /usr/lib/pkgconfig
