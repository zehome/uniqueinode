PLATFORMS := darwin/amd64 linux/amd64
GO := go
GOFLAGS = CGO_ENABLED=0
OUTDIR := .

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

all: $(PLATFORMS)

$(PLATFORMS):
	$(GOFLAGS) GOOS=$(os) GOARCH=$(arch) $(GO) build -ldflags='-s -w' -o "$(OUTDIR)/uniqueinode_$(os)_$(arch)" main.go

.PHONY: $(PLATFORMS)
