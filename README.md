# Amur

Go-Gin Rest API CDN with WebP conversion, image resizing, video processing, and redis capabilities.

## Roadmap
- Add FFmpeg support
- Integrate optional use of redis

## ENV
These are all default values so you don't have to define them all yourself.
```
PORT=3000
DOMAIN=localhost
PUBLIC_FOLDER=/public/
UPLOAD_KEY=o86bo86b // Must be over 8 characters long to upload files

USE_REDIS=false
REDIS_HOST=localhost:6379
REDIS_PASS= //blank
REDIS_DB=1
```

To properly work the following are REQUIRED:
- ``DOMAIN`` To be able to get a return link from uploading to use for ShareX or saving the file URL.
- ``UPLOAD_KEY`` MUST be over 8 characters to unlock upload capabilities.
- ``PUBLIC_FOLDER`` If you are running on host, make sure to define something like ``./public/``


## Credits
- [Logging](github.com/rs/zerolog)
- [Gin](github.com/gin-gonic/gin)
- [WebP encoding](github.com/nickalie/go-webpbin)
- [MimeType](github.com/gabriel-vasile/mimetype)
- [UUID](github.com/google/uuid)
- [ENV](github.com/joho/godotenv)
