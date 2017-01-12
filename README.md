# GoMailRuCloud
Implementation of [mail.ru cloud](https://cloud.mail.ru/) API written in golang.

A trivial command line client (**tmrc**) based on the api also supplied.

# Installation & Usage
    git clone https://github.com/sergevs/GoMailRuCloud
    cd GoMailRuCloud
    go get github.com/sergevs/GoMailRuCloud
    go get github.com/gosuri/uiprogress
    go build tmrc.go

    export MAILRU_USER=<your mail.ru username>
    export MAILRU_PASSWORD=<your mail.ru password>
    ./tmrc -help

# Documentation
Most of API documented using godoc. To view the inline documentation use:

    godoc github.com/sergevs/GoMailRuCloud/Api
    
