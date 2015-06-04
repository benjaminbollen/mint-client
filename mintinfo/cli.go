package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/tendermint/tendermint/binary"
	"github.com/tendermint/tendermint/types"
)

func prettyPrint(o interface{}) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, binary.JSONBytes(o), "", "\t")
	if err != nil {
		return "", err
	}
	return string(prettyJSON.Bytes()), nil
}

func FieldFromTag(v reflect.Value, field string) (string, error) {
	iv := v.Interface()
	st := reflect.TypeOf(iv)
	for i := 0; i < v.NumField(); i++ {
		tag := st.Field(i).Tag.Get("json")
		if tag == field {
			return st.Field(i).Name, nil
		}
	}
	return "", fmt.Errorf("Invalid field name")
}

func formatOutput(c *cli.Context, i int, o interface{}) (string, error) {
	args := c.Args()
	if len(args) < i+1 {
		return prettyPrint(o)
	}
	arg0 := args[i]
	v := reflect.ValueOf(o).Elem()
	name, err := FieldFromTag(v, arg0)
	if err != nil {
		return "", err
	}
	f := v.FieldByName(name)
	return prettyPrint(f.Interface())
}

func cliStatus(c *cli.Context) {
	r, err := client.Status()
	ifExit(err)
	s, err := formatOutput(c, 0, r)
	ifExit(err)
	fmt.Println(s)
}

func cliNetInfo(c *cli.Context) {
	r, err := client.NetInfo()
	ifExit(err)
	s, err := formatOutput(c, 0, r)
	ifExit(err)
	fmt.Println(s)
}

func cliGenesis(c *cli.Context) {
	r, err := client.Genesis()
	ifExit(err)
	s, err := formatOutput(c, 0, r)
	ifExit(err)
	fmt.Println(s)
}

func cliValidators(c *cli.Context) {
	r, err := client.ListValidators()
	ifExit(err)
	s, err := formatOutput(c, 0, r)
	ifExit(err)
	fmt.Println(s)
}

func cliConsensus(c *cli.Context) {
	r, err := client.DumpConsensusState()
	ifExit(err)
	rs := r.RoundState
	prss := r.PeerRoundStates
	// TODO ... get fields
	fmt.Println("round_state:", rs)
	fmt.Println("peer_round_states:")
	for _, prs := range prss {
		fmt.Println(prs)
	}
}

func cliUnconfirmed(c *cli.Context) {
	r, err := client.ListUnconfirmedTxs()
	ifExit(err)
	s, err := formatOutput(c, 0, r)
	ifExit(err)
	fmt.Println(s)
}

func cliAccounts(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		r, err := client.ListAccounts()
		ifExit(err)
		s, err := formatOutput(c, 0, r)
		ifExit(err)
		fmt.Println(s)
	} else {
		addr := args[0]
		addrBytes, err := hex.DecodeString(addr)
		if err != nil {
			exit(fmt.Errorf("Addr %s is improper hex: %v", addr, err))
		}
		r, err := client.GetAccount(addrBytes)
		ifExit(err)
		s, err := formatOutput(c, 1, r)
		ifExit(err)
		fmt.Println(s)
	}
}

func cliNames(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		r, err := client.ListNames()
		ifExit(err)
		s, err := formatOutput(c, 1, r)
		ifExit(err)
		fmt.Println(s)
	} else {
		name := args[0]
		r, err := client.GetName(name)
		ifExit(err)
		s, err := formatOutput(c, 1, r)
		ifExit(err)
		if len(args) > 1 {
			if args[1] == "data" {
				s, err = strconv.Unquote(s)
				ifExit(err)
			}
		}
		fmt.Println(s)
	}
}

func cliBlocks(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		exit(fmt.Errorf("must specify a height to get a single block, or two heights to get all blocks between them"))
	} else if len(args) == 1 {
		height := args[0]
		i, err := strconv.ParseUint(height, 10, 32)
		ifExit(err)
		r, err := client.GetBlock(uint(i))
		ifExit(err)
		s, err := formatOutput(c, 1, r)
		ifExit(err)
		fmt.Println(s)
	} else {
		minHeightS, maxHeightS := args[0], args[1]
		minHeight, err := strconv.ParseUint(minHeightS, 10, 32)
		ifExit(err)
		maxHeight, err := strconv.ParseUint(maxHeightS, 10, 32)
		ifExit(err)
		if maxHeight <= minHeight {
			exit(fmt.Errorf("maxHeight must be greater than minHeight"))
		}
		r, err := client.BlockchainInfo(uint(minHeight), uint(maxHeight))
		ifExit(err)
		s, err := formatOutput(c, 2, r)
		ifExit(err)
		fmt.Println(s)
	}
}

func cliStorage(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		exit(fmt.Errorf("must specify an address to dump all storage, and an optional key to get just that storage"))
	} else if len(args) == 1 {
		addr := args[0]
		addrBytes, err := hex.DecodeString(addr)
		ifExit(err)
		r, err := client.DumpStorage(addrBytes)
		ifExit(err)
		s, err := formatOutput(c, 1, r)
		ifExit(err)
		fmt.Println(s)
	} else {
		addr, key := args[0], args[1]
		addrBytes, err := hex.DecodeString(addr)
		ifExit(err)
		keyBytes, err := hex.DecodeString(key)
		ifExit(err)
		r, err := client.GetStorage(addrBytes, keyBytes)
		ifExit(err)
		s, err := formatOutput(c, 2, r)
		ifExit(err)
		fmt.Println(s)
	}
}

func cliCall(c *cli.Context) {
	args := c.Args()
	if len(args) < 2 {
		exit(fmt.Errorf("must specify an address and data to send"))
	}
	addr, data := args[0], args[1]
	addrBytes, err := hex.DecodeString(addr)
	ifExit(err)
	dataBytes, err := hex.DecodeString(data)
	ifExit(err)
	r, err := client.Call(addrBytes, dataBytes)
	ifExit(err)
	s, err := formatOutput(c, 2, r)
	ifExit(err)
	fmt.Println(s)
}

func cliCallCode(c *cli.Context) {
	args := c.Args()
	if len(args) < 2 {
		exit(fmt.Errorf("must specify code to execute and data to send"))
	}
	code, data := args[0], args[1]
	codeBytes, err := hex.DecodeString(code)
	ifExit(err)
	dataBytes, err := hex.DecodeString(data)
	ifExit(err)
	r, err := client.CallCode(codeBytes, dataBytes)
	ifExit(err)
	s, err := formatOutput(c, 2, r)
	ifExit(err)
	fmt.Println(s)
}

func cliBroadcast(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		exit(fmt.Errorf("must specify transaction bytes to broadcast"))
	}
	txS := args[0]
	// TODO: we should switch over a hex vs. json flag
	txBytes, err := hex.DecodeString(txS)
	ifExit(err)
	var tx types.Tx
	n := new(int64)
	buf := bytes.NewBuffer(txBytes)
	binary.ReadBinary(tx, buf, n, &err)
	ifExit(err)
	r, err := client.BroadcastTx(tx)
	ifExit(err)
	s, err := formatOutput(c, 1, r)
	ifExit(err)
	fmt.Println(s)
}