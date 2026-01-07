//go:build uint256

package tests

//go:generate go run ../cmd -var Uint256TestABI -output=uint256.abi.go -package=tests -uint256

var Uint256TestABI = []string{
	"function transfer(address to, uint256 amount) returns (bool)",
	"function balanceOf(address account) view returns (uint256)",
	"function multiTransfer(address[] recipients, uint256[] amounts)",
	"event Transfer(address indexed from, address indexed to, uint256 value)",
}
