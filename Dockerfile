from golang:1.22.2-bookworm
SHELL ["/bin/bash", "-c"]

# Install necessary packages to build things. Download CodeQL
RUN apt update && apt upgrade -y && apt install ssh python3-pip python3-setuptools python3-wheel ninja-build build-essential flex bison git cmake libsctp-dev libgnutls28-dev libgcrypt-dev libssl-dev libidn11-dev libmongoclient0 libbson-dev libyaml-dev libnghttp2-dev libmicrohttpd-dev libcurl4-gnutls-dev libnghttp2-dev libtins-dev libtalloc-dev meson sudo lsb-release libspdlog-dev openssl libcurl4-gnutls-dev libevent-dev libssl-dev libtool libxml2 libxml2-dev nettle-dev libcurl4 net-tools pkg-config libboost-thread-dev libyaml-cpp-dev nlohmann-json3-dev libconfig++-dev libcpp-jwt-dev libmongoc-dev libmongo-client-dev libmongoclient-dev -y && wget https://dev.mysql.com/get/mysql-apt-config_0.8.29-1_all.deb && echo "4" | sudo dpkg -i mysql-apt-config_0.8.29-1_all.deb && apt update && apt install mysql-server mysql-client libmysqlclient21 libmysqlclient-dev mysql-common -y && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts && mkdir codeql-home && cd codeql-home && git clone https://github.com/github/codeql.git codeql-repo && cd .. && wget https://github.com/github/codeql-action/releases/download/codeql-bundle-v2.17.0/codeql-bundle.tar.gz && tar -xf codeql-bundle.tar.gz && mkdir /dbs

ENV PATH="${PATH}:/go/codeql"
# Clone repos we want to analyze
RUN git clone https://github.com/free5gc/free5gc.git free5gc && cd free5gc && git submodule init && git submodule update --recursive --remote
COPY build_free5gc.sh build_free5gc.sh
RUN ./build_free5gc.sh


COPY build_sdcore.sh build_sdcore.sh
COPY build_all_sd_core_helper.sh build_all_sd_core_helper.sh
RUN mkdir sd-core && ./build_sdcore.sh

# Open5GS
RUN git clone https://github.com/open5gs/open5gs.git open5gs_r17 && cp open5gs_r17 open5gs_r16 -r && cd open5gs_r17 && codeql database create /dbs/open5gs_r17 --language=c-cpp --threads=0
RUN cd open5gs_r16 && git checkout 2ccd19e3f5af7cf553ec0c9ca04b17d4203d8a9f && codeql database create /dbs/open5gs_r16 --language=c-cpp --threads=0

# OAI
RUN git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-amf.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-smf.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-nrf.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-ausf.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-udm.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-nssf.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-nef.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-udr.git && git clone https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-pcf.git

# Create database for OAI here
RUN mkdir clones && cd clones && git clone https://github.com/oktal/pistache.git && cd pistache && git checkout e18ed9baeb2145af6f9ea41246cf48054ffd9907 && mkdir _build && cd _build && cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release \
        -DPISTACHE_BUILD_EXAMPLES=false \
        -DPISTACHE_BUILD_TESTS=false \
        -DPISTACHE_BUILD_DOCS=false \
        .. \
 && make -j $(nproc) && make install 

RUN cd clones && git clone https://github.com/nghttp2/nghttp2.git && cd nghttp2 && git submodule update --init && autoreconf -i && automake && ./configure --enable-asio-lib --enable-lib-only && make -j $(nproc) && make install && ldconfig 
RUN cd clones && git clone https://github.com/nghttp2/nghttp2-asio.git && cd nghttp2-asio && autoreconf -i && automake && ./configure && make -j $(nproc) && make install && ldconfig
RUN cd clones && git clone https://github.com/mongodb/mongo-cxx-driver.git && cd mongo-cxx-driver && cd build && cmake .. -DCMAKE_BUILD_TYPE=Release \
      -DCMAKE_INSTALL_PREFIX=/usr/local \
      -DBSONCXX_POLY_USE_IMPLS=ON \
#      -DBSONCXX_POLY_USE_BOOST=1 \
#      -DBSONCXX_POLY_USE_MNMLSTC=0 \
      -DBSONCXX_POLY_USE_STD_EXPERIMENTAL=0 \
      -DENABLE_TESTS=OFF \ 
&& cmake --build . -j $(nproc) && sudo cmake --build . --target install

#RUN cd /usr/include/nghttp2 && wget https://raw.githubusercontent.com/nghttp2/nghttp2-asio/main/lib/includes/nghttp2/asio_http2_server.h && wget https://raw.githubusercontent.com/nghttp2/nghttp2-asio/main/lib/includes/nghttp2/asio_http2.h && wget https://raw.githubusercontent.com/nghttp2/nghttp2-asio/main/lib/includes/nghttp2/asio_http2_client.h
COPY build_oai.sh build_oai.sh
RUN codeql database create /dbs/oai --overwrite --db-cluster --language=c-cpp --command=/go/build_oai.sh --threads=0

WORKDIR /repo

ENTRYPOINT ["/repo/docker_run_all.sh"]
