set -e
(cd cmd/cfops/ && GOOS=linux GOARCH=amd64 godep go build -ldflags "-X main.VERSION=${NEXT_VERSION}" && mkdir -p ${WERCKER_OUTPUT_DIR}/${BUILD_DIR}/linux64 && mv cfops ${WERCKER_OUTPUT_DIR}/${BUILD_DIR}/linux64)
(cd cmd/cfops/ && GOOS=darwin GOARCH=amd64 godep go build -ldflags "-X main.VERSION=${NEXT_VERSION}" && mkdir -p ${WERCKER_OUTPUT_DIR}/${BUILD_DIR}/osx && mv cfops ${WERCKER_OUTPUT_DIR}/${BUILD_DIR}/osx)
(cd cmd/cfops/ && GOOS=windows GOARCH=amd64 godep go build -ldflags "-X main.VERSION=${NEXT_VERSION}" && mkdir -p ${WERCKER_OUTPUT_DIR}/${BUILD_DIR}/win64 && mv cfops.exe ${WERCKER_OUTPUT_DIR}/${BUILD_DIR}/win64)
