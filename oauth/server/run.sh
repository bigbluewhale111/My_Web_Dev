#! /bin/sh

export OAUTH_DB=postgres://pg:pg_pass123@localhost:5433/oauth
export REDIS=redis://myoauth:myoauth123@localhost:6379/1
export SECRET=secret
export CLIENT_ID=12345abcd
export CLIENT_SECRET=1a2b3c4d5e6f7g8h9i0j


./main