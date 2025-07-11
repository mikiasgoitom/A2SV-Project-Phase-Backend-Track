/*
Fundamentals of Go Tasks
Task: Student Grade Calculator
Create a Go console application that allows students to calculate their average grade based on different subjects. The application should prompt the student to enter their name and the number of subjects they have taken. For each subject, the student should enter the subject name and the grade obtained (numeric value). After entering all subjects and grades, the application should display the student's name, individual subject grades, and the calculated average grade.

Requirements:
Use variables and data types to store student data.
Use conditional statements to validate input (e.g., ensure grade values are within a valid range).
Implement loops to handle multiple subjects and grades.
Utilize collections (e.g., List, Dictionary) to store subject names and corresponding grades.
Define a method to calculate the average grade based on the entered grades.
Use string interpolation to display the results in a user-friendly format.
Write test for your code [Optional]
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Student struct{
	name string
	numSubjects int8
	subject map[string]float32
}

func main()  {
	reader := bufio.NewReader(os.Stdin)
	var s1 Student
	s1.subject = make(map[string]float32)
	
	fmt.Println("GRADE CALCULATOR")
	fmt.Println("=================")
	fmt.Print("Enter your name: ")
	nameInput, _ := reader.ReadString('\n')
	s1.name = strings.TrimSpace(nameInput)
	
	fmt.Print("Enter the number of subject: ")
	numInput, _ := reader.ReadString('\n')
	num, err := strconv.Atoi(strings.TrimSpace(numInput))
	if err != nil {
		panic("Invalid input! expected an integer")
	}
	s1.numSubjects = int8(num)

	for i := 0; i < int(s1.numSubjects); i++ {
		fmt.Println("")
		var subjectName string
		fmt.Print("Enter the subject's name: ")
		subjectRaw, _ := reader.ReadString('\n')
		subjectName = strings.TrimSpace(subjectRaw)

		fmt.Printf("Enter your grade on %s: ", subjectName)
		gradeRaw, _ := reader.ReadString('\n')
		grade, err := strconv.ParseFloat(strings.TrimSpace(gradeRaw), 32)
		
		if err != nil {
			panic("Invalid input! expected a number")
		}

		if grade > 100. {
			panic("Grade can't be above a 100")
		} else if grade < 0. {
			panic("Grade can't be less than 0")
		}
		s1.subject[subjectName] = float32(grade)
	}
	fmt.Println("")
	fmt.Printf("Student name: %s \n", s1.name)
	var totalGrade float32
	for sub, grd := range s1.subject{
		fmt.Printf("Subject name: %s \t grade: %.2f \n", sub, grd)
		totalGrade += grd
	}
	fmt.Println("")	
	fmt.Printf("Average grade: %.2f", totalGrade / float32(s1.numSubjects))
}