#!/bin/bash

# MedHash Tools
# Copyright (c) 2021 GHIFARI160
# MIT License

source config

echo "Cleaning artifact directory"
rm -rf $artifactDir

echo "Cleaning release directory"
rm -rf $releaseDir

echo "Done"
