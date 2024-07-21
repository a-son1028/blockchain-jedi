const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");
const { expect } = require("chai");
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

describe("KYC contract test", async function () {
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




  const createUser = async (userInfo) => {
    const options = {
      pinataMetadata: {
        name: "metadata.json",
      },
      pinataOptions: {
        cidVersion: 0,
      },
    };
  
    const payload = {
      data: userInfo
    }
    const { IpfsHash } = await pinata.pinJSONToIPFS(payload, options);
    await contract.mint(TOKEN_ID++, IpfsHash);

    // Generate an RSA key pair
    const key = new NodeRSA({ b: 512 }); // You can adjust the key size as needed
    const rsaPublicKey = key.exportKey('public');
    const rsaPrivateKey = key.exportKey('private');

    return { nftURI: IpfsHash, rsaPublicKey, rsaPrivateKey }
  }



  it("Create user A and delegate B", async function () {
    const userA = await createUser({
      fullName: "User A",
      address: "Address 1",
    })
    // generate jedi private key
    const jediPKResA = await JediInstance.get(`/jedi-private-key/`)
    const jediPrivateKeyA = jediPKResA.data.result
    expect(userA).to.be.a('object');
    console.log({userA: {
      ...userA,
      jediPrivateKey: jediPrivateKeyA
    }})


    const userB = await createUser({
      fullName: "User B",
      address: "Address B",
    })

    // delete user B
    const jediPKResB = await JediInstance.get(`/jedi-private-key/?parent=${jediPrivateKeyA}`)
    const jediPrivateKeyB = jediPKResB.data.result
    
    const publicKeyEncryption = new NodeRSA(userB.rsaPublicKey);
    const privateKeyDecryption = new NodeRSA(userB.rsaPrivateKey);

    const encryptedMessage = publicKeyEncryption.encrypt(jediPrivateKeyB)
    const decryptedMessage = privateKeyDecryption.decrypt(encryptedMessage, 'utf8');


    console.log({userB: {
      ...userB,
      jediPrivateKey: decryptedMessage
    }})
  });
});



