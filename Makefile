GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BUILD_DIR=$(shell pwd)
APP=$(BUILD_DIR)/monitoring
GENERATOR=$(BUILD_DIR)/generator

default: clean
	$(GO_BUILD_ENV) go build -v -o $(APP) ./cmd/monitoring-system
	$(GO_BUILD_ENV) go build -v -o $(GENERATOR) ./cmd/test-generator

clean:
	rm -f $(APP)
	rm -f $(GENERATOR)
	rm -f $(BUILD_DIR)/email.data
	rm -f $(BUILD_DIR)/billing.data
	rm -f $(BUILD_DIR)/sms.data
	rm -f $(BUILD_DIR)/voice.data

