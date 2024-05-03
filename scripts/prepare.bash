#!/usr/bin/env bash

setupenv() {
    cp $PWD/example/.env.example $PWD/.env;
}

setupenv && printf "[SUCCESSED] default '.env' file created\n" || printf "[FAILED] fail to create default '.env' file";
