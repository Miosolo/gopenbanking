echo "Starting up Gopenbanking Project"
echo "Server Name: anz.italktoyou.cn"

# First setting all environment path
domain="anz.italktoyou.cn"
export set CORE_PEER_MSPCONFIGPATH=/opt/hyperledger/gopenbanking/crypto-config/peerOrganizations/$domain/users/peer0.$domain/msp
export set CORE_PEER_LOCALMSPID=ANZBankMSP
export set FABRIC_CFG_PATH=/opt/hyperledger/gopenbanking/config

# Then starting all peer and orderer nodes
if test -e /opt/hyperledger/gopenbanking/config/orderer.yaml
then
    echo "----------------------------------------"
    echo "Find orderer.yaml, starting orderer node"
    echo "----------------------------------------"
    # if ps aux | grep -c "orderer" is larger than or equal 1
    # it should stop starting orderer
    if [ $(ps aux | grep -c "orderer") -le 1 ]
    then
        nohup orderer start >> /opt/hyperledger/gopenbanking/log/log_orderer.log 2>&1 &
        echo "---------------------------------"
        echo "orderer node started successfully"
        echo "---------------------------------"
    else
        echo "-----------------------------------------------------------"
        echo "orderer node has already started, end starting orderer node"
        echo "-----------------------------------------------------------"
    fi
else
    echo "---------------------------------------------------------------------------------------------------"
    echo "Cannot find orderer.yaml in /opt/hyperledger/goenbanking/config document, end starting orderer node"
    echo "---------------------------------------------------------------------------------------------------"
fi

if test -e /opt/hyperledger/gopenbanking/config/core.yaml
then
    echo "----------------------------------"
    echo "Find core.yaml, starting peer node"
    echo "----------------------------------"
    if [ $(ps aux | grep -c "peer") -le 1 ]
    then
        nohup peer node start >> /opt/hyperledger/gopenbanking/log/log_peer.log 2>&1 &
        echo "------------------------------"
        echo "peer node started successfully"
        echo "------------------------------"
    else
        echo "-----------------------------------------------------"
        echo "peer node has already started, end starting peer node"
        echo "-----------------------------------------------------"
    fi
else
    echo "----------------------------------------------------------------------------------------------"
    echo "Cannot find core.yaml in /opt/hyperledger/gopenbanking/config document, end starting peer node"
    echo "----------------------------------------------------------------------------------------------"
fi

echo "All peers and orderer has started successfully"
