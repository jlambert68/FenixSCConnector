#!/bin/bash

cd resources
fyne bundle Fenix_83_green_32x32.png > bundledIcons.go
fyne bundle -append Fenix_83_red_32x32.png >> bundledIcons.go
