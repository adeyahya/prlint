const path = require("node:path");
const yaml = require("js-yaml");
const c = require("ansi-colors");
const { minimatch } = require("minimatch");
const { spawn } = require("node:child_process");
const { readFile } = require("node:fs/promises");

const argv = process.argv ?? [];
const isDebug = argv.some(n => n === "--debug");

/**
 * 
 * @param {string} target target branch
 * @param {string} source source branch
 * @returns {Promise<string>}
 */
const getDiffFiles = (target, source = "HEAD") => new Promise((resolve, reject) => {
  if (isDebug) console.log("branch", "source:" + source, "target:" + target);
  let isError = false;
  let result = '';

  const gitDiff = spawn("git", ["diff", "--name-only", `${source}...${target}`], {
    cwd: process.cwd()
  })

  gitDiff.stdout.on("data", (data) => {
    result += data;
  })
  gitDiff.stderr.on("data", (data) => {
    isError = true;
    result += data;
  })
  gitDiff.on("close", () => {
    if (isDebug) console.log("diff", JSON.stringify(result, null, 2))
    if (!isError) return resolve(result.split("\n"));
    reject(new Error(result));
  })
});

(async () => {
  try {
    const configStr = await readFile(path.resolve(process.cwd(), ".prlint"), { encoding: "utf-8" });
    const config = yaml.load(configStr);
    if (isDebug) console.log("config", JSON.stringify(config, null, 2));
    const diffFiles = await getDiffFiles(process.env.TARGET_BRANCH, process.env.SOURCE_BRANCH);
    for (const key in config) {
      const rule = config[key];
      if (rule.files && !Array.isArray(rule.files)) throw new Error(`Invalid config, files should be an array in "${key}"`);
      if (!Array.isArray(rule.rules)) throw new Error(`Invalid config, rules should be an array in "${key}"`);
      if (typeof rule.envar !== "string") throw new Error(`Invalid config, eval should be a string "${key}"`);
      const description = rule.description || "";

      const isRuleMatch = !rule.files ? true : diffFiles.some(fname => {
        return rule.files.some(r => minimatch(fname, r));
      });

      if (isRuleMatch) {
        const isValid = rule.rules.some(ruleRe => {
          const re = new RegExp(ruleRe, "i");
          const targetEvaluation = process.env[`${rule.envar}`];
          return re.exec(targetEvaluation) !== null;
        })
        if (!isValid)
          throw new Error(`${c.red(key)} ${c.yellow(rule.description ? "(" + rule.description + ")" : "")}`);
      }
    }
    process.exit(0);
  } catch (error) {
    console.log(error.message);
    process.exit(1);
  }
})()
