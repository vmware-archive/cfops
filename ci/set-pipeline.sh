#!/bin/bash -eu

SSH_KEY=$(lpass show "Shared-London Services"/london-ci/git-ssh-key --notes)

lpass show Shared-PCF-Backup-and-Restore/concourse-secrets --notes > \
  secrets.yml

fly -t london set-pipeline \
  --pipeline cfops \
  --config pipeline.yml \
  --load-vars-from secrets.yml \
  --var git-private-key="${SSH_KEY}"
