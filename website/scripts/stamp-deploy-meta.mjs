import { execFile } from "node:child_process";
import { mkdir, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

// Stamp the deploy build with its commit + build time so the post-deploy
// smoke can assert the CDN serves THIS deploy (freshness), not the previous
// one riding out its propagation window. Written to public/ so it ships as
// /deploy-meta.json. Usage: node scripts/stamp-deploy-meta.mjs [commit]
// (CI passes $GITHUB_SHA; local runs fall back to `git rev-parse HEAD`).

const root = import.meta.dirname;
const outPath = join(root, "..", "public", "deploy-meta.json");

async function resolveCommit() {
	const [commit] = process.argv.slice(2);
	if (commit) {
		return commit;
	}

	try {
		const { stdout } = await execFileAsync("git", ["rev-parse", "HEAD"], {
			cwd: join(root, ".."),
		});
		return stdout.trim();
	} catch {
		return "unknown";
	}
}

async function main() {
	const commit = await resolveCommit();
	const meta = {
		built_at: new Date().toISOString(),
		commit,
	};

	await mkdir(join(root, "..", "public"), { recursive: true });
	await writeFile(outPath, `${JSON.stringify(meta, null, 2)}\n`);
	console.log(`deploy-meta: commit ${meta.commit}`);
}

await main();
