#!/bin/bash

protoc internal/pkg/greetpb/greet.proto --go_out=plugins=grpc:.
protoc internal/pkg/calculatorpb/calculator.proto --go_out=plugins=grpc:.