#!/bin/bash

go install go.uber.org/mock/mockgen@latest

mockgen -destination=mocks/spf13/afero_fs.go -package=mock_afero github.com/spf13/afero Fs
mockgen -destination=mocks/spf13/afero_file.go -package=mock_afero github.com/spf13/afero File
mockgen -destination=mocks/os/types.go -package=mock_os os FileInfo

