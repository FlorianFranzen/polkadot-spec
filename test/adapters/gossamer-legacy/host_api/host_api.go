// Copyright (c) 2019 Web3 Technologies Foundation

// This file is part of Polkadot Host Test Suite

// Polkadot Host Test Suite is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Polkadot Host Tests is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Foobar.  If not, see <https://www.gnu.org/licenses/>.

package host_api

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ChainSafe/gossamer/dot/state"
	"github.com/ChainSafe/gossamer/lib/keystore"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmtime"
	"github.com/ChainSafe/gossamer/lib/trie"

	database "github.com/ChainSafe/chaindb"
)

// #include <errno.h>
import "C"

var RELATIVE_WASM_ADAPTER_PATH = "bin/hostapi_runtime_legacy.compact.wasm"

func GetRuntimePath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get current working directory")
	}
	return path.Join(dir, RELATIVE_WASM_ADAPTER_PATH)
}

func GetTestStorage() *state.TrieState {
	store, err := state.NewTrieState(database.NewMemDatabase(), trie.NewEmptyTrie())
	if err != nil {
		fmt.Println("Failed initialize storage: ", err)
		os.Exit(1)
	}
	return store
}

func ProcessHostApiCommand(args []string) {

	// List of expected flags
	functionTextPtr := flag.String("function", "", "Function to call (required).")
	inputTextPtr := flag.String("input", "", "Input to pass on call (required).")
	
	wasmtimeBoolPtr := flag.Bool("wasmtime", false, "Use wasmtime instead of wasmer.")

	// Parse provided argument list
	flag.CommandLine.Parse(args)

	if !flag.Parsed() {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Verify that all required flags are provided
	if (*functionTextPtr == "") || (*inputTextPtr == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	function := *functionTextPtr
	inputs   := strings.Split(*inputTextPtr, ",")

	// Initialize runtime environment...
	var rtm runtime.Instance
	if *wasmtimeBoolPtr {
		// ... using wasmtime
		cfg := &wasmtime.Config{
			Imports: wasmtime.ImportsHostAPITester,
		}
		cfg.Storage = GetTestStorage()
		cfg.Keystore = keystore.NewGenericKeystore("test")
		cfg.LogLvl = 2

		r, err := wasmtime.NewInstanceFromFile(GetRuntimePath(), cfg)
		if err != nil {
			fmt.Println("Failed initialize runtime: ", err)
			os.Exit(1)
		}
		rtm = r
	} else {
		// ... using wasmer
		cfg := &wasmer.Config{
			Imports: wasmer.RegisterImports_NodeRuntime,
		}
		cfg.Storage = GetTestStorage()
		cfg.Keystore = keystore.NewGenericKeystore("test")
		cfg.LogLvl = 2

		r, err := wasmer.NewInstanceFromFile(GetRuntimePath(), cfg)
		if err != nil {
			fmt.Println("Failed initialize runtime: ", err)
			os.Exit(1)
		}
		rtm = r
	}

	// Run requested test function
	switch function {

	// test crypto api
	case "test_blake2_128",
	     "test_blake2_256",
	     "test_twox_64",
	     "test_twox_128",
	     "test_twox_256",
	     "test_keccak_256":
		test_crypto_hash(rtm, function[5:], inputs[0])
	//case "test_blake2_256_enumerated_trie_root":
	//case "test_ed25519":
	//case "test_sr25519":
	//case "test_secp256k1_ecdsa_recover":

	// test storage api
	//case "test_clear_prefix":
	//case "test_clear_storage":
	//case "test_exists_storage":
	//case "test_set_get_local_storage":
	case "test_set_get_storage":
		test_set_get_storage(rtm, inputs[0], inputs[1])
	//case "test_set_get_storage_into":
	//case "test_storage_root":
	//case "test_storage_changes_root":
	//case "test_local_storage_compare_and_set":

	// test child storage api
	//case "test_clear_child_prefix":
	//case "test_clear_child_storage":
	//case "test_exists_child_storage":
	//case "test_kill_child_storage":
	//case "test_set_get_child_storage":
	//case "test_get_child_storage_into":
	//case "test_child_storage_root":

	default:
		fmt.Println("Not implemented: ", function)
		os.Exit(C.EOPNOTSUPP)
	}

	os.Exit(0)
}
