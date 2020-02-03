# bad-tar

> Bad Tar, Bad Tar! They ride across the nation, the thoroughbred of sin!

This is a small tool to build an "evil" gzipped tarball to demonstrate [zip zlip](https://snyk.io/research/zip-slip-vulnerability) vulnerabilities.

Libraries and applications that are vulnerable will deposit `evil.txt` into `/tmp`.
