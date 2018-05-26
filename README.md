# Tao 

[![Build Status](https://travis-ci.org/GaruGaru/Tao.svg?branch=master)](https://travis-ci.org/GaruGaru/Tao)
[![Go Report Card](https://goreportcard.com/badge/github.com/GaruGaru/Tao)](https://goreportcard.com/report/github.com/GaruGaru/Tao)
![license](https://img.shields.io/github/license/GaruGaru/Tao.svg)
![Docker Pulls](https://img.shields.io/docker/pulls/garugaru/tao.svg)
[![Docker image](https://images.microbadger.com/badges/image/garugaru/tao.svg)](https://microbadger.com/images/garugaru/tao "Get your own image badge on microbadger.com")

![Logo](https://github.com/GaruGaru/Tao/blob/master/res/Tao.png)

### Fast, scalable, container ready aggregator for coderdojo events

## Supported platforms 

- Eventbrite 
- Zen
- Custom api feed (wip)

# Run locally

    docker run -p 80:80 -e EVENT_BRITE_TOKEN=A_VALID_TOKEN garugaru/tao

# Powered with ♥ by

- [GoLang](https://golang.org/)
- [Gin Gonic - The fastest full-featured web framework for Golang.](https://gin-gonic.github.io/gin/) 
