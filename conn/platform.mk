PLATS = linux windows
.PHONY : none $(PLATS)

CGO_ENABLED := 1
GOOS := windows
GOARCH := amd64

linux : CGO_ENABLED := 0
linux : GOOS := linux

windows linux :
	$(MAKE) build CGO_ENABLED="$(CGO_ENABLED)" GOOS="$(GOOS)"

