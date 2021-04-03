docker run --rm -it -v gitlab-runner-config:/etc/gitlab-runner gitlab/gitlab-runner:latest register \
  # The name of this runner
  --description=buildy-docker-privileged \
  # A list of tags that allows us to define what jobs will be executed by which runner.
  --tag-list=docker,privileged \
  # By default new runner are locked. This unlocks them without having to go into the GitLab UI.
  --locked=false \
  # Allow the runner to run build jobs that have no tags defined
  --run-untagged=true \
  # Make this a Docker runner
  --executor=docker \
  # The default image used for all build jobs (can be override on a per job basis)
  --docker-image=ubuntu \
  # Allow this Docker container to access the host it is running on in order to start new Docker
  # container on the host.
  --docker-privileged=true