.PHONY: default clean 
	request-remote-public-list request-remote-public-info request-remote-public-download request-remote-public-delete
	request-remote-personal-list request-remote-personal-info request-remote-personal-download request-remote-personal-delete
	request-local-public-list request-local-public-info request-local-public-download request-local-public-upload request-local-public-delete
	request-local-personal-list request-local-personal-info request-local-personal-download request-local-personal-upload request-local-personal-delete
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean

request-remote-public-list:
	./_request/remote_public_list.sh
request-remote-public-info:
	./_request/remote_public_info.sh
request-remote-public-download:
	./_request/remote_public_download.sh
request-remote-public-delete:
	./_request/remote_public_delete.sh

request-remote-personal-list:
	./_request/remote_personal_list.sh
request-remote-personal-info:
	./_request/remote_personal_info.sh
request-remote-personal-download:
	./_request/remote_personal_download.sh
request-remote-personal-delete:
	./_request/remote_personal_delete.sh

request-local-public-list:
	./_request/local_public_list.sh
request-local-public-info:
	./_request/local_public_info.sh
request-local-public-download:
	./_request/local_public_download.sh
request-local-public-upload:
	./_request/local_public_upload.sh
request-local-public-delete:
	./_request/local_public_delete.sh

request-local-personal-list:
	./_request/local_personal_list.sh
request-local-personal-info:
	./_request/local_personal_info.sh
request-local-personal-download:
	./_request/local_personal_download.sh
request-local-personal-upload:
	./_request/local_personal_upload.sh
request-local-personal-delete:
	./_request/local_personal_delete.sh
