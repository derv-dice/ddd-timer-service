# Для Windows
SHELL := powershell.exe
.SHELLFLAGS := -Command

IMAGE_BUILD = ddd-timer-builder:latest
IMAGE_SYNC = ddd-timer-sync:latest

deploy:
	docker build -t $(IMAGE_BUILD) -f Dockerfile.build .

sync:
	docker build -t $(IMAGE_SYNC) -f Dockerfile.sync .
	docker run --name db_sync_container --rm -v "$(shell pwd):/mnt" ddd-timer-sync:latest cp db.sqlite /mnt/db.sqlite
