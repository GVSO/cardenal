FROM cardenal-server-prod

RUN go get github.com/pilu/fresh

CMD ["./.docker/app_scripts/dev_initialize.sh"]