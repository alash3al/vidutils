About
======
> a very simple, tiny and intuitive ffmpeg wrapper with a cli interface for inspecting & transforming media files supported by the original ffmpeg software.

Why
===
> I wanted to learn more about ffmpeg and dig deep into its use cases, so I tried to build a simple wrapper that is consdered a shortcut for ffmpeg.

Features
========
- [x] Inspecting a video file to get its info including (duration, size, thumbnail, ... etc).
- [x] Transforming a video file (converting, resizing and compressing).
- [x] Generating a HLS playlist from a video file.

Downloads
=========
> go to the [releases page](https://github.com/alash3al/vidutils/releases).

Usage
=====
> for now, execute the downloaded binary with the `--help` flag, example:
```shell
$ ./vidutils_linux_amd64 --help
```

Thanks To
==========
- [H.264 Video Encoding Guide](https://trac.ffmpeg.org/wiki/Encode/H.264)
- [CRF Guide (Constant Rate Factor in x264, x265 and libvpx)](https://slhck.info/video/2017/02/24/crf-guide.html)
- [Creating A Production Ready Multi Bitrate HLS VOD stream](https://docs.peer5.com/guides/production-ready-hls-vod/)
- [How can I reduce a video's size with ffmpeg?](https://unix.stackexchange.com/questions/28803/how-can-i-reduce-a-videos-size-with-ffmpeg)
- [Video Quality – What is Video Bitrate(kbps), Pixels (p) & Aspect Ratios](https://www.vdocipher.com/blog/2020/09/video-quality-bitrate-pixels/)
- [How do I reduze the size of a video to a target size?](https://unix.stackexchange.com/questions/520597/how-do-i-reduze-the-size-of-a-video-to-a-target-size?rq=1)
- [Resize/Scale/Change Resolution of a Video using FFmpeg Easily](https://ottverse.com/change-resolution-resize-scale-video-using-ffmpeg/)
- [How to consider bitrate, -maxrate and -bufsize of a video for web](https://superuser.com/questions/945413/how-to-consider-bitrate-maxrate-and-bufsize-of-a-video-for-web)
- [How much data does YouTube actually use?](https://www.androidauthority.com/how-much-data-does-youtube-use-964560/)
- [YouTube Audio Quality Bitrate](https://www.h3xed.com/web-and-internet/youtube-audio-quality-bitrate-240p-360p-480p-720p-1080p)
