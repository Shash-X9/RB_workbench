package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	type Report struct {
		segment_idx int `json:"segment_id"`
		Triggered   int `json:"Triggered"`
		Opened      int `json:Opened`
	}
	type Employees struct {
		Reports []Report `json:"report"`
	}

	db, err := sql.Open("mysql", "rbrideuser:ridep@s$w0rd@tcp(rb-platform-db.ccoqqognomd8.ap-southeast-1.rds.amazonaws.com:3306)/ride")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	//e.POST("/employee", func(c echo.Context) error {
	//	emp := new(Employee)
	//	if err := c.Bind(emp); err != nil {
	//		return err
	//	}
	//	//
	//	sql := "INSERT INTO employee(employee_name, employee_age, employee_salary) VALUES(?,?,?)"
	//	stmt, err := db.Prepare(sql)
	//
	//	if err != nil {
	//		fmt.Print(err.Error())
	//	}
	//	defer stmt.Close()
	//	result, err2 := stmt.Exec(emp.Name, emp.Age, emp.Salary)
	//
	//	// Exit if we get an error
	//	if err2 != nil {
	//		panic(err2)
	//	}
	//	fmt.Println(result.LastInsertId())
	//
	//	return c.JSON(http.StatusCreated, emp.Name)
	//})

	//e.DELETE("/employee/:id", func(c echo.Context) error {
	//	requested_id := c.Param("id")
	//	sql := "Delete FROM employee Where id = ?"
	//	stmt, err := db.Prepare(sql)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	result, err2 := stmt.Exec(requested_id)
	//	if err2 != nil {
	//		panic(err2)
	//	}
	//	fmt.Println(result.RowsAffected())
	//	return c.JSON(http.StatusOK, "Deleted")
	//})

	e.GET("/report/:id", func(c echo.Context) error {
		requested_id := c.Param("id")
		fmt.Println(requested_id)
		var triggered sql.NullInt32
		var segment_idx sql.NullInt32
		var opened sql.NullInt32

		err = db.QueryRow("SELECT segment_idx, opened,triggered FROM report WHERE segment_idx = ?", requested_id).Scan(&segment_idx, &opened, &triggered)

		if err != nil {
			fmt.Println(err)
		}

		//response := Report{segment_idx: segment_idx, Triggered: triggered, Opened: opened}
		response := Report{}
		if triggered.Valid {
			response.Triggered = int(triggered.Int32)
		}
		if segment_idx.Valid {
			response.segment_idx = int(segment_idx.Int32)
		}
		if opened.Valid {
			response.Opened = int(opened.Int32)
		}
		return c.JSON(http.StatusOK, response)
	})

	e.Logger.Fatal(e.Start(":3308"))
}
