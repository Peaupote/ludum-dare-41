BIN = a.out

SRCS = physics.go simulation.go universe.go player.go main.go

run:
	go run $(SRCS)