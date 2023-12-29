build:
    go build -v -o tfvars . 
    du -sh tfvars

run: build
    ./tfvars

install: build
    cp tfvars /usr/local/bin/tfvars

    