// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

<<<<<<< HEAD
import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
=======
import { ERC721 } from "@rari-capital/solmate/src/tokens/ERC721.sol";
>>>>>>> v0.5.23

contract TestERC721 is ERC721 {
    constructor() ERC721("TEST", "TST") {}

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }
<<<<<<< HEAD
=======

    function tokenURI(uint256) public pure virtual override returns (string memory) {}
>>>>>>> v0.5.23
}
