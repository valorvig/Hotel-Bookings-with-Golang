#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -dbname=bookings -dbuser=postgres -dbpass=1234 -cache=false -production=false