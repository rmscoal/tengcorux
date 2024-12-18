# go-redis Tracing Plugin

This package is inspired by [redisotel](https://github.com/redis/go-redis/blob/master/extra/redisotel/README.md) however re-written for tengcorux's own [tracer](https://github.com/rmscoal/tengcorux/tree/main/tracer) package.

The go-redis tracing plugin adds hooks that creates span and injects attributes when dialing redis, processing commands, and processing pipeline commands. Such that, we are able to see what query was executed to redis and so on.