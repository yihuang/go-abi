package examples

//go:generate go run ../cmd -var ERC20ABI -output erc20.abi.go

// ERC20ABI contains the standard ERC20 interface
var ERC20ABI = []string{
	"function name() view returns (string)",
	"function symbol() view returns (string)",
	"function decimals() view returns (uint8)",
	"function totalSupply() view returns (uint256)",
	"function balanceOf(address account) view returns (uint256)",
	"function transfer(address to, uint256 amount) returns (bool)",
	"function allowance(address owner, address spender) view returns (uint256)",
	"function approve(address spender, uint256 amount) returns (bool)",
	"function transferFrom(address from, address to, uint256 amount) returns (bool)",
	"event Transfer(address indexed from, address indexed to, uint256 value)",
	"event Approval(address indexed owner, address indexed spender, uint256 value)",
}

//go:generate go run ../cmd -var SimpleABI -output simple.abi.go

// SimpleABI contains a single function definition
var SimpleABI = "function send(address to, uint256 amount)"
