#!/bin/sh

cd $HOME/ドキュメント/PandorA
cd cmd/pandora && GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui"
cd ../form && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=$(which gcc) go build -ldflags "-H=windowsgui"
#mv form.exe ../pandora/.