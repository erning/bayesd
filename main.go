package main

import (
	"flag"
	"fmt"

	"github.com/iris-contrib/middleware/logger"
	"github.com/jbrukh/bayesian"
	"github.com/kataras/iris"
)

var listenOn string
var dataFile string
var newClassifier bool
var classifier *bayesian.Classifier

func containsClass(class bayesian.Class) bool {
	for _, a := range classifier.Classes {
		if a == class {
			return true
		}
	}
	return false
}

func handleLearn(c *iris.Context) {
	class := bayesian.Class(c.Param("class"))
	if !containsClass(class) {
		c.Text(iris.StatusForbidden, "not such class")
		c.JSON(400, map[string]string{"status": "error"})
		return
	}

	tokens := []string{}
	c.ReadJSON(&tokens)
	classifier.Learn(tokens, class)
	classifier.WriteToFile(dataFile)
	c.JSON(iris.StatusOK, map[string]string{"status": "ok"})
}

func handleGuess(c *iris.Context) {
	tokens := []string{}
	c.ReadJSON(&tokens)

	var scores []float64
	var likely int
	var strict bool

	args := c.RequestCtx.QueryArgs()
	use := string(args.Peek("use"))
	if use == "log" {
		scores, likely, strict = classifier.LogScores(tokens)
	} else if use == "safe" {
		scores, likely, strict, _ = classifier.SafeProbScores(tokens)
	} else {
		scores, likely, strict = classifier.ProbScores(tokens)
	}

	result := make(map[string]interface{})
	result["likely"] = string(classifier.Classes[likely])
	result["strict"] = strict

	scoreMap := make(map[string]float64)
	for i, v := range scores {
		scoreMap[string(classifier.Classes[i])] = v
	}
	result["scores"] = scoreMap

	err := c.JSON(iris.StatusOK, result)
	if err != nil {
		c.JSON(iris.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
}

func main() {
	var verbose bool
	flag.StringVar(&listenOn, "listen", ":8080", "HTTPd listen on")
	flag.StringVar(&dataFile, "data", "bayesd.data", "Data file name")
	flag.BoolVar(&newClassifier, "new", false, "Create new classifier")
	flag.BoolVar(&verbose, "verbose", true, "Enable logger or not")
	flag.Parse()

	args := flag.Args()
	if newClassifier {
		if len(args) < 2 {
			fmt.Println("Error: provide at least two classes")
			return
		}
		classes := make([]bayesian.Class, len(args), len(args))
		for i, v := range args {
			classes[i] = bayesian.Class(v)
		}
		classifier = bayesian.NewClassifier(classes...)
		classifier.WriteToFile(dataFile)
		fmt.Printf("New classifier %s\n", dataFile)
		fmt.Printf("Classes: %v\n", classes)
	} else {
		var err error
		classifier, err = bayesian.NewClassifierFromFile(dataFile)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Load classifier from %s\n", dataFile)
		fmt.Printf("Classes: %v\n", classifier.Classes)
	}

	if verbose {
		fmt.Println("Verbose: YES")
		iris.Use(logger.New())
	}

	iris.Post("/learn/:class", handleLearn)
	iris.Post("/guess", handleGuess)

	iris.Listen(listenOn)
}
