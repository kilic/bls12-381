module github.com/drand/bls12-381

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/drand/kyber v1.0.1-0.20200110225416-8de27ed8c0e2
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/sys v0.0.0-20191025090151-53bf42e6b339
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

//replace github.com/drand/bls12381rs => ../map12381/
