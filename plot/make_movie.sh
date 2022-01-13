#!/bin/sh

ffmpeg -framerate 1/5 -i plots/%03d_x50x50x.png -c:v libx264 -r 30 -pix_fmt yuv420p out.mp4
