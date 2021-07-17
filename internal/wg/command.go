package wg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const commandWG = `wg`
const argPrivateKey = `genkey`
const argPublicKey = `pubkey`
const argSyncConf = `syncconf`

// GenPrivateKey generate a new private key using the command wg
func GenPrivateKey() string {
	//! Not need for root
	out, err := executeWGNoInput([]string{argPrivateKey})
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(out)
}

// GenPublicKey generate a new private key using the command wg
func GenPublicKey(privateKey string) string {
	//! Not need for root
	out, err := executeWGWithInput([]string{argPublicKey}, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(out)
}

// SyncPeerConf will sync the peer configuration file with the interface file
func SyncPeerConf(inter string, fileName string) error {
	//! Need for root
	out, err := executeWGNoInput([]string{argSyncConf, inter, fileName})
	if err != nil {
		return fmt.Errorf("fail to sync peer config: %s - %w", strings.TrimSpace(out), err)
	}
	return nil
}

func executeWGWithInput(args []string, input string) (output string, err error) {
	/* #nosec */
	cmd := exec.Command(commandWG, args...)
	cmd.Stdin = bytes.NewBuffer([]byte(input))
	return executeWG(cmd)
}

func executeWGNoInput(args []string) (output string, err error) {
	/* #nosec */
	cmd := exec.Command(commandWG, args...)
	return executeWG(cmd)
}

func executeWG(cmd *exec.Cmd) (output string, err error) {
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// CheckWGExists will check if it can find the wg bin on the user path
func CheckWGExists() bool {
	//! No need for root
	_, err := exec.LookPath(commandWG)
	if err != nil {
		return false
	}
	return true
}
