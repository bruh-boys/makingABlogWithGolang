package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type document struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Mineatura string `json:"mineatura"`
	Body      string `json:"bodyOfDocument"`
}
type publications struct {
	Size         int
	Publications []document
}

func getConnection() *sql.DB {
	dsn := "postgres://ranon:ranon@127.0.0.1:5432/publications?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
	}
	return db
}
func addPublication(e document) error {
	q := `INSERT INTO 
	publ(titulo,mineatura,body) 
	values($1,$2,$3)
	
	`
	db := getConnection()
	defer db.Close()
	stm, err := db.Prepare(q)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stm.Close()
	r, err := stm.Exec(&e.Title, &e.Mineatura, &e.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	i, _ := r.RowsAffected()
	if i != 1 {
		return errors.New("se esperaba una sola fila omg")
	}
	return nil
}
func getPublications(min, max int) (publications, error) {
	// este es el consultorio croe que se llamaba asi , ya no me acuerdo xd
	q := fmt.Sprintf("SELECT * FROM publ WHERE id <=%d AND id>=%d", max, min)
	db := getConnection()
	// aqui lo que hace es conectarse a la base de datos
	defer db.Close()
	//espera a cerrarse para evitar ciertos problemas de seguridad
	m, err := db.Query(q) // envia esto y la salida deb de ser la siguiente
	if err != nil {
		fmt.Println(err) // solo por si hay un error xd

		return publications{}, err
	}
	defer m.Close() // espera a cerrar el canal ( por razones de seguridad)

	var pubs publications
	for m.Next() {
		// repasa la informacion,
		var d document
		// cambia los valores de publication
		err := m.Scan(&d.ID, &d.Title, &d.Mineatura, &d.Body)
		if err != nil {
			// en caso de que haya un error

			log.Println("fuck", err)
			return publications{}, err
		}
		pubs.Publications = append(pubs.Publications, d)
		// los agrega a una listaa
	}

	return pubs, nil
}
func getTheSizeOfTheQuery() (int, error) {
	q := `SELECT id FROM publ
	ORDER BY id DESC
	LIMIT 1;`
	var dataSize int
	db := getConnection()
	defer db.Close()
	m, err := db.Query(q)
	if err != nil {
		return 0, err
	}
	defer m.Close()
	for m.Next() {

		err = m.Scan(&dataSize)
		if err != nil {
			return 0, err
		}

	}
	return dataSize, nil
}
