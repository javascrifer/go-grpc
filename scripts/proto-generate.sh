#!/bin/bash

protoc internal/pkg/greetpb/greet.proto --go_out=plugins=grpc:.