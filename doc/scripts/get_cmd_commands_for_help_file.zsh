#! /usr/env/bin bash

# Copy to clipboard the arguments that can be found in the cmd file in order to help build out the help function
grep "case" cmd/cmd.go \
  | sed -E 's/.*if[[:space:]]*case[[:space:]]"([^"]*)".*/aphrodite.PrintInfo("\1")/' 
