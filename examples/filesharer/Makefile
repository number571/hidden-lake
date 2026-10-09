PERSONAL ?= false

.PHONY: default clean _clean_files
	request-test-download-public-file request-test-info-public-file request-test-list-public
	request-test-download-personal-file request-test-info-personal-file request-test-list-personal
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean

request-test-download-file: clean-files
# Bob upload file to (PUBLIC || PERSONAL/Alice) dir
# Alice download file from Bob (PUBLIC || PERSONAL/Alice) dir
	@if [ "$(PERSONAL)" = "false" ]; then 		\
		echo "download public..." && echo;		\
		./_request/local_public_upload.sh; 		\
		./_request/remote_public_download.sh; 	\
	else 										\
		echo "download personal..." && echo;	\
		./_request/local_personal_upload.sh; 	\
		./_request/remote_personal_download.sh; \
	fi

request-test-download-proc: clean-files
# Bob upload file to (PUBLIC || PERSONAL/Alice) dir
# Alice download file from Bob (PUBLIC || PERSONAL/Alice) dir
# Bob delete file from (PUBLIC || PERSONAL/Alice) dir
# Alice delete file from local (PUBLIC/Bob || PERSONAL/Alice) dir
	@if [ "$(PERSONAL)" = "false" ]; then 							\
		echo "download/proc public..." && echo;						\
		./_request/local_public_upload.sh; 							\
		./_request/remote_public_download.sh & 						\
		for i in {1..100}; do 										\
			sleep 1 && echo $$i;									\
			./_request/remote_public_proc.sh || exit 1;				\
		done;														\
	else 															\
		echo "download/proc personal..." && echo;					\
		./_request/local_personal_upload.sh; 						\
		./_request/remote_personal_download.sh & 					\
		for i in {1..100}; do 										\
			sleep 1 && echo $$i;									\
			./_request/remote_personal_proc.sh || exit 1;			\
		done;														\
	fi

request-test-download-proc-cancel: clean-files
# Bob upload file to (PUBLIC || PERSONAL/Alice) dir
# Alice download file from Bob (PUBLIC || PERSONAL/Alice) dir
# Bob delete file from (PUBLIC || PERSONAL/Alice) dir
# Alice delete file from local (PUBLIC/Bob || PERSONAL/Alice) dir
	@if [ "$(PERSONAL)" = "false" ]; then 							\
		echo "download/proc public..." && echo;						\
		./_request/local_public_upload.sh; 							\
		./_request/remote_public_download.sh & 						\
		for i in {1..25}; do 										\
			sleep 1 && echo $$i;									\
			./_request/remote_public_proc.sh || exit 1;				\
		done;														\
		./_request/remote_public_proc_cancel.sh;					\
	else 															\
		echo "download/proc personal..." && echo;					\
		./_request/local_personal_upload.sh; 						\
		./_request/remote_personal_download.sh & 					\
		for i in {1..25}; do 										\
			sleep 1 && echo $$i;									\
			./_request/remote_personal_proc.sh || exit 1;			\
		done;														\
		./_request/remote_personal_proc_cancel.sh;					\
	fi

request-test-info-public-file: clean-files
# Bob upload file to (PUBLIC || PERSONAL/Alice) dir
# Alice get file info from Bob (PUBLIC || PERSONAL/Alice) dir
	@if [ "$(PERSONAL)" = "false" ]; then 		\
		echo "get info public..." && echo;		\
		./_request/local_public_upload.sh; 		\
		./_request/remote_public_info.sh; 		\
	else 										\
		echo "get info personal..." && echo;	\
		./_request/local_personal_upload.sh; 	\
		./_request/remote_personal_info.sh; 	\
	fi

request-test-list-public: clean-files
# Bob upload files to (PUBLIC || PERSONAL/Alice) dir
# Alice get list of files from Bob (PUBLIC || PERSONAL/Alice) dir
# Bob delete files from (PUBLIC || PERSONAL/Alice) dir
	@if [ "$(PERSONAL)" = "false" ]; then 		\
		echo "get list public..." && echo;		\
		./_request/local_public_upload.sh;		\
		./_request/remote_public_list.sh;		\
	else 										\
		echo "get list personal..." && echo;	\
		./_request/local_personal_upload.sh;	\
		./_request/remote_personal_list.sh;		\
	fi

clean-files:
	./_request/remote_public_proc_cancel.sh		2>/dev/null || true
	./_request/remote_personal_proc_cancel.sh 	2>/dev/null || true
	./_request/local_public_delete.sh			2>/dev/null || true
	./_request/local_personal_delete.sh			2>/dev/null || true
	./_request/remote_public_delete.sh			2>/dev/null || true
	./_request/remote_personal_delete.sh		2>/dev/null || true
