#!/bin/sh

cd cmd/pandora && go build
cd ../form && go build
mv form ../pandora/.