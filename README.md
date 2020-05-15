# Perun mobile bindings
This project provides Android bindings for [go-perun](https://github.com/perun-network/go-perun) called *prnm*.  
Right now it only provides on two-party-payment channels.  

## Security Disclaimer
The authors take no responsibility for any loss of digital assets or other damage caused by the use of this software.  
**Do not use this software with real funds**.

### Getting Started
```sh
# Get perun-eth-mobile
go get -d https://github.com/perun-network/perun-eth-mobile
# Install gomobile, see https://godoc.org/golang.org/x/mobile/cmd/gomobile
go get golang.org/x/mobile/cmd/gomobile
gomobile init
# Generate the bindings
gomobile bind -o prnm.aar -target=android github.com/perun-network/perun-eth-mobile/
```

`prnm.aar` can then be included with Android studio.  
The easiest way of getting started with an app is to try out the *go-mobile* [example app](https://github.com/golang/go/wiki/Mobile#sdk-applications-and-generating-bindings) and instead of using the `hello.aar` replace it with `prnm.aar` that you generated with the command above.

### In Java
Go to the `MainActivity.java` of your app and import `prnm.*`. Just to see whether or not the compilation and linking is working, you can try to call `Prnm.contextBackground()` which should return a non-null object.  
A sample setup for a two-party-payment channel could look like this:  
```java
// You can add a sleep here to ensure that the Android studio logger is attached.
// Thread.sleep(2000);
// This seems to be the correct path for persistent files.
String appDir = getApplicationContext().getFilesDir().getAbsolutePath();
String dbPath = appDir +"/database";
// 10.0.2.2 is the IP of the host PC when using Android Simulator and the host is running a ganache-cli.
// 8545 is the standart port of ganache-cli.
String ethUrl = "ws://10.0.2.2:8545";
// Create a wallet.
Wallet wallet = new Wallet(appDir +"/wallet", "password");
// We are alice in this example and this is our on-chain secret key holding the ETH.
Address alice = wallet.importAccount("0x69cb97043e56883d66627e8f7a828877a56022d0fb05ae6197e6e16fb56282d0");
// Listen on 127.0.0.1:5750 for channel Proposals.
Config cfg = new Config("Alice ", alice, dbPath, ethUrl, "127.0.0.1", 5750);
// Create a client, currently skipping the ProposalHandler.
Client client = new Client(cfg, wallet);
// PerunId (currently an Address) of the peer that we want to open a channel with.
Address bob = new Address("0xA298Fc05bccff341f340a11FffA30567a00e651f");
// Tell the client what `bob`s ip and port are.
client.addPeer(bob, "10.0.2.2", 5750);
// Create the initial balances of the channel, me starting with 2000 and bob with 1000.
BigInts initBals = Prnm.newBalances(new BigInt(2000), new BigInt(1000));
// The ongoing channel proposal can be cancelled with ctx.cancel().
Context ctx = Prnm.contextWithCancel();
// Propose a channel to `bob` with `initBals` and 60s challenge duration.
PaymentChannel channel = client.proposeChannel(ctx, bob, 60, initBals);
// ctx.cancel() must always be called or it will leak resources.
// Better put it in the `finally` block of the surounding `try`.
ctx.cancel();
```  
This code must be wrapped in a `try` as the compiler will let you know.  
You can always `Ctrl+Click` on a *prnm* function to use the Android Studio decompiler, wich is really helpful to see all available Java functions.

The implementation of incoming channel proposals and channel updates will be added shortly.  


### Android App permissions
The `AndroidManifest.xml` file must contain at least the following permissions for the app.  
```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
<uses-permission android:name="android.permission.WRITE_EXTERNAL_STORAGE" />
<uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" />
```
Place this above the `<application>` section.

## Copyright
Copyright &copy; 2020 Chair of Applied Cryptography, Technische Universität Darmstadt, Germany.
All rights reserved.
Use of the source code is governed by the Apache 2.0 license that can be found in the [LICENSE file](LICENSE).

Contact us at [info@perun.network](mailto:info@perun.network).