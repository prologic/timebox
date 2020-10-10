module banking

go 1.13

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/gorilla/mux v1.8.0
	github.com/kode4food/timebox v0.0.0
	github.com/stretchr/testify v1.6.1
	github.com/vektah/gqlparser/v2 v2.1.0
)

replace banking => ./

replace github.com/kode4food/timebox => ../..
