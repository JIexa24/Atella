#!/bin/bash

test -f /pkg/deb/${SERVICE}_${VERSION_RELEASE}-1_${ARCH}.deb
if [[ $? == 0 ]]; then rm -f /pkg/deb/${SERVICE}_${VERSION_RELEASE}-1_${ARCH}.deb; fi

fpm -t deb \
	-s "tar" \
	--description "${DESCRIPTION}" \
	--vendor "${VENDOR}" \
	--url "${URL}" \
	--license "${LICENSE}" \
	--name "${SERVICE}" \
	--version "${VERSION_RELEASE}" \
	--iteration "1" \
	--config-files "/etc/${SERVICE}/${SERVICE}.yml" \
	--before-install "/pkg/deb/pre-install.sh" \
	--after-install "/pkg/deb/post-install.sh" \
	--before-remove "/pkg/deb/pre-remove.sh" \
	--after-remove "/pkg/deb/post-remove.sh" \
	-p /pkg/deb \
	/pkg/tar/${SERVICE}-${VERSION_RELEASE}.tar.gz