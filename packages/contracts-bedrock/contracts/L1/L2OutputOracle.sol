<<<<<<< HEAD
//SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title L2OutputOracle
 * @notice The L2 state is committed to in this contract
 * The payable keyword is used on appendL2Output to save gas on the msg.value check.
 * This contract should be deployed behind an upgradable proxy
 */
// slither-disable-next-line locked-ether
contract L2OutputOracle is Ownable {
    /**********
     * Events *
     **********/

    /// @notice Emitted when an output is appended.
    event l2OutputAppended(
        bytes32 indexed _l2Output,
        uint256 indexed _l1Timestamp,
        uint256 indexed _l2timestamp
    );

    /// @notice Emitted when an output is deleted.
    event l2OutputDeleted(
        bytes32 indexed _l2Output,
        uint256 indexed _l1Timestamp,
        uint256 indexed _l2timestamp
    );

    /**********************
     * Contract Variables *
     **********************/

    /// @notice The interval in seconds at which checkpoints must be submitted.
    uint256 public immutable SUBMISSION_INTERVAL;

    /// @notice The time between blocks on L2.
    uint256 public immutable L2_BLOCK_TIME;

    /// @notice The number of blocks in the chain before the first block in this contract.
    uint256 public immutable HISTORICAL_TOTAL_BLOCKS;

    /// @notice The timestamp of the first L2 block recorded in this contract.
    uint256 public immutable STARTING_BLOCK_TIMESTAMP;

    /// @notice The timestamp of the most recent L2 block recorded in this contract.
    uint256 public latestBlockTimestamp;

    /// @notice A mapping from L2 timestamps to the output root for the block with that timestamp.
    mapping(uint256 => OutputProposal) internal l2Outputs;

    /// @notice OutputProposal represents a commitment to the L2 state.
    /// The timestamp is the L1 timestamp that the output root is posted.
    /// This timestamp is used to verify that the finalization period
    /// has passed since the output root was submitted.
    struct OutputProposal {
        bytes32 outputRoot;
        uint256 timestamp;
    }

    /***************
     * Constructor *
     ***************/

    /**
     * @notice Initialize the L2OutputOracle contract.
     * @param _submissionInterval The desired interval in seconds at which
     *        checkpoints must be submitted.
     * @param _l2BlockTime The desired L2 inter-block time in seconds.
     * @param _genesisL2Output The initial L2 output of the L2 chain.
     * @param _historicalTotalBlocks The number of blocks that preceding the
     *        initialization of the L2 chain.
     * @param _startingBlockTimestamp The timestamp to start L2 block at.
     */
    constructor(
        uint256 _submissionInterval,
        uint256 _l2BlockTime,
        bytes32 _genesisL2Output,
        uint256 _historicalTotalBlocks,
        uint256 _startingBlockTimestamp,
        address sequencer
    ) {
        require(
            _submissionInterval % _l2BlockTime == 0,
            "Submission Interval must be a multiple of L2 Block Time"
        );

        SUBMISSION_INTERVAL = _submissionInterval;
        L2_BLOCK_TIME = _l2BlockTime;
        // solhint-disable-next-line not-rely-on-time
        l2Outputs[_startingBlockTimestamp] = OutputProposal(_genesisL2Output, block.timestamp);
        HISTORICAL_TOTAL_BLOCKS = _historicalTotalBlocks;
        // solhint-disable-next-line not-rely-on-time
        latestBlockTimestamp = _startingBlockTimestamp;
        // solhint-disable-next-line not-rely-on-time
        STARTING_BLOCK_TIMESTAMP = _startingBlockTimestamp;

        _transferOwnership(sequencer);
    }

    /*********************************
     * External and Public Functions *
     *********************************/

    /**
     * @notice Accepts an L2 outputRoot and the timestamp of the corresponding L2 block. The
     * timestamp must be equal to the current value returned by `nextTimestamp()` in order to be
     * accepted.
     * This function may only be called by the Sequencer.
     * @param _l2Output The L2 output of the checkpoint block.
     * @param _l2timestamp The L2 block timestamp that resulted in _l2Output.
     * @param _l1Blockhash A block hash which must be included in the current chain.
     * @param _l1Blocknumber The block number with the specified block hash.
     */
    function appendL2Output(
        bytes32 _l2Output,
        uint256 _l2timestamp,
        bytes32 _l1Blockhash,
        uint256 _l1Blocknumber
    ) external payable onlyOwner {
        require(_l2timestamp < block.timestamp, "Cannot append L2 output in future");
        require(_l2timestamp == nextTimestamp(), "Timestamp not equal to next expected timestamp");
        require(_l2Output != bytes32(0), "Cannot submit empty L2 output");
=======
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {
    OwnableUpgradeable
} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { Semver } from "../universal/Semver.sol";
import { Types } from "../libraries/Types.sol";

/**
 * @custom:proxied
 * @title L2OutputOracle
 * @notice The L2 state is committed to in this contract
 *         The payable keyword is used on proposeL2Output to save gas on the msg.value check.
 *         This contract should be deployed behind an upgradable proxy
 */
// slither-disable-next-line locked-ether
contract L2OutputOracle is OwnableUpgradeable, Semver {
    /**
     * @notice The interval in L2 blocks at which checkpoints must be submitted.
     */
    uint256 public immutable SUBMISSION_INTERVAL;

    /**
     * @notice The number of blocks in the chain before the first block in this contract.
     */
    uint256 public immutable HISTORICAL_TOTAL_BLOCKS;

    /**
     * @notice The number of the first L2 block recorded in this contract.
     */
    uint256 public immutable STARTING_BLOCK_NUMBER;

    /**
     * @notice The timestamp of the first L2 block recorded in this contract.
     */
    uint256 public immutable STARTING_TIMESTAMP;

    /**
     * @notice The time between L2 blocks in seconds.
     */
    uint256 public immutable L2_BLOCK_TIME;

    /**
     * @notice The address of the proposer;
     */
    address public proposer;

    /**
     * @notice The number of the most recent L2 block recorded in this contract.
     */
    uint256 public latestBlockNumber;

    /**
     * @notice A mapping from L2 block numbers to the respective output root. Note that these
     *         outputs should not be considered finalized until the finalization period (as defined
     *         in the Optimism Portal) has passed.
     */
    mapping(uint256 => Types.OutputProposal) internal l2Outputs;

    /**
     * @notice Emitted when an output is proposed.
     *
     * @param outputRoot    The output root.
     * @param l1Timestamp   The L1 timestamp when proposed.
     * @param l2BlockNumber The L2 block number of the output root.
     */
    event OutputProposed(
        bytes32 indexed outputRoot,
        uint256 indexed l1Timestamp,
        uint256 indexed l2BlockNumber
    );

    /**
     * @notice Emitted when an output is deleted.
     *
     * @param outputRoot    The output root.
     * @param l1Timestamp   The L1 timestamp when proposed.
     * @param l2BlockNumber The L2 block number of the output root.
     */
    event OutputDeleted(
        bytes32 indexed outputRoot,
        uint256 indexed l1Timestamp,
        uint256 indexed l2BlockNumber
    );

    /**
     * @notice Emitted when the proposer address is changed.
     *
     * @param previousProposer The previous proposer address.
     * @param newProposer      The new proposer address.
     */
    event ProposerChanged(address indexed previousProposer, address indexed newProposer);

    /**
     * @notice Reverts if called by any account other than the proposer.
     */
    modifier onlyProposer() {
        require(proposer == msg.sender, "L2OutputOracle: function can only be called by proposer");
        _;
    }

    /**
     * @custom:semver 0.0.1
     *
     * @param _submissionInterval    Interval in blocks at which checkpoints must be submitted.
     * @param _genesisL2Output       The initial L2 output of the L2 chain.
     * @param _historicalTotalBlocks Number of blocks preceding this L2 chain.
     * @param _startingBlockNumber   The number of the first L2 block.
     * @param _startingTimestamp     The timestamp of the first L2 block.
     * @param _l2BlockTime           The time per L2 block, in seconds.
     * @param _proposer              The address of the proposer.
     * @param _owner                 The address of the owner.
     */
    constructor(
        uint256 _submissionInterval,
        bytes32 _genesisL2Output,
        uint256 _historicalTotalBlocks,
        uint256 _startingBlockNumber,
        uint256 _startingTimestamp,
        uint256 _l2BlockTime,
        address _proposer,
        address _owner
    ) Semver(0, 0, 1) {
        require(
            _startingTimestamp <= block.timestamp,
            "L2OutputOracle: starting L2 timestamp must be less than current time"
        );

        SUBMISSION_INTERVAL = _submissionInterval;
        HISTORICAL_TOTAL_BLOCKS = _historicalTotalBlocks;
        STARTING_BLOCK_NUMBER = _startingBlockNumber;
        STARTING_TIMESTAMP = _startingTimestamp;
        L2_BLOCK_TIME = _l2BlockTime;

        initialize(_genesisL2Output, _proposer, _owner);
    }

    /**
     * @notice Initializer.
     *
     * @param _genesisL2Output     The initial L2 output of the L2 chain.
     * @param _proposer            The address of the proposer.
     * @param _owner               The address of the owner.
     */
    function initialize(
        bytes32 _genesisL2Output,
        address _proposer,
        address _owner
    ) public initializer {
        require(_proposer != _owner, "L2OutputOracle: proposer cannot be the same as the owner");
        l2Outputs[STARTING_BLOCK_NUMBER] = Types.OutputProposal(_genesisL2Output, block.timestamp);
        latestBlockNumber = STARTING_BLOCK_NUMBER;
        __Ownable_init();
        changeProposer(_proposer);
        _transferOwnership(_owner);
    }

    /**
     * @notice Deletes the most recent output. This is used to remove the most recent output in the
     *         event that an erreneous output is submitted. It can only be called by the contract's
     *         owner, not the proposer. Longer term, this should be replaced with a more robust
     *         mechanism which will allow deletion of proposals shown to be invalid by a fault
     *         proof.
     *
     * @param _proposal Represents the output proposal to delete
     */
    // solhint-disable-next-line ordering
    function deleteL2Output(Types.OutputProposal memory _proposal) external onlyOwner {
        Types.OutputProposal memory outputToDelete = l2Outputs[latestBlockNumber];

        require(
            _proposal.outputRoot == outputToDelete.outputRoot,
            "L2OutputOracle: output root to delete does not match the latest output proposal"
        );

        require(
            _proposal.timestamp == outputToDelete.timestamp,
            "L2OutputOracle: timestamp to delete does not match the latest output proposal"
        );

        emit OutputDeleted(outputToDelete.outputRoot, outputToDelete.timestamp, latestBlockNumber);

        delete l2Outputs[latestBlockNumber];
        latestBlockNumber = latestBlockNumber - SUBMISSION_INTERVAL;
    }

    /**
     * @notice Accepts an outputRoot and the timestamp of the corresponding L2 block. The
     *         timestamp must be equal to the current value returned by `nextTimestamp()` in order
     *         to be accepted. This function may only be called by the Proposer.
     *
     * @param _outputRoot    The L2 output of the checkpoint block.
     * @param _l2BlockNumber The L2 block number that resulted in _outputRoot.
     * @param _l1Blockhash   A block hash which must be included in the current chain.
     * @param _l1BlockNumber The block number with the specified block hash.
     */
    function proposeL2Output(
        bytes32 _outputRoot,
        uint256 _l2BlockNumber,
        bytes32 _l1Blockhash,
        uint256 _l1BlockNumber
    ) external payable onlyProposer {
        require(
            _l2BlockNumber == nextBlockNumber(),
            "L2OutputOracle: block number must be equal to next expected block number"
        );

        require(
            computeL2Timestamp(_l2BlockNumber) < block.timestamp,
            "L2OutputOracle: cannot propose L2 output in the future"
        );

        require(
            _outputRoot != bytes32(0),
            "L2OutputOracle: L2 output proposal cannot be the zero hash"
        );
<<<<<<< HEAD
        require(_l2Output != bytes32(0), "OutputOracle: Cannot submit empty L2 output.");
>>>>>>> v0.5.23
=======
>>>>>>> v0.5.24

        if (_l1Blockhash != bytes32(0)) {
            // This check allows the proposer to propose an output based on a given L1 block,
            // without fear that it will be reorged out.
            // It will also revert if the blockheight provided is more than 256 blocks behind the
            // chain tip (as the hash will return as zero). This does open the door to a griefing
            // attack in which the proposer's submission is censored until the block is no longer
            // retrievable, if the proposer is experiencing this attack it can simply leave out the
            // blockhash value, and delay submission until it is confident that the L1 block is
            // finalized.
            require(
<<<<<<< HEAD
                blockhash(_l1Blocknumber) == _l1Blockhash,
                "Blockhash does not match the hash at the expected height."
            );
        }

        l2Outputs[_l2timestamp] = OutputProposal(_l2Output, block.timestamp);
        latestBlockTimestamp = _l2timestamp;

        emit l2OutputAppended(_l2Output, block.timestamp, _l2timestamp);
    }

    /**
     * @notice Deletes the most recent output.
     * @param _proposal Represents the output proposal to delete
     */
    function deleteL2Output(OutputProposal memory _proposal) external onlyOwner {
        OutputProposal memory outputToDelete = l2Outputs[latestBlockTimestamp];

        require(
            _proposal.outputRoot == outputToDelete.outputRoot,
            "Can only delete the most recent output."
        );
        require(_proposal.timestamp == outputToDelete.timestamp, "");

        emit l2OutputDeleted(
            outputToDelete.outputRoot,
            outputToDelete.timestamp,
            latestBlockTimestamp
        );

        delete l2Outputs[latestBlockTimestamp];
        latestBlockTimestamp = latestBlockTimestamp - SUBMISSION_INTERVAL;
    }

    /**
     * @notice Computes the timestamp of the next L2 block that needs to be checkpointed.
     */
    function nextTimestamp() public view returns (uint256) {
        return latestBlockTimestamp + SUBMISSION_INTERVAL;
    }

    /**
     * @notice Returns the L2 output proposal given a target L2 block timestamp.
     * Returns a null output proposal if none is found.
     * @param _l2Timestamp The L2 block timestamp of the target block.
     */
    function getL2Output(uint256 _l2Timestamp) external view returns (OutputProposal memory) {
        return l2Outputs[_l2Timestamp];
    }

    /**
     * @notice Computes the L2 block number given a target L2 block timestamp.
     * @param _l2timestamp The L2 block timestamp of the target block.
     */
    function computeL2BlockNumber(uint256 _l2timestamp) external view returns (uint256) {
        require(
            _l2timestamp >= STARTING_BLOCK_TIMESTAMP,
            "Timestamp prior to startingBlockTimestamp"
        );
        // For the first block recorded (ie. _l2timestamp = STARTING_BLOCK_TIMESTAMP), the
        // L2BlockNumber should be HISTORICAL_TOTAL_BLOCKS + 1.
        unchecked {
            return
                HISTORICAL_TOTAL_BLOCKS +
                ((_l2timestamp - STARTING_BLOCK_TIMESTAMP) / L2_BLOCK_TIME);
        }
=======
                blockhash(_l1BlockNumber) == _l1Blockhash,
                "L2OutputOracle: blockhash does not match the hash at the expected height"
            );
        }

        l2Outputs[_l2BlockNumber] = Types.OutputProposal(_outputRoot, block.timestamp);
        latestBlockNumber = _l2BlockNumber;

        emit OutputProposed(_outputRoot, block.timestamp, _l2BlockNumber);
    }

    /**
     * @notice Returns the L2 output proposal associated with a target L2 block number. If the
     *         L2 block number provided is between checkpoints, this function will rerutn the next
     *         proposal for the next checkpoint.
     *         Reverts if the output proposal is either not found, or predates
     *         the STARTING_BLOCK_NUMBER.
     *
     * @param _l2BlockNumber The L2 block number of the target block.
     */
    function getL2Output(uint256 _l2BlockNumber)
        external
        view
        returns (Types.OutputProposal memory)
    {
        require(
            _l2BlockNumber >= STARTING_BLOCK_NUMBER,
            "L2OutputOracle: block number cannot be less than the starting block number."
        );

        // Find the distance between _l2BlockNumber, and the checkpoint block before it.
        uint256 offset = (_l2BlockNumber - STARTING_BLOCK_NUMBER) % SUBMISSION_INTERVAL;

        // If the offset is zero, then the _l2BlockNumber should be checkpointed.
        // Otherwise, we'll look up the next block that will be checkpointed.
        uint256 lookupBlockNumber = offset == 0
            ? _l2BlockNumber
            : _l2BlockNumber + (SUBMISSION_INTERVAL - offset);

        Types.OutputProposal memory output = l2Outputs[lookupBlockNumber];
        require(
            output.outputRoot != bytes32(0),
            "L2OutputOracle: No output found for that block number."
        );
        return output;
    }

    /**
     * @notice Overrides the standard implementation of transferOwnership
     *         to add the requirement that the owner and proposer are distinct.
     *         Can only be called by the current owner.
     */
    function transferOwnership(address _newOwner) public override onlyOwner {
        require(_newOwner != proposer, "L2OutputOracle: owner cannot be the same as the proposer");
        super.transferOwnership(_newOwner);
    }

    /**
     * @notice Transfers the proposer role to a new account (`newProposer`).
     *         Can only be called by the current owner.
     */
    function changeProposer(address _newProposer) public onlyOwner {
        require(
            _newProposer != address(0),
            "L2OutputOracle: new proposer cannot be the zero address"
        );

        require(
            _newProposer != owner(),
            "L2OutputOracle: proposer cannot be the same as the owner"
        );

        emit ProposerChanged(proposer, _newProposer);
        proposer = _newProposer;
    }

    /**
     * @notice Computes the block number of the next L2 block that needs to be checkpointed.
     */
    function nextBlockNumber() public view returns (uint256) {
        return latestBlockNumber + SUBMISSION_INTERVAL;
    }

    /**
     * @notice Returns the L2 timestamp corresponding to a given L2 block number.
     *         If the L2 block number provided is between checkpoints, this function will return the
     *         timestamp of the previous checkpoint.
     *
     * @param _l2BlockNumber The L2 block number of the target block.
     */
    function computeL2Timestamp(uint256 _l2BlockNumber) public view returns (uint256) {
        require(
            _l2BlockNumber >= STARTING_BLOCK_NUMBER,
            "L2OutputOracle: block number must be greater than or equal to starting block number"
        );

        return STARTING_TIMESTAMP + ((_l2BlockNumber - STARTING_BLOCK_NUMBER) * L2_BLOCK_TIME);
    }
<<<<<<< HEAD

    /**
     * @notice Transfers the sequencer role to a new account (`newSequencer`).
     *         Can only be called by the current owner.
     */
    function changeSequencer(address _newSequencer) public onlyOwner {
        require(_newSequencer != address(0), "OutputOracle: new sequencer is the zero address");
        require(_newSequencer != owner(), "OutputOracle: sequencer cannot be same as the owner");
        emit SequencerChanged(sequencer, _newSequencer);
        sequencer = _newSequencer;
>>>>>>> v0.5.23
    }
=======
>>>>>>> v0.5.24
}
