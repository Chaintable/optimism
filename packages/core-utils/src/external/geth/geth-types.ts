// Types explicitly related to dealing with Geth.

/**
 * Represents the Ethereum state, in the format that Geth expects it.
 */
export interface State {
  [address: string]: {
<<<<<<< HEAD
    nonce: number
    balance: string
=======
    nonce?: string
    balance?: string
>>>>>>> v0.5.23
    codeHash?: string
    root?: string
    code?: string
    storage?: {
      [key: string]: string
    }
    secretKey?: string
  }
}

/**
 * Represents Geth's ChainConfig
 */
export interface ChainConfig {
  chainId: number
  homesteadBlock: number
  eip150Block: number
<<<<<<< HEAD
=======
  eip150Hash?: string
>>>>>>> v0.5.23
  eip155Block: number
  eip158Block: number
  byzantiumBlock: number
  constantinopleBlock: number
  petersburgBlock: number
  istanbulBlock: number
  muirGlacierBlock: number
  berlinBlock: number
  londonBlock?: number
  arrowGlacierBlock?: number
<<<<<<< HEAD
=======
  grayGlacierBlock?: number
<<<<<<< HEAD
>>>>>>> v0.5.23
  mergeForkBlock?: number
=======
  mergeNetsplitBlock?: number
>>>>>>> v0.5.24
  terminalTotalDifficulty?: number
  clique?: {
    period: number
    epoch: number
  }
  ethash?: {}
}

/**
 * Represents Geth's genesis file format.
 */
export interface Genesis {
  config: ChainConfig
<<<<<<< HEAD
  nonce?: number
  timestamp?: number
  difficulty: string
  mixHash?: string
  coinbase?: string
=======
  nonce?: string
  timestamp?: string
  difficulty: string
  mixHash?: string
  coinbase?: string
  number?: string
>>>>>>> v0.5.23
  gasLimit: string
  gasUsed?: string
  parentHash?: string
  extraData: string
  baseFeePerGas?: string
  alloc: State
}

/**
 * Represents the chain config for an Optimism chain
 */
export interface OptimismChainConfig extends ChainConfig {
  optimism: {
    baseFeeRecipient: string
    l1FeeRecipient: string
  }
}

/**
 * Represents the Genesis file format for an Optimism chain
 */
export interface OptimismGenesis extends Genesis {
  config: OptimismChainConfig
}
