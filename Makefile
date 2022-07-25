PARAM = CGO_ENABLED=0
BUILD = go build
FUNDUS = cmd/fundus-service/main.go
FUNDUS_OUT = fundus-service
MEMBER = cmd/member-service/main.go
MEMBER_OUT = member-service
ITEM = cmd/item-service/main.go
ITEM_OUT = item-service
ENTRY = cmd/entry-service/main.go
ENTRY_OUT = entry-service
ARCHIVE = bkc-fundus-management.zip

fundus: Makefile $(FUNDUS)
	$(PARAM) $(BUILD) $(FUNDUS)
	mv main $(FUNDUS_OUT)

member: Makefile $(MEMBER)
	$(PARAM) $(BUILD) $(MEMBER)
	mv main $(MEMBER_OUT)

item: Makefile $(ITEM)
	$(PARAM) $(BUILD) $(ITEM)
	mv main $(ITEM_OUT)

entry: Makefile $(ENTRY)
	$(PARAM) $(BUILD) $(ENTRY)
	mv main $(ENTRY_OUT)

bundle: Makefile all
	mkdir bin
	mv $(MEMBER_OUT) $(ITEM_OUT) $(ENTRY_OUT) $(FUNDUS_OUT) bin/
	cp -r pkg bin
	cp -r scripts bin
	rm -frv bin/pkg/constants bin/pkg/models
	cd bin/
	zip -r $(ARCHIVE) bin/

clean:
	rm -rvf $(MEMBER_OUT) $(ITEM_OUT) $(ENTRY_OUT) $(FUNDUS_OUT) bin/ $(ARCHIVE)

all: Makefile fundus member item entry