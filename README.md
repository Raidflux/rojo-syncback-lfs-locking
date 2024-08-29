# Getting started

## Installation

```bash
rokit add Raidflux/rojo-syncback-lfs-locking syncback-lfs
```

## Usage

```bash
syncback-lfs watch --input project.rbxl
```

## ⚠️ Important notes

Wire your brain to Ctrl+S to save your changes in Roblox Studio as much as possible, as this will trigger a syncback and a request to lock the file. If a lock is not granted, the syncback will not proceed and you will have to reconnect Rojo to revert your changes. Otherwise, future changes can't be synced back.
