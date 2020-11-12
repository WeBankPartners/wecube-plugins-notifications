current_dir=$(shell pwd)
version=$(PLUGIN_VERSION)
project_name=$(shell basename "${current_dir}")

clean:
	rm -rf wecube-plugins-notifications

build: clean
	chmod +x ./build/*.sh
	docker run --rm -v $(current_dir):/go/src/github.com/WeBankPartners/$(project_name) --name build_$(project_name) golang:1.12.5 /bin/bash /go/src/github.com/WeBankPartners/$(project_name)/build/build.sh

image: build
	docker build -t $(project_name):$(version) .

package: image
	cp build/register.xml ./
	sed -i "s~{{PLUGIN_VERSION}}~$(version)~g" ./register.xml
	docker save -o image.tar $(project_name):$(version)
	zip  wecube-plugins-notification-$(version).zip image.tar register.xml
	rm -f register.xml
	rm -rf ./*.tar
	docker rmi $(project_name):$(version)

upload: package
	$(eval container_id:=$(shell docker run -v $(current_dir):/package -itd --entrypoint=/bin/sh minio/mc))
	docker exec $(container_id) mc config host add wecubeS3 $(s3_server_url) $(s3_access_key) $(s3_secret_key) wecubeS3
	docker exec $(container_id) mc cp /package/wecube-plugins-notification-$(version).zip wecubeS3/wecube-plugin-package-bucket
	docker rm -f $(container_id)
	rm -rf wecube-plugins-notification-$(version).zip