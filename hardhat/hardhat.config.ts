import { HardhatUserConfig } from 'hardhat/config'
import '@typechain/hardhat'
import '@nomicfoundation/hardhat-toolbox'
import '@nomicfoundation/hardhat-ethers'
import '@nomicfoundation/hardhat-chai-matchers'
import 'hardhat-deploy'
import * as dotenv from 'dotenv'

import {
  ACCOUNT_ADDRESSES,
  PRIVATE_KEYS,
} from './utils/accounts'

const ENV_FILE = process.env.DOTENV_CONFIG_PATH || '../.env'
dotenv.config({ path: ENV_FILE })

const INFURA_KEY = process.env.INFURA_KEY || "";

const config: HardhatUserConfig = {
  solidity: '0.8.21',
  defaultNetwork: 'geth',
  namedAccounts: ACCOUNT_ADDRESSES,
  networks: {
    hardhat: {},
    geth: {
      url: 'http://localhost:8545',
      chainId: 1337,
      accounts: PRIVATE_KEYS,
    },
    sepolia: {
      url: `https://sepolia.infura.io/v3/${INFURA_KEY}`,
      accounts: PRIVATE_KEYS,
    },
  },
};

module.exports = config
