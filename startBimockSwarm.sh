#!/usr/bin/env bash

sh -c "./bin/bimock -p=9000 &\
./bin/bimock -p=9001 &\
./bin/bimock -p=9002 &\
./bin/bimock -p=9003 &\
./bin/bimock -p=9004 &\
./bin/bimock -p=9005 &\
./bin/bimock -p=9006 &\
./bin/bimock -p=9007 &\
./bin/bimock -p=9009 &\
wait"