FROM google/golang:stable

MAINTAINER Leo Correa <lcorr005@gmail.com>

ENV APP_DIR $GOPATH/src/github.com/Tonkpils/goosse

WORKDIR $APP_DIR
ADD . $APP_DIR

RUN go install

CMD exec goosse

