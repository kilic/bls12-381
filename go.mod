module github.com/drand/bls12-381

go 1.14

require (
	github.com/drand/bls12381rs v0.0.0-20200326192142-9d237da36846
	github.com/drand/kyber v1.0.1-0.20200110225416-8de27ed8c0e2
	github.com/filecoin-project/filecoin-ffi v0.0.0-20200326140424-59e1ffb992c1 // indirect
	github.com/kilic/bls12-381 v0.0.0-20191103193557-038659eaa189
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/sys v0.0.0-20191025090151-53bf42e6b339
	gopkg.in/yaml.v2 v2.2.4
)

//replace github.com/drand/bls12381rs => ../map12381/
