// CMPT 315 (Winter 2018)
//
// Lab 03: PostgreSQL and jmoiron/sql
// Author: Nicholas Boers
package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var Connect string = "dbname=lab03 user=postgres host=localhost port=5432 sslmode=disable"

// student represents all of the data stored for a single student
type student struct {
	StudentId int    `db:"student_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

// course represents all of the data stored for a single course
type course struct {
	CourseId     int    `db:"course_id"`
	Department   string `db:"department"`
	CourseNumber int    `db:"course_number"`
	Section      string `db:"section"`
}

// enrollment represents the relationship between students and courses;
// when a student is enrolled in a course, both the student ID and
// course ID appear in such a record
type enrollment struct {
	EnrollmentId int `db:"enrollment_id"`
	StudentId    int `db:"student_id"`
	CourseId     int `db:"course_id"`
}

// isUniqueViolation returns true if the supplied error resulted from
// a unique constraint violation
func isUniqueViolation(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		return err.Code == "23505"
	}

	return false
}

// createCourse inserts a new course into the database
func createCourse(db *sqlx.DB) error {
	// check the arguments
	args := os.Args[2:]
	if len(args) != 3 {
		return fmt.Errorf("three arguments required: department, course number, and section")
	}

	// insert the data
	q := `INSERT INTO course (department, course_number, section)
                   VALUES ($1, $2, $3)`
	result, err := db.Exec(q, args[0], args[1], args[2])
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("course already exists")
		}

		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("%d course(s) created.\n", count)

	return nil
}

// printStudents prints all of the students appearing in the slice of
// students
func printCourses(courses []course) {
	fmt.Printf("ID   Department, Course Number, Section\n")
	fmt.Printf("--------------------------------\n")
	for _, course := range courses {
		fmt.Printf("%-4d %s, %d, %s\n", course.CourseId, course.Department, course.CourseNumber, course.Section)
	}
}

// listCourses displays a table of all the courses in the database
func listCourses(db *sqlx.DB) error {
	courses := []course{}
	q := `SELECT *
                FROM course`
	if err := db.Select(&courses, q); err != nil {
		return err
	}

	printCourses(courses)

	return nil
}

// createStudent inserts a new student into the database
func createStudent(db *sqlx.DB) error {
	args := os.Args[2:]
	if len(args) != 2 {
		return fmt.Errorf("two arguments required: first name, last name")
	}

	q := `INSERT INTO student (first_name, last_name)
                   VALUES ($1, $2)`
	result, err := db.Exec(q, args[0], args[1])
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("course already exists")
		}

		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("%d student(s) created.\n", count)

	return nil
}

// printStudents prints all of the students appearing in the slice of
// students
func printStudents(students []student) {
	fmt.Printf("ID   Last name, First name\n")
	fmt.Printf("--------------------------------\n")
	for _, student := range students {
		fmt.Printf("%-4d %s, %s\n", student.StudentId, student.LastName, student.FirstName)
	}
}

// listStudents displays a table of all the students in the database
func listStudents(db *sqlx.DB) error {
	students := []student{}
	q := `SELECT *
                FROM student`
	if err := db.Select(&students, q); err != nil {
		return err
	}

	printStudents(students)

	return nil
}

// enrollStudent enrolls a student in the indicated course
func enrollStudent(db *sqlx.DB) error {
	args := os.Args[2:]
	if len(args) != 2 {
		return fmt.Errorf("two arguments required: studentId, courseId")
	}

	q := `INSERT INTO enrollment (student_id, course_id)
                   VALUES ($1, $2)`
	result, err := db.Exec(q, args[0], args[1])
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("student already enrolled in course")
		}

		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("%d student(s) enrolled.\n", count)

	return nil
}

// dropStudent drops a student from the indicated course
func dropStudent(db *sqlx.DB) error {
	args := os.Args[2:]
	if len(args) != 2 {
		return fmt.Errorf("two arguments required: studentId, courseId")
	}

	q := `DELETE FROM enrollment 
                    WHERE student_id = $1
                      AND course_id = $2`
	result, err := db.Exec(q, args[0], args[1])
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("%d student(s) dropped.\n", count)

	return nil
}

// courseRoster displays the course roster for the indicated course
func courseRoster(db *sqlx.DB) error {
	args := os.Args[2:]
	if len(args) != 1 {
		return fmt.Errorf("one argument required: courseId")
	}

	students := []student{}
	q := `SELECT *
                FROM student
               WHERE student_id IN (
                     SELECT student_id
                       FROM enrollment
                      WHERE course_id = $1)`
	if err := db.Select(&students, q, args[0]); err != nil {
		return err
	}

	printStudents(students)

	return nil
}

// displayUsage displays a usage message
func displayUsage() {
	fmt.Fprintf(os.Stderr, `usage: %s command [args...]

Commands:
  help
  createCourse  department number   section
  listCourses
  createStudent firstName  lastName
  listStudents
  enrollStudent studentId  courseId
  dropStudent   studentId  courseId
  courseRoster  courseId
`, path.Base(os.Args[0]))
}

// executeCommand executes the specified command (specified on the
// command line)
func executeCommand(db *sqlx.DB) error {
	commands := map[string]func(*sqlx.DB) error{
		"help": func(db *sqlx.DB) error {
			displayUsage()
			return nil
		},
		"createCourse":  createCourse,
		"listCourses":   listCourses,
		"createStudent": createStudent,
		"listStudents":  listStudents,
		"enrollStudent": enrollStudent,
		"dropStudent":   dropStudent,
		"courseRoster":  courseRoster,
	}

	if len(os.Args) == 1 {
		return errors.New("required command missing")
	}

	// attempt to lookup and then execute the command
	if f, ok := commands[os.Args[1]]; ok {
		err := f(db)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown command (%v)", os.Args[1])
	}

	return nil
}

// connectToDB connects to an SQL database and provides a pointer to
// the connected database
func connectToDB() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", Connect)
}

func main() {
	db, err := connectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to database: %v\n", err)
		os.Exit(1)
	}

	err = executeCommand(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot execute command: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "run '%s help' for help\n", path.Base(os.Args[0]))
		os.Exit(1)
	}
}
