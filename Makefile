BIN = a.out

SRCS = menu.go physics.go simulation.go universe.go player.go main.go

run:
	go run $(SRCS)