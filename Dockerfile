FROM golang:1.18.3 as builder

WORKDIR /app
COPY go.mod . 
COPY go.sum . 
RUN go mod download

COPY main.go .
COPY morse/* ./morse/

RUN go build -o main .

FROM ubuntu:kinetic-20220602 as runner

RUN apt-get update -y && apt-get install -y sox libssl-dev pkg-config ca-certificates openssl ffmpeg

WORKDIR /app
COPY --from=builder /app/main main

COPY ./sounds /app/sounds

ENV DASH_SOUND_FILE /app/sounds/dash.wav
ENV DOT_SOUND_FILE /app/sounds/dot.wav
ENV DOT_DASH_BREAK_SOUND_FILE /app/sounds/dot_dash_break.wav
ENV LETTER_BREAK_SOUND_FILE /app/sounds/letter_break.wav
ENV WORD_BREAK_SOUND_FILE /app/sounds/word_break.wav

ENTRYPOINT [ "/app/main" ]
