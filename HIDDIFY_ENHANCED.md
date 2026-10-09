# Hiddify Enhanced

Hiddify Enhanced is a side-by-side desktop fork of [Hiddify App](https://github.com/hiddify/hiddify-app). It keeps the upstream Flutter interface, subscriptions, profile management, and proxy/TUN controls while using a distinct application identity on macOS and Windows.

## PasarGuard protocol support

The pinned Hiddify Core source includes Hiddify's sing-box fork and the `ray2sing` subscription-link converter. The fork has a VLESS client implementation for the `mlkem768x25519plus` encryption format (hybrid ML-KEM-768 and X25519) and XHTTP transport. The converter carries the VLESS `encryption` value and XHTTP mode, host, and path into the generated sing-box outbound. This does not require routing these links through Xray.

I fixed a parser edge case in the pinned sing-box fork: padding segments such as `100-111-1111` can decode as Base64, so the old parser mistook them for key material and rejected the link. The parser now recognizes only supported key sizes as keys, preserves earlier segments as padding, and still rejects invalid key sizes after key parsing starts.

The `ray2sing` regression test `TestVlessEncryptionWithXHTTP` checks a link containing XHTTP and the Base64-like padding segment. A Hiddify Core config test also runs that profile through sing-box option validation and outbound initialization. The formerly skipped VLESS parser test for Base64-like padding is now enabled and passes. These tests do not connect to a live PasarGuard server.

Subscriptions, profile refresh, and the existing TUN controls remain in the upstream Flutter app. TUN operation has not been tested on real macOS or Windows machines in this environment.

## Separate app identity

- macOS bundle identifier: `app.hiddify.enhanced`
- macOS app name and URL scheme: `Hiddify Enhanced` / `hiddify-enhanced`
- Windows executable and mutex: `HiddifyEnhanced.exe` / `HiddifyEnhancedMutex`
- Windows install directory and package identity are distinct from the official app
- Windows TUN service, local control port, and interface: `HiddifyEnhancedTunnelService`, `18021`, and `HiddifyEnhTun`

The installer cleanup only targets the Enhanced process and data directory. The MSIX manifest uses a development publisher identity; producing a distributable signed MSIX requires a signing certificate whose publisher matches that value.

Update feeds are unset in this source snapshot, so it will not offer the official app's updates as updates for Enhanced. Set `release_api_url` and `appcast_url` Dart defines to the fork's own release endpoints when publishing builds.

## Build

This source tree pins Flutter 3.38.5 or newer as permitted by `pubspec.yaml` and Go 1.26.3 for Hiddify Core. Build from the target operating system with the platform's SDKs installed. The upstream Makefile contains the dependency and package steps:

```sh
# macOS Apple Silicon
make macos-prepare
make macos-release

# Windows x64, from Windows with Flutter, Visual Studio C++, and packaging tools
make windows-prepare
make windows-release
```

The release targets package files under `dist/`. macOS distribution signing/notarization requires an Apple Developer identity. Windows MSIX signing requires the matching certificate. Run and verify each package on the actual destination operating system before distributing it.

For a cloud macOS build without installing Xcode locally, the source includes a manually triggered GitHub Actions workflow at `.github/workflows/build-enhanced-macos.yml`. Publish this source tree to a GitHub repository, then choose **Actions → Build Hiddify Enhanced for macOS → Run workflow**. A public repository uses standard GitHub-hosted runners at no charge. The workflow builds the patched universal Hiddify Core and uploads unsigned DMG/PKG artifacts; it has not yet been run in this source snapshot.

## Build status for this source snapshot

- Passed: `go test ./ray2sing_test -run 'TestVless' -count=1` in `hiddify-core/ray2sing`, including the new VLESS Encryption + XHTTP conversion regression test.
- Passed: full Hiddify sing-box VLESS package tests, including parser, Base64-like padding, and encryption handshake tests.
- Passed: `go test ./v2/config -run 'TestPanelURIsConvert/vless_encryption_xhttp$' -count=1` in `hiddify-core`.
- Passed: the complete Hiddify Core panel URI conversion and field-preservation tests with the desktop build tags from its Makefile.
- Passed: Go compilation of the shared tunnel-service package on this Mac (it has no package-local tests); this does not exercise the Windows service runtime.
- Not built: the desktop app or installers. Flutter/Dart are not installed, and this Mac only has Xcode Command Line Tools, not full Xcode. Windows packaging and runtime testing require a Windows host.
- Not verified: live PasarGuard connectivity, TUN behavior, macOS signing/notarization, Windows signing, or system-level coexistence with an installed Hiddify release.

## Source and license

The app source is based on the official Hiddify App repository. Hiddify Core, its sing-box fork, and `ray2sing` are included as pinned Git submodules. Keep their license and attribution files with any redistributed copy. See `LICENSE.md` and the license files in each submodule.
