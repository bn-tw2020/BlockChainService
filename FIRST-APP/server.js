// 1. 모듈추가
const express = require('express');
var bodyParser = require('body-parser');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');

// 2. 하이퍼레저 속성 설정
const { FileSystemWallet, Gateway } = require('fabric-network');

// 3. 웹 서버 설정
const app = express();
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const PORT = 8080;
const HOST = '0.0.0.0';
// app.use
app.use(express.static(path.join(__dirname, 'views')));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// 4. 웹페이지 라우팅
// 4.1 /GET 라우팅 index.html
// 4.2 /key-create.html GET
// 4.3 /key-query.html GET
app.get('/', (req, res)=>{
    res.sendFile(__dirname + '/index.html');
})


// 5. REST API 라우팅
// 5.1 /asset POST -> sacc set param: key, value
app.post('/asset', async(req, res)=>{
    // 클라이언트 파라미터 받는 부분
    const {key, value} = req.body;

    // user1 인증서 가져오기 및 검사
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // 체인코드 연결
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('sacc');
    await contract.submitTransaction('set', key, value);
    console.log('Transaction has been submitted');
    await gateway.disconnect();

    res.status(200).send('Transaction has been submitted');

})

// 5-2 /asset GET -> sacc get param: key
app.get('/asset', async(req, res)=>{
    const {key} = req.query;
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('sacc');
    const result = await contract.evaluateTransaction('get',key);
    console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

    // 클라이언트에 JSON 형태로 보내기
    var obj = JSON.parse(result);
    res.status(200).json(obj);
    
})

// 5.3 /assethistory GET -> sacc getHistoryForKey param: key
app.get('/assethistory', async(req, res)=>{
    const {key} = req.query;
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('sacc');
    const result = await contract.evaluateTransaction('getHistoryForKey',key);
    console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

    // 클라이언트에 JSON 형태로 보내기
    var obj = JSON.parse(result);
    res.status(200).json(obj);
    
})

// 5.4 /assets GET -> sacc getAllKeys
app.get('/assets', async(req, res)=>{
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('sacc');
    const result = await contract.evaluateTransaction('getAllKeys');
    console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

    // 클라이언트에 JSON 형태로 보내기
    var obj = JSON.parse(result);
    res.status(200).json(obj);
    
})
// 6. 서버시작
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);