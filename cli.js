#!/usr/bin/env node
const { spawn } = require("child_process");
const { join } = require("path");

const bin = join(__dirname, "resolvebench");
spawn(bin, process.argv.slice(2), { stdio: "inherit" }).on("exit", process.exit);
