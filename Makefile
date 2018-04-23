BIN = a.out

SRCS = menu.go physics.go simulation.go universe.go player.go main.go

install:
	go get github.com/Peaupote/ludum-dare-41

run:
	go run $(SRCS)