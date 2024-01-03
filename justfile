build:
    cp VERSION ./cmd/svart
    go build -v -o $GOPATH/bin/svart ./cmd/svart
    du -sh bin/svart
    rm ./cmd/svart/VERSION

run: build
    exec ./bin/svart

install: build
    exec cp $GOPATH/bin/svart /usr/local/bin/svart

bump-version arg:
    version {{arg}}

release: 
    git push
    git push --tags