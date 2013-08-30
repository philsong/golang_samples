package main;
 
import (
    "fmt"
    "database/sql"
    _"odbc/driver"
)
 
func main(){
    conn,err := sql.Open("odbc","driver={Microsoft Access Driver (*.mdb)};dbq=d:\\test.mdb");
    if(err!=nil){
        fmt.Println("Connecting Error");
        return;
    }
    defer conn.Close();
    stmt,err := conn.Prepare("select * from test");
    if(err!=nil){
        fmt.Println("Query Error");
        return;
    }
    defer stmt.Close();
    row,err := stmt.Query();
    if err!=nil {
        fmt.Println("Query Error");
        return;
    }
    defer row.Close();
    for row.Next() {
        var id int;
        var name string;
        if err := row.Scan(&id,&name);err==nil {
            fmt.Println(id,name);
        }
    }
    fmt.Printf("%s\n","finish");
    return;
}
