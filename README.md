# DiscordGifSmaller

## Currently this only works on windows

DiscordGifSmaller is a wrapper around gifsicle.exe

Just pass in the path to a file, and it will keep resizing until the size is less than 256KB

You can drag multiple files onto `run.bat` to have them run in a row.

The aspect ratio is not kept, so if an image is 100x50(2:1), and goes through 20 resizes, it will be 80x30(8:3)
