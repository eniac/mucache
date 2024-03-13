#!/bin/bash

## This file does not need to be run (it has been run once in the beginning)
## It is only here for provenance to determine how the input data was acquired.

## Download and unzip the dataset
mkdir socfb
curl -C - -o socfb/socfb-Reed98.zip https://nrvis.com/download/data/socfb/socfb-Reed98.zip
( cd socfb; unzip socfb-Reed98.zip )
