package examples

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/sirupsen/logrus"
	"testing"
)

const (
	CallFactFuncDRL = `
rule RaiseHandToClap "Raise the hand to be able to clap" {
	when 
		Clapper.ClapCount < 10 &&
		Clapper.HandsUp == false
	then
		Clapper.HandsUp = true;
}

rule CheckIfCanClap "If hands is up we can clap" {
	when
		Clapper.HandsUp &&
		Clapper.CanClap == false
	then
		Clapper.CanClap = true;
}

rule LetsClap "If hands are up and can clap, lets clap" {
	when
		Clapper.HandsUp && Clapper.CanClap
	then
		Clapper.Clap();
		Changed("Clapper.CanClap");
		Changed("Clapper.HandsUp");
		Changed("Clapper.ClapCount");
		Log("Clapped " + Clapper.ClapCount + " times");
}
`
)

type Clapper struct {
	CanClap   bool
	HandsUp   bool
	ClapCount int64
}

func (c *Clapper) Clap() {
	fmt.Println("CLAP !!!")
	c.ClapCount++
	c.HandsUp = false
	c.CanClap = false
}

func TestCallingFactFunction(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	c := &Clapper{
		CanClap:   false,
		HandsUp:   false,
		ClapCount: 0,
	}

	dataContext := ast.NewDataContext()
	err := dataContext.Add("Clapper", c)
	if err != nil {
		panic(err)
	}

	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	err = ruleBuilder.BuildRuleFromResource("CallingFactFunction", "0.1.1", pkg.NewBytesResource([]byte(CallFactFuncDRL)))
	knowledgeBase := lib.NewKnowledgeBaseInstance("CallingFactFunction", "0.1.1")
	if err != nil {
		panic(err)
	} else {
		eng1 := &engine.GruleEngine{MaxCycle: 500}
		err := eng1.Execute(dataContext, knowledgeBase)
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		if c.ClapCount != 10 {
			t.Logf("Expected 10 clap, but %d", c.ClapCount)
			t.Fail()
		}
	}
}
