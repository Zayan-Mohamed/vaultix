# Vaultix v2.0.1 - Disk Space Optimization

This is a patch release fixing an important issue causing "no disk space left" errors when encrypting or decrypting extremely large folders. 

## 🔧 Bug Fixes
- **Disk Space Exhaustion on Init:** Optimized `Initialize` (`cd directory && vaultix init`) to securely delete each original unencrypted file immediately after it has been encrypted, rather than waiting for all files to be processed. This vastly reduces the peak disk space required.
- **Disk Space Exhaustion on DropAllFiles:** Optimized `DropAllFiles` to immediately delete the encrypted vault object as soon as the corresponding file is successfully decrypted and restored, keeping disk space overhead near zero for massive unencryption operations. 

## 🔄 Upgrading
Since this release does not change the core encryption schema or metadata structures introduced in v2.0.0, upgrading is simple: just replace your binary. Existing `v2.x` vaults are fully compatible with this release.