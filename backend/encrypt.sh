#!/bin/bash

INPUT=$1
FILENAME=$2
ID=$3
KEY=$4
lcpencrypt -input $INPUT -filename $FILENAME -output uploads/ -contentid $ID -contentkey $KEY