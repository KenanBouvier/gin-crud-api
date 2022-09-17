package main

import(
    "net/http"
    "github.com/gin-gonic/gin"
    "errors"
)

type book struct{
    ID string `json:"id"`
    Title string `json:"title"`
    Author string `json:"author"`
    Quantity int `json:"quantity"`
}


var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}


func getBooks(c *gin.Context){
    c.IndentedJSON(http.StatusOK,books);
}

func bookById(c *gin.Context){
    id := c.Param("id");
    book,error := getBookById(id);
    if error != nil{
        c.IndentedJSON(http.StatusNotFound,gin.H{"error":"No such book exists"})
        return
    }

    c.IndentedJSON(http.StatusOK,book);

}

func getBookById(id string) (*book,error){
    for index,book := range books{
        if book.ID == id{
            return &books[index],nil;
        }
    }
    return nil,errors.New("Book not found");
}


func createBook(c *gin.Context){
    var newBook book;

    if err := c.BindJSON(&newBook); err!=nil{
        return;
    }

    books = append(books,newBook);
    
    c.IndentedJSON(http.StatusCreated,newBook);

}

// checkin and checkout books

func checkOutBook(c *gin.Context){
    id:= c.Param("id");
    for index,item := range books{
        if item.ID == id{
            if books[index].Quantity>0{
                books[index].Quantity -=1;
                c.IndentedJSON(http.StatusAccepted,books[index]);
                return;
            }
            c.IndentedJSON(http.StatusNotFound,gin.H{"Error":"No more books to take :("})
            return;
        }
    } 
    c.IndentedJSON(http.StatusNotFound,gin.H{"Error":"No such book exists"})
}

func checkInBook(c *gin.Context){
    id := c.Param("id");

    for index, book := range books{
        if book.ID == id{
            books[index].Quantity +=1;
            c.IndentedJSON(http.StatusAccepted,books[index]);
            return;
        }
    }
    c.IndentedJSON(http.StatusNotFound,gin.H{"Error":"No such book exists"})
}

func main(){
    router := gin.Default();
    router.GET("/books",getBooks);
    router.POST("/books",createBook);
    router.GET("/books/:id",bookById);

    //checkout
    router.PATCH("/books/out/:id",checkOutBook);
    router.PATCH("/books/in/:id",checkInBook);

    router.Run("localhost:8080");
}


