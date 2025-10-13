import { $ } from 'bun';

const currentDir = process.cwd();
console.info('当前执行命令的目录:', currentDir);
$.cwd(currentDir);

try {
  await $`git clean -fdX`;
} catch (error) {
  console.error(`rm err code: ${error.exitCode}`);
  console.info(error.stdout.toString());
  console.info(error.stderr.toString());
  process.exit(1);
}
process.exit(0);
