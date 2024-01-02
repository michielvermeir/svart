build:
    cp VERSION ./cmd/svart
    go build -v -o $GOPATH/bin/svart ./cmd/svart
    du -sh bin/svart
    rm ./cmd/svart/VERSION

run: build
    exec ./bin/svart

install: build
    exec cp svart /usr/local/bin/svart

bump-version arg:
    version {{arg}}
    git push
    git push --tags