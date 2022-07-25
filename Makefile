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
DIRECTORY = bkc-fundus-management
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
	mkdir $(DIRECTORY)
	mv $(MEMBER_OUT) $(ITEM_OUT) $(ENTRY_OUT) $(FUNDUS_OUT) $(DIRECTORY)
	cp -r pkg $(DIRECTORY)
	cp -r scripts $(DIRECTORY)
	rm -frv $(DIRECTORY)/pkg/constants $(DIRECTORY)/pkg/models
	zip -r $(ARCHIVE) $(DIRECTORY)

clean:
	rm -rvf $(MEMBER_OUT) $(ITEM_OUT) $(ENTRY_OUT) $(FUNDUS_OUT) $(DIRECTORY) $(ARCHIVE)

all: Makefile fundus member item entry