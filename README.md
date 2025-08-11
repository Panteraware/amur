# Amur
Go Echo Rest API CDN with WebP conversion, image resizing, video processing, and redis capabilities. Amur aims to make image and video content delivery more ease-of-use to offload work time and resource consumptions. Through images conversions to webp and thumbnails and video conversions to HLS and lower resolution files, Amur can drastically reduce load times and bandwidth use.

## Roadmap
- Add video dimensions check 
- Limit re-processing
- Cronjob for re-processing

## Functionality

### Images
- Thumbnail creation (currently half the size of the original image)
- WebP conversion
- WebP thumbnail (There is a weird bug where the webp thumbnail has a slight white haze)

### Videos
- HLS conversion
- Dimension resizing for smaller file sizes

### All capabilities
- Re-process failed optimizations
- Redis queue
- File watching, optimize on new image uploads
- Process files on interval
- Upload files via REST
- Complete structured logs (compatible with alloy, promtail, etc.)
- Healthcheck route
- Graceful shutdown
- Prometheus metrics

## ENV
These are all default values so you don't have to define them all yourself.

| Variable               | Default Value                  | Data type | Description                                                                  |
|------------------------|--------------------------------|-----------|------------------------------------------------------------------------------|
| PORT                   | 3000                           | *int      | Port that the server runs on                                                 |
| DOMAIN                 | localhost                      | string    | Domain the server runs on, useful for custom CORS configurations and cookies |
| PUBLIC_FOLDER          | /public                        | *string   | Folder with all files to be served by server                                 |
| UPLOAD_KEY             |                                | *string   | Upload key for upload authentication                                         |
| USE_REDIS              | false                          | *boolean  | Rather to use the queue processor, redis has to be used and connected        |
| CAN_CONVERT_HLS        | false                          | *boolean  | Convert videos to m3u8                                                       |
| CAN_SCALE_VIDEO        | false                          | *boolean  | Scale video to smaller and larger resolutions defined in ``VIDEO_SCALE``     |
| VIDEO_SCALE            | 720                            | *string   | Comma-separated string with common resolutions (2560,1440,1080,720, etc.)    |
| REDIS_HOST             | localhost:6379                 | *string   | Host + Port for redis                                                        |
| REDIS_PASS             |                                | *string   | Redis password for authentication (optional)                                 |
| REDIS_DB               | 1                              | *int      | Which redis database to use                                                  |
| CORS_ALLOW_ORIGINS     | localhost                      | *string   | Which hosts should be allowed through CORS                                   |
| CORS_ALLOW_HEADERS     | GET,HEAD,POST,PUT,PATCH,DELETE | *string   | Request types allowed through CORS                                           |
| CORS_ALLOW_CREDENTIALS | true                           | *boolean  | Allow cookies                                                                |
| PROM_USERNAME          | admin                          | *string   | Username to use with prometheus                                              |
| PROM_PASSWORD          |                                | *string   | Password for prometheus, if left empty then prometheus is disabled           |

``*`` Purely optional variables that aren't needed for basic functionality.

To properly work the following are **REQUIRED**:
- ``DOMAIN`` **MUST** be set to your serving domain for security header purposes.
- ``UPLOAD_KEY`` **MUST** be over 8 characters to unlock upload capabilities.
- ``PUBLIC_FOLDER`` If you are running on host, make sure to define something like ``./public/``, otherwise leave the default in the docker container.


## Credits
- [Logging](https://github.com/rs/zerolog)
- [Echo](https://echo.labstack.com/)
- [Webp](https://developers.google.com/speed/webp/download)
- [UUID](https://github.com/google/uuid)
- [ENV](https://github.com/joho/godotenv)
- [FFmpeg Wrapper](https://github.com/u2takey/ffmpeg-go)
