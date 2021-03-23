all: demo

.PHONY: deploy
deploy: 
	flow project deploy

.PHONY: demo
demo: deploy
	go run ./demo/main.go

.PHONY: mint
mint:
	go run ./mint/main.go

.PHONY: event
event:
	go run ./event/main.go


.PHONY: script
script:
	go run ./script/main.go

.PHONY: emulator
emulator:
	flow project start-emulator -v
