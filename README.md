
###  Overview
NEAR Protocol is a decentralized application platform that is secure enough to manage high value assets like money or identity and performant enough to make them useful for everyday people, putting the power of the Open Web in their hands.


# NEAR Protocol Warchest Bot Stakewars Warchest bot  for monitoring and adapting the stake Ⓝ
It is warchest bot written in Go language which is used for monitoring and adapting the stake


### Requirements 
This warchest bot requires:
* Go language version 1.14.6 and higher
* near-shell install in the system where script will run.
    * To install near-shell run below command:
    `
    npm install -g near-shell
    `
* near need to be logged in shell (command)
 export NODE_ENV=betanet
 export NODE_ENV=testnet

To Build the script you need to execute the following command
`
### Build

go build main.go
`

### Run
To run the command you need to execute the command
`
./main
`



# Guildnet Network Only!!!!!.
Install near-cli
git clone https://github.com/near-guildnet/near-cli.git


cd near-cli


npm install -g

# Setting up your environment
To use the guildnet network you need to update the environment via the command line.


  export NODE_ENV=guildnet
  
  
 echo 'export NODE_ENV=guildnet' >> ~/.bashrc 
------------------------------

Install Nearup
The Prerequisite has python3, git and curl toolset, which have been installed

curl --proto '=https' --tlsv1.2 -sSfL https://raw.githubusercontent.com/near-guildnet/nearup/master/nearup | python3

Launch validator node
We recommand to use Officially Compiled Binary to lauch validator node,
###  source ~/.profile
###  then
nearup guildnet --nodocker

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
