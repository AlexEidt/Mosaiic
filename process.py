#!/usr/local/bin/python3
"""
Alex Eidt

Converts image frames into GIFs.
"""

import imageio
import os
import argparse


FRAMES = 'Frames'
GIF = 'GIFs'


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('out', help='Output GIF File name.')
    parser.add_argument('fps', help='Frames per second.')
    parser.add_argument('io', help='Zoom in/out.')

    args = parser.parse_args()

    writer = imageio.get_writer(
        os.path.join(GIF, f'{args.out}.gif'),
        fps=float(args.fps)
    )

    frames = sorted(os.listdir(FRAMES), key=lambda x: int(x.split('.')[0]))
    for frame in frames:
        writer.append_data(imageio.imread(os.path.join(FRAMES, frame)))
    
    if args.io == 'true':
        for frame in frames[1:-1][::-1]:
            writer.append_data(imageio.imread(os.path.join(FRAMES, frame)))

    writer.close()


if __name__ == '__main__':
    main()