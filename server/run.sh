#! /bin/sh

export DB=postgres://pg:pg_pass123@localhost:5434/task
export REDIS=redis://mytaskapi:mytaskapi123@localhost:6379/0
export OAUTH_URL=http://localhost:7000
export JWTSECRET=secret@jwt123
export OAUTH_CLIENT_URL=http://localhost:7003
export OAUTH_URL=http://localhost:7000

./main