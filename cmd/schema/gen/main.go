package main

import (
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

type Name struct {
	Name string `json:"name"`
	Sex  string `json:"sex,omitempty"`
}

func main() {
	jsch := makeJSONSchema[Name]()

	explore(jsch)

	person := Name{Name: "Siu Yin", Sex: "M"}
	person2 := Name{Name: "Kit Siew"}

	jschResolved, err := jsch.Resolve(nil)
	if err != nil {
		log.Fatal("resolve: ", err)
	}

	if err := jschResolved.Validate(&person); err != nil {
		log.Fatal("validate error: ", err)
	}
	fmt.Println("person validated")

	if err := jschResolved.Validate(&person2); err != nil {
		log.Fatal("validate error: ", err)
	}
	fmt.Println("person2 validated")

	dat, err := jsch.MarshalJSON()
	if err != nil {
		log.Fatal("marshal: ", err)
	}
	fmt.Printf("%s\n", dat)

	if err := jsch.UnmarshalJSON(dat); err != nil {
		log.Fatal("unmarshal: ", err)
	}
}

func makeJSONSchema[T any]() *jsonschema.Schema {
	jsch, err := jsonschema.For[T]()
	if err != nil {
		log.Fatal("for: ", err)
	}
	return jsch
}

func explore(jsch *jsonschema.Schema) {
	fmt.Printf("%v %v:\n", jsch.Type, jsch.Properties)
	for _, r := range jsch.Required {
		fmt.Printf("\t%s: %s\n", r, jsch.Properties[r].Type)
	}
}
