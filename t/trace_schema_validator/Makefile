RM = rm -f
PROGRAMS = trace_schema_validator
GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get

.PHONY: fetch_deps clean

all: $(PROGRAMS)

trace_schema_validator: fetch_deps trace_schema_validator.go
	$(GOBUILD) -o trace_schema_validator

fetch_deps:
	$(GOGET) github.com/xeipuuv/gojsonschema

clean:
	$(RM) $(PROGRAMS)
