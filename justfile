build:
    go build -v -o svart . 
    du -sh svart

run: build
    ./svart

install: build
    cp svart /usr/local/bin/svart

bump-version arg:
    version {{arg}}
    git push
    git push --tags