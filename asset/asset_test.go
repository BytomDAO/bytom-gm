package asset

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	dbm "github.com/tendermint/tmlibs/db"

	"github.com/bytom-gm/consensus"
	"github.com/bytom-gm/crypto/sm2/chainkd"
	"github.com/bytom-gm/database/leveldb"
	"github.com/bytom-gm/protocol"
	"github.com/bytom-gm/testutil"
)

func TestDefineAssetWithLowercase(t *testing.T) {
	reg := mockNewRegistry(t)
	alias := "lower"
	asset, err := reg.Define([]chainkd.XPub{testutil.TestXPub}, 1, nil, alias, nil)
	if err != nil {
		t.Fatal(err)
	}
	if *asset.Alias != strings.ToUpper(alias) {
		t.Fatal("created asset alias should be uppercase")
	}
}

func TestDefineAssetWithSpaceTrimed(t *testing.T) {
	reg := mockNewRegistry(t)
	alias := " WITH SPACE "
	asset, err := reg.Define([]chainkd.XPub{testutil.TestXPub}, 1, nil, alias, nil)
	if err != nil {
		t.Fatal(err)
	}
	if *asset.Alias != strings.TrimSpace(alias) {
		t.Fatal("created asset alias should be uppercase")
	}
}

func TestDefineAsset(t *testing.T) {
	ctx := context.Background()
	reg := mockNewRegistry(t)
	asset, err := reg.Define([]chainkd.XPub{testutil.TestXPub}, 1, nil, "asset-alias", nil)
	if err != nil {
		testutil.FatalErr(t, err)
	}

	found, err := reg.FindByID(ctx, &asset.AssetID)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if !testutil.DeepEqual(asset, found) {
		t.Errorf("expected asset %v to be recorded as %v", asset, found)
	}
}

func TestDefineBtmAsset(t *testing.T) {
	reg := mockNewRegistry(t)
	_, err := reg.Define([]chainkd.XPub{testutil.TestXPub}, 1, nil, consensus.BTMAlias, nil)
	if err == nil {
		testutil.FatalErr(t, err)
	}
}

func TestFindAssetByID(t *testing.T) {
	ctx := context.Background()
	reg := mockNewRegistry(t)
	keys := []chainkd.XPub{testutil.TestXPub}
	asset, err := reg.Define(keys, 1, nil, "TESTASSET", nil)
	if err != nil {
		testutil.FatalErr(t, err)

	}
	found, err := reg.FindByID(ctx, &asset.AssetID)
	if err != nil {
		testutil.FatalErr(t, err)
	}

	if !testutil.DeepEqual(asset, found) {
		t.Errorf("expected %v and %v to match", asset, found)
	}
}

func TestUpdateAssetAlias(t *testing.T) {
	reg := mockNewRegistry(t)

	oldAlias := "OLD_ALIAS"
	newAlias := "NEW_ALIAS"

	asset, err := reg.Define([]chainkd.XPub{testutil.TestXPub}, 1, nil, oldAlias, nil)
	if err != nil {
		testutil.FatalErr(t, err)
	}

	if reg.UpdateAssetAlias(asset.AssetID.String(), newAlias) != nil {
		testutil.FatalErr(t, err)
	}

	asset1, err := reg.FindByAlias(newAlias)
	if err != nil {
		testutil.FatalErr(t, err)
	}

	gotAlias := *asset1.Alias
	if !reflect.DeepEqual(gotAlias, newAlias) {
		t.Fatalf("alias:\ngot:  %v\nwant: %v", gotAlias, newAlias)
	}
}

func TestListAssets(t *testing.T) {
	reg := mockNewRegistry(t)

	firstAlias := "FIRST_ALIAS"
	secondAlias := "SECOND_ALIAS"

	firstAsset, err := reg.Define([]chainkd.XPub{testutil.TestXPub}, 1, nil, firstAlias, nil)
	if err != nil {
		testutil.FatalErr(t, err)
	}

	secondAsset, err := reg.Define([]chainkd.XPub{testutil.TestXPub}, 1, nil, secondAlias, nil)
	if err != nil {
		testutil.FatalErr(t, err)
	}

	wantAssets := []*Asset{DefaultNativeAsset, firstAsset, secondAsset}

	gotAssets, err := reg.ListAssets("")
	if err != nil {
		testutil.FatalErr(t, err)
	}

	if !testutil.DeepEqual(gotAssets, wantAssets) {
		t.Fatalf("got:\ngot:  %v\nwant: %v", gotAssets, wantAssets)
	}
}

func mockChain(testDB dbm.DB) (*protocol.Chain, error) {
	store := leveldb.NewStore(testDB)
	txPool := protocol.NewTxPool(store)
	chain, err := protocol.NewChain(store, txPool)
	if err != nil {
		return nil, err
	}
	return chain, nil
}

func mockNewRegistry(t *testing.T) *Registry {
	dirPath, err := ioutil.TempDir(".", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirPath)

	testDB := dbm.NewDB("testdb", "leveldb", "temp")
	defer os.RemoveAll("temp")

	chain, err := mockChain(testDB)
	if err != nil {
		t.Fatal(err)
	}

	return NewRegistry(testDB, chain)
}
