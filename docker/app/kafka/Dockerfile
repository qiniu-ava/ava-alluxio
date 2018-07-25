FROM ubuntu:16.04 
ENV KAFKA_USER=kafka \
KAFKA_DATA_DIR=/var/lib/kafka/data \
JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64 \
KAFKA_HOME=/opt/kafka \
PATH=$PATH:/opt/kafka/bin

ARG KAFKA_VERSION=0.10.2.1
ARG KAFKA_DIST=kafka_2.11
RUN set -x \
    && apt-get update \
    && apt-get install -y openjdk-8-jre-headless wget \
    && wget -q "http://mirrors.shu.edu.cn/apache/kafka/${KAFKA_VERSION}/${KAFKA_DIST}-${KAFKA_VERSION}.tgz" \
    && wget -q "http://www.apache.org/dist/kafka/${KAFKA_VERSION}/${KAFKA_DIST}-${KAFKA_VERSION}.tgz.asc" \
    && wget -q "http://kafka.apache.org/KEYS" \
    && export GNUPGHOME="$(mktemp -d)" \
    && gpg --import KEYS \
    && gpg --batch --verify "${KAFKA_DIST}-${KAFKA_VERSION}.tgz.asc" "${KAFKA_DIST}-${KAFKA_VERSION}.tgz" \
    && tar -xzf "${KAFKA_DIST}-${KAFKA_VERSION}.tgz" -C /opt \
    && rm -r "$GNUPGHOME" "${KAFKA_DIST}-${KAFKA_VERSION}.tgz" "${KAFKA_DIST}-${KAFKA_VERSION}.tgz.asc" \
    && apt-get autoremove -y wget \
    && rm -rf /var/lib/apt/lists/*

COPY log4j.properties /opt/${KAFKA_DIST}-${KAFKA_VERSION}/config/

RUN set -x \
    && ln -s /opt/${KAFKA_DIST}-${KAFKA_VERSION} $KAFKA_HOME \
    && useradd $KAFKA_USER \
    && [ `id -u $KAFKA_USER` -eq 1000 ] \
    && [ `id -g $KAFKA_USER` -eq 1000 ] \
    && mkdir -p $KAFKA_DATA_DIR \
    && chown -R "$KAFKA_USER:$KAFKA_USER"  /opt/${KAFKA_DIST}-${KAFKA_VERSION} \
    && chown -R "$KAFKA_USER:$KAFKA_USER"  $KAFKA_DATA_DIR
 
