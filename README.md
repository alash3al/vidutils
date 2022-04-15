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
- [ ] Generating a HLS playlist from a video file.
- [ ] Supporting S3 as a valid file source/output beside the currently and only supported "local filesystem".
- [ ] Supporting FTP as a valid file source/output beside the currently and only supported "local filesystem".
- [ ] Supporting HTTP as a valid file source beside the currently and only supported "local filesystem".
- [ ] Implementing a "Distributed Queue" engine to simplify the integration with the real world problems/apps.

Downloads
=========
> go to the [releases page](https://github.com/alash3al/vidutils/releases).

Thanks To
==========
- [H.264 Video Encoding Guide](https://trac.ffmpeg.org/wiki/Encode/H.264)
- [CRF Guide (Constant Rate Factor in x264, x265 and libvpx)](https://slhck.info/video/2017/02/24/crf-guide.html)
- [Creating A Production Ready Multi Bitrate HLS VOD stream](https://docs.peer5.com/guides/production-ready-hls-vod/)
- [How can I reduce a video's size with ffmpeg?](https://unix.stackexchange.com/questions/28803/how-can-i-reduce-a-videos-size-with-ffmpeg)
- [Video Quality – What is Video Bitrate(kbps), Pixels (p) & Aspect Ratios](https://www.vdocipher.com/blog/2020/09/video-quality-bitrate-pixels/)
- [How do I reduze the size of a video to a target size?](https://unix.stackexchange.com/questions/520597/how-do-i-reduze-the-size-of-a-video-to-a-target-size?rq=1)
- [Resize/Scale/Change Resolution of a Video using FFmpeg Easily](https://ottverse.com/change-resolution-resize-scale-video-using-ffmpeg/)