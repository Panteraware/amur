# Amur
Go-Gin Rest API CDN with WebP conversion, image resizing, video processing, and redis capabilities. Amur aims to make image and video content delivery more ease-of-use to offload work time and resource consumptions. Through images conversions to webp and thumbnails and video conversions to HLS and lower resolution files, Amur can drastically reduce load times and bandwidth use.

## Roadmap
- Integrate optional use of redis (might be cancelled)
- Add video dimensions check 
- Limit re-processing
- Processing queue
- Cronjob for re-processing

## Functionality
All processing is done off the main-thread to ensure faster serving times. Currently, processing is always redone on every data fetch.

### Images
- Thumbnail creation (currently half the size of the original image)
- WebP conversion
- WebP thumbnail (There is a weird bug where the webp thumbnail has a slight white haze)

### Videos
- HLS conversion
- Dimension resizing for smaller file sizes

## ENV
These are all default values so you don't have to define them all yourself.
```
// General Variables
PORT=3000
DOMAIN=localhost
PUBLIC_FOLDER=/public/
UPLOAD_KEY=o86bo86b // Must be over 8 characters long to upload files

// Redis Variables
USE_REDIS=false
REDIS_HOST=localhost:6379
REDIS_PASS= //blank
REDIS_DB=1

// Video Processing Variables
VIDEO_SCALE=1080,720 // Default is: 720
CONVERT_HLS=false // Set true to convert videos to .m3u8 format
```

To properly work the following are **REQUIRED**:
- ``DOMAIN`` **MUST** be set to your serving domain for security header purposes.
- ``UPLOAD_KEY`` **MUST** be over 8 characters to unlock upload capabilities.
- ``PUBLIC_FOLDER`` If you are running on host, make sure to define something like ``./public/``.


## Credits
- [Logging](github.com/rs/zerolog)
- [Gin](github.com/gin-gonic/gin)
- [WebP Encoding](github.com/nickalie/go-webpbin)
- [MimeType](github.com/gabriel-vasile/mimetype)
- [UUID](github.com/google/uuid)
- [ENV](github.com/joho/godotenv)
- [FFmpeg Wrapper](https://github.com/u2takey/ffmpeg-go)
