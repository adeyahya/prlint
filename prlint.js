const path = require("node:path");
const yaml = require("js-yaml");
const c = require("ansi-colors");
const { minimatch } = require("minimatch");
const { spawn } = require("node:child_process");
const { readFile } = require("node:fs/promises");

/**
 * 
 * @param {string} source source branch
 * @param {string} target target branch
 * @returns {Promise<string>}
 */
const getDiffFiles = (source, target) => new Promise((resolve, reject) => {
  let isError = false;
  let result = '';

  const gitDiff = spawn("git", ["diff", "--name-only", `${target}...${source}`], {
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
    if (!isError) return resolve(result.split("\n"));
    reject(result);
  })
});

(async () => {
  try {
    const configStr = await readFile(path.resolve(process.cwd(), ".prlint"), { encoding: "utf-8" });
    const config = yaml.load(configStr);
    const diffFiles = await getDiffFiles("development", "master");
    for (const key in config) {
      const rule = config[key];
      if (rule.files && !Array.isArray(rule.files)) throw new Error(`Invalid config, files should be an array in "${key}"`);
      if (!Array.isArray(rule.rules)) throw new Error(`Invalid config, rules should be an array in "${key}"`);
      if (typeof rule.eval !== "string") throw new Error(`Invalid config, eval should a string "${key}"`);
      const description = rule.description || "";

      const isRuleMatch = !rule.files ? true : diffFiles.some(fname => {
        return rule.files.some(r => minimatch(fname, r));
      });

      if (isRuleMatch) {
        const isValid = rule.rules.some(ruleRe => {
          const re = new RegExp(ruleRe);
          const targetEvaluation = process.env[`${rule.eval}`];
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
