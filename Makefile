

cmd/relaymail/relaymail:
	cd cmd/relaymail && go build -gcflags="-trimpath=${PWD}"

