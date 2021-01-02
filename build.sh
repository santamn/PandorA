#!/bin/sh


cd $HOME/ドキュメント/PandorA
cd cmd/pandora && go build
cd ../form && go build
mv form ../pandora/.