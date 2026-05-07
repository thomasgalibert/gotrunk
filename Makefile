BINARY := gotrunk
DIST := dist
LDFLAGS := -s -w

.PHONY: all build build-all clean run mac-intel mac-apple windows

all: build

build:
	go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY) .

build-all: mac-intel mac-apple windows

mac-intel:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-mac-intel .

mac-apple:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-mac-apple .

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY).exe .

# Renomme le binaire Mac en .command pour qu'un double-clic dans le Finder
# l'ouvre automatiquement dans Terminal.
package-mac: mac-apple mac-intel
	cp $(DIST)/$(BINARY)-mac-apple $(DIST)/BerryerSetup-mac-apple.command
	cp $(DIST)/$(BINARY)-mac-intel $(DIST)/BerryerSetup-mac-intel.command
	chmod +x $(DIST)/BerryerSetup-mac-apple.command $(DIST)/BerryerSetup-mac-intel.command

package-windows: windows
	cp $(DIST)/$(BINARY).exe $(DIST)/BerryerSetup.exe

run: build
	./$(DIST)/$(BINARY)

clean:
	rm -rf $(DIST)
