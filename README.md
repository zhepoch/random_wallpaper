# macOS random wallpaper 

[![Build Status](https://travis-ci.com/zhepoch/random_wallpaper.svg?branch=master)](https://travis-ci.com/zhepoch/random_wallpaper)

## Features
+ runtime change wallpaper style
+ backend clean extra download wallpaper

## To use
+ Register a [Unsplash](https://unsplash.com/developers) developer account
+ Create you application, and get this access_key.
+ Download this repository,and build `go build -o /usr/local/bin .`
+ Run this command: `random_wallpaper -a your_access_key`

## To runtime change wallpaper style
> open browser: `http:://127.0.0.1:16606`, change wallpaper style

## Support Parameters
| shorthand | longhand | type | usage |
| --------- |:-------- | ---- | ----- |
| -a        | --access_key | string | Access key of unsplash. |
| -f        | --file_path | string | save download wallpaper path. (default "/tmp/random_wallpaper/"). |
| -p        | --listen_port | int | change get unsplash query key http server listen port. (default 16606). |
| -v        | --log_level | uint | debug level 0-5, 0:panic, 1:Fatal, 2:Error, 3:Warn, 4:Info 5:debug (default 4). |
| -q        | --photo_query_key | string | Limit selection to photos matching a search term. |
| -t        | --replace_time | int | Change wallpaper every few minutes. (default 5). |

## And
... Nothing ...
