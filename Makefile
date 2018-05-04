TARGET = gotime
APP_NAME = Gotime.app
RELEASE_DIR = build/
APP_TEMPLATE = templates/mac/$(APP_NAME)
APP_DIR = $(RELEASE_DIR)
APP_BINARY = $(RELEASE_DIR)/$(TARGET)
APP_BINARY_DIR  = $(APP_DIR)/$(APP_NAME)/Contents/MacOS

DMG_NAME = Gotime.dmg
DMG_DIR = $(RELEASE_DIR)

vpath $(TARGET) $(RELEASE_DIR)
vpath $(APP_NAME) $(APP_DIR)
vpath $(DMG_NAME) $(APP_DIR)

run:
	go run cmd/gotime/gotime.go

build: data
	go build -o build/gotime cmd/gotime/gotime.go 

data:
	go-bindata -nocompress -nometadata -pkg timer -o bindata.go assets/...

# Build release binary
$(TARGET): clean-build build icon-mac icon-win

icon-win:
	convert assets/win/icon.png assets/win/icon.ico

icon-mac:
	mkdir app.iconset
	sips -z 16 16     assets/icon128.png --out app.iconset/icon_16x16.png
	sips -z 32 32     assets/icon128.png --out app.iconset/icon_16x16@2x.png
	sips -z 32 32     assets/icon128.png --out app.iconset/icon_32x32.png
	sips -z 64 64     assets/icon128.png --out app.iconset/icon_32x32@2x.png
	sips -z 128 128   assets/icon128.png --out app.iconset/icon_128x128.png
	sips -z 256 256   assets/icon128.png --out app.iconset/icon_128x128@2x.png
	sips -z 256 256   assets/icon128.png --out app.iconset/icon_256x256.png
	sips -z 512 512   assets/icon128.png --out app.iconset/icon_256x256@2x.png
	sips -z 512 512   assets/icon128.png --out app.iconset/icon_512x512.png
	sips -z 1024 102  assets/icon128.png --out app.iconset/icon_512x512@2x.png
	iconutil -c icns app.iconset
	rm -R app.iconset
	mv app.icns templates/mac/Gotime.app/Contents/Resources

 # Clone macOS app template and mount binary
app: | $(APP_NAME)
$(APP_NAME): $(TARGET) $(APP_TEMPLATE)
	mkdir -p $(APP_BINARY_DIR)
	cp -fRp $(APP_TEMPLATE) $(APP_DIR)
	cp -fp $(APP_BINARY) $(APP_BINARY_DIR)
	@echo "Created '$@' in '$(APP_DIR)'"

# Pack macOS app into .dmg
dmg: | $(DMG_NAME)
$(DMG_NAME): $(APP_NAME)
	@echo "Packing disk image..."
	hdiutil create $(DMG_DIR)/$(DMG_NAME) \
		-volname "Gotime" \
		-fs HFS+ \
		-srcfolder $(APP_DIR) \
		-ov -format UDZO
	@echo "Packed '$@' in '$(APP_DIR)'"

# Mount disk image
install: $(DMG_NAME)
	@open $(DMG_DIR)/$(DMG_NAME)

clean-build:
	rm -rf $(APP_DIR)

update-deps:
	rm -rf vendor
	dep ensure
	dep ensure -update
