#!/usr/bin/env bash

sh -c "./bin/simulator -p=10000 &\
./bin/simulator -p=10001 &\
./bin/simulator -p=10002 &\
./bin/simulator -p=10003 &\
./bin/simulator -p=10004 &\
./bin/simulator -p=10005 &\
./bin/simulator -p=10006 &\
./bin/simulator -p=10007 &\
./bin/simulator -p=10009 &\
wait"