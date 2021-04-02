docker run --rm -it -v gitlab-runner-config:/etc/gitlab-runner gitlab/gitlab-runner:latest register \
  --description=buildy-docker-privileged \
  --tag-list=docker,privileged \
  --locked=false \
  --run-untagged=true \
  --executor=docker \
  --docker-image=ubuntu \
  --docker-privileged=true