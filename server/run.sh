#! /bin/sh
export DB=postgres://pg:pg_pass123@localhost:5434/task
export OAUTH_URL=http://localhost:7000
export SECRET=secret
export REDIS=redis://mytaskapi:mytaskapi123@localhost:6379/0

./main