# GRPC Streaming Server Usando GoLang
## Descripción

Este es un ejemplo básico de como funciona Streaming Server RPC en gRPC usango GoLang.

## Como probarlo

### Servidor

Ejecuta `go run server/maing.go`, el cual será el encargado de recibir *actualizaciones de un partido* y el encargado de compartir estas actualizaciones a los clientes que se han suscrito por streaming.

### Cliente (Suscriptor)

Ejecuta `go run client/client2.go` en tantas terminales como quieras y pueden ser en tiempo diferentes, es decir, después de que otras ya hayan recibido alguna actualización.

### Cliente (Actualizador)

Ejecuta `go run/client3.go -args` para enviar una nueva actualización..

## Authors
* Santiago Yepes Tamayo <zetogk@gmail.com>