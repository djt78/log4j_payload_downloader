# log4j_payload_downloader  
log4shell CVE-2021-44228

Quick and Dirty Utility to download the classfile payload from log4j ldap exploitation attempt url bassed on the project https://github.com/Adikso/minecraft-log4j-honeypot

1. Install go
2. get code from git
3. Enter code directory
4. `bash build.sh`
5. `l4jdl ldap://somedomain/Exploit`

Should get the class URL from the LDAP server, then download the class and save it to downloaded dir (with filename based of md5sum of class)

Once you download the class file you can decompile them  in seconds online I used https://jdec.app/ when testing, which seemed to work fine.
