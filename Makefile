APP_NAME = megajam
APP_ID = com.example.megajam
PKG_CONFIG_PATH_LINUX = /usr/lib/x86_64-linux-gnu/pkgconfig/
FYNE_CROSS_BIN = /home/megalith/go/bin/fyne-cross
CUSTOM_IMAGE = custom-fyne-cross

.PHONY: all windows linux macos clean build-docker-image

all: build-docker-image windows linux macos

build-docker-image:
	@echo "Building custom Docker image with mounted pkg-config..."
	docker build -t $(CUSTOM_IMAGE) -f Dockerfile .

windows:
	@echo "$(FYNE_CROSS_BIN) windows -output $(APP_NAME) -app-id $(APP_ID)"
	$(FYNE_CROSS_BIN) windows -output $(APP_NAME) -app-id $(APP_ID)

linux: build-docker-image
	@echo "$(FYNE_CROSS_BIN) linux -output $(APP_NAME) -app-id $(APP_ID) -env PKG_CONFIG_PATH=/usr/lib/pkgconfig/ -image $(CUSTOM_IMAGE)"
	$(FYNE_CROSS_BIN) linux -output $(APP_NAME) -app-id $(APP_ID) \
		-env PKG_CONFIG_PATH=/usr/lib/pkgconfig/ \
		-image $(CUSTOM_IMAGE)

macos:
	@echo "$(FYNE_CROSS_BIN) darwin -output $(APP_NAME) -app-id $(APP_ID)"
	$(FYNE_CROSS_BIN) darwin -output $(APP_NAME) -app-id $(APP_ID)

clean:
	@echo "Cleaning up build artifacts and Docker images..."
	rm -rf fyne-cross
	docker rmi $(CUSTOM_IMAGE) || true
