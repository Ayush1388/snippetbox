Remember: The internal directory is being used to hold ancillary non-application-
specific code, which could potentially be reused. A database model which could be
used by other applications in the future (like a command line interface application) fits
the bill here.
Let’s open the internal/models/snippets.go file and add a new Snippet struct to
represent the data for an individual snippet, along with a SnippetModel type with methods
on it to access and manipulate the snippets in our database. Like so:
File: internal/models/snippets.go
package models
import (
"database/sql"
"time"
)
// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
ID int
Title string
Content string
Created time.Time
Expires time.Time
}
// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
DB *sql.DB
}
// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
return 0, nil
}
// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
return Snippet{}, nil
}
// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
return nil, nil
}