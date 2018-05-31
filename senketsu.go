package senketsu
// senketsu - a dialect of io with a persistent memory-mapped lobby & a MUMPS-style global persistent list-keyed database

import (
//	"strconv"
	"fmt"
	"strings"
)

type object struct {
	typeName string
	name	string
	parent *object
	supportedMessages map[string]string
	mailbox chan message
}
type message struct {
	messageName string
	args []object
	response chan object
	expectsResponse bool
}

var Object object
var Future object
var Nil object

var Lobby map[string]*object

func initializeEnvironment() {
	Object := object{"Object", "Object", nil, nil, make(chan message)}
	Future := object{"Future", "Future", &Object, nil, make(chan message)}
	Lobby["Object"]=&Object
	Lobby["Future"]=&Future
	Lobby["Nil"]=&Nil

}
func nextChunk(s string) (string, string) {
	var head string
	var tail string
	s=strings.TrimSpace(s)
	nextSpace := strings.IndexAny(s, " \n\t")
	nextParen := strings.Index(s, "(")
	if(nextSpace==-1) { nextSpace=len(s) }
	if(nextParen==-1 || nextSpace<nextParen) {
		head = s[0:nextSpace-1]
		if(nextSpace+1<len(s)-1) {
			tail = s[nextSpace+1:len(s)-1]
		} else {
			tail = ""
		}
		return head, tail
	} else if(nextParen!=0) {
		head = s[0:nextParen-1]
		if(nextParen<len(s)-1) {
			tail = s[nextParen:len(s)-1]
		} else {
			tail = ""
		}
		return head, tail
	} else { // Our parenthetical expression begins the chunk
		depth := 1
		i:=0
		for i=0; i<len(s); i++ {
			if (s[i]=='(') {
				depth++
			} else if(s[i]==')') {
				depth--
				if(depth==0) {
					head = s[0:i]
					if(i+1<len(s)-1) {
						tail = s[i+1:len(s)]
					} else {
						tail = ""
					}
					return head, tail
				}
			}
		}
	}
	return head, tail // only necessary because Go is stupid
}

func (self object) clone(newName string) object {
	var ret=object{}
	if (newName=="") {
		if(self.typeName=="") { self.typeName="Object" }
		ret.typeName=self.typeName
		newName=self.typeName+fmt.Sprintf("_%x", &ret)
	} else {
		if(newName[0]>='A' && newName[0]<='Z') {
			ret.typeName=newName
		} else {
			ret.typeName=self.typeName
		}
	}
	ret.name=newName
	ret.parent=&self
	ret.run()
	return ret
}
func (self object) handle(m message) object {
	// TODO implement error stack and error condition
	var impl string
	if(self.supportedMessages != nil) {
		impl=self.supportedMessages[m.messageName]
	} else {
		impl=""
	}
	var ret object
	if(impl!="") {
		// TODO call out to parser/executor, set ret
	} else {
		if(self.parent!=nil) {
			ret=self.parent.handle(m)
		} else {
			// TODO throw error
		}
	}
	if(m.expectsResponse) {
		m.response<-ret
	}
	return ret
}
func (self object) run() {
	self.mailbox=make(chan message)
	go func() {
		select{
		case msg := <-self.mailbox:
			self.handle(msg)
		}
	}()
}
func (self object) call(other object, msg message) object{
	other.mailbox <- msg
	if(msg.expectsResponse) {
		ret := <-msg.response
		return ret
	}
	return Nil 
}
func (self object) callAsync(other object, msg message) object {
	var ret=Future.clone("")
	go func() {
		ret=self.call(other, msg)
	}()
	return ret
}
