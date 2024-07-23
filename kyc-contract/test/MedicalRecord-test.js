const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");
const { expect, use } = require("chai");
const crypto = require('crypto');
const {ethers, deployments, getNamedAccounts, getUnnamedAccounts} = require("hardhat")
const axios = require('axios');
const NodeRSA = require('node-rsa');

const ENCRYPTION_KEY = 'WbHPUuDNQLpFtSSGiBtTVSeqPdAHBZzw'; // Replace this with your actual encryption key
const IV_LENGTH = 16;
let TOKEN_ID = 1


const JediInstance = axios.default.create({
  baseURL: "http://localhost:8080"
})

describe("JEDI Blockchain for Medical record TEST", async function () {
  var owner, otherAccount, contract, pinata
  before(async () => {
    await loadFixture(
      deployOneYearLockFixture
    );
  })
 
  async function deployOneYearLockFixture() {
    const pinataSDK = require("@pinata/sdk");
    pinata = new pinataSDK({
      pinataApiKey: "5cf60e7146bc19a4e813",
      pinataSecretApiKey:
        "e0c23a5a169760c651f595778611245375008974b25bd79fc1e8718653cd4780",
    });

    // Contracts are deployed using the first signer/account by default
    const [ownerTemp, otherAccountTemp] = await ethers.getSigners();

    const Contract = await ethers.getContractFactory("KYC");
    contract = await Contract.deploy();

    owner = ownerTemp
    otherAccount = otherAccountTemp
  }



  const pinMetadataToIPFS = async (metadata) => {
    const options = {
      pinataMetadata: {
        name: "metadata.json",
      },
      pinataOptions: {
        cidVersion: 0,
      },
    };

    const { IpfsHash } = await pinata.pinJSONToIPFS(metadata, options);

    return IpfsHash
  }
  // pin the metadata to IPFS create a nft on blockchain and 
  const createNFT = async (IpfsHash) => {
    const nft = await contract.mint(TOKEN_ID++, IpfsHash);

    return nft
  }


  const userInfoA = {
    fullName: "patient A",
    address: "address 1",
    age: 20,
  }

  it("Step 1: Generate RSA keys and create NFT for storing KYC on blockchain", async function () {
    const key = new NodeRSA({ b: 512 }); 

    // Generate an RSA key pair
    const rsaPrivateKey = key.exportKey('private');
    const rsaPublicKey = key.exportKey('public');

    // Pin data to IPFS
    const IpfsHash = await pinMetadataToIPFS({
      data: {
        ...userInfoA,
        rsaPublicKey: rsaPublicKey
      }
    });

    expect(IpfsHash).to.be.a('string')

    // Create NFT
    const nft = await createNFT(IpfsHash)
    
    expect(nft).to.be.a('object')
    expect(nft).to.have.property('hash')
    expect(nft.hash).to.be.a('string')


    // save info for next steps
    userInfoA.rsaPublicKey = rsaPublicKey
    userInfoA.rsaPrivateKey = rsaPrivateKey
 

    // return data to user
    return {
      rsaPrivateKey
    }
  });

  it("Step 2.1: Delegate JEDI key for user A and encrypt it with RSA key", async function () {
    // generate jedi private key
    const jediPKResA = await JediInstance.get(`/jedi-private-key/`)
    const jediKeys = jediPKResA.data.data

    expect(jediKeys).to.be.a('string')

   

    const publicKeyEncryption = new NodeRSA(userInfoA.rsaPublicKey);
    const encryptedPublickey = publicKeyEncryption.encrypt(jediKeys)

    // save info for next steps
    userInfoA.encryptedJediKey = encryptedPublickey

    // return data to user
    return {
      encryptedPublickey
    }
  })

  it("Step 2.2: Decrypt encrypted jedi keys", async function () {
    const privateKeyDecryption = new NodeRSA(userInfoA.rsaPrivateKey);
    const jediKeys = privateKeyDecryption.decrypt(userInfoA.encryptedJediKey, 'utf8');

     // save info for next steps
      userInfoA.jediKeys = jediKeys

    // return data to user
    return {
      jediKeys
    }
  })

  it("Step 3: Delegate JEDI key for user B from user A", async function () {
    const jediKeyRes = await JediInstance.get(`/jedi-private-key?parent=${userInfoA.jediKeys}`)

    expect(jediKeyRes.data).to.be.a('object')
    expect(jediKeyRes.data).to.have.property('data')
    expect(jediKeyRes.data.data).to.be.a('string')

    const jediKeys = jediKeyRes.data.data
    return {
      jediKeys
    }
  })
  

  it("Encrypt and Decrypt message with JEDI", async function () {
    const message = "Hello world"
    
    const encryptedMsgRes = await JediInstance.post(`/encrypt`, {
      message,
      uri: "a/b/c/d",
    })
    
    expect(encryptedMsgRes).to.be.a('object')
    expect(encryptedMsgRes).to.have.property('data')
    expect(encryptedMsgRes.data.data).to.be.a('string')

    const encryptedMsg = encryptedMsgRes.data.data

    const decryptedMsgRes = await JediInstance.post(`/decrypt`, {
      encryptedMessage: encryptedMsg,
      uri: "a/b/c/d",
      key: userInfoA.jediKeys
    })
    
    expect(decryptedMsgRes).to.be.a('object')
    expect(decryptedMsgRes).to.have.property('data')
    expect(decryptedMsgRes.data.data).to.be.a('string')

    const decryptedMsg = decryptedMsgRes.data.data
    expect(decryptedMsg).to.be.eq(message)
  })
});



