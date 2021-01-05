#!/bin/sh


cd $HOME/Documents/PandorA
cd cmd/pandora && go build
cd ../form && go build
mv form ../pandora/.