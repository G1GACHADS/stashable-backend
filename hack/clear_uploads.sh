#!/bin/bash

set -ex

ls public/uploads | grep -xv ".gitkeep" | parallel rm