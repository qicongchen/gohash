include $(GOROOT)/src/Make.inc

TARG=gohash.googlecode.com/hg/hashlessset
GOFILES=\
	hashlessset.go\

GC+= -I ../_obj
LD+= -L ../_obj

../_obj/$(TARG).a: _obj/$(TARG).a
	mkdir -p $(dir $@); cp -f $< $@
package: ../_obj/$(TARG).a

include $(GOROOT)/src/Make.pkg
