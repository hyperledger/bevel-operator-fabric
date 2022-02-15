
const FabricCAServices = require('fabric-ca-client');
const { Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
let mspId = "Org1MSP";
const admin = 'enroll';
const adminpw = 'enrollpw';
const caUrl = 'org1-ca.default' // refer this from the connection profile
const yaml = require("js-yaml");
async function enroll() {
    try {
        
        const ccpPath = path.resolve(__dirname,  "connection-org.yaml");
        if (ccpPath.includes(".yaml")) {
            ccp = yaml.load(fs.readFileSync(ccpPath, "utf8"));
        } else {
            ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));
        }

        const caInfo = ccp.certificateAuthorities[caUrl];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        const walletPath = path.join(process.cwd(), 'wallet', mspId);
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        const identity = await wallet.get(admin);
        if (identity) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        const enrollment = await ca.enroll({ enrollmentID: admin, enrollmentSecret: adminpw });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: mspId,
            type: 'X.509',
        };
        await wallet.put(admin, x509Identity);
        console.log('Successfully enrolled admin user "enroll" and imported it into the wallet');

    } catch (error) {
        console.error(`Failed to enroll admin user "enroll": ${error}`);
        process.exit(1);
    }
}

enroll();