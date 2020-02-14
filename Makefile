GOBUILD=go build .

all: mailheaders

.PHONY: buddy
mailheaders:
	$(GOBUILD)


.PHONY: sass
sass:
	sassc static/scss/mailheaders.scss > static/css/mailheaders.css

clean:
	rm go-mail-headers
