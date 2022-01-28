.PHONY: help
help:
	$(info ---------------------------------------------)
	$(info                 pwnkit-go)
	$(info ---------------------------------------------)
	$(info - build: Build the binary as `pwnkit-go`)
	$(info - run: Run the exploit)
	$(info - vm: Start the vulnerable vagrant box)
	$(info - ssh: SSH onto the vagrant box)
	$(info - scp: SCP the binary to the vulnerable vagrant box in the `/tmp` dir)

build:
	go build -o pwnkit-go

run:
	go run main.go

vm:
	vagrant destroy -f && vagrant up

ssh:
	vagrant ssh

scp: build
	vagrant scp ./pwnkit-go /tmp/pwnkit-go

